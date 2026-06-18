package bucketx

import (
	"sync"
	"time"
)

// TokenBucket 令牌桶限流器接口
// 以恒定速率填充令牌，调用方消耗令牌，令牌不足时拒绝请求
type TokenBucket interface {
	// TryTake 尝试获取 n 个令牌，返回是否成功
	TryTake(n int) bool
	// Available 返回当前可用令牌数
	Available() int
	// Rate 返回填充速率（令牌/秒）
	Rate() int
	// Capacity 返回桶的最大容量
	Capacity() int
	// Reset 重置令牌桶到满容量状态
	Reset()
}

// bucket 令牌桶实现
type bucket struct {
	mu       sync.Mutex
	rate     int       // 每秒填充的令牌数
	capacity int       // 桶最大容量
	tokens   float64   // 当前令牌数
	lastTime time.Time // 上次填充时间
}

// NewTokenBucket 创建一个新的令牌桶
// rate: 每秒填充的令牌数
// capacity: 桶的最大容量（初始令牌数 = capacity）
func NewTokenBucket(rate, capacity int) TokenBucket {
	return &bucket{
		rate:     rate,
		capacity: capacity,
		tokens:   float64(capacity),
		lastTime: time.Now(),
	}
}

// refill 根据经过的时间补充令牌（内部方法，调用前需持有锁）
func (b *bucket) refill() {
	now := time.Now()
	elapsed := now.Sub(b.lastTime).Seconds()
	b.tokens += elapsed * float64(b.rate)
	if b.tokens > float64(b.capacity) {
		b.tokens = float64(b.capacity)
	}
	b.lastTime = now
}

// TryTake 尝试从桶中获取 n 个令牌
// 返回 true 表示获取成功，false 表示令牌不足
func (b *bucket) TryTake(n int) bool {
	if n <= 0 {
		return true
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	b.refill()

	if b.tokens >= float64(n) {
		b.tokens -= float64(n)
		return true
	}
	return false
}

// Available 返回当前可用令牌数（会先执行填充）
func (b *bucket) Available() int {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.refill()
	return int(b.tokens)
}

// Rate 返回令牌填充速率
func (b *bucket) Rate() int {
	return b.rate
}

// Capacity 返回桶的最大容量
func (b *bucket) Capacity() int {
	return b.capacity
}

// Reset 重置令牌桶，清空已有令牌并填满
func (b *bucket) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.tokens = float64(b.capacity)
	b.lastTime = time.Now()
}
