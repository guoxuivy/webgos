---
name: "backend-api-creator"
description: "自动创建完整的后端API接口，采用接口化设计模式，包括模型、DTO、服务接口与实现、处理器和路由。当用户需要创建新的后端接口或API端点时调用此技能。"
---

# 后端 CRUD 接口代码模板

> **本文件是纯代码模板，不定义规范。**
> 后端规范的唯一权威源：[`../backend-conventions.md`](../backend-conventions.md)
> （分层、Context 传递、读写分离选库、分页返回格式、JSONB 序列化等一律以它为准）。
> 本文件只提供可直接复制的代码骨架，**不再复述任何规范**；两者冲突时以权威源为准，并回来修正本文件。

## 功能特性

- 自动生成符合项目规范的完整 API 接口
- 采用接口化设计模式，实现服务层解耦
- 包含模型（Model）、DTO、服务接口与实现、处理器（Handler）和路由（Route）
- 支持 RESTful 设计、Swagger 文档注释、GORM 数据库操作、入参校验与统一响应
- 支持服务层 Context 传递

## 使用场景

需要新增一个后端接口 / API 端点时使用。需提供：

1. **实体名称**（如 Product、User、Order）
2. **实体字段**（字段名、类型、验证规则）
3. **业务功能需求**（创建订单、查询用户列表、更新产品状态等）
4. **服务方法名称**（按业务语义命名，如 `CreateOrder`、`GetUserList`、`UpdateProductStatus`）

## 生成的代码结构

### 1. 模型文件（internal/models/）

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

### 2. DTO 文件（internal/dto/）

```go
// EntityName 实体数据传输对象
// @description 实体数据传输对象，用于API请求和响应
type EntityName struct {
    ID     *int     `json:"id" validate:"omitempty,gte=0" label:"实体ID"`
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
// 注意：仅做字段映射，不查库；如需更新，由服务层用 ctxDB(ctx) 写入
func (dto *EntityName) ToModel() *models.EntityName {
    var model models.EntityName

    if dto.ID != nil {
        model.ID = *dto.ID
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

    return &model
}
```

### 3. 服务文件（internal/services/）

服务层通过 `ctxDB(ctx)`（主库/写）与 `ctxSDB(ctx)`（从库/读）获取带上下文的 `*gorm.DB`，
两个辅助函数定义在 `internal/services/helper.go`。

> **选库原则、主从复制延迟边界、事务闭包用法见
> [`../backend-conventions.md` §读写分离与选库](../backend-conventions.md)，本模板不重复。**

```go
import (
    "context"

    "webgos/internal/dto"
    "webgos/internal/models"
)

// EntityNameService 实体服务接口
// 方法按业务语义命名，避免千篇一律的 Save/Get/Delete 前缀
type EntityNameService interface {
    // 统一保存方法（创建和更新）
    SaveEntity(ctx context.Context, entity *models.EntityName) error

    // 查询方法
    GetEntityByID(ctx context.Context, id int) (*models.EntityName, error)
    GetEntityPage(ctx context.Context, query dto.EntityNameQuery) ([]models.EntityName, int64)

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
    // 手动序列化 JSONB 字段（如有）
    if err := entity.Serialize(); err != nil {
        return err
    }

    db := ctxDB(ctx)

    if entity.ID > 0 {
        return db.Select("*").Updates(entity).Error
    }

    // 创建逻辑
    return db.Create(entity).Error
}

// GetEntityByID 根据ID获取实体详情
func (s *entityNameService) GetEntityByID(ctx context.Context, id int) (*models.EntityName, error) {
    var entity models.EntityName
    err := ctxSDB(ctx).First(&entity, id).Error
    if err != nil {
        return nil, err
    }

    // 手动反序列化 JSONB 字段（如有）
    if err := entity.Deserialize(); err != nil {
        return nil, err
    }

    return &entity, nil
}

// GetEntityPage 分页查询实体列表
func (s *entityNameService) GetEntityPage(ctx context.Context, query dto.EntityNameQuery) ([]models.EntityName, int64) {
    db := ctxSDB(ctx).Model(&models.EntityName{})

    // 构建查询条件（指针类型判断参数是否存在）
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
    // 使用 models.Page Scope 统一处理分页（替代手动计算 offset）
    db = db.Scopes(models.Page(query.Page, query.PageSize))
    if err := db.Find(&entities).Error; err != nil {
        return []models.EntityName{}, 0
    }

    // 手动反序列化所有记录的 JSONB 字段（如有）
    for i := range entities {
        if err := entities[i].Deserialize(); err != nil {
            continue
        }
    }

    return entities, total
}

// DeleteEntity 删除实体
func (s *entityNameService) DeleteEntity(ctx context.Context, id int) error {
    return ctxDB(ctx).Delete(&models.EntityName{}, id).Error
}
```

### 4. 处理器文件（internal/handlers/）

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
    entities, total := entityService.GetEntityPage(c, query) // c *gin.Context 实现了 context.Context

    response.Success(c, "获取成功", gin.H{
        "items": entities,  // 必须为 items，禁止 list / data / rows
        "total": total,     // 必须为 total，禁止 count / sum
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

### 5. 路由文件（internal/routes/）

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

## 示例用法

提供具体业务需求即可：

### 示例 1：楼栋管理 API

"创建一个楼栋管理API，包含以下功能：

- 统一创建/更新接口：SaveBuilding
- 获取楼栋详情：GetBuildingByID
- 分页查询楼栋列表：GetBuildingPage
- 删除楼栋：DeleteBuilding

实体字段包括：楼栋编号、园区ID、楼层数、面积、用途类型、建筑类型等"

### 示例 2：产品管理 API

"创建一个产品管理API，包含以下功能：

- 统一创建/更新接口：SaveProduct
- 获取产品详情：GetProductByID
- 分页查询产品列表：GetProductPage
- 删除产品：DeleteProduct

实体字段包括：产品名称、编码、类型、状态、价格、库存等"

### 示例 3：订单管理 API

"创建一个订单管理API，包含以下功能：

- 统一创建/更新接口：SaveOrder
- 获取订单详情：GetOrderByID
- 分页查询订单列表：GetOrderPage
- 删除订单：DeleteOrder

实体字段包括：订单号、用户ID、商品信息、金额、状态等"

## 重要说明

**服务重启策略**：

- 每次修改后端代码后，本模板不会自动重启服务
- 需要手动重启服务以应用代码变更（`make run`）
- 这样可以避免意外的服务中断和数据丢失

**代码生成原则**：

- 遵循项目现有架构模式
- 保持代码风格一致性
- 优先使用项目已有的工具和库
- 确保生成的代码可以直接运行
- **生成后必须对照 [`../backend-conventions.md`](../backend-conventions.md) 复核**，尤其是分页返回格式、Context 传递与选库
