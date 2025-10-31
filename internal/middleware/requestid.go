package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 为每个请求生成唯一的RequestID
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查请求头中是否已存在RequestID
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			// 如果不存在，则生成一个新的UUID
			requestID = uuid.New().String()
		}

		// 将RequestID存储在上下文中
		c.Set("request_id", requestID)

		// 将RequestID添加到响应头中
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}
