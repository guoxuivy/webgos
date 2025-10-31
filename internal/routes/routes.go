package routes

import (
	"hserp/internal/config"
	"hserp/internal/handlers"

	"hserp/internal/middleware"
	"hserp/internal/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置路由
func SetupRoutes(config *config.Config) *gin.Engine {
	// 设置Gin模式
	gin.SetMode(config.Server.Mode)
	// 创建不带默认中间件的路由引擎
	r := gin.New()

	// 应用通用中间件
	middleware.ApplyMiddlewares(r, config)

	// 测试相关路由
	testGroup := r.Group("/test")
	{
		testGroup.GET("/Test", handlers.Test)
		testGroup.GET("/TestCurlList", handlers.TestCurlList)
		testGroup.GET("/TestTransaction", handlers.TestTransaction)
		testGroup.GET("/TestTransaction2", handlers.TestTransaction2)
		testGroup.GET("/TestCurlUpdate/:id/:age", handlers.TestCurlUpdate)

	}
	// 登录相关路由（公开）
	publicUserGroup := r.Group("/auth")
	{
		publicUserGroup.POST("/register", handlers.RegisterUser)
		publicUserGroup.POST("/login", handlers.Login)
		publicUserGroup.POST("/logout", handlers.Logout)
		publicUserGroup.POST("/reset-password", handlers.ResetPassword)

	}

	// 需要认证的路由组
	api := r.Group("/api")
	api.Use(middleware.JWT())
	// api.Use(middleware.Rbac())
	{
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

		user := WrapRouter(api.Group("/user"))
		{
			user.GET("/info", "当前用户", handlers.UserInfo)
			user.POST("/list", "创建角色", handlers.UsersList)
			user.POST("/edit", "创建修改用户", handlers.UserEdit)
		}

		products := WrapRouter(api.Group("/products"))
		{
			products.POST("/add", "创建商品", handlers.AddProduct)
			products.GET("/:id", "获取商品详情", handlers.GetProductByID)
		}

		// 库存相关路由
		inventory := WrapRouter(api.Group("/inventory"))
		{
			// 为入库操作添加防抖中间件，防止重复提交 demo
			inventory.POST("/in", "入库测试", middleware.Debounce(), handlers.ProductIn)
			inventory.POST("/out", "出库测试", handlers.ProductOut)
		}

	}

	// 404处理
	r.NoRoute(handleNotFound)

	return r
}

// handleNotFound 处理404错误
func handleNotFound(c *gin.Context) {
	response.Error(c, "请求的资源不存在", http.StatusNotFound)
}
