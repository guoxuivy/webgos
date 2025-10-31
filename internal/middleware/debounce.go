package middleware

import (
	"hserp/internal/utils/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

// 防抖中间件使用的缓存  默认过期时间5分钟，每10分钟清理一次
// 这里使用了第三方库 go-cache 来实现简单的内存缓存
var debounceCache = cache.New(5*time.Minute, 10*time.Minute)

// DebounceMiddleware 防抖中间件
// duration: 可选参数，指定防抖时间窗口，默认500毫秒
// key: user_id + ":" + URL.path
func Debounce(timeout ...time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		var duration time.Duration
		if len(timeout) > 0 {
			duration = timeout[0]
		} else {
			duration = 500 * time.Millisecond // 默认防抖时间窗口为500ms
		}
		// 获取路由路径
		path := c.Request.URL.Path

		user_id := c.GetString("user_id")
		if user_id == "" {
			user_id = c.ClientIP() // 匿名用户用IP
		}
		key := user_id + "@" + path

		// 检查是否已存在相同请求
		if _, found := debounceCache.Get(key); found {
			response.Error(c, "请求过于频繁，请稍后再试", http.StatusTooManyRequests)
			c.Abort()
			return
		}
		// 将请求ID存入缓存，设置防抖时间窗口
		debounceCache.Set(key, true, duration)
		// 继续处理请求
		c.Next()
	}
}
