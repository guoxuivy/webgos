package middleware

import (
	"net/http"
	"sync"
	"webgos/common/bucketx"
	"webgos/internal/xlog"

	"github.com/gin-gonic/gin"
)

// IPLimiter 基于令牌桶的IP限流中间件
// rate: 每秒填充的令牌数（即每秒允许的请求数）
// capacity: 桶容量，决定允许的瞬时突发量（一般和 rate 一致表示不允许突发）
// 每个IP独立维护一个令牌桶，请求消耗1个令牌，令牌不足时返回 500
func IPLimiter(rate, capacity int) gin.HandlerFunc {
	var mu sync.Mutex
	buckets := make(map[string]bucketx.TokenBucket)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		mu.Lock()
		b, ok := buckets[ip]
		if !ok {
			b = bucketx.NewTokenBucket(rate, capacity)
			buckets[ip] = b
		}
		mu.Unlock()

		if b.TryTake(1) {
			c.Next()
		} else {
			xlog.Warn("[SECURITY] IP %s 超过请求频率限制", ip)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "服务器内部错误",
			})
		}
	}
}
