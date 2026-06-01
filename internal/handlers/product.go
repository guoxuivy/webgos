package handlers

import (
	"webgos/internal/models"
	"webgos/internal/services"
	"webgos/internal/utils/response"

	"github.com/gin-gonic/gin"
)

// AddProduct 创建商品
// @Summary 创建商品
// @Description 新增商品
// @Tags 商品
// @Accept json
// @Produce json
// @Param body body models.Product true "商品参数"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/products/add [post]
// @Security BearerAuth
func AddProduct(c *gin.Context) {
	var product models.Product

	if err := c.ShouldBind(&product); err != nil {
		response.Error(c, err.Error())
		return
	}

	productService := services.NewProductService()
	if err := productService.CreateProduct(&product); err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, "产品添加成功", nil)
}

// GetProductByID 获取商品详情
// @Summary 获取商品详情
// @Description 根据ID获取商品详情
// @Tags 商品
// @Produce json
// @Param id path int true "商品ID"
// @Success 200 {object} response.Response{data=models.Product{}}
// @Failure 400 {object} response.Response
// @Router /api/products/{id} [get]
// @Security BearerAuth
func GetProductByID(c *gin.Context) {
	idStr := c.Param("id")
	productService := services.NewProductService()
	product, err := productService.GetProductByID(idStr)
	if err != nil {
		response.Error(c, "产品不存在")
		return
	}
	response.Success(c, "获取产品成功", product)
}
