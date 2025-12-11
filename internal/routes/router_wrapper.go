package routes

import (
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"webgos/internal/config"
	"webgos/internal/models"
)

// 存储注册的路由信息
type RouteInfo struct {
	Method      string
	Path        string
	Name        string
	Description string
}

// 存储所有路由信息
var routeInfos []RouteInfo

// RouterWrapper 包装gin的RouterGroup，用于收集路由信息
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
	// 添加路由信息到routeInfos（只记录路径和方法，不记录中间件）
	// 仅在配置中启用自动同步RBAC权限点时才收集路由信息
	if config.GlobalConfig.AutoRBACPoint {
		fullPath := strings.ToLower(w.calculateFullPath(relativePath))
		routeInfos = append(routeInfos, RouteInfo{
			Path:        fullPath,
			Method:      method,
			Description: description,
			Name:        fullPath + "#" + method,
		})
	}

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

// SyncPermissions 将收集的路由信息同步到数据库作为权限点
func SyncPermissions(db *gorm.DB) error {
	for _, route := range routeInfos {
		// 查找是否已存在该权限
		var existingPermission models.RBACPermission
		result := db.Where("name = ?", route.Name).First(&existingPermission)

		if result.Error != nil {
			// 权限不存在，创建新权限
			permission := models.RBACPermission{
				Path:        route.Path,
				Method:      route.Method,
				Description: route.Description,
				Name:        route.Name,
			}
			if err := db.Create(&permission).Error; err != nil {
				return err
			}
		} else {
			// 权限已存在，更新描述信息
			existingPermission.Description = route.Description
			if err := db.Save(&existingPermission).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
