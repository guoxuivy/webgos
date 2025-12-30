package syncx

import "sync"

// LockedCalls 提供基于 key 的锁定调用，确保同一时间只有一个协程执行相同 key 的函数
// 比如 用户账户操作串行化、同一文件的多个写操作串行化 等场景
type LockedCalls interface {
	Do(key string, fn func() (any, error)) (any, error)
}
type lockedGroup struct {
	mu sync.Mutex
	m  map[string]*sync.WaitGroup
}

func NewLockedCalls() LockedCalls {
	return &lockedGroup{
		m: make(map[string]*sync.WaitGroup),
	}
}

func (lg *lockedGroup) Do(key string, fn func() (any, error)) (any, error) {
begin:
	lg.mu.Lock()
	if wg, ok := lg.m[key]; ok {
		// case 1: 该 key 正在被其他协程处理
		lg.mu.Unlock()
		wg.Wait()  // 等待其完成
		goto begin // 重新尝试获取锁
	}

	// case 2: 该 key 没有被处理，当前协程获得执行权
	return lg.makeCall(key, fn)
}

func (lg *lockedGroup) makeCall(key string, fn func() (any, error)) (any, error) {
	var wg sync.WaitGroup
	wg.Add(1)
	lg.m[key] = &wg // 标记该 key 正在处理
	lg.mu.Unlock()

	defer func() {
		// 注意顺序：先删除 key，再 Done()
		// 如果反过来，可能有协程 Wait() 返回但 key 还在 map 中
		lg.mu.Lock()
		delete(lg.m, key)
		lg.mu.Unlock()
		wg.Done()
	}()

	return fn() // 执行实际函数
}
