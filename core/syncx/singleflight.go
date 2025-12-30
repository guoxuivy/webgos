package syncx

import "sync"

// 多个请求只执行一次，共享第一个请求的结果
// 不适用 执行时间很长的操作（>5秒）
type SingleFlight interface {
	Do(key string, fn func() (any, error)) (any, error)
}

func NewSingleFlight() SingleFlight {
	return &flightGroup{
		calls: make(map[string]*call),
	}
}

type call struct {
	wg  sync.WaitGroup // 用于等待第一个请求完成
	val any            // 缓存结果值
	err error          // 缓存错误
}

type flightGroup struct {
	calls map[string]*call // key -> 正在执行的调用
	lock  sync.Mutex       // 保护 calls map
}

func (g *flightGroup) Do(key string, fn func() (any, error)) (any, error) {
	// 1. 尝试创建或获取已存在的 call
	c, done := g.createCall(key)
	if done {
		// 已经有其他协程在执行了，直接返回结果
		return c.val, c.err
	}

	// 2. 第一个协程执行实际的函数
	g.makeCall(c, key, fn)
	return c.val, c.err
}

func (g *flightGroup) createCall(key string) (c *call, done bool) {
	g.lock.Lock()
	if c, ok := g.calls[key]; ok {
		// case 1: 已经有请求在执行中
		g.lock.Unlock()
		c.wg.Wait() // 等待第一个请求完成
		return c, true
	}

	// case 2: 第一个请求，创建新的 call
	c = new(call)
	c.wg.Add(1)      // 设置等待计数
	g.calls[key] = c // 注册到 map
	g.lock.Unlock()

	return c, false
}

func (g *flightGroup) makeCall(c *call, key string, fn func() (any, error)) {
	defer func() {
		g.lock.Lock()
		delete(g.calls, key) // 执行完成，从 map 中移除
		g.lock.Unlock()
		c.wg.Done() // 通知等待的协程
	}()

	c.val, c.err = fn() // 执行实际函数
}
