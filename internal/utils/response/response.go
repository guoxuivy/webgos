package response

import (
	"net/http"
	"webgos/internal/utils/code"
	"webgos/internal/xlog"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code      int    `json:"code"`
	Msg       string `json:"message"`
	Data      any    `json:"data,omitempty"`
	RequestID string `json:"request_id"`
}

func Resp(c *gin.Context, code int, msg string, data any) {
	requestID := GetRequestID(c)
	resp := Response{}
	resp.Code = code
	resp.Msg = msg
	resp.Data = data
	resp.RequestID = requestID
	if code != 0 {
		xlog.Error("request err: requestID %s url %s method %s code %d msg %s", requestID, c.Request.URL, c.Request.Method, code, msg)
	}
	c.AbortWithStatusJSON(http.StatusOK, resp)
}

// Success 成功响应
func Success(c *gin.Context, msg string, data any) {
	Resp(c, code.OK, msg, data)
}

func Error(c *gin.Context, msg string) {
	Resp(c, code.Error, msg, "")
}

func ErrorWithCode(c *gin.Context, msg string, errCode int) {
	Resp(c, errCode, msg, "")
}

// 用户登录错误
func AuthError(c *gin.Context, msg string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, msg)
}

// 权限不足错误
func Forbidden(c *gin.Context, msg string) {
	c.AbortWithStatusJSON(http.StatusForbidden, msg)
}

// GetRequestID 从上下文中获取RequestID
func GetRequestID(c *gin.Context) string {
	requestID, exists := c.Get("request_id")
	if exists {
		return requestID.(string)
	}
	return ""
}
