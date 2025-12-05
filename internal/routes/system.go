package routes

import (
	"webgos/internal/handlers"
	"webgos/internal/middleware"

	"github.com/gin-gonic/gin"
)

// init 自动注册路由
func init() {
	Register(func(router *gin.Engine) {

		// 登录相关路由（公开）
		publicGroup := router.Group("/auth")
		{
			publicGroup.POST("/register", handlers.RegisterUser)
			publicGroup.POST("/login", handlers.Login)
			publicGroup.POST("/logout", handlers.Logout)
			publicGroup.POST("/reset-password", handlers.ResetPassword)
		}

		// 需要认证的路由组
		api := router.Group("/api")
		api.Use(middleware.JWT())

		// 菜单管理路由
		menu := WrapRouter(api.Group("/menu"))
		{
			menu.POST("/add", "创建菜单", handlers.AddMenu)
			menu.POST("/edit", "编辑菜单", handlers.EditMenu)
			menu.POST("/delete/:id", "删除菜单", handlers.DeleteMenu)
			menu.GET("/info/:id", "获取菜单详情", handlers.GetMenuByID)
			menu.GET("/list", "获取菜单列表", handlers.GetMenus)
			menu.GET("/tree", "获取菜单树", handlers.GetMenuTree)
			menu.GET("/name_exists", "检查菜单名称是否存在", handlers.NameExists)
			menu.GET("/path_exists", "检查菜单路径是否存在", handlers.PathExists)
			menu.GET("/user_menus", "获取当前用户目录", handlers.GetUserMenus)
		}

		// 角色管理路由(勿动)~
		rbac := WrapRouter(api.Group("/rbac"))
		{
			rbac.POST("/role", "创建角色", handlers.AddRole)
			rbac.POST("/edit_role", "编辑角色", handlers.EditRole)
			rbac.GET("/roles", "角色列表", handlers.GetRoles)
			rbac.POST("/assign_roles", "分配角色给用户", handlers.AssignRoles)
			rbac.POST("/assign_permissions", "分配权限给角色", handlers.AssignPermissions)
			rbac.GET("/permissions", "全部权限项", handlers.GetPermissions)
			rbac.GET("/role_permissions/:id", "角色权限项", handlers.GetRolePermissions)
			rbac.GET("/role/:id", "角色详情", handlers.GetRoleByID)
			rbac.GET("/user_roles/:id", "用户角色", handlers.GetUserRoles)
		}

		// 用户管理路由
		user := WrapRouter(api.Group("/user"))
		{
			user.GET("/info", "当前用户", handlers.UserInfo)
			user.POST("/list", "获取用户列表", handlers.UsersList)
			user.POST("/edit", "修改用户", handlers.UserEdit)
		}

	})
}
