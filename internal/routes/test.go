package routes

import (
	"webgos/internal/handlers"
	"webgos/internal/middleware"

	"github.com/gin-gonic/gin"
)

// init 自动注册路由
func init() {
	Register(func(router *gin.Engine) {
		// 测试相关路由
		testGroup := router.Group("/test")
		// testGroup.Use(middleware.Gzip())
		{
			testGroup.GET("/Test", middleware.Gzip(), handlers.Test)
			testGroup.GET("/TestCurlList", handlers.TestCurlList)
			testGroup.GET("/TestTransaction", handlers.TestTransaction)
			testGroup.GET("/TestTransaction2", handlers.TestTransaction2)
			testGroup.GET("/TestCurlUpdate/:id/:age", handlers.TestCurlUpdate)
		}
	})
}
