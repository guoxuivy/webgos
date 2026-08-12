// main.go
package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
	"webgos/internal/bootstrap"
	"webgos/internal/config"
	"webgos/internal/routes"
	"webgos/internal/xlog"

	_ "webgos/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title webgos API
// @version 1.0
// @description webgos 企业资源计划系统 API 文档

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT认证方式，值为"Bearer {token}"

// @BasePath /
func main() {
	defer bootstrap.Close()
	// 解析命令行参数
	configPath := flag.String("c", "./config/config.yaml", "Specify the config file path")
	flag.Parse()

	// 初始化项目
	if err := bootstrap.Initialize(*configPath); err != nil {
		panic("failed to initialize project: " + err.Error())
	}
	globalConfig := config.GlobalConfig

	// 创建 http.Server
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", globalConfig.Server.Port),
		Handler:           routes.REngine,
		ReadTimeout:       30 * time.Second,  // 读取请求超时
		WriteTimeout:      60 * time.Second,  // 写入响应超时
		IdleTimeout:       120 * time.Second, // 连接空闲超时（keep-alive 连接保持时间）
		ReadHeaderTimeout: 10 * time.Second,  // 读取请求头超时
		MaxHeaderBytes:    1 << 20,           // 最大请求头大小（1MB）
	}

	// 按需开启Swagger文档
	if globalConfig.Server.Swag {
		routes.REngine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger/doc.json")))
	}

	quit := make(chan os.Signal, 1)
	// kill -SIGINT 或 kill -SIGTERM 会触发优雅关闭 kill <pid> 或 kill -2 <pid>
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 按需开启 pprof 性能分析（独立 debug 端口，不影响业务路由）
	if globalConfig.Server.Pprof {
		pprofServer(quit)
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		<-quit
		xlog.Access("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			xlog.Access("Server forced to shutdown: %v", err)
		}
		xlog.Access("Server exiting")
		close(idleConnsClosed)
	}()

	xlog.Access("Server started on port %d", globalConfig.Server.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("listen: %s\n", err)
	}
	<-idleConnsClosed
}

func pprofServer(quit chan os.Signal) {
	pprofMux := http.NewServeMux()
	pprofMux.HandleFunc("/debug/pprof/", pprof.Index)
	pprofMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	pprofMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	pprofMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	pprofMux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	pprofMux.Handle("/debug/pprof/block", pprof.Handler("block"))
	pprofMux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	pprofMux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	pprofMux.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	pprofMux.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))

	pprofSrv := &http.Server{
		Addr:    fmt.Sprintf(":%d", 6060),
		Handler: pprofMux,
	}
	go func() {
		xlog.Access("pprof server started on port 6060")
		if err := pprofSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			xlog.Access("pprof server listen error: %v", err)
		}
	}()
	go func() {
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = pprofSrv.Shutdown(ctx)
	}()
}
