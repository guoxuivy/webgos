package middleware

import (
	"fmt"
	"net/http"
	"time"
	"webgos/internal/config"
	"webgos/internal/database"
	"webgos/internal/models"
	"webgos/internal/utils/response"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

// 全局权限缓存实例
// 缓存用户权限，键为 user_id，值为权限映射 map[string]bool
// 缓存过期时间：10分钟，清理间隔：30分钟
var permissionCache = cache.New(10*time.Minute, 30*time.Minute)

// RBAC 中间件用于检查用户权限 RBAC 通过路由节点自动检查
// 该中间件会在请求上下文中查找用户ID，并根据用户ID查询用户的角色和权限
// 如果用户没有权限访问当前请求的资源，将返回 403 Forbidden 错误
// 如果用户是管理员角色，则跳过权限检查
// 注意：此中间件应该在 JWT 中间件之后使用，以确保用户信息已被正确设置到上下文中
// 例如：在 JWT 中间件中可以设置用户ID到上下文中：c.Set("user_id", user.ID)
// 这样在 Rbac 中间件中就可以直接获取到用户ID
// 如果用户ID不存在或为空，将返回 401 Unauthorized 错误
func RBAC() gin.HandlerFunc {
	return func(c *gin.Context) {

		user_id, exists := c.Get("user_id")
		if !exists || user_id == "" {
			response.Error(c, "缺少用户信息", http.StatusUnauthorized)
			c.Abort()
			return
		}

		// 从缓存获取用户权限
		cacheKey := fmt.Sprintf("permissions:%v", user_id)
		userPermissions, found := permissionCache.Get(cacheKey)
		if !found {
			// 缓存未命中，查询数据库
			var user models.User
			if err := database.DB.Preload("Roles.Permissions").Where("id = ?", user_id).First(&user).Error; err != nil {
				response.Error(c, "用户不存在", http.StatusUnauthorized)
				c.Abort()
				return
			}
			// 超管跳过权限检查
			if user.Username == config.GlobalConfig.SuperAccount {
				c.Next()
				return
			}
			// 收集用户所有权限
			permissions := make(map[string]bool)
			for _, role := range user.Roles {
				for _, perm := range role.Permissions {
					key := perm.Path + ":" + perm.Method
					permissions[key] = true
				}
			}

			// 将权限存入缓存
			permissionCache.Set(cacheKey, permissions, cache.DefaultExpiration)
			userPermissions = permissions
		} else {
			// 缓存命中，直接使用缓存的权限
			permissionsMap, ok := userPermissions.(map[string]bool)
			if !ok {
				// 缓存类型错误，重新查询
				permissionCache.Delete(cacheKey)
				c.Next()
				return
			}
			userPermissions = permissionsMap
		}

		// 类型断言，确保userPermissions是map[string]bool类型
		permissionsMap, _ := userPermissions.(map[string]bool)

		// 检查当前请求是否有权限
		currentPath := c.FullPath()
		currentMethod := c.Request.Method
		requiredPermission := currentPath + ":" + currentMethod

		if permissionsMap[requiredPermission] {
			c.Next()
		} else {
			response.Error(c, "没有访问权限", http.StatusForbidden)
			c.Abort()
			return
		}
	}
}
