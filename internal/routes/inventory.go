package routes

import (
	"webgos/internal/handlers"
	"webgos/internal/middleware"

	"github.com/gin-gonic/gin"
)

// init 自动注册路由
func init() {
	Register(func(router *gin.Engine) {

		// 库存相关路由
		inventory := WrapRouter(router.Group("/api/inventory"))
		inventory.Use(middleware.JWT())
		inventory.Use(middleware.RBAC())
		{
			// 为入库操作添加防抖中间件，防止重复提交 demo
			inventory.POST("/in", "入库测试", middleware.Debounce(), handlers.ProductIn)
			inventory.POST("/out", "出库测试", handlers.ProductOut)
		}
	})
}
