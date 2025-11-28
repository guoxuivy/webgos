package routes

import (
	"webgos/internal/handlers"
	"webgos/internal/middleware"

	"github.com/gin-gonic/gin"
)

// init 自动注册路由
func init() {
	Register(func(router *gin.Engine) {
		// 需要认证的路由组
		api := router.Group("/api")
		api.Use(middleware.JWT())

		products := WrapRouter(api.Group("/products"))
		{
			products.POST("/add", "创建商品", handlers.AddProduct)
			products.GET("/:id", "获取商品详情", handlers.GetProductByID)
		}
	})
}
