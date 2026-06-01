## Context

当前系统已实现用户管理、角色管理和权限管理功能，但缺少部门管理模块。为完善组织架构管理能力，需要实现部门的树形层级管理功能，支持用户与部门的关联。

## Goals / Non-Goals

**Goals:**
- 实现部门树形层级结构管理
- 支持部门的增删改查操作
- 用户与部门建立一对一关联
- 部门设置唯一负责人
- 部门列表以树形结构展示

**Non-Goals:**
- 部门编码管理（用户明确不需要）
- 部门列表分页功能（树形结构展示，不需要分页）
- 用户多部门归属（一个用户只属于一个部门）

## Decisions

### 1. 数据模型设计

**部门模型 (Department)**
- 表名：`departments`
- 字段：
  - `id`: 主键ID
  - `parent_id`: 父部门ID（支持树形结构）
  - `name`: 部门名称（必填，唯一）
  - `leader_id`: 部门负责人ID（关联用户表）
  - `remark`: 部门备注（可选）
  - `status`: 状态（0-禁用，1-启用）
  - `order`: 排序字段
  - `created_at`, `updated_at`, `deleted_at`: 时间戳

**用户模型修改 (User)**
- 新增字段：`department_id`: 所属部门ID

### 2. 数据库关系

```
┌─────────────────────┐          ┌─────────────────────┐
│     departments     │          │        users        │
├─────────────────────┤          ├─────────────────────┤
│ id (PK)             │◄───────►│ id (PK)             │
│ parent_id (FK) ─────┼───┐     │ department_id (FK)  │
│ name                │   │     │ leader_of (virtual) │
│ leader_id (FK) ─────┼───┘     │ ...                 │
│ remark              │          └─────────────────────┘
│ status              │
│ order               │
└─────────────────────┘
```

### 3. API接口设计

| 接口 | 方法 | 路径 | 描述 |
|------|------|------|------|
| 创建部门 | POST | `/api/department` | 创建新部门 |
| 更新部门 | PUT | `/api/department/{id}` | 更新部门信息 |
| 删除部门 | DELETE | `/api/department/{id}` | 删除部门（级联删除子部门） |
| 获取部门详情 | GET | `/api/department/{id}` | 获取单个部门详情 |
| 获取部门树 | GET | `/api/department/tree` | 获取部门树形结构列表 |
| 获取部门用户 | GET | `/api/department/{id}/users` | 获取部门下的用户列表 |
| 分配负责人 | PUT | `/api/department/{id}/leader` | 设置部门负责人 |

### 4. 服务层设计

**DepartmentService接口**
- `Create(dto)`: 创建部门
- `Update(dto)`: 更新部门
- `Delete(id)`: 删除部门
- `GetByID(id)`: 获取部门详情
- `GetTree()`: 获取部门树形结构
- `GetUsers(id)`: 获取部门用户
- `SetLeader(id, leaderID)`: 设置负责人

### 5. 树形结构实现

采用递归查询方式构建树形结构：
1. 先查询所有部门（按parent_id和order排序）
2. 使用Map建立ID到部门的映射
3. 遍历所有部门，将子部门添加到父部门的Children字段
4. 返回根部门列表（parent_id=0或null）

## Risks / Trade-offs

| 风险 | 缓解措施 |
|------|----------|
| 删除部门时子部门处理 | 采用级联删除，删除父部门时同时删除所有子部门 |
| 部门负责人不存在 | 在设置负责人时验证用户是否存在 |
| 用户所属部门不存在 | 创建用户时验证部门ID有效性 |
| 树形结构深度过大影响性能 | 限制部门层级深度（建议不超过5层） |