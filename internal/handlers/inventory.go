package handlers

import (
	"webgos/internal/models"
	"webgos/internal/services"
	"webgos/internal/utils/response"

	"github.com/gin-gonic/gin"
)

// ProductIn 产品入库
// @Summary 产品入库
// @Description 产品入库操作
// @Tags 库存
// @Accept json
// @Produce json
// @Param body body models.InventoryRecord true "入库参数"
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

	inventoryService := services.NewInventoryService()
	if err := inventoryService.ProductIn(&record); err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, "产品入库成功", nil)
}

// ProductOut 产品出库
// @Summary 产品出库
// @Description 产品出库操作
// @Tags 库存
// @Accept json
// @Produce json
// @Param body body models.InventoryRecord true "出库参数"
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

	inventoryService := services.NewInventoryService()
	if err := inventoryService.ProductOut(&record); err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, "产品出库成功", nil)
}
