// main.go
package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
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
		Addr:    fmt.Sprintf(":%d", globalConfig.Server.Port),
		Handler: routes.REngine,
	}

	// 按需开启Swagger文档
	routes.REngine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger/doc.json")))

	quit := make(chan os.Signal, 1)
	// kill -SIGINT 或 kill -SIGTERM 会触发优雅关闭 kill <pid> 或 kill -2 <pid>
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
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
