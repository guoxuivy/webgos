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

		// 角色管理路由(勿动)~
		rbac := WrapRouter(api.Group("/rbac"))
		{
			rbac.POST("/add/role", "创建角色", handlers.AddRole)
			rbac.POST("/edit/role", "编辑角色", handlers.EditRole)
			rbac.POST("/assign/roles", "分配角色给用户", handlers.AssignRoles)
			rbac.POST("/assign/permissions", "分配权限给角色", handlers.AssignPermissions)
			rbac.GET("/roles", "角色列表", handlers.GetRoles)
			rbac.GET("/permissions", "全部权限项", handlers.GetPermissions)
			rbac.GET("/role/:id/permissions", "角色权限项", handlers.GetRolePermissions)
			rbac.GET("/role/:id", "角色详情", handlers.GetRoleByID)
			rbac.GET("/user/:id/roles", "用户角色", handlers.GetUserRoles)
		}
	})
}
