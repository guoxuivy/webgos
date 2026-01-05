package xlog

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"
	"webgos/internal/config"

	"gorm.io/gorm/logger"
)

// 定义颜色常量（文本色）
const (
	ColorReset  = "\033[0m"  // 重置
	ColorRed    = "\033[31m" // 红色
	ColorGreen  = "\033[32m" // 绿色
	ColorYellow = "\033[33m" // 黄色
	ColorBlue   = "\033[34m" // 蓝色
)

var levelColor = map[string]string{
	"INFO":  ColorGreen,
	"ERROR": ColorRed,
	"WARN":  ColorYellow,
	"DEBUG": ColorYellow,
	"SQL":   ColorBlue,
}

type logMsg struct {
	level   string
	message string
}

// Xlog 自定义日志结构体
type log struct {
	logFiles  map[string]*os.File // 按日志级别存储日志文件句柄
	mu        sync.Mutex          // 添加互斥锁 以确保并发安全
	isDebug   bool                // 新增字段，用于判断是否处于调试模式
	buffer    chan *logMsg        // 日志缓冲通道，用于传递结构化的日志消息
	closeChan chan struct{}       // 关闭信号通道
	doneChan  chan struct{}       // 通知Close() flushLoop已完成
	logDir    string              // 存储日志目录路径
}

var Xlogger *log
var logLevel logger.LogLevel = logger.Info

// 是否开启访问日志
var logAccess bool

// 使用Xlogger创建一个GORM日志实例
func NewGormLogger() logger.Interface {
	if Xlogger == nil {
		panic("Xlogger is not initialized. Please call InitLogger first.")
	}

	globalConfig := config.GlobalConfig
	sqllogLevel := logger.Info
	switch globalConfig.Log.LevelSQL {
	case "Error":
		sqllogLevel = logger.Error
	case "Warn":
		sqllogLevel = logger.Warn
	case "Info":
		sqllogLevel = logger.Info
	case "Silent":
		sqllogLevel = logger.Silent
	default:
		sqllogLevel = logger.Info
	}

	gormLogger := logger.New(Xlogger, logger.Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  sqllogLevel,
		IgnoreRecordNotFoundError: false,
		Colorful:                  true,
	})
	return gormLogger
}

// InitLogger 创建一个新的日志实例
func InitLogger() error {
	logDir := "./logs" // 默认日志目录
	config := config.GlobalConfig
	logAccess = true // 默认开启访问日志
	isDebug := true
	if config != nil {
		logDir = config.Log.Dir
		logAccess = config.Log.Access
		isDebug = config.Server.Mode == "debug"
		switch config.Log.Level {
		case "Error":
			logLevel = logger.Error
		case "Warn":
			logLevel = logger.Warn
		case "Info":
			logLevel = logger.Info
		default:
			logLevel = logger.Info
		}
	}

	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	Xlogger = &log{
		logFiles:  make(map[string]*os.File),
		isDebug:   isDebug,
		buffer:    make(chan *logMsg, 200),
		closeChan: make(chan struct{}),
		doneChan:  make(chan struct{}),
		logDir:    logDir,
	}

	go Xlogger.flushLoop() // 启动后台刷新协程
	return nil
}

// Info 记录信息级别日志
func Info(format string, v ...any) {
	if logLevel >= logger.Info {
		Xlogger.enqueue("INFO", format, v...)
	}
}

// Warn 记录警告级别日志
func Warn(format string, v ...any) {
	if logLevel >= logger.Warn {
		Xlogger.enqueue("WARN", format, v...)
	}
}

// Error 记录错误级别日志
func Error(format string, v ...any) {
	if logLevel >= logger.Error {
		Xlogger.enqueue("ERROR", format, v...)
	}
}

// Debug 记录调试信息
func Debug(format string, v ...any) {
	Xlogger.enqueue("DEBUG", format, v...)
}

// Access 记录访问日志
func Access(format string, v ...any) {
	if logAccess {
		Xlogger.enqueue("ACCESS", format, v...)
	}
}

// gorm SQL专用打印函数
func (l *log) Printf(format string, v ...any) {
	l.enqueue("SQL", format, v...)
}

// enqueue 将日志条目加入缓冲通道
func (l *log) enqueue(level, format string, v ...any) {
	message := fmt.Sprintf(format, v...)
	// now := time.Now().Format("2006-01-02 15:04:05")
	now := time.Now().Format("15:04:05")
	switch level {
	case "ACCESS", "SQL":
		message = fmt.Sprintf("%s %s", now, message)
	default:
		_, file, line, _ := runtime.Caller(2)
		fileName := filepath.Base(file)
		message = fmt.Sprintf("%s %s:%d - %s", now, fileName, line, message)
	}

	msg := &logMsg{
		level:   level,
		message: message,
	}

	if Xlogger.isDebug {
		fmt.Printf("%v[%s]%v %s\n", levelColor[msg.level], msg.level, ColorReset, msg.message)
	}
	// ANSI颜色代码的正则表达式
	var ansiColorRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)
	msg.message = ansiColorRegex.ReplaceAllString(msg.message, "")

	// 将消息发送到缓冲通道
	select {
	case l.buffer <- msg:
	default:
		// 缓冲区满时直接写入（防止阻塞）
		go l.flushBuffer([]*logMsg{msg})
	}

}

// writeToFile 封装日志消息写入文件的逻辑，包括文件的创建、打开和写入操作。
// 参数:
//
//	level: 日志级别(INFO, ERROR等)
//	message: 需要写入的日志消息
func (l *log) writeToFile(level, message string) {
	today := time.Now().Format("2006-01-02")
	logFilePath := filepath.Join(l.logDir, fmt.Sprintf("%s-%s.log", strings.ToLower(level), today))

	file := l.logFiles[level]
	if file != nil {
		// 获取当前文件名以比较日期
		if !strings.HasSuffix(file.Name(), today) {
			file.Close() // 关闭旧文件
			file = nil
		}
	}

	if file == nil {
		var err error
		file, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open log file for level %s: %v\n", level, err)
			return
		}
		l.logFiles[level] = file
	}

	_, err := file.WriteString(message + "\n")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write %s log: %v\n", level, err)
	}
	file.Sync()
}

// flushLoop 后台刷新协程，定时5秒或按量刷新缓冲区
func (l *log) flushLoop() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	var buffer []*logMsg
	for {
		select {
		case msg := <-l.buffer:
			buffer = append(buffer, msg)
			if len(buffer) >= 10 { // 达到10条立即刷新
				l.flushBuffer(buffer)
				buffer = nil
			}
		case <-ticker.C:
			if len(buffer) > 0 {
				l.flushBuffer(buffer)
				buffer = nil
			}
		case <-l.closeChan:
			if len(buffer) > 0 {
				l.flushBuffer(buffer)
			}
			close(l.doneChan) // 通知Close()：已处理完所有日志
			return
		}
	}
}

// flushBuffer 实际执行日志写入
func (l *log) flushBuffer(buffer []*logMsg) {
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, msg := range buffer {
		l.writeToFile(msg.level, msg.message)
	}
}

// Close 关闭日志文件并刷新缓冲区
func (l *log) Close() {
	close(l.closeChan)
	// 等待flushLoop处理完所有日志
	<-l.doneChan

	for _, file := range l.logFiles {
		if file != nil {
			file.Close()
		}
	}
}
