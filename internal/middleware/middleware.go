package middleware

import (
	"hserp/internal/config"

	"github.com/gin-gonic/gin"
)

// 应用全局中间件包含请求ID、恢复、日志记录和跨域
func ApplyMiddlewares(r *gin.Engine, config *config.Config) {
	// 请求ID中间件
	r.Use(RequestID())

	// 自定义恢复中间件，用于捕获panic
	r.Use(Recovery())

	// 自定义日志记录中间件
	r.Use(Logging())

	// 跨域中间件
	r.Use(CORS())

}
