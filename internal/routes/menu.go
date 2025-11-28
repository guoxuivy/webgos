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
	})
}
