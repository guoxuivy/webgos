package xlog

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm/logger"
)

type gormLog struct {
	logLevel logger.LogLevel
}

func NewGormLogger() logger.Interface {
	return &gormLog{
		logLevel: logger.Info, // 默认日志级别
	}
}

func (l *gormLog) LogMode(level logger.LogLevel) logger.Interface {
	// 返回一个新的日志实例，设置日志级别
	newLogger := *l
	newLogger.logLevel = level
	return &newLogger
}

func (l *gormLog) Info(ctx context.Context, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	SQL("[Info] %v", message)
}

func (l *gormLog) Warn(ctx context.Context, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	SQL("[Warn] %v", message)
}

func (l *gormLog) Error(ctx context.Context, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	SQL("%s[Error]%s %v", ColorRed, ColorReset, message)
}

func (l *gormLog) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, rows := fc()
	elapsed := time.Since(begin)
	message := fmt.Sprintf("SQL Query sql %v rows %v duration %v error %v", sql, rows, elapsed.String(), err)
	if err != nil {
		l.Error(ctx, "%v", message)
		return
	}
	SQL("[Trace] %v", message)
}
