---
name: "backend-api-creator"
description: "自动创建完整的后端API接口，采用接口化设计模式，包括模型、DTO、服务接口与实现、处理器和路由。当用户需要创建新的后端接口或API端点时调用此技能。"
---

# 后端API创建器

这个技能帮助您快速创建完整的后端API接口，遵循项目的现有架构模式和代码规范。

## 功能特性

- 自动生成符合项目规范的完整API接口
- 采用接口化设计模式，实现服务层解耦
- 包含模型（Model）、数据传输对象（DTO）、服务接口与实现、处理器（Handler）和路由（Route）
- 支持RESTful API设计模式
- 自动添加Swagger文档注释
- 集成GORM数据库操作
- 包含输入验证和错误处理
- 支持依赖注入和单元测试

## 使用场景

当您需要创建新的后端接口时，调用此技能并提供以下信息：

1. **实体名称**（如：Product, User, Order等）
2. **实体字段**（字段名、类型、验证规则等）
3. **业务功能需求**（具体的业务操作，如：创建订单、查询用户列表、更新产品状态等）
4. **服务方法名称**（根据业务语义命名，如：CreateOrder、GetUserList、UpdateProductStatus等）

## 核心架构原则：职责分离

### 分层架构

| 层级 | 目录 | 职责 | 允许的操作 | 禁止的操作 |
|------|------|------|------------|------------|
| **Handlers** | `internal/handlers/` | 参数接收、参数验证、响应格式化 | 解析请求参数、调用服务层、返回响应 | 构造模型、数据库操作、业务逻辑判断 |
| **Services** | `internal/services/` | 业务逻辑、模型构造、数据库操作 | 构造模型、执行业务逻辑、调用数据库 | 直接访问请求上下文 |
| **Models** | `internal/models/` | 数据结构定义 | 定义字段、定义方法 | 包含业务逻辑 |
| **DTO** | `internal/dto/` | 数据传输对象、参数验证规则 | 定义字段、验证标签、转换方法 | 包含业务逻辑 |
| **Database** | `internal/xdb/` | 数据库连接管理 | 获取数据库连接 | 业务逻辑 |

### Handlers 层职责（重要）

**Handlers 层只负责以下工作：**

1. **接收参数** - 从请求中解析参数（路径参数、查询参数、请求体）
2. **参数格式验证** - 使用 `param.Validate` 进行参数格式验证（如必填字段、数据类型、范围限制）
3. **调用服务** - 将 DTO 传递给服务层处理
4. **响应输出** - 使用 `response.Success/Error` 返回结果

**Handlers 层禁止以下操作：**

❌ 构造模型对象（如 `&models.User{...}`）
❌ 直接访问数据库（如 `xdb.GetDB().First(...)`）
❌ 业务规则验证（如名称唯一性、库存是否充足、余额是否足够）
❌ 数据转换和计算

### Services 层职责（重要）

**Services 层负责以下工作：**

1. **接收 DTO** - 从 handler 接收已通过格式验证的 DTO
2. **模型构造** - 将 DTO 转换为模型对象
3. **业务规则验证** - 验证业务规则（如名称不能为空、名称唯一性检查、库存是否充足、余额是否足够）
4. **数据库操作** - 使用 `xdb.GetDB()` 进行 CRUD 操作
5. **返回结果** - 返回模型或错误

### Database 层职责（重要）

**Database 层负责以下工作：**

1. **数据库连接管理** - 管理主库和备库连接
2. **连接获取** - 通过 `GetDB()` 获取主库连接，通过 `GetSlaveDB()` 获取备库连接
3. **读写分离** - 支持读写分离配置

**Database 层禁止以下操作：**

❌ 直接暴露 `MasterDB` 或 `SlaveDBs` 变量
❌ 在业务代码中直接访问连接变量

## 生成的代码结构

技能将创建以下文件，遵循项目的最佳实践模式：

### 1. 模型文件 (internal/models/)
```go
type EntityName struct {
    BaseFields
    Name     string `gorm:"size:100;not null;comment:名称" json:"name"`
    Code     string `gorm:"size:50;unique;comment:编码" json:"code"`
    Status   int    `gorm:"comment:状态" json:"status"`
}
```

### 2. DTO文件 (internal/dto/)
```go
type EntityName struct {
    ID     int      `json:"id" validate:"omitempty,gte=0" label:"实体ID"`
    Name   *string  `json:"name" validate:"omitempty,max=100" label:"实体名称"`
    Code   *string  `json:"code" validate:"omitempty,max=50" label:"实体编码"`
    Status *int     `json:"status" validate:"omitempty,oneof=1 2 3 4" label:"状态"`
}

type EntityNameQuery struct {
    Page     int     `form:"page" validate:"omitempty,gte=1"`
    PageSize int     `form:"pageSize" validate:"omitempty,gte=1,lte=100"`
    Name     *string `form:"name" validate:"omitempty,max=100"`
    Status   *int    `form:"status" validate:"omitempty,oneof=1 2 3 4"`
}
```

### 3. 服务文件 (internal/services/)
```go
import (
    "errors"
    "webgos/internal/xdb"
    "webgos/internal/dto"
    "webgos/internal/models"
)

type EntityNameService interface {
    AddEntity(dtoModel dto.EntityName) (*models.EntityName, error)
    EditEntity(id int, dtoModel dto.EntityName) error
    GetEntityByID(id int) (*models.EntityName, error)
    GetEntityPage(query dto.EntityNameQuery) ([]models.EntityName, int)
    DeleteEntity(id int) error
}

type entityNameService struct{}

func NewEntityNameService() EntityNameService {
    return &entityNameService{}
}

func (s *entityNameService) AddEntity(dtoModel dto.EntityName) (*models.EntityName, error) {
    entity := &models.EntityName{
        Name:   *dtoModel.Name,
        Code:   *dtoModel.Code,
        Status: *dtoModel.Status,
    }
    if err := xdb.GetDB().Create(entity).Error; err != nil {
        return nil, err
    }
    return entity, nil
}

func (s *entityNameService) EditEntity(id int, dtoModel dto.EntityName) error {
    var entity models.EntityName
    if err := xdb.GetDB().First(&entity, id).Error; err != nil {
        return errors.New("实体不存在")
    }
    
    if dtoModel.Name != nil {
        entity.Name = *dtoModel.Name
    }
    if dtoModel.Code != nil {
        entity.Code = *dtoModel.Code
    }
    if dtoModel.Status != nil {
        entity.Status = *dtoModel.Status
    }
    
    return xdb.GetDB().Select("*").Updates(&entity).Error
}

func (s *entityNameService) GetEntityByID(id int) (*models.EntityName, error) {
    var entity models.EntityName
    err := xdb.GetDB().First(&entity, id).Error
    return &entity, err
}

func (s *entityNameService) GetEntityPage(query dto.EntityNameQuery) ([]models.EntityName, int64) {
    var entities []models.EntityName
    var total int64
    
    db := xdb.GetDB().Model(&models.EntityName{})
    if query.Name != nil {
        db = db.Where("name LIKE ?", "%"+*query.Name+"%")
    }
    if query.Status != nil {
        db = db.Where("status = ?", *query.Status)
    }
    
    if err := db.Count(&total).Error; err != nil {
        return []models.EntityName{}, 0
    }
    
    db = db.Scopes(models.Page(query.Page, query.PageSize))
    if err := db.Find(&entities).Error; err != nil {
        return []models.EntityName{}, 0
    }
    
    return entities, total
}

func (s *entityNameService) DeleteEntity(id int) error {
    return xdb.GetDB().Delete(&models.EntityName{}, id).Error
}
```

### 4. 处理器文件 (internal/handlers/)
```go
// EntityNameList 获取实体列表
// @Summary 获取实体列表
// @Description 分页获取实体列表
// @Tags 实体管理
// @Accept json
// @Produce json
// @Param body body dto.EntityNameQuery true "查询参数"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/entityname/list [post]
// @Security BearerAuth
func EntityNameList(c *gin.Context) {
    var query dto.EntityNameQuery
    if err := param.Validate(c, &query); err != nil {
        response.Error(c, err.Error())
        return
    }

    entityService := services.NewEntityNameService()
    entities, total := entityService.GetEntityPage(query)

    response.Success(c, "获取成功", gin.H{
        "list":  entities,
        "total": total,
    })
}

// EntityNameDetail 获取实体详情
// @Summary 获取实体详情
// @Description 根据ID获取实体详情
// @Tags 实体管理
// @Produce json
// @Param id path int true "实体ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/entityname/{id} [get]
// @Security BearerAuth
func EntityNameDetail(c *gin.Context) {
    id := convert.S2Int(c.Param("id"))
    if id == 0 {
        response.Error(c, "无效的实体 ID")
        return
    }

    entityService := services.NewEntityNameService()
    entity, err := entityService.GetEntityByID(id)
    if err != nil {
        response.Error(c, "实体不存在")
        return
    }

    response.Success(c, "获取成功", entity)
}

// EntityNameEdit 创建或更新实体
// @Summary 创建或更新实体
// @Description 创建或更新实体信息
// @Tags 实体管理
// @Accept json
// @Produce json
// @Param body body dto.EntityName true "实体信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/entityname/edit [post]
// @Security BearerAuth
func EntityNameEdit(c *gin.Context) {
    var postDTO dto.EntityName
    if err := param.Validate(c, &postDTO); err != nil {
        response.Error(c, err.Error())
        return
    }

    entityService := services.NewEntityNameService()
    
    if postDTO.ID > 0 {
        if err := entityService.EditEntity(postDTO.ID, postDTO); err != nil {
            response.Error(c, "更新失败："+err.Error())
            return
        }
    } else {
        _, err := entityService.AddEntity(postDTO)
        if err != nil {
            response.Error(c, "创建失败："+err.Error())
            return
        }
    }

    response.Success(c, "保存成功", nil)
}

// EntityNameDelete 删除实体
// @Summary 删除实体
// @Description 删除指定实体
// @Tags 实体管理
// @Produce json
// @Param id path int true "实体ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/entityname/{id} [delete]
// @Security BearerAuth
func EntityNameDelete(c *gin.Context) {
    id := convert.S2Int(c.Param("id"))
    if id == 0 {
        response.Error(c, "无效的实体 ID")
        return
    }

    entityService := services.NewEntityNameService()
    if err := entityService.DeleteEntity(id); err != nil {
        response.Error(c, "删除失败："+err.Error())
        return
    }

    response.Success(c, "删除成功", nil)
}
```

### 5. 路由文件 (internal/routes/)
```go
func init() {
    Register(func(router *Router) {
        wr := router.Group("/api/entityname")
        wr.Use(middleware.JWT())
        {
            wr.POST("/list", "获取列表", handlers.EntityNameList)
            wr.GET("/:id", "详情", handlers.EntityNameDetail)
            wr.POST("/edit", "创建/更新", handlers.EntityNameEdit)
            wr.DELETE("/:id", "删除", handlers.EntityNameDelete)
        }
    })
}
```

## 代码规范

### 架构模式规范
- **职责分离**：Handler 只处理参数和响应，Service 处理所有业务逻辑
- **参数验证**：使用 `param.Validate` 进行统一的参数验证
- **响应格式**：使用统一的 `response.Response` 格式返回数据
- **路由注册**：使用 `init()` 函数和 `WrapRouter` 进行路由注册
- **中间件**：默认添加 JWT 认证中间件
- **数据库访问**：统一使用 `xdb.GetDB()` 获取数据库连接

### Handlers 层规范（重要）
- **参数接收**：解析路径参数、查询参数、请求体
- **参数格式验证**：使用 `param.Validate` 进行参数格式验证（如必填字段、数据类型、范围限制）
- **响应处理**：调用 `response.Success/Error` 返回结果
- **禁止构造模型**：不要在 handler 中创建 `&models.Xxx{}`
- **禁止数据库操作**：不要在 handler 中调用 `xdb.GetDB()`
- **禁止业务验证**：不要在 handler 中进行业务规则验证（名称唯一性、库存检查等）

### Services 层规范
- **接口化设计**：定义接口 + 实现结构体
- **接口命名**：使用 `EntityNameService`，实现结构体使用 `entityNameService`
- **工厂函数**：返回接口类型：`NewEntityNameService() EntityNameService`
- **接收 DTO**：方法参数使用 DTO 类型，如 `AddEntity(dtoModel dto.EntityName)`
- **构造模型**：在 service 层将 DTO 转换为模型对象
- **数据库访问**：使用 `xdb.GetDB()` 进行数据库操作，**禁止**直接访问 `xdb.MasterDB`
- **业务规则验证**：在 service 层进行业务规则验证（如名称唯一性、库存是否充足、余额是否足够）
- **分页逻辑**：使用 `db.Scopes(models.Page(page, pageSize))` 进行分页

### 模型层规范
- **嵌入 BaseFields**：使用 `BaseFields` 结构体嵌入基础字段（ID、CreatedAt、UpdatedAt、DeletedAt）
- **ID 字段类型**：使用 `int` 类型
- **JSON 字段名**：使用小驼峰命名法
- **软删除**：`DeletedAt` 使用指针类型 `*gorm.DeletedAt`，配合 `omitempty` 标签

### DTO 规范
- **统一 DTO 设计**：使用单个 DTO 结构体同时支持创建和更新操作
- **指针类型字段**：所有字段使用指针类型支持部分更新
- **验证标签**：使用 `validate` 标签进行参数验证

### Database 层规范
- **包名**：使用 `xdb` 作为包名，位于 `internal/xdb/`
- **连接获取**：使用 `GetDB()` 获取主库连接，使用 `GetSlaveDB()` 获取备库连接
- **禁止暴露变量**：禁止直接访问 `MasterDB` 和 `SlaveDBs` 变量
- **读写分离**：通过配置控制是否启用读写分离

## 示例用法

### 示例：产品管理API
"创建一个产品管理API，包含以下功能：
- 创建产品：AddProduct
- 更新产品：EditProduct
- 获取产品详情：GetProductByID
- 分页查询产品列表：GetProductPage
- 删除产品：DeleteProduct

实体字段包括：产品名称、编码、类型、状态、价格、库存等"

技能将根据您提供的业务需求，生成符合项目最佳实践的完整API接口。

## 重要说明

**服务重启策略**：
- 每次修改后端代码后，技能不会自动重启服务
- 您需要手动重启服务以应用代码变更

**代码生成原则**：
- 遵循项目现有架构模式
- 保持代码风格一致性
- 优先使用项目已有的工具和库
- 确保生成的代码可以直接运行
- 严格遵守职责分离原则：Handler 只处理参数和响应，Service 处理业务逻辑
- 严格遵守数据库访问规范：统一使用 `xdb.GetDB()`，禁止直接访问 `xdb.MasterDB`