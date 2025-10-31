package middleware

import (
	"hserp/internal/utils/response"
	"hserp/internal/xlog"
	"time"

	"github.com/gin-gonic/gin"
)

// 访问日志记录中间件
func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求ID
		requestID := response.GetRequestID(c)

		// 记录请求开始时间
		start := time.Now()

		c.Next()

		// 请求处理完成后记录一条包含所有信息的日志，运行时间精确到毫秒
		duration := time.Since(start) / time.Millisecond
		xlog.Access("RequestID=%s [%s] %s %s %d %dms", requestID, c.Request.Method, c.Request.URL.Path, getClientIP(c), c.Writer.Status(), duration)

		// 如果发生错误，记录错误日志
		if len(c.Errors) > 0 {
			xlog.Error("RequestID=%s ERROR [%s] %s: %s", requestID, c.Request.Method, c.Request.URL.Path, c.Errors.String())
		}
	}
}

// getClientIP 获取客户端IP
func getClientIP(c *gin.Context) string {
	ip := c.ClientIP()
	if ip == "" {
		ip = c.GetHeader("X-Forwarded-For")
	}
	if ip == "" {
		ip = c.GetHeader("X-Real-IP")
	}
	return ip
}
