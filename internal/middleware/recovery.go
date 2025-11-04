package middleware

import (
	"net/http"
	"runtime/debug"
	"webgos/internal/utils/response"
	"webgos/internal/xlog"

	"github.com/gin-gonic/gin"
)

// Recovery 中间件用于捕获和处理 panic
// 如果发生 panic，将记录错误并返回 500 错误响应
// 该中间件应该在所有其他中间件之后使用，以确保能够捕获所有的 panic
// 注意：此中间件不应该重新生成 request_id，因为 RequestID 中间件已经处理了 request_id 的生成和存储
// 如果需要在 panic 时获取 request_id，可以直接从上下文中获取
// 例如：requestID, _ := c.Get("request_id")
// 这样可以确保在 panic 时仍然能够获取到正确的 request_id
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否已存在request_id，如果不存在则不重新生成
		// RequestIDMiddleware应该已经设置了request_id
		defer func() {
			if err := recover(); err != nil {
				requestID, _ := c.Get("request_id")
				stackInfo := debug.Stack()

				xlog.Error("%v Recovered from panic: %v", requestID, err)
				xlog.Error("异常堆栈：%v\n", string(stackInfo))

				response.Error(c, "内部服务器错误", http.StatusInternalServerError)
				c.Abort()
			}
		}()
		c.Next()
	}
}
