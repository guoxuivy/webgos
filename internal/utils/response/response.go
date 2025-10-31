package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code      int    `json:"code"`
	Msg       string `json:"message"`
	Data      any    `json:"data,omitempty"`
	RequestID string `json:"request_id"`
}

// Success 成功响应
func Success(c *gin.Context, message string, data any) {
	requestID := GetRequestID(c)
	c.JSON(http.StatusOK, Response{
		Code:      http.StatusOK,
		Msg:       message,
		Data:      data,
		RequestID: requestID,
	})
}

// Error 错误响应
// 参数:
//   - c: gin 上下文
//   - message: 错误消息
//   - code: 可变参数，第一个是HTTP状态码，第二个是 业务 状态码
func Error(c *gin.Context, message string, code ...int) {
	requestID := GetRequestID(c)

	// 业务状态码，默认为 400
	bizCode := 400
	// HTTP 状态码，默认为 400
	statusCode := http.StatusBadRequest

	if len(code) > 0 {
		bizCode = code[0]
		statusCode = code[0] // 默认业务码和HTTP状态码相同
	}

	if len(code) > 1 {
		bizCode = code[1] // 可以分别指定业务码和HTTP状态码
	}

	c.JSON(statusCode, Response{
		Code:      bizCode,
		Msg:       message,
		RequestID: requestID,
	})
	c.Abort()
}

// ErrorWithData 带数据的错误响应
// 参数:
//   - c: gin 上下文
//   - message: 错误消息
//   - data: 返回的数据
//   - code: 可变参数，第一个是HTTP状态码，第二个是 业务 状态码
func ErrorWithData(c *gin.Context, message string, data any, code ...int) {
	requestID := GetRequestID(c)

	// 业务状态码，默认为 400
	bizCode := 400
	// HTTP 状态码，默认为 400
	statusCode := http.StatusBadRequest

	if len(code) > 0 {
		bizCode = code[0]
		statusCode = code[0] // 默认业务码和HTTP状态码相同
	}

	if len(code) > 1 {
		bizCode = code[1] // 可以分别指定业务码和HTTP状态码
	}

	c.JSON(statusCode, Response{
		Code:      bizCode,
		Msg:       message,
		Data:      data,
		RequestID: requestID,
	})
	c.Abort()
}

// GetRequestID 从上下文中获取RequestID
func GetRequestID(c *gin.Context) string {
	requestID, exists := c.Get("request_id")
	if exists {
		return requestID.(string)
	}
	return ""
}
