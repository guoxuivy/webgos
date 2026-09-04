package routes

import (
	"path"

	"github.com/gin-gonic/gin"

	"webgos/internal/models"
)

// RouterWrapper 包装gin的RouterGroup
type RouterWrapper struct {
	*gin.RouterGroup
}

// WrapRouter 包装一个RouterGroup，返回自定义的RouterWrapper
func WrapRouter(group *gin.RouterGroup) *RouterWrapper {
	return &RouterWrapper{group}
}

// 以下是对各种HTTP方法的包装，自动收集路由信息作为权限点

// GET 包装GET方法，自动收集路由信息 支持中间件注入
func (w *RouterWrapper) GET(relativePath string, description string, handlers ...gin.HandlerFunc) {
	w.addRouteInfoWithHandlers(relativePath, "GET", description, handlers...)
}

// POST 包装POST方法，自动收集路由信息 支持中间件注入
func (w *RouterWrapper) POST(relativePath string, description string, handlers ...gin.HandlerFunc) {
	w.addRouteInfoWithHandlers(relativePath, "POST", description, handlers...)
}

// PUT 包装PUT方法，自动收集路由信息 支持中间件注入
func (w *RouterWrapper) PUT(relativePath string, description string, handlers ...gin.HandlerFunc) {
	w.addRouteInfoWithHandlers(relativePath, "PUT", description, handlers...)
}

// DELETE 包装DELETE方法，自动收集路由信息 支持中间件注入
func (w *RouterWrapper) DELETE(relativePath string, description string, handlers ...gin.HandlerFunc) {
	w.addRouteInfoWithHandlers(relativePath, "DELETE", description, handlers...)
}

// addRouteInfoWithHandlers 是一个内部方法，用于添加路由信息并注册处理函数
// 支持多个处理函数（包括中间件）
func (w *RouterWrapper) addRouteInfoWithHandlers(relativePath, method, description string, handlers ...gin.HandlerFunc) {

	// 注册处理函数到路由组
	switch method {
	case "GET":
		w.RouterGroup.GET(relativePath, handlers...)
	case "POST":
		w.RouterGroup.POST(relativePath, handlers...)
	case "PUT":
		w.RouterGroup.PUT(relativePath, handlers...)
	case "DELETE":
		w.RouterGroup.DELETE(relativePath, handlers...)
	}
	// 将路由信息写入内存描述表（key = path#method），供管理端实时列出权限点。
	// 权限点不再持久化到数据库，始终等于真实注册的路由。
	fullPath := w.calculateFullPath(relativePath)
	models.RouteDescriptions[models.BuildPermKey(fullPath, method)] = description

}
func lastChar(str string) uint8 {
	if str == "" {
		panic("The length of the string can't be 0")
	}
	return str[len(str)-1]
}

// 计算完整的路由路径
func (w *RouterWrapper) calculateFullPath(relativePath string) string {
	absolutePath := w.BasePath()
	if relativePath == "" {
		return absolutePath
	}

	finalPath := path.Join(absolutePath, relativePath)
	if lastChar(relativePath) == '/' && lastChar(finalPath) != '/' {
		return finalPath + "/"
	}
	return finalPath
}
