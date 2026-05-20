package middleware

import (
	"strconv"
	"time"
	"webgos/internal/config"
	"webgos/internal/database"
	"webgos/internal/models"
	"webgos/internal/utils/response"

	"github.com/gin-gonic/gin"

	"webgos/internal/cache"
)

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

		userID := c.GetInt("user_id")
		if userID == 0 {
			response.AuthError(c, "缺少用户信息")
			return
		}

		// 从缓存获取用户权限
		cacheKey := "permissions:" + strconv.Itoa(userID)
		userPermissions, found := cache.GetCache().Get(cacheKey)
		var permissions map[string]bool

		if !found {
			// 缓存未命中，查询数据库
			var user models.User
			if err := database.GetDB().Preload("Roles.Permissions").Where("id = ?", userID).First(&user).Error; err != nil {
				response.AuthError(c, "用户不存在")
				return
			}
			// 超管跳过权限检查
			if user.Username == config.GlobalConfig.SuperAccount {
				c.Next()
				return
			}
			// 收集用户所有权限
			permissions = make(map[string]bool)
			for _, role := range user.Roles {
				for _, perm := range role.Permissions {
					key := perm.Path + ":" + perm.Method
					permissions[key] = true
				}
			}
			// 将权限存入缓存
			cache.GetCache().Set(cacheKey, permissions, 5*time.Minute)
		} else {
			// 缓存命中，直接使用缓存的权限
			permissionsMap, ok := userPermissions.(map[string]bool)
			if !ok {
				// 缓存类型错误，重新查询或拒绝请求
				cache.GetCache().Delete(cacheKey)
				response.AuthError(c, "权限验证失败")
				return
			}
			permissions = permissionsMap
		}

		// 检查当前请求是否有权限
		currentPath := c.FullPath()
		currentMethod := c.Request.Method
		requiredPermission := currentPath + ":" + currentMethod

		if permissions[requiredPermission] {
			c.Next()
		} else {
			response.Forbidden(c, "没有访问权限")
			return
		}
	}
}
