package database

import (
	"fmt"
	"time"
	"webgos/internal/config"
	"webgos/internal/xlog"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() (*gorm.DB, error) {
	dialector, err := dialector()
	if err != nil {
		return nil, fmt.Errorf("failed to create dialector: %w", err)
	}

	sqlDB, err := gorm.Open(dialector, &gorm.Config{
		Logger: xlog.NewGormLogger(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// 设置连接池参数
	sqlDBInstance, err := sqlDB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}
	dbConfig := config.GlobalConfig.Database
	// 设置连接池参数
	sqlDBInstance.SetMaxOpenConns(dbConfig.MaxOpenConns)                                // 设置最大打开连接数
	sqlDBInstance.SetMaxIdleConns(dbConfig.MaxIdleConns)                                // 设置最大空闲连接数
	sqlDBInstance.SetConnMaxLifetime(time.Duration(dbConfig.MaxLifetime) * time.Minute) // 设置连接的最大生命周期

	DB = sqlDB
	return DB, nil
}

func dialector() (gorm.Dialector, error) {
	dbConfig := config.GlobalConfig.Database
	// 根据配置的 dialect 选择相应的数据库驱动
	var dialector gorm.Dialector
	switch dbConfig.Dialect {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			dbConfig.Username, dbConfig.Password, dbConfig.Host,
			dbConfig.Port, dbConfig.DBName)
		dialector = mysql.Open(dsn)
	case "postgres":
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
			dbConfig.Host, dbConfig.Port, dbConfig.Username,
			dbConfig.Password, dbConfig.DBName)
		dialector = postgres.Open(dsn)
	// case "sqlite":
	// 	dialector = sqlite.Open(dbConfig.DBName)
	// case "sqlserver":
	// 	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
	// 		dbConfig.Username, dbConfig.Password, dbConfig.Host,
	// 		dbConfig.Port, dbConfig.DBName)
	// 	dialector = sqlserver.Open(dsn)
	default:
		return nil, fmt.Errorf("unsupported database dialect: %s", dbConfig.Dialect)
	}
	return dialector, nil
}

func CloseDB() {
	if DB == nil {
		return
	}
	sqlDB, err := DB.DB()
	if err != nil {
		fmt.Println("Failed to get database instance:", err)
		return
	}
	if err := sqlDB.Close(); err != nil {
		fmt.Println("Failed to close database connection:", err)
	}
}
