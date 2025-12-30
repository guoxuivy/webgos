package middleware

import (
	"strings"
	"webgos/internal/services"
	"webgos/internal/utils/response"

	"github.com/gin-gonic/gin"
)

// JWT中间件
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.AuthError(c, "缺少认证令牌")
			return
		}

		// 检查Bearer token格式
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.AuthError(c, "令牌格式错误")
			return
		}

		// 解析token
		tokenString := parts[1]
		c.Set("tokenString", tokenString)

		// token验证
		service := services.NewAuthService()
		claims, err := service.ValidateToken(tokenString)
		if err != nil {
			response.AuthError(c, err.Error())
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", int((*claims)["user_id"].(float64)))
		c.Set("username", (*claims)["username"].(string))
		c.Next()
	}
}
