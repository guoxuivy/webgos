package middleware

import (
	"net/http"
	"webgos/internal/database"
	"webgos/internal/models"
	"webgos/internal/utils/response"

	"github.com/gin-gonic/gin"
)

// Rbac 中间件用于检查用户权限 RBAC 通过路由节点自动检查
// 该中间件会在请求上下文中查找用户ID，并根据用户ID查询用户的角色和权限
// 如果用户没有权限访问当前请求的资源，将返回 403 Forbidden 错误
// 如果用户是管理员角色，则跳过权限检查
// 注意：此中间件应该在 JWT 中间件之后使用，以确保用户信息已被正确设置到上下文中
// 例如：在 JWT 中间件中可以设置用户ID到上下文中：c.Set("user_id", user.ID)
// 这样在 Rbac 中间件中就可以直接获取到用户ID
// 如果用户ID不存在或为空，将返回 401 Unauthorized 错误
func Rbac() gin.HandlerFunc {
	return func(c *gin.Context) {

		user_id, exists := c.Get("user_id")
		if !exists || user_id == "" {
			response.Error(c, "缺少用户信息", http.StatusUnauthorized)
			c.Abort()
			return
		}

		// 查询用户及其角色权限
		// todo 优化：可以将用户角色权限添加缓存，减少数据库查询
		var user models.User
		if err := database.DB.Preload("Roles.Permissions").Where("id = ?", user_id).First(&user).Error; err != nil {
			response.Error(c, "用户不存在", http.StatusUnauthorized)
			c.Abort()
			return
		}

		// 检查是否为管理员角色(可以跳过权限检查)
		isAdmin := false
		for _, role := range user.Roles {
			if role.Name == "Super" {
				isAdmin = true
				break
			}
		}
		if isAdmin {
			// 管理员拥有所有权限
			c.Next()
			return
		}
		//  收集用户所有权限
		userPermissions := make(map[string]bool)
		for _, role := range user.Roles {
			for _, perm := range role.Permissions {
				key := perm.Path + ":" + perm.Method
				// key := perm.Name //Name不作为权限验证依据，只为方便作权限点的分组管理
				userPermissions[key] = true
			}
		}

		// 检查当前请求是否有权限
		currentPath := c.FullPath()
		currentMethod := c.Request.Method
		requiredPermission := currentPath + ":" + currentMethod

		if userPermissions[requiredPermission] {
			c.Next()
		} else {
			response.Error(c, "没有访问权限", http.StatusForbidden)
			c.Abort()
			return
		}
	}
}
