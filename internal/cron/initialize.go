package cron

import (
	"context"
	"sync"
	"time"

	"webgos/internal/xlog"
)

// Job 定时任务定义。
// 每个 Job 独占一个 goroutine，由调度器按间隔串行驱动：
// 上一次执行结束后才重新计时，因此既不会并发重入，也不会因单次执行超时而背靠背连跑。
type Job struct {
	Name       string                          // 任务名，用于日志排障
	Interval   time.Duration                   // 执行间隔（从上次执行结束开始计时）
	RunOnStart bool                            // 是否在启动时立即执行一次
	Fn         func(ctx context.Context) error // 任务逻辑，需同步执行完成
}

// jobs 已注册的定时任务列表
var jobs []Job

// Register 注册定时任务（在各任务文件的 init 中调用）
func Register(j Job) {
	jobs = append(jobs, j)
}

var (
	stopOnce sync.Once
	stopCh   = make(chan struct{})
	wg       sync.WaitGroup
)

// SetUp 启动所有已注册的定时任务
func SetUp() {
	for _, j := range jobs {
		wg.Add(1)
		go func(j Job) {
			defer wg.Done()
			timer := time.NewTimer(j.Interval)
			defer timer.Stop()

			// 启动后立即执行一次（可选项），随后进入定时循环
			if j.RunOnStart {
				runJob(j)
			}

			for {
				select {
				case <-stopCh:
					return
				case <-timer.C:
					// 同步调用：本次执行完才可能进入下一次，天然串行、无并发重入
					runJob(j)
					// 执行结束后重新计时：保证两次执行至少间隔 Interval，
					// 避免执行耗时超过 Interval 时 Ticker 补跑导致的连续执行
					timer.Reset(j.Interval)
				}
			}
		}(j)
	}
	xlog.Info("cron setup done, jobs=%d", len(jobs))
}

// ShutDown 优雅停止所有定时任务，并等待在途任务执行结束（进程退出前调用）。
// 注意是"等待"而非"打断"：正在执行的任务会跑完，保证聚合类任务的完整性。
func ShutDown() {
	stopOnce.Do(func() {
		close(stopCh)
	})
	wg.Wait()
	xlog.Info("cron stopped")
}

// runJob 执行单个任务：panic 恢复 + 耗时与错误日志。
// recover 必须在此处（而非 goroutine 顶部）：否则 panic 后 goroutine 直接退出，
// 该任务将永远失去调度。
func runJob(j Job) {
	start := time.Now()
	defer func() {
		if err := recover(); err != nil {
			xlog.Error("cron job=%s panic: %v", j.Name, err)
		}
		xlog.Info("cron job=%s finished cost=%s", j.Name, time.Since(start))
	}()
	if err := j.Fn(context.Background()); err != nil {
		xlog.Error("cron job=%s err: %v", j.Name, err)
	}
}
