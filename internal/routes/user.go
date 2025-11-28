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
		user := WrapRouter(router.Group("/api/user"))
		user.Use(middleware.JWT())
		{
			user.GET("/info", "当前用户", handlers.UserInfo)
			user.POST("/list", "获取用户列表", handlers.UsersList)
			user.POST("/edit", "修改用户", handlers.UserEdit)
		}
	})
}
