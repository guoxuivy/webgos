package bootstrap

import (
	"fmt"
	"webgos/internal/config"
	"webgos/internal/database"
	"webgos/internal/database/migrate"
	"webgos/internal/routes"
	"webgos/internal/xlog"

	"github.com/gin-gonic/gin"
)

var R *gin.Engine

func Initialize(config *config.Config) error {
	// 使用传入的配置
	// 构建DSN字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Database.Username, config.Database.Password, config.Database.Host,
		config.Database.Port, config.Database.DBName)

	// 初始化日志
	logDir := config.Server.LogDir
	if logDir == "" {
		logDir = "./logs" // 默认日志目录
	}
	err := xlog.InitLogger(logDir, config.Server.Mode == "debug") // 将logger替换为xlog
	if err != nil {
		panic(fmt.Sprintf("failed to initialize logger: %v", err))
	}

	// 初始化数据库
	_, err = database.InitDB(dsn)
	if err != nil {
		panic(fmt.Sprintf("Database initialization error: %v", err))
	}

	// 自动迁移模型（根据配置决定是否执行）
	if config.AutoMigrate {
		xlog.Access("Starting auto migration...")
		if err := migrate.AutoMigrate(); err != nil {
			panic(fmt.Sprintf("Model migration error: %v", err))
		}
		xlog.Access("Auto migration completed")
	} else {
		xlog.Access("Auto migration is disabled")
	}

	// 注册路由
	R = routes.SetupRoutes(config)

	// 同步权限到数据库（根据配置决定是否收集）
	if err := routes.SyncPermissions(database.DB); err != nil {
		panic(fmt.Sprintf("Failed to sync permissions: %v", err))
	}

	return nil
}

func Close() {
	xlog.Access("Closing resources...")
	database.CloseDB()
	xlog.Xlogger.Close()
}
