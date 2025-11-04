package middleware

import (
	"net/http"
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
			response.Error(c, "缺少认证令牌", http.StatusUnauthorized)
			c.Abort()
			return
		}

		// 检查Bearer token格式
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Error(c, "令牌格式错误", http.StatusUnauthorized)
			c.Abort()
			return
		}

		// 解析token
		tokenString := parts[1]
		c.Set("tokenString", tokenString)

		// token验证
		service := services.NewAuthService()
		claims, err := service.ValidateToken(tokenString)
		if err != nil {
			response.Error(c, err.Error(), http.StatusUnauthorized)
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		// JWT中的数字默认是float64类型 转换为int类型
		c.Set("user_id", int(claims["user_id"].(float64)))
		c.Set("username", claims["username"])

		c.Next()
	}
}
