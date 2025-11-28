package routes

import (
	"webgos/internal/config"

	"net/http"
	"webgos/internal/middleware"
	"webgos/internal/utils/response"

	"github.com/gin-gonic/gin"
)

// RouteRegister 路由注册器函数类型
type RouteRegister func(router *gin.Engine)

// 全局路由注册器列表
var routeRegisters []RouteRegister

// Register 注册路由注册器
func Register(register RouteRegister) {
	routeRegisters = append(routeRegisters, register)
}

// SetupRoutes 设置路由
func SetupRoutes(config *config.Config) *gin.Engine {
	// 设置Gin模式
	gin.SetMode(config.Server.Mode)
	// 创建不带默认中间件的路由引擎
	r := gin.New()

	// 应用通用中间件
	middleware.ApplyMiddlewares(r, config)

	// 注册所有路由
	for _, register := range routeRegisters {
		register(r)
	}

	// 404处理
	r.NoRoute(handleNotFound)

	return r
}

// handleNotFound 处理404错误
func handleNotFound(c *gin.Context) {
	response.Error(c, "请求的资源不存在", http.StatusNotFound)
}
