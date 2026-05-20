package middleware

import (
	"net/http"
	"strconv"
	"time"
	"webgos/internal/utils/response"

	"github.com/gin-gonic/gin"

	"webgos/internal/cache"
)

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

		// user_id := c.GetString("user_id")
		userID := c.GetInt("user_id")
		user_id := strconv.Itoa(userID)
		if userID == 0 {
			user_id = c.ClientIP() // 匿名用户用IP
		}
		key := user_id + "@" + path

		// 检查是否已存在相同请求
		if _, found := cache.GetCache().Get(key); found {
			response.ErrorWithCode(c, "请求过于频繁，请稍后再试", http.StatusTooManyRequests)
			return
		}
		// 将请求ID存入缓存，设置防抖时间窗口
		cache.GetCache().Set(key, true, duration)
		// 继续处理请求
		c.Next()
	}
}
