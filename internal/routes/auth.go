package routes

import (
	"webgos/internal/handlers"

	"github.com/gin-gonic/gin"
)

// init 自动注册路由
func init() {
	Register(func(router *gin.Engine) {
		// 登录相关路由（公开）
		publicUserGroup := router.Group("/auth")
		{
			publicUserGroup.POST("/register", handlers.RegisterUser)
			publicUserGroup.POST("/login", handlers.Login)
			publicUserGroup.POST("/logout", handlers.Logout)
			publicUserGroup.POST("/reset-password", handlers.ResetPassword)
		}
	})
}
