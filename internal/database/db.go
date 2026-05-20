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

var (
	MasterDB *gorm.DB   // 主库连接
	SlaveDBs []*gorm.DB // 备库连接池
)

// InitDB 初始化数据库连接
func InitDB() error {
	dialector, err := dialector()
	if err != nil {
		return fmt.Errorf("failed to create dialector: %w", err)
	}

	sqlDB, err := gorm.Open(dialector, &gorm.Config{
		Logger: xlog.NewGormLogger(),
	})
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	// 设置连接池参数
	sqlDBInstance, err := sqlDB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}
	dbConfig := config.GlobalConfig.Database
	// 设置连接池参数
	sqlDBInstance.SetMaxOpenConns(dbConfig.MaxOpenConns)                                // 设置最大打开连接数
	sqlDBInstance.SetMaxIdleConns(dbConfig.MaxIdleConns)                                // 设置最大空闲连接数
	sqlDBInstance.SetConnMaxLifetime(time.Duration(dbConfig.MaxLifetime) * time.Minute) // 设置连接的最大生命周期

	MasterDB = sqlDB
	return nil
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
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
			dbConfig.Host, dbConfig.Port, dbConfig.Username,
			dbConfig.Password, dbConfig.DBName)
		// 1. 连接数据库（v1.6.0+ 驱动配置）
		dialector = postgres.New(postgres.Config{
			DSN:                  dsn,
			PreferSimpleProtocol: true, // 依然需要启用文本协议
		})
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
	if MasterDB == nil {
		return
	}
	sqlDB, err := MasterDB.DB()
	if err != nil {
		fmt.Println("Failed to get database instance:", err)
		return
	}
	if err := sqlDB.Close(); err != nil {
		fmt.Println("Failed to close database connection:", err)
	}
	for _, slaveDB := range SlaveDBs {
		sqlDB, err = slaveDB.DB()
		if err != nil {
			fmt.Println("Failed to get database instance:", err)
			continue
		}
		if err := sqlDB.Close(); err != nil {
			fmt.Println("Failed to close slave database connection:", err)
		}
	}
}
