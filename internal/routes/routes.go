package routes

import (
	"webgos/internal/config"

	"net/http"
	"webgos/internal/middleware"
	"webgos/internal/utils/response"

	"github.com/gin-gonic/gin"
)

// 全局路由引擎实例
var REngine *gin.Engine

// RouteRegister 路由注册器函数类型
type RouteRegister func(router *gin.Engine)

// 全局路由注册器列表
var routeRegisters []RouteRegister

// Register 注册路由注册器
func Register(register RouteRegister) {
	routeRegisters = append(routeRegisters, register)
}

// 创建路由引擎
func New(config *config.Config) *gin.Engine {
	// 设置Gin模式
	gin.SetMode(config.Server.Mode)
	// 创建不带默认中间件的路由引擎
	REngine = gin.New()

	// 应用通用中间件
	middleware.ApplyMiddlewares(REngine, config)

	// 注册所有路由
	for _, register := range routeRegisters {
		register(REngine)
	}

	// 404处理
	REngine.NoRoute(handleNotFound)
	return REngine
}

// handleNotFound 处理404错误
func handleNotFound(c *gin.Context) {
	response.ErrorWithCode(c, "请求的资源不存在", http.StatusNotFound)
}
