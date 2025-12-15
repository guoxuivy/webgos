package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// 压缩中间件 需要处理函数前调用
func Gzip() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查客户端是否支持 gzip 压缩
		if !strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
			c.Next()
			return
		}

		// 创建 gzip writer
		gz := gzip.NewWriter(c.Writer)
		defer gz.Close()
		// 设置响应头
		c.Header("Content-Encoding", "gzip")
		c.Header("Vary", "Accept-Encoding")

		// 替换 Writer
		c.Writer = &gzipWriter{c.Writer, gz}
		// 处理请求
		c.Next()

	}
}

type gzipWriter struct {
	gin.ResponseWriter
	gz *gzip.Writer
}

func (w *gzipWriter) Write(data []byte) (int, error) {
	// 如果没有设置 Content-Type，则自动检测
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", http.DetectContentType(data))
	}
	return w.gz.Write(data)
}
