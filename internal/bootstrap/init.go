package bootstrap

import (
	"fmt"
	"webgos/internal/config"
	"webgos/internal/database"
	"webgos/internal/database/migrate" // 添加数据库迁移包导入
	"webgos/internal/routes"
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
	if _, err = database.InitDB(); err != nil {
		return fmt.Errorf("Database initialization error: %v", err)
	}

	// 自动迁移模型
	if err = migrate.AutoMigrate(); err != nil {
		return fmt.Errorf("Model migration error: %v", err)
	}

	// 注册路由
	routes.New(globalConfig)

	// 同步权限到数据库（根据配置决定是否收集）
	if err := routes.SyncPermissions(database.DB); err != nil {
		return fmt.Errorf("Failed to sync permissions: %v", err)
	}

	return nil
}

func Close() {
	xlog.Access("Closing resources...")
	database.CloseDB()
	if xlog.Xlogger != nil {
		xlog.Xlogger.Close()
	}
}
