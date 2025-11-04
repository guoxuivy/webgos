package database

import (
	"fmt"
	"time"
	"webgos/internal/xlog"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB(dsn string) (*gorm.DB, error) {
	gormLogger := xlog.NewGormLogger()

	sqlDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// 设置连接池参数
	sqlDBInstance, err := sqlDB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// 设置连接池参数
	sqlDBInstance.SetMaxOpenConns(10)           // 设置最大打开连接数
	sqlDBInstance.SetMaxIdleConns(5)            // 设置最大空闲连接数
	sqlDBInstance.SetConnMaxLifetime(time.Hour) // 设置连接的最大生命周期

	DB = sqlDB
	return DB, nil
}

func CloseDB() {
	sqlDB, err := DB.DB()
	if err != nil {
		fmt.Println("Failed to get database instance:", err)
		return
	}
	if err := sqlDB.Close(); err != nil {
		fmt.Println("Failed to close database connection:", err)
	}
}
