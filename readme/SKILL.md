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
- **支持服务层 Context 传递**

## 使用场景

当您需要创建新的后端接口时，调用此技能并提供以下信息：

1. **实体名称**（如：Product, User, Order等）
2. **实体字段**（字段名、类型、验证规则等）
3. **业务功能需求**（具体的业务操作，如：创建订单、查询用户列表、更新产品状态等）
4. **服务方法名称**（根据业务语义命名，如：CreateOrder、GetUserList、UpdateProductStatus等）

## 生成的代码结构

技能将创建以下文件，遵循项目的最佳实践模式：

### 1. 模型文件 (internal/models/)
```go
type EntityName struct {
    BaseFields
    // 字段定义 - 命名尽量简化，JSON标签使用小驼峰
    Name     string `gorm:"size:100;not null;comment:名称" json:"name"`
    Code     string `gorm:"size:50;unique;comment:编码" json:"code"`
    Status   int    `gorm:"comment:状态" json:"status"`
    // 其他简化字段...
}

// 复杂类型字段使用JSONB存储
// 例如位置信息、配置信息等
```

### 2. DTO文件 (internal/dto/)
```go
// EntityName 实体数据传输对象
// @description 实体数据传输对象，用于API请求和响应
type EntityName struct {
    ID     int      `json:"id" validate:"omitempty,gte=0" label:"实体ID"`
    Name   *string  `json:"name" validate:"omitempty,max=100" label:"实体名称"`
    Code   *string  `json:"code" validate:"omitempty,max=50" label:"实体编码"`
    Status *int     `json:"status" validate:"omitempty,oneof=1 2 3 4" label:"状态"`
    // 其他字段...
}

// EntityNameQuery 实体查询参数
// @description 实体列表查询参数
// Page和PageSize 只能是int类型，不能是指针类型
type EntityNameQuery struct {
    Page     int     `form:"page" validate:"omitempty,gte=1"`
    PageSize int     `form:"pageSize" validate:"omitempty,gte=1,lte=100"`
    Name     *string `form:"name" validate:"omitempty,max=100"`
    Status   *int    `form:"status" validate:"omitempty,oneof=1 2 3 4"`
}

// ToModel 将 DTO 转换为模型（用于创建和更新操作）
func (dto *EntityName) ToModel() *models.EntityName {
    var model *models.EntityName
    if dto.ID == 0 {
        model = &models.EntityName{}
    } else {
        var existing models.EntityName
        if err := database.MasterDB.First(&existing, dto.ID).Error; err != nil {
            return nil
        }
        model = &existing
    }

    if dto.Name != nil {
        model.Name = *dto.Name
    }
    if dto.Code != nil {
        model.Code = *dto.Code
    }
    if dto.Status != nil {
        model.Status = *dto.Status
    }
    // 其他字段转换...

    return model
}
```

### 3. 服务文件 (internal/services/)
```go
import (
    "context"
    "cyp/internal/database"
    "cyp/internal/dto"
    "cyp/internal/models"
)

// EntityNameService 实体服务接口
type EntityNameService interface {
    // 统一保存方法（创建和更新）
    SaveEntity(ctx context.Context, data *models.EntityName) error
    
    // 查询方法
    GetEntityByID(ctx context.Context, id int) (*models.EntityName, error)
    GetEntityPage(ctx context.Context, query dto.EntityNameQuery) ([]models.EntityName, int)
    
    // 删除方法
    DeleteEntity(ctx context.Context, id int) error
}

// entityNameService 实现 EntityNameService 接口
type entityNameService struct{}

// NewEntityNameService 创建实体服务实例
func NewEntityNameService() EntityNameService {
    return &entityNameService{}
}

// SaveEntity 统一保存实体（创建和更新）
func (s *entityNameService) SaveEntity(ctx context.Context, entity *models.EntityName) error {
    // 手动序列化 JSONB 字段
    if err := entity.Serialize(); err != nil {
        return err
    }
    
    if entity.ID > 0 {
        return database.MasterDB.Select("*").Updates(entity).Error
    }
    
    // 创建逻辑
    return database.MasterDB.Create(entity).Error
}

// GetEntityByID 根据ID获取实体详情
func (s *entityNameService) GetEntityByID(ctx context.Context, id int) (*models.EntityName, error) {
    var entity models.EntityName
    err := database.MasterDB.First(&entity, id).Error
    if err != nil {
        return nil, err
    }
    
    // 手动反序列化 JSONB 字段
    if err := entity.Deserialize(); err != nil {
        return nil, err
    }
    
    return &entity, nil
}

// GetEntityPage 分页查询实体列表
func (s *entityNameService) GetEntityPage(ctx context.Context, query dto.EntityNameQuery) ([]models.EntityName, int) {
    db := database.MasterDB.Model(&models.EntityName{})

    // 构建查询条件
    if query.Name != nil {
        db = db.Where("name LIKE ?", "%"+*query.Name+"%")
    }
    if query.Status != nil {
        db = db.Where("status = ?", *query.Status)
    }

    var total int64
    if err := db.Count(&total).Error; err != nil {
        return []models.EntityName{}, 0
    }

    var entities []models.EntityName
    offset := (query.Page - 1) * query.PageSize
    if err := db.Offset(offset).Limit(query.PageSize).Find(&entities).Error; err != nil {
        return []models.EntityName{}, 0
    }
    
    // 手动反序列化所有记录的 JSONB 字段
    for i := range entities {
        if err := entities[i].Deserialize(); err != nil {
            continue
        }
    }
    
    return entities, int(total)
}

// DeleteEntity 删除实体
func (s *entityNameService) DeleteEntity(ctx context.Context, id int) error {
    return database.MasterDB.Delete(&models.EntityName{}, id).Error
}
```

### 4. 处理器文件 (internal/handlers/)
```go
// EntityNameList 获取列表
// @Summary 实体列表
// @Description 获取实体列表接口
// @Tags 实体
// @Accept json
// @Produce json
// @Param body body dto.EntityNameQuery true "实体查询参数"
// @Success 200 {array} models.EntityName
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
    entities, total := entityService.GetEntityPage(c, query)

    response.Success(c, "获取成功", gin.H{
        "list":  entities,
        "total": total,
    })
}

// EntityNameDetail 获取详情
// @Summary 实体详情
// @Description 获取实体详情接口
// @Tags 实体
// @Accept json
// @Produce json
// @Param id path int true "实体ID"
// @Success 200 {object} models.EntityName
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
    entity, err := entityService.GetEntityByID(c, id)
    if err != nil {
        response.Error(c, "实体不存在")
        return
    }

    response.Success(c, "获取成功", entity)
}

// EntityNameEdit 创建或更新（统一接口）
// @Summary 实体创建/更新
// @Description 创建或更新实体接口，不传 ID 则为创建，传 ID 则为更新
// @Tags 实体
// @Accept json
// @Produce json
// @Param body body dto.EntityName true "实体信息 修改时需要包含ID字段"
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

    // 转换为模型
    entity := postDTO.ToModel()
    if entity == nil {
        response.Error(c, "实体不存在")
        return
    }

    entityService := services.NewEntityNameService()
    if err := entityService.SaveEntity(c, entity); err != nil {
        response.Error(c, "保存失败："+err.Error())
        return
    }

    response.Success(c, "保存成功", nil)
}

// EntityNameDelete 删除
// @Summary 实体删除
// @Description 删除实体接口
// @Tags 实体
// @Accept json
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
    if err := entityService.DeleteEntity(c, id); err != nil {
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
- **统一保存接口**：服务层使用`SaveEntity`方法统一处理创建和更新操作
- **参数验证**：使用`param.Validate`进行统一的参数验证
- **响应格式**：使用统一的`response.Response`格式返回数据
- **路由注册**：使用`init()`函数和`WrapRouter`进行路由注册
- **中间件**：默认添加JWT认证中间件
- **Context 传递**：服务层方法统一接收 `context.Context` 参数
- **GORM原生操作**：直接使用`database.MasterDB`进行数据库操作，不再使用BaseModel封装

### 服务层规范
- **接口化设计**：定义接口 + 实现结构体
- **接口命名**：使用`EntityNameService`，实现结构体使用`entityNameService`
- **工厂函数**：返回接口类型：`NewEntityNameService() EntityNameService`
- **Context 支持**：所有方法第一个参数为 `ctx context.Context`
- **方法命名**：
  - `SaveEntity(ctx context.Context, data *models.EntityName) error` - 统一保存方法
  - `GetEntityByID(ctx context.Context, id int) (*models.EntityName, error)` - 获取详情
  - `GetEntityPage(ctx context.Context, query dto.EntityNameQuery) ([]models.EntityName, int)` - 分页查询
  - `DeleteEntity(ctx context.Context, id int) error` - 删除操作
- **数据库操作**：直接使用`database.MasterDB`进行GORM原生操作

### 处理器层规范
- **统一编辑接口**：使用`EntityNameEdit`方法统一处理创建和更新
- **参数处理**：使用`convert.S2Int`进行ID参数转换
- **错误处理**：统一的错误消息格式
- **API文档**：包含完整的Swagger注释
- **Context 传递**：调用服务方法时直接传递 `c *gin.Context`（`gin.Context` 实现了 `context.Context` 接口）

### 模型层规范
- **BaseFields嵌入**：模型通过嵌入`BaseFields`结构体获得ID、CreatedAt、UpdatedAt、DeletedAt字段
- **ID字段类型**：所有与ID相关的参数必须使用`int`类型
- **字段命名简化**：避免冗余和过长的名称
- **JSON字段名**：使用小驼峰命名法，如：`landType`、`isActive`
- **复杂类型存储**：复杂字段使用JSONB存储，配合序列化方法

### DTO规范
- **统一DTO设计**：使用单个DTO结构体同时支持创建和更新操作
- **指针类型字段**：所有字段使用指针类型支持部分更新
- **验证标签**：使用`validate`标签进行参数验证，包含`label`标签用于错误消息
- **转换方法**：实现`ToModel()`方法进行DTO到模型的转换，使用GORM原生查询获取现有记录
- **查询DTO**：查询参数使用指针类型支持可选参数

### 序列化规范
- **JSONB字段**：复杂类型字段需要实现`Serialize()`和`DeSerialize()`方法
- **手动调用**：在服务层手动调用序列化方法，确保可靠性
- **创建/更新**：在`SaveEntity`方法中手动调用`Serialize()`
- **查询操作**：在`GetEntityByID`和`GetEntityPage`方法中手动调用`DeSerialize()`
- **避免钩子**：不依赖GORM钩子，采用明确的手动调用方式

### 查询构建规范
- **链式查询**：使用GORM的链式查询构建器
- **条件判断**：使用指针类型判断参数是否存在
- **模糊查询**：字符串字段使用`LIKE`进行模糊匹配
- **分页处理**：手动计算offset，使用`Offset()`和`Limit()`方法进行分页查询

## 重要优化点

### 1. DTO设计优化
- **统一DTO**：不再区分Create和Update DTO，使用单个DTO结构体
- **指针类型**：所有字段使用指针类型，支持部分更新
- **智能转换**：`ToModel()`方法根据ID自动判断是创建还是更新

### 2. 查询构建优化
- **链式构建**：使用GORM的链式查询构建器
- **条件过滤**：通过指针判断是否添加查询条件
- **灵活查询**：支持精确查询和模糊查询
- **GORM原生**：直接使用`database.MasterDB`进行查询

### 3. 错误处理优化
- **统一格式**：使用`response.Error()`和`response.Success()`
- **详细错误**：包含具体的错误信息
- **参数验证**：使用`param.Validate`进行统一验证

### 4. 序列化优化
- **明确调用**：在服务层明确调用序列化方法
- **批量处理**：分页查询时批量反序列化
- **错误处理**：序列化失败时继续处理其他记录

### 5. Context 支持优化
- **统一传递**：所有服务方法接收 `context.Context` 参数
- **请求上下文**：支持传递请求级别的数据和超时控制
- **接口兼容**：`gin.Context` 实现了 `context.Context` 接口，可直接传递

### 6. 模型层重构优化
- **移除BaseModel**：不再继承BaseModel泛型基类
- **BaseFields嵌入**：使用`BaseFields`结构体嵌入获得基础字段
- **GORM原生操作**：服务层直接使用`database.MasterDB`进行数据库操作
- **简化设计**：降低代码复杂度，提高可维护性

## 示例用法

当您需要创建新的API接口时，请提供具体的业务需求：

### 示例1：楼栋管理API（基于优化后的模式）
"创建一个楼栋管理API，包含以下功能：
- 统一创建/更新接口：SaveBuilding
- 获取楼栋详情：GetBuildingByID
- 分页查询楼栋列表：GetBuildingPage
- 删除楼栋：DeleteBuilding

实体字段包括：楼栋编号、园区ID、楼层数、面积、用途类型、建筑类型等"

### 示例2：产品管理API
"创建一个产品管理API，包含以下功能：
- 统一创建/更新接口：SaveProduct
- 获取产品详情：GetProductByID
- 分页查询产品列表：GetProductPage
- 删除产品：DeleteProduct

实体字段包括：产品名称、编码、类型、状态、价格、库存等"

### 示例3：订单管理API
"创建一个订单管理API，包含以下功能：
- 统一创建/更新接口：SaveOrder
- 获取订单详情：GetOrderByID
- 分页查询订单列表：GetOrderPage
- 删除订单：DeleteOrder

实体字段包括：订单号、用户ID、商品信息、金额、状态等"

技能将根据您提供的业务需求，生成符合项目最佳实践的完整API接口。

## 重要说明

**服务重启策略**：
- 每次修改后端代码后，技能不会自动重启服务
- 您需要手动重启服务以应用代码变更
- 这样可以避免意外的服务中断和数据丢失
- 您可以根据自己的开发节奏控制重启时机

**代码生成原则**：
- 遵循项目现有架构模式
- 保持代码风格一致性
- 优先使用项目已有的工具和库
- 确保生成的代码可以直接运行
