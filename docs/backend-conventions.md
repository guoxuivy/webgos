# webgos 后端开发规范（唯一权威源）

> 本文件是 `webgos` 后端代码的**唯一权威规范**，适用于 `internal/` 目录下的全部代码。
>
> - 历史入口 `readme/backend_rules.md` 已迁移至此，原文仅保留存根指向本文件。
> - [`docs/templates/crud-api.md`](templates/crud-api.md) 是**代码生成模板**，不定义规范；模板与本文件冲突时，**一律以本文件为准**。
> - 修改本文件即修改团队约定，请谨慎；代码与本文件冲突时，改代码。

---

## 1. 技术栈（不可协商）

| 层 | 选型 |
| --- | --- |
| 语言 | Go 1.25 |
| Web 框架 | Gin v1.11 |
| ORM | GORM v1.31（原生 API，不使用泛型 BaseModel 封装） |
| 数据库 | MySQL / PostgreSQL（`go.mod` 同时引入两个驱动，实际由配置决定） |
| 认证 | golang-jwt/jwt/v4 |
| 参数校验 | go-playground/validator/v10 |
| API 文档 | Swaggo（swag + gin-swagger） |
| 缓存 | patrickmn/go-cache |
| 测试 | testify |

引入任何新依赖前，先确认现有工具/库已无法满足，并在提交信息中说明理由。

---

## 2. 目录结构

```
internal/
├── bootstrap/     # 项目启动初始化
├── cache/         # 缓存封装
├── config/        # 配置加载与校验
├── dto/           # 数据传输对象（请求参数 / 响应数据）
├── handlers/      # HTTP 处理器（参数校验 → 调 Service → 统一响应）
├── middleware/    # 中间件（JWT / CORS / 日志 / 限流 / 安全等）
├── models/        # 数据模型
├── routes/        # 路由注册（自动同步为 RBAC 权限点）
├── services/      # 业务逻辑（事务边界、DB 操作）
├── utils/         # 工具类（code 业务码 / param 校验 / response 响应）
├── xdb/           # 数据库句柄（主库 GetDB / 从库 GetSlaveDB）
└── xlog/          # 日志
```

Swagger 生成物位于 `internal/swagger/`（由 `make swagger` 生成，**勿手改**）。
手写规范文档位于 `docs/`，与生成物物理隔离。

---

## 3. 分层架构

```
Routes → Handlers → Services → Models
              ↘ DTO ↗
```

- **Routes**：注册路由与中间件，路由自动同步为权限点。
- **Handlers**：只做参数绑定与校验、调用 Service、输出统一响应；**不写业务逻辑、不直接操作 DB**。
- **Services**：业务逻辑、事务边界、数据库操作；**不引用 `*gin.Context`**。
- **Models**：数据结构，嵌入 `BaseFields`，复杂字段 JSONB + 手动序列化。
- **DTO**：入参与出参结构及校验标签。

---

## 4. 中间件

### 全局中间件（`middleware.ApplyMiddlewares` 注册，**顺序固定**）

`IPBlacklist` → `RequestID` → `Recovery` → `Logging` → `CORS` → `Gzip`

### 安全防范

- `CheckSensitivePath` 在 404 handler 中调用，检测 `.env` / `.git` / `phpmyadmin` 等敏感路径与 `/shell`、`/exec` 等危险关键字。
- **1 小时内命中 5 次**自动将该 IP 加入黑名单，持久化到 `blacklist.json`。

### IP 限流

- `IPLimiter(rate, capacity)` 基于令牌桶，每个 IP 独立限流，用于登录等高风险路由防爆破（如 `IPLimiter(1, 1)`）。

### 路由分组中间件

| 中间件 | 作用 |
| --- | --- |
| `JWT` | 登录态校验 |
| `Auth` | RBAC 权限验证 |
| `Debounce` | 防重复提交 |

---

## 5. 模型与数据库操作

### 模型：嵌入 `BaseFields`

```go
type Park struct {
    models.BaseFields
    Name   string `gorm:"size:100" json:"name"`
    Status int64  `json:"status"`
}
```

- 通过嵌入 `BaseFields` 获得 `ID` / `CreatedAt` / `UpdatedAt` / `DeletedAt`。
- **不要修改 `BaseFields`**，嵌入即可。
- ID 相关参数统一使用 `int` 类型。
- JSON 标签使用**小驼峰**（如 `landType`、`isActive`）。
- **软删除默认开启**；物理删除显式使用 `Unscoped()`。

### 数据库操作：Service 层直接用 GORM 原生 API

不再经过泛型 `BaseModel` 封装，统一通过 `ctxDB(ctx)` / `ctxSDB(ctx)` 获取带上下文的 `*gorm.DB`（定义见 `internal/services/helper.go`）：

- 查询：`ctxSDB(ctx).Where(...).Order(...).Find(&items)`
- 单条：`ctxSDB(ctx).First(&item, id)`
- 分页：`ctxSDB(ctx).Where(...).Count(&total)` 配合分页 Scope
- 创建：`ctxDB(ctx).Create(&item)`
- 更新：`ctxDB(ctx).Select("*").Updates(&item)`（默认仅更新非零值字段）
- 删除（软删除）：`ctxDB(ctx).Delete(&item, id)`
- 事务：`ctxDB(ctx).Transaction(func(tx *gorm.DB) error { ... })`

### Context 约定

- **所有 Service 方法第一个参数必须是 `ctx context.Context`**，用于请求级超时与取消。
- Handler 调用 Service 时**直接传 `c *gin.Context`** —— `gin.Context` 实现了 `context.Context` 接口，无需 `c.Request.Context()`。
- **Service 层只依赖标准库 `context.Context`**，不 import gin、不使用 `*gin.Context` 专属方法，保证可脱离 Web 框架单独测试。

---

## 6. 读写分离与选库（强制）

`internal/services/helper.go` 提供两个入口，**由 Service 自行决定**走哪个库，不做自动推断：

| 入口 | 含义 |
| --- | --- |
| `ctxDB(ctx)` | 主库（写库） |
| `ctxSDB(ctx)` | 从库（读库）；当 `xdb.GetSlaveDB()` 返回 `nil`（未开启读写分离）时**自动回退主库**，调用方无需判断 |

### 选库原则

| 场景 | 入口 |
| --- | --- |
| `Create` / `Update` / `Delete` / `Transaction` | `ctxDB` |
| 写前校验查询（先查父级 / 查重名再写） | `ctxDB` |
| 写操作之后紧跟读取自己刚写入的数据 | `ctxDB` |
| 强一致性读（支付状态、库存扣减前余额校验） | `ctxDB` |
| 普通只读查询（`First` / `Find` / `Count` / `Pluck` 等无写入） | `ctxSDB` |
| **事务闭包内** | **只用闭包参数 `tx`，禁止调用 `ctxDB` / `ctxSDB`** |

### 关键边界：从库复制延迟

1. **不要**在写入后立即依赖从库读取自身刚写入的数据，否则可能读到旧值或空值。
2. **不要**把「写前校验 / 强一致校验」放到从库，否则可能因延迟导致校验失效（如并发重复名、超卖）。
3. 未开启读写分离时 `ctxSDB` 自动回退主库，行为与 `ctxDB` 一致。

### 性能说明

`WithContext` 在 GORM 中无副作用（仅基于全局连接池复制实例并设置 ctx 字段），同一函数内多次调用 `ctxDB(ctx)` / `ctxSDB(ctx)` 没有性能问题，也可在开头取一次局部变量连续使用。

### 正确示例

```go
// 写操作 + 写后读自己 → 全程 ctxDB
func (s *inventoryService) ProductOut(ctx context.Context, record *models.InventoryRecord) error {
    db := ctxDB(ctx)
    var product models.Product
    if err := db.First(&product, record.ProductID).Error; err != nil { // 强一致读
        return err
    }
    if product.Stock < record.Quantity {
        return errors.New("库存不足")
    }
    if err := db.Create(record).Error; err != nil {
        return err
    }
    product.Stock -= record.Quantity
    return db.Select("*").Updates(&product).Error
}

// 纯查询 → ctxSDB（未开启分离时自动回退主库）
func (s *productService) GetProductByID(ctx context.Context, id int) (*models.Product, error) {
    var product models.Product
    err := ctxSDB(ctx).First(&product, id).Error
    return &product, err
}
```

---

## 7. Service 接口模式

```go
type ParkService interface {
    SavePark(ctx context.Context, park *models.Park) error
    GetParkByID(ctx context.Context, id int) (*models.Park, error)
    ParksPage(ctx context.Context, query dto.ParkQuery) ([]models.Park, int64)
    DeletePark(ctx context.Context, id int) error
}

type parkService struct{}

func NewParkService() ParkService { return &parkService{} }
```

- **接口化设计**：`XxxService` 接口 + `xxxService` 实现结构体。
- **工厂函数返回接口类型**：`NewXxxService() XxxService`。
- **方法命名按业务语义**，避免千篇一律的 Save/Get/Delete 前缀。
- **分页方法返回 `int64`**：GORM `Count()` 签名要求 `*int64`，统一使用 `int64`（`backend_rules.md` 旧示例中的 `int` 已废弃）。

---

## 8. DTO 模式

```go
// 请求 DTO —— 字段用指针类型表示可选，支持部分更新
type Park struct {
    ID     *int    `json:"id" validate:"omitempty,gte=0" label:"园区ID"`
    Name   *string `json:"name" validate:"omitempty,min=2,max=100" label:"园区名称"`
    Status *int    `json:"status" validate:"omitempty,oneof=1 2 3" label:"状态"`
}

// 查询参数 —— Page / PageSize 只能是 int，不能是指针
type ParkQuery struct {
    Page     int     `form:"page" validate:"omitempty,gte=1"`
    PageSize int     `form:"pageSize" validate:"omitempty,gte=1,lte=100"`
    Name     *string `form:"name" validate:"omitempty,max=100"`
    Status   *int    `form:"status" validate:"omitempty,oneof=1 2 3"`
}

// ToModel —— 仅做字段映射，不查库
func (dto *Park) ToModel() *models.Park { ... }
```

- **统一 DTO**：不区分 Create 与 Update，单个结构体同时支持创建与更新，靠 `ID` 是否为零值判断。
- **指针字段**：支持部分更新；查询参数也用指针判断条件是否存在。
- **`label` 标签**：给中文字段名，用于错误消息。
- 响应 DTO 可嵌入模型并扩展字段：`type ParkResp struct { models.Park; EnterpriseCount int \`json:"enterpriseCount"\` }`。

---

## 9. 统一响应

```go
response.Success(c, "消息", data)                     // {"code":0,"message":"消息","data":...}
response.Error(c, "错误信息")                          // {"code":1,"message":"错误信息"}
response.Unauthorized(c, "未认证")                      // HTTP 401
response.Forbidden(c, "无权限")                         // HTTP 403
response.ErrorWithCode(c, "消息", http.StatusNotFound)  // 自定义状态码
```

- **错误统一用 `response` 返回，不要直接 `return`**。
- 业务状态码定义在 `internal/utils/code/bizcode.go`。

---

## 10. 分页接口返回格式（强制）

**所有分页列表接口必须使用 `items` + `total` 格式！**

```go
func XxxList(c *gin.Context) {
    var query dto.XxxQuery
    if err := param.Validate(c, &query); err != nil {
        response.Error(c, err.Error())
        return
    }

    items, total := services.NewXxxService().XxxsPage(c, query)

    response.Success(c, "获取成功", gin.H{
        "items": items,  // 必须为 items
        "total": total,  // 必须为 total
    })
}
```

返回 JSON 结构：

```json
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "items": [],
    "total": 100
  }
}
```

**禁止**：用 `list`、`data`、`rows` 代替 `items`；用 `count`、`sum` 代替 `total`。

> 这是本仓库最常见的回归点。旧 `readme/SKILL.md`（现为 [`docs/templates/crud-api.md`](templates/crud-api.md)）模板曾写成 `list`，已修正——**生成代码后请复核此处**。

---

## 11. JSONB 字段序列化

复杂类型字段（位置信息、配置等）用 JSONB 存储，**手动序列化，不依赖 GORM 钩子**：

- 保存前在 `SaveXxx` 中调用 `entity.Serialize()`。
- 查询后在 `GetXxxByID` / `GetXxxPage` 中调用 `entity.Deserialize()`。
- 分页查询需**批量**反序列化；单条反序列化失败时 `continue`，不中断整批。

---

## 12. Swagger 注解与路由注册

### Swagger 注解放在 Handler 函数上方

```go
// XxxList 获取列表
// @Summary 实体列表
// @Description 获取实体列表接口
// @Tags 实体
// @Accept json
// @Produce json
// @Param body body dto.XxxQuery true "查询参数"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/xxx/list [post]
// @Security BearerAuth
func XxxList(c *gin.Context) { ... }
```

注解修改后执行 `make swagger` 重新生成（产物在 `internal/swagger/`，勿手改）。

### 路由注册

```go
func init() {
    Register(func(router *Router) {
        wr := router.Group("/api/xxx")
        wr.Use(middleware.JWT())
        {
            wr.POST("/list", "获取列表", handlers.XxxList)
            wr.GET("/:id", "详情", handlers.XxxDetail)
            wr.POST("/edit", "创建/更新", handlers.XxxEdit)
            wr.DELETE("/:id", "删除", handlers.XxxDelete)
        }
    })
}
```

- 使用 `init()` + `Register` 注册。
- **路由自动同步为 RBAC 权限点**，权限标识为 `路径:HTTP方法`（如 `/api/products:GET`）。新增路由即新增权限点，需在菜单/角色中授权后才能访问。

---

## 13. 参数验证标签

| 标签 | 说明 | 示例 |
| --- | --- | --- |
| `required` | 必填 | `validate:"required"` |
| `omitempty` | 允许为空 | `validate:"omitempty,min=3"` |
| `min/max` | 长度 / 值范围 | `validate:"min=3,max=100"` |
| `gte/lte` | 大于等于 / 小于等于 | `validate:"gte=18,lte=100"` |
| `email/url` | 格式验证 | `validate:"email"` |
| `oneof` | 枚举值 | `validate:"oneof=1 2 3"` |
| `phone` | 手机号 | `validate:"phone"` |
| `label` | 字段中文名（用于错误消息） | `label:用户名` |

统一使用 `param.Validate(c, &dto)` 校验。

---

## 14. 新增 CRUD 模块清单

```
□ internal/models/xxx.go       # BaseFields 嵌入 + 字段 + Serialize/Deserialize
□ internal/dto/xxx.go          # 实体 DTO（指针字段） + Query + ToModel()
□ internal/services/xxx.go     # 接口定义 + 实现（ctx 首参，ctxDB/ctxSDB 选库）
□ internal/handlers/xxx.go     # Validate + Service 调用 + Response（含 Swagger 注解）
□ internal/routes/xxx.go       # init() + Register + Group + CRUD 路由
```

---

## 15. 注意事项

1. **DTO 可选字段用指针类型**：`*string`、`*int64`。
2. **Service 方法始终接收 `context.Context`**（首参数），Handler 直接传 `c` 即可。
3. **JSONB 字段**：保存前调 `Serialize()`，读取后调 `Deserialize()`。
4. **增删改后清除缓存**：`cache.GetCache().DeleteByPrefix(...)`。
5. **错误统一用 `response` 返回**，不要直接 `return`。
6. **不要修改 `BaseFields`**，嵌入即可。
7. **软删除默认开启**，物理删除用 `Unscoped()`。
8. **Swagger 注解**放在 Handler 函数上方。
9. 用户 ID 等身份信息**从 context 获取**，不要在 Handler 中重复设置。

---

## 16. 常见错误

1. **分页返回格式错误**：用 `list` 代替 `items`（最高频）。
2. **缺少 `Preload`**：关联数据为空。
3. **Service 层缺少 `ctx` 参数**：上下文丢失，无法传递超时与取消。
4. **事务闭包内误用 `ctxDB` / `ctxSDB`**：脱离事务另开会话。
5. **写后立即从从库读**：因复制延迟读到旧值或空值。
6. **Handler 中重复设置 UserID**：应从 context 获取。

---

## 17. 项目初始化规范

- `cmd/main.go` **只负责程序入口与启动**，保持简洁：解析 `-c` 配置路径 → 调用
  `bootstrap.Initialize(*configPath)` → 构造 `http.Server` 并启动 → 优雅关闭。
- 初始化逻辑集中在 `internal/bootstrap`，`Initialize()` **接收配置路径参数**，依次完成：
  日志系统初始化 → 数据库连接初始化 → 模型自动迁移 → Gin 路由初始化 → 路由注册 → 权限点同步。
- 初始化后的 Gin 引擎通过 **`routes.REngine`** 获取（`main.go` 用它注册业务路由与 Swagger 路由）。
- 配置依赖**显式传递**，避免隐式全局状态；统一初始化错误处理，`Initialize` 返回 error 由 `main` 处理。
- 优雅关闭：`defer bootstrap.Close()`，并监听 `SIGINT` / `SIGTERM`，用带超时的 `context`
  调用 `srv.Shutdown`。
