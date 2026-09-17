# webgos 统一命令入口
#
# 说明：本仓库没有 CI，也没有 git 钩子，make check 靠自觉执行。
# 跳过它，docs/ 下的规范就只是一堆文字。

BIN := webgos
GO := go
# 用 @version 隔离模式运行 swag：不污染 go.mod / go.sum，也不需要全局安装 swag CLI
SWAG_PKG := github.com/swaggo/swag/cmd/swag@v1.16.6
SWAG_OUT := internal/swagger

ifeq ($(OS),Windows_NT)
	BIN := webgos.exe
endif

.DEFAULT_GOAL := help
.PHONY: help fmt vet test build run swagger swag check

help: ## 显示可用命令
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2}'

fmt: ## gofmt 格式化
	$(GO) fmt ./...

vet: ## go vet 静态检查
	$(GO) vet ./...

test: ## 运行单元测试
	$(GO) test ./... -count=1

build: ## 编译二进制
	$(GO) build -o $(BIN) cmd/main.go

run: ## 本地运行（配置默认 config/config.yaml）
	$(GO) run cmd/main.go -c ./config/config.yaml

swag: ## 重新生成 Swagger 文档到 internal/swagger（改了注解必须跑）
	$(GO) run $(SWAG_PKG) init -g cmd/main.go -o $(SWAG_OUT)

swagger: swag ## 兼容旧命令，等价于 make swag

check: fmt vet test build ## 一键体检：fmt + vet + test + build
