# 权限点绑定菜单 —— 开发计划文档（Spec）

## 1. 背景与现状

当前 RBAC 系统采用经典的 **三层模型**：

```
User ──< rbac_user_roles >── Role ──< rbac_role_permissions >── Permission
                                     │
                                     └── MenuIDs (字符串 "3,4,5")
```

- 权限点（`RBACPermission`，含 `Path` + `Method`）**只绑定角色**，通过 `rbac_role_permissions` 中间表。
- 菜单（`Menu`）是**另一条独立链路**，挂在 `RBACRole.MenuIDs` 逗号分隔字符串上，与权限点平行、互不相干。
- `Menu` 表已存在 `auth_code` 字段，但目前仅作标记字符串，**未被任何逻辑使用**。
- 权限校验中间件 `middleware/rbac.go` 走 `perm.Path + ":" + perm.Method`，与 `auth_code` 无关。

**问题**：角色既要管菜单又要管权限点，双轨冗余，维护成本高，且无法从菜单自然推断出其包含的权限。

## 2. 目标

将架构改造为 **权限点归属于菜单（多对多），角色绑定菜单后间接获得权限**：

```
User ──< rbac_user_roles >── Role ──< rbac_role_menus >── Menu ──< rbac_menu_permissions >── RBACPermission
```

- 权限点与菜单是 **多对多** 关系：同一个 `RBACPermission` 可以同时绑定到多个菜单上（例如一个公共接口权限点被多个菜单共享）。
- 角色只绑定菜单，通过 `rbac_role_menus` 关联；角色经由菜单经由 `rbac_menu_permissions` 间接获得权限。

管理员只需给角色分配菜单，角色自动拥有所绑菜单下所有按钮/接口的权限。权限校验逻辑对前端透明，API 调用方式不变。

## 3. 改造范围

| 层 | 文件 | 改动 |
| --- | --- | --- |
| 模型 | `internal/models/rbac.go` | `RBACPermission` 加 `Menus []Menu`（many2many `rbac_menu_permissions`），**不使用单外键**；`RBACRole` 加 `Menus []Menu`（many2many `rbac_role_menus`）；移除 `RBACRolePermission` 结构体 |
| 模型 | `internal/models/menu.go` | `Menu` 加 `Permissions []RBACPermission`（many2many `rbac_menu_permissions`） |
| 服务 | `internal/services/rbac.go` | 移除 `AssignPermissionsToRole`，新增 `AssignMenusToRole`；`AddRole`/`EditRole`/`GetRoles` 使用菜单关联；`SyncPermissions` 仅同步权限点本身（创建/更新描述），不自动归属菜单 |
| 服务 | `internal/services/menu.go` | `GetUserMenus` 改用 GORM 关联替代 `GetMenuIDs()` 字符串解析；新增 `AssignPermissionsToMenu` 维护菜单-权限多对多 |
| 中间件 | `internal/middleware/rbac.go` | 校验链路改为 `user → roles → menus → permissions`，三层 Preload（含菜单-权限多对多） |
| DTO | `internal/dto/rbac.go` | `AssignPermissionsDTO`（role_id + perm_ids）改为 `AssignMenusDTO`（role_id + menu_ids）；新增 `AssignPermissionsToMenuDTO`（menu_id + perm_ids） |
| Handler | `internal/handlers/rbac.go` | 权限分配接口改为菜单分配；菜单编辑接口可返回/编辑其隶属的多个权限点 |
| 路由 | `internal/routes/system.go` | 调整/新增菜单分配相关路由；新增菜单-权限绑定路由 |
| 数据库 | `readme/sys_menus_init.sql` | 新增 `rbac_role_menus` 表；新增 `rbac_menu_permissions(menu_id, permission_id)` 表（取代原 `menu_id` 列方案） |
| 迁移 | 新增迁移脚本 | 三步迁移历史数据 |

## 4. 核心功能

- **角色绑定菜单**：通过 `rbac_role_menus` 规范化中间表，取代 `MenuIDs` 字符串。
- **权限点归属菜单（多对多）**：`RBACPermission` 与 `Menu` 通过 `rbac_menu_permissions(menu_id, permission_id)` 关联；**同一个权限点可绑定到多个菜单**，菜单也可包含多个权限点。
- **路由权限点同步**：启动时 `SyncPermissions` 仅将注册的路由同步为 `rbac_permissions` 权限点（创建缺失项、更新描述），**不做菜单归属**；菜单与权限点的绑定由 `AssignPermissionsToMenu` 显式维护（菜单 path 为前端路由、接口 path 为后端 API，二者无必然前缀关系，按前缀猜测会污染关联数据）。
- **用户权限校验链路**：`User → Role → Menu → Permission → 接口鉴权`。
- **超级管理员**：跳过权限检查，保持现有逻辑不变。
- **缓存**：键保持 `permissions:` + userID，内容仍为 `map[string]bool`，仅查询链路改变。

## 5. 技术实现方案

### 5.1 三层链路穿透查询

中间件中用 GORM 链式预加载一次性取出用户的所有菜单及权限点，避免 N+1。由于权限-菜单是多对多，Preload 走 `rbac_menu_permissions` 关联表：

```go
var user models.User
db.Preload("Roles.Menus.Permissions").First(&user, userID)
permissions := make(map[string]bool)
for _, role := range user.Roles {
    for _, menu := range role.Menus {
        for _, perm := range menu.Permissions {
            permissions[perm.Path+":"+perm.Method] = true
        }
    }
}
```

> 说明：GORM 的 `many2many` 标签会自动使用 `rbac_menu_permissions` 作为桥表，一个菜单的 `Permissions` 是其关联的所有权限点集合，天然支持同一权限出现在多个菜单下。

### 5.2 权限点菜单归属策略（多对多）

菜单与权限点的绑定**不**由 `SyncPermissions` 自动推断。`SyncPermissions` 仅负责将注册的路由同步为 `rbac_permissions` 权限点（创建缺失项、更新描述），因为菜单 `Path`（前端路由，如 `/system/dept`）与接口 `Path`（后端 API，如 `/api/department/tree`）没有必然的前缀关系，按前缀猜测会污染 `rbac_menu_permissions` 关联数据、导致 RBAC 权限判定失真。

归属关系统一由手动/前端维护入口 `AssignPermissionsToMenu(menuID, permIDs)` 控制，支持把同一个权限点批量绑定到多个菜单；调用前先 `Replace` 该菜单已有关联，幂等安全。

### 5.3 数据迁移策略

**Step 1 — 建新表结构**
- 创建 `rbac_role_menus(role_id, menu_id)` 中间表
- 创建 `rbac_menu_permissions(menu_id, permission_id)` 中间表（取代原 `rbac_permissions.menu_id` 列方案）

**Step 2 — 迁移历史数据**
- 解析每个 `RBACRole.MenuIDs` 字符串（如 `"3,4,5"`），拆解后逐条插入 `rbac_role_menus`
- 运行 `SyncPermissions` 仅同步权限点本身（`rbac_permissions`），菜单-权限绑定由 `AssignPermissionsToMenu` 单独维护

**Step 3 — 清理旧结构**
- 从 `RBACRole` 删除 `menu_ids` 列
- 删除 `rbac_role_permissions` 中间表（权限-角色直连关系不再使用）
- 从 AutoMigrate 中移除 `RBACRolePermission`

**可执行迁移脚本（已生成）**
- Go 一次性命令：`cmd/migrate_rbac/main.go`
  - 用法：`go run ./cmd/migrate_rbac -c ./config/config.yaml`
  - 流程：AutoMigrate 建中间表 → 注册路由 + `SyncPermissions` 仅同步权限点（`rbac_permissions`）→ 迁移 `rbac_roles.menu_ids` 到 `rbac_role_menus` → 删除 `rbac_role_permissions` 表与 `menu_ids` 列。幂等、可重入。菜单-权限绑定由 `AssignPermissionsToMenu` 维护。
- SQL 脚本：`internal/xdb/migrate/rbac_menu_migration.sql`
  - 负责建表 + 角色菜单数据搬运（含 MySQL/PostgreSQL 分支）+ 清理旧结构；权限-菜单归属由应用 `AssignPermissionsToMenu` 维护，不依赖 `SyncPermissions` 前缀匹配。

## 6. 关键决策与权衡

| 决策 | 选择 | 理由 |
| --- | --- | --- |
| 权限-菜单关系 | **多对多**（`rbac_menu_permissions` 中间表） | 同一权限点需可绑定到多个菜单，单外键 `menu_id` 无法满足 |
| 权限点归属方式 | 由 `AssignPermissionsToMenu` 显式维护，不依赖路径前缀自动匹配 | 菜单 path（前端路由）与接口 path（后端 API）无必然前缀关系，自动匹配会污染 `rbac_menu_permissions`；`SyncPermissions` 仅同步权限点本身 |
| 角色绑菜单方式 | `rbac_role_menus` 中间表 | 规范化多对多，优于逗号分隔字符串 |
| 是否保留 `RBACRolePermission` | 彻底移除 | 消除双轨冗余，角色权限来源唯一（经菜单间接获得） |
| 旧数据兼容 | 提供迁移脚本，不保留 `MenuIDs` 列 | 一次迁移干净，避免遗留技术债 |
| GORM 预加载 | `Preload("Roles.Menus.Permissions")` | 三层嵌套 Preload，多 LEFT JOIN 一次取齐，天然支持一权限多菜单 |

## 7. 实现注意事项

- **性能**：三层 Preload 生成多表 JOIN，但角色数（≤5）、菜单数（≤50）、权限数（≤200）均在可控范围，不构成瓶颈；缓存命中时完全跳过查询。
- **日志**：Preload 失败需记录错误但不暴露内部细节（返回通用"权限验证失败"），用 `xlog` 记录完整堆栈。
- **兼容性**：超级管理员跳过逻辑保持不变。
- **前端透明**：权限校验对前端透明，API 调用方式不变；角色管理页改为只分配菜单，菜单编辑页可展示隶属权限点。

## 8. 验收标准

1. 启动后权限点经 `rbac_menu_permissions` 正确归属菜单，且同一权限点可同时出现在多个菜单下；`rbac_permissions` 不再使用单 `menu_id` 列。
2. 为角色分配菜单后，用户登录即可访问该菜单下所有接口，无需单独分配权限点。
3. 未绑定菜单的角色无法访问其接口（超管除外）。
4. 旧 `MenuIDs` 数据完整迁移至 `rbac_role_menus`，原 `rbac_role_permissions` 逻辑不再使用。
5. 缓存与现有 5 分钟过期策略一致，权限变更后最多 5 分钟内生效。
