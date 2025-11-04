package handlers

import (
	"webgos/internal/models"
	"webgos/internal/services"
	"webgos/internal/utils/response"

	"github.com/gin-gonic/gin"
)

// @Summary 产品入库
// @Description 产品入库操作
// @Tags 库存
// @Accept json
// @Produce json
// @Param data body models.InventoryRecord true "入库参数"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/inventory/in [post]
// @Security BearerAuth
func ProductIn(c *gin.Context) {
	var record models.InventoryRecord

	if err := c.ShouldBind(&record); err != nil {
		response.Error(c, err.Error())
		return
	}

	// 验证必需字段
	if record.ProductID == 0 || record.Quantity <= 0 {
		response.Error(c, "产品ID和数量必须大于0")
		return
	}

	// 强制设置类型为"in"
	record.Type = "in"

	inventoryService := services.NewInventoryService()
	if err := inventoryService.ProductIn(&record); err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, "产品入库成功", nil)
}

// @Summary 产品出库
// @Description 产品出库操作
// @Tags 库存
// @Accept json
// @Produce json
// @Param data body models.InventoryRecord true "出库参数"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/inventory/out [post]
// @Security BearerAuth
func ProductOut(c *gin.Context) {
	var record models.InventoryRecord

	if err := c.ShouldBind(&record); err != nil {
		response.Error(c, err.Error())
		return
	}

	// 验证必需字段
	if record.ProductID == 0 || record.Quantity <= 0 {
		response.Error(c, "产品ID和数量必须大于0")
		return
	}

	// 强制设置类型为"out"
	record.Type = "out"

	inventoryService := services.NewInventoryService()
	if err := inventoryService.ProductOut(&record); err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, "产品出库成功", nil)
}
