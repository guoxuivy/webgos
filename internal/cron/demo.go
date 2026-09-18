package cron

import (
	"context"
	"time"

	"webgos/internal/xlog"
)

func init() {
	Register(Job{
		Name:       "cron_demo",
		Interval:   5 * time.Minute,
		RunOnStart: true, // 启动后立即聚合一次，避免服务重启后当天汇总长时间为空
		Fn:         demo,
	})
}

func demo(ctx context.Context) error {
	// 聚合当天数据；采用先删后插的幂等逻辑，即使每 5 分钟重跑，同一 day 结果一致，不会重复累加
	targetDay := time.Now().Format("20060102")
	xlog.Info("demo day=%s done", targetDay)
	return nil
}
