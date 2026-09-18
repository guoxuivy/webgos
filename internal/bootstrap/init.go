package bootstrap

import (
	"fmt"
	"webgos/internal/config"
	"webgos/internal/cron"
	"webgos/internal/routes"
	"webgos/internal/xdb"
	"webgos/internal/xdb/migrate"
	"webgos/internal/xlog"
)

func Initialize(configPath string) error {
	// 配置初始化,必须最先执行！
	globalConfig, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// 初始化日志
	if err = xlog.InitLogger(); err != nil {
		return fmt.Errorf("failed to initialize logger: %v", err)
	}

	// 初始化数据库
	if err = xdb.InitDB(); err != nil {
		return fmt.Errorf("Database initialization error: %v", err)
	}

	// 自动迁移模型
	if err = migrate.AutoMigrate(); err != nil {
		return fmt.Errorf("Model migration error: %v", err)
	}

	// 注册路由
	routes.New(globalConfig)

	// 启动定时任务：任务在各文件 init 中 Register，此处统一拉起
	cron.SetUp()

	return nil
}

func Close() {
	xlog.Access("Closing resources...")
	xdb.CloseDB()
	// 等待在途定时任务跑完再关日志组件，否则 ShutDown 的收尾日志会打到已关闭的 logger
	cron.ShutDown()
	if xlog.Xlogger != nil {
		xlog.Xlogger.Close()
	}
}
