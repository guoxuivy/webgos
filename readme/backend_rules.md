# Go 后端项目开发规则

本文档定义了后端项目的开发规范和约定，适用于 `internal/` 目录下的代码，所有后端代码必须遵循这些规则。

## 技术栈

Go 1.25 + Gin v1.11 + GORM v1.31 + PostgreSQL + JWT + Swagger

## 目录结构

```
internal/
├── dto/           # 数据传输对象（请求参数/响应数据）
├── handlers/      # HTTP处理器
├── middleware/    # 中间件（JWT/CORS/日志等）
├── models/        # 数据模型
├── routes/        # 路由注册
├── services/      # 业务逻辑
└── utils/         # 工具类（response/param/cache）
```

## 分层架构：Routes → Handlers → Services → Models → DTO

## 核心模式

### 1. BaseModel 泛型模型（Active Record）

```go
type BaseModel[T any] struct {
    ID int `gorm:"primarykey" json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    DeletedAt *time.Time `json:"deleted_at,omitempty"`
    queryHandler *gorm.DB `gorm:"-"`
}
```

**核心方法**：

- `Create(item)` / `Read(id)` / `Update(item)` / `Delete(id)`
- `More()` / `One()` / `Page(page, pageSize)` → `(items, total, err)`
- `Count()` / `Exist()`
- 链式查询：`Where()` / `Order()` / `Select()` / `Preload()` / `Limit()`
- 事务：`WithTx(tx)` / `WithCtx(ctx)`

### 2. Service 接口模式

```go
type ParkService interface {
    SavePark(ctx context.Context, park *models.Park) error
    GetParkByID(ctx context.Context, id int) (*models.Park, error)
    ParksPage(ctx context.Context, query dto.ParkQuery) ([]models.Park, int)
}

type parkService struct{}
func NewParkService() ParkService { return &parkService{} }
func (s *parkService) SavePark(ctx context.Context, park *models.Park) error { ... }
```

### 3. DTO 模式

```go
// 请求DTO - 字段用指针类型表示可选
type ParkQuery struct {
    Page     int     `form:"page" validate:"omitempty,min=1"`
    Name     *string `form:"name" validate:"omitempty,min=2,max=100"`
    Status   *int64  `form:"status" validate:"omitempty"`
}

// 创建DTO
type ParkCreate struct {
    Name   string `form:"name" validate:"required,min=2,max=100"`
    Status int64  `form:"status" validate:"required,oneof=1 2 3"`
    func ToModel() *models.Park { return &models.Park{Name: p.Name, Status: int(p.Status)} }
}

// 响应DTO - 可扩展字段
type ParkResp struct { models.Park; EnterpriseCount int `json:"enterpriseCount"` }
```

## 统一响应

```go
response.Success(c, "消息", data)          // {"code":0,"message":"消息","data":...}
response.Error(c, "错误信息")               // {"code":1,"message":"错误信息"}
response.AuthError(c, "未认证")             // HTTP 401
response.Forbidden(c, "无权限")              // HTTP 403
response.ErrorWithCode(c, "消息", http.StatusNotFound)  // 自定义状态码
```

## ⚠️ 分页接口返回格式（强制要求）

**所有分页列表接口必须使用 `items` + `total` 格式！**

### ✅ 正确格式

```go
func XxxList(c *gin.Context) {
    var query dto.XxxQuery
    if err := param.Validate(c, &query); err != nil {
        response.Error(c, err.Error())
        return
    }

    items, total := services.NewXxxService().XxxsPage(c, query)

    response.Success(c, "获取成功", gin.H{
        "items": items,  // ✅ 必须使用 items
        "total": total,  // ✅ 必须使用 total
    })
}
```

**返回 JSON 结构**：

```json
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "items": [...],
    "total": 100
  }
}
```

### ❌ 禁止格式

- 使用 `list`、`data`、`rows` 代替 `items`
- 使用 `count`、`sum` 代替 `total`

## 参数验证标签

| 标签 | 说明 | 示例 |
| --- | --- | --- |
| `required` | 必填 | `validate:"required"` |
| `omitempty` | 允许为空 | `validate:"omitempty,min=3"` |
| `min/max` | 长度/值范围 | `validate:"min=3,max=100"` |
| `gte/lte` | 大小等于 | `validate:"gte=18,lte=100"` |
| `email/url` | 格式验证 | `validate:"email"` |
| `oneof` | 枚举值 | `validate:"oneof=1 2 3"` |
| `phone` | 手机号 | `validate:"phone"` |
| `label` | 字段中文名 | `label:用户名` |

## 新增CRUD模块清单

```
□ internal/models/xxx.go       # BaseModel + 字段 + Serialize/Deserialize
□ internal/dto/xxx.go          # Query/Create + ToModel()
□ internal/services/xxx.go     # 接口定义 + 实现
□ internal/handlers/xxx.go     # Validate + Service调用 + Response
□ internal/routes/xxx.go       # init() + WrapRouter + CRUD路由
```

## 注意事项

1. **DTO可选字段用指针类型**：`*string`, `*int64`
2. **Service方法始终接收context.Context**
3. **JSONB字段**：保存前调`Serialize()`，读取后调`Deserialize()`
4. **增删改后清除缓存**：`cache.GetCache().DeleteByPrefix(...)`
5. **错误统一用response返回**，不要直接return
6. **不要修改BaseModel**，继承即可
7. **软删除默认开启**，物理删除用`Unscoped()`
8. **Swagger注解**放在Handler函数上方

## 常见错误

1. **分页接口返回格式错误**：使用 `list` 代替 `items`
2. **缺少 Preload**：关联数据为空
3. **Service层缺少 WithCtx(ctx)**：上下文丢失
4. **Handler中重复设置UserID**：应从context获取