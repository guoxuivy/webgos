README.md
# webgos 项目文档

## 项目概述
webgos 是一个基于 Go 的企业级 Web 系统快速开发脚手架，基于 Gin 和 GORM，面向可扩展的业务开发。项目目标是提供一套工程化、可测试、可扩展的模板，包含常见的鉴权、日志、请求追踪、事务、分页、统一响应等能力，帮助团队快速落地业务。

## 项目特点
- **技术先进**：使用 Go 1.25 语言开发，基于 Gin 框架和 GORM ORM 工具
- **架构清晰**：采用 MVC 架构模式，分层设计清晰
- **配置灵活**：支持 YAML 格式的配置文件
- **易于扩展**：分层设计，便于新增功能集成
- **维护性强**：统一的响应格式和完善的错误处理机制
- **安全可靠**：基于 RBAC 的权限管理机制，自动化注册路由节点为权限点
- **文档完善**：集成 Swagger API 文档，便于接口调试和使用
- **日志系统**：完善的自定义日志记录系统，支持请求追踪和问题排查
- **优雅关闭**：支持服务的优雅启动和关闭
- **事务支持**：提供便捷的事务处理能力，支持事务嵌套
- **模型封装**：模型通过嵌入 `BaseFields` 获得通用基础字段，数据库操作使用 GORM 原生 API

## 技术栈
- **编程语言**：Go 1.25
- **Web 框架**：Gin
- **ORM 框架**：GORM
- **数据库**：MySQL（通过 gorm.io/driver/mysql）
- **缓存**：github.com/patrickmn/go-cache
- **配置管理**：YAML（gopkg.in/yaml.v3）
- **数据验证**：github.com/go-playground/validator/v10
- **API 文档**：Swaggo（swag + gin-swagger）
- **JWT 认证**：github.com/golang-jwt/jwt/v4
- **UUID 生成**：github.com/google/uuid
- **测试框架**：testify

## 目录结构
```
webgos/
├── cmd/                            # 可执行文件相关代码
│   └── main.go                     # 程序入口文件
├── common/                         # 常用工具代码收集
│   ├── convert/                    # 数据类型转换工具
│   │   └── convert.go              # 类型转换函数
│   ├── json/                       # JSON处理工具
│   │   └── json.go                 # JSON编解码封装
│   ├── syncx/                      # 并发安全工具
│   │   ├── lockedcalls.go          # 并发安全调用工具
│   │   └── singleflight.go         # 防止缓存击穿工具
│   └── time/                       # 时间处理工具
│       └── time.go                 # 时间格式化和解析函数
├── config/                         # 配置管理
│   └── config.yaml                 # 主配置文件
├── internal/                       # 核心业务逻辑代码
│   ├── bootstrap/                  # 项目启动初始化
│   │   └── init.go                 # 项目初始化逻辑
│   ├── config/                     # 配置加载和验证
│   │   └── config.go               # 配置管理实现
│   ├── database/                   # 数据库相关代码
│   │   ├── migrate/                # 数据库迁移
│   │   │   └── migrate.go          # 模型自动迁移逻辑
│   │   └── db.go                   # 数据库连接和迁移逻辑
│   ├── dto/                        # 数据传输对象
│   │   ├── inventory.go            # 库存相关DTO
│   │   ├── menu.go                 # 菜单相关DTO
│   │   ├── rbac.go                 # RBAC相关DTO
│   │   └── user.go                 # 用户相关DTO
│   ├── handlers/                   # HTTP请求处理器
│   │   ├── auth.go                 # 认证相关请求处理
│   │   ├── inventory.go            # 库存相关请求处理
│   │   ├── menu.go                 # 菜单相关请求处理
│   │   ├── product.go              # 产品相关请求处理
│   │   ├── rbac.go                 # RBAC相关请求处理
│   │   ├── test.go                 # 测试相关请求处理
│   │   └── user.go                 # 用户相关请求处理
│   ├── middleware/                 # Gin框架中间件
│   │   ├── auth.go                 # rbac权限认证中间件
│   │   ├── cors.go                 # 跨域中间件
│   │   ├── debounce.go             # 防抖中间件
│   │   ├── gzip.go                 # Gzip压缩中间件
│   │   ├── jwt.go                  # JWT登录认证中间件
│   │   ├── limiter.go              # IP令牌桶限流中间件
│   │   ├── logging.go              # 日志记录中间件
│   │   ├── middleware.go           # 中间件注册与接口定义
│   │   ├── recovery.go             # 恢复中间件
│   │   ├── requestid.go            # 请求ID中间件
│   │   └── security.go             # 安全防范：敏感路径检测、恶意IP自动封禁、IP黑名单
│   ├── models/                     # 数据访问层
│   │   ├── base_fields.go          # 基础字段结构体
│   │   ├── inventory_record.go     # 库存记录数据模型
│   │   ├── menu.go                 # 菜单数据模型
│   │   ├── product.go              # 产品数据模型
│   │   ├── rbac.go                 # RBAC权限数据模型
│   │   └── user.go                 # 用户数据模型
│   ├── routes/                     # 路由注册
│   │   ├── router_wrapper.go       # 路由注册rbac包装器
│   │   └── routes.go               # 路由注册和管理
│   ├── services/                   # 业务逻辑层
│   │   ├── inventory.go            # 库存业务逻辑
│   │   ├── menu.go                 # 菜单业务逻辑
│   │   ├── product.go              # 产品业务逻辑
│   │   ├── rbac.go                 # RBAC业务逻辑
│   │   └── user.go                 # 用户业务逻辑
│   ├── utils/                      # 公共工具函数
│   │   ├── code/                   # 业务状态码
│   │   │   └── bizcode.go          # 业务状态码定义
│   │   ├── param/                  # 参数处理
│   │   │   └── validator.go        # 参数验证器
│   │   └── response/               # 响应处理
│   │       └── response.go         # 统一响应格式
│   └── xlog/                       # 日志工具
│       └── xlog.go                 # 日志系统实现
├── readme/                         # 详细功能说明文档
│   ├── 商品管理.md                  # 商品管理功能说明
│   ├── 库存管理.md                  # 库存管理功能说明
│   ├── 权限管理.md                  # 权限管理功能说明
│   └── 菜单管理.md                  # 菜单管理功能说明
├── docs/                           # API文档
│   └── swagger/                    # Swagger自动生成的API文档
├── tests/                          # 测试代码
│   ├── integration/                # 集成测试
│   └── unit/                       # 单元测试
├── go.mod                          # Go模块定义文件
├── go.sum                          # Go依赖校验文件
├── README.md                       # 项目说明文档
└── doc.go                          # Swagger文档入口文件
```

## 主要功能模块
- **用户管理**：用户注册、登录、登出、JWT 认证
- **产品管理**：产品信息的增删改查
- **库存管理**：库存记录的查询与更新、出入库操作
- **权限管理**：基于 RBAC 的权限控制系统，路由自动注册为权限点

## 系统架构

项目采用分层架构（类似 MVC）：

1. **Handlers（表现层）**：处理 HTTP 请求、参数验证与响应（Gin）
2. **Services（业务层）**：组织业务逻辑、事务边界、调用模型层
3. **Models（数据访问层）**：定义数据结构并嵌入 `BaseFields` 获取基础字段，数据库操作通过 GORM 原生 API 在 Service 层完成
4. **DTO（数据传输对象）**：处理入参与出参结构定义与验证
5. **Infrastructure（基础设施）**：配置、数据库连接、日志、middleware 等

各层通过接口与注入解耦，便于单元测试与替换实现。

## RBAC 权限管理系统

项目实现基于角色的访问控制（RBAC）。路由在注册时会被收集并同步为权限点，权限标识采用 `路径:HTTP方法`（例如 `/api/products:GET`）。

### 核心概念
- **用户（User）**：系统的使用者，可以被分配一个或多个角色
- **角色（Role）**：一组权限的集合，可以分配给一个或多个用户
- **权限（Permission）**：系统中最小的不可再分的访问控制单元，由路由节点自动生成

### 权限自动生成机制
系统实现了基于路由的权限点自动生成机制：
1. 在路由注册时收集路由信息
2. 系统启动时将收集到的路由信息同步到数据库作为权限点
3. 如果权限点已存在，则更新其描述信息；如果不存在，则创建新权限点

权限标识采用 `路径:HTTP方法` 的格式，例如：
- `/api/products:GET` - 查看商品列表
- `/api/products:POST` - 创建商品

### 权限验证流程
1. 使用JWT进行用户身份认证
2. 通过RBAC中间件进行权限检查
3. 查询当前用户是否拥有访问当前路径和方法的权限
4. 如果有权限，则继续处理请求；否则返回403错误

## 数据验证机制

系统使用 [go-playground/validator](https://github.com/go-playground/validator) 库进行数据验证，提供以下特性：

### 验证功能特点
- 支持结构体字段验证
- 支持自定义验证规则
- 支持友好的错误消息显示
- 支持字段标签自定义（使用label标签作为字段名）

### 自定义验证规则
系统已实现以下自定义验证规则：
- **手机号验证**：使用`phone`标签验证中国手机号格式

### 错误消息处理
验证器支持通过`label`标签来自定义字段显示名称，并提供友好的错误提示信息：
- required: "为必填项"
- min: "长度不能少于{n}个字符" 或 "不能小于{n}"
- max: "长度不能超过{n}个字符" 或 "不能大于{n}"
- email: "格式不正确"
- gte: "必须大于等于{n}"
- lte: "必须小于等于{n}"
- oneof: "必须是{n}中的一个"

## 项目初始化规范
- `main.go`负责程序入口和启动，调用`bootstrap.Initialize()`函数
- 将项目初始化逻辑集中到`bootstrap.Initialize()`函数中
- `main.go`应保持简洁，只负责程序入口和启动
- `bootstrap.Initialize()`函数接收配置参数，处理以下初始化任务:
  - 日志系统初始化
  - 数据库连接初始化
  - 模型自动迁移
  - Gin路由初始化
  - 路由注册
  - 权限点同步
- 配置依赖应显式传递，避免使用全局变量
- Initialize()函数应能接收不同的配置参数，提高灵活性
- 数据库初始化和路由设置的依赖关系要明确
- 统一项目初始化错误处理机制
- main.go应通过`bootstrap.R`获取初始化后的路由引擎
- 更清晰的错误日志输出方式

## 中间件系统

系统实现了多种中间件来处理请求的前置和后置逻辑：

### 核心中间件

**全局中间件（在 `middleware.ApplyMiddlewares` 中注册，按以下顺序执行）**
1. **IPBlacklist中间件**：最高优先级拦截，拒绝黑名单中的 IP（精确 IP 与 CIDR 网段），黑名单持久化到 `blacklist.json` 并每 5 分钟自动保存
2. **RequestID中间件**：为每个请求生成唯一标识，用于日志追踪
3. **Recovery中间件**：捕获系统 panic，防止服务崩溃
4. **Logging中间件**：记录请求日志，便于问题追踪
5. **CORS中间件**：处理跨域请求
6. **Gzip中间件**：对响应进行 Gzip 压缩，降低传输体积

**安全防范中间件**
- **敏感路径检测（CheckSensitivePath）**：在全局 404 handler 中调用，匹配 `.env`、`.git`、`phpmyadmin`、`wp-admin`、`.sql`、备份/压缩包等敏感路径与 `/shell`、`/exec` 等危险关键字；命中按时间窗口（1 小时）计数，达到阈值（5 次）自动将该 IP 加入黑名单
- **IP 限流（IPLimiter）**：基于令牌桶（每个 IP 独立桶）的限流，用于 `/api/auth/login` 等高风险路由防爆破，例如 `IPLimiter(1, 1)` 表示每秒 1 个请求、桶容量 1（不允许突发），超限返回 500

**路由分组中间件**
- **JWT中间件**：处理用户身份认证（登录态校验）
- **Auth中间件**：处理 RBAC 权限验证（路由即权限点）
- **Debounce中间件**：防止重复提交

### 中间件执行顺序
```
全局：IPBlacklist -> RequestID -> Recovery -> Logging -> CORS -> Gzip
路由组（如 /api/auth）：JWT -> Auth -> IPLimiter -> 业务处理
404 处理：CheckSensitivePath（命中则记录并可能触发 IP 自动封禁）
```

## 数据库连接池配置

项目中使用 GORM 管理数据库连接，并配置了连接池参数：

### 连接池配置
- **最大打开连接数**：10
- **最大空闲连接数**：5
- **连接的最大生命周期**：1小时

这些参数可以根据实际需求在 [database/db.go](file:///d:/Goroot/webgos/internal/database/db.go) 文件中进行调整。

## 统一响应格式

系统采用统一的 JSON 响应格式：

```json
{
  "code": 200,
  "msg": "success",
  "data": {}
}
```

- `code`：优先表示业务状态码（项目约定），同时可映射为 HTTP 状态码；请参阅 `internal/utils/response/response.go` 的实现。
- `msg`：描述信息
- `data`：返回的数据体

注意：项目中已经有 `response` 工具用于统一封装响应，调用方应按照库提供的方法传入明确的 HTTP 状态码与业务码。

## 模型层与数据库操作说明

模型定义放在 `internal/models/` 中，各模型通过嵌入 `BaseFields` 结构体获得 `ID`、`CreatedAt`、`UpdatedAt`、`DeletedAt` 等通用基础字段，自身只声明业务字段。

数据库操作不再经过泛型 `BaseModel` 封装，而是由 **Service 层直接使用 GORM 原生 API**（如 `xdb.GetDB().WithContext(ctx).Where(...).Find(&items)`）完成，事务通过 `xdb.GetDB().Transaction(func(tx *gorm.DB) error { ... })` 处理。

Service 层所有业务方法首参数统一为 `ctx context.Context`，用于请求级超时与取消。Handler 调用时直接传入 `c *gin.Context` 即可（`gin.Context` 实现了 `context.Context` 接口），无需取 `c.Request.Context()`；Service 内部通过 `xdb.GetDB().WithContext(ctx)` 使用上下文，不直接引用 `*gin.Context`，因此 Service 与 Gin 保持解耦，可独立测试。

模型定义示例：

```go
type User struct {
	models.BaseFields
	Username string `gorm:"size:64;uniqueIndex" json:"username"`
	Password string `json:"-"`
}
```

Service 层使用 GORM 原生操作示例：

```go
func (s *userService) GetUserInfo(ctx context.Context, id int) (*models.User, error) {
	var u models.User
	if err := xdb.GetDB().WithContext(ctx).First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}
```

## 开发与部署

### 开发环境要求
- Go 1.25 或更高版本
- MySQL 5.7 或更高版本
- Git 版本管理工具

### 配置文件
在 [config/config.yaml](file:///d:/Goroot/webgos/config/config.yaml) 中配置数据库连接、服务端口、JWT密钥等参数。

### 构建和运行
```bash
# 克隆项目
git clone <项目地址>

# 进入项目目录
cd webgos

# 安装依赖
go mod tidy

# 构建项目
go build -o webgos cmd/main.go

# 运行项目
./webgos
```

### 命令行参数
```bash
# 指定配置文件路径
./webgos -c ./config/config.yaml
```

### 优雅关闭
项目支持优雅关闭，使用 `kill <pid>` 或 `Ctrl+C` 可以安全关闭服务。

## 测试

项目包含单元测试和集成测试，目录为 `tests/unit` 与 `tests/integration`。注意：直接对 `./tests` 顶层运行 `go test ./tests` 会失败（因为顶层目录没有 Go 包文件），应使用子包路径或 `./...` 模式。

运行建议：

```bash
# 运行所有测试（包括 tests 下的子包）
go test ./... -v

# 仅运行所有单元测试
go test ./tests/unit/... -v

# 仅运行所有集成测试
go test ./tests/integration/... -v

# 生成覆盖率
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

集成测试会访问数据库，请务必使用独立的测试数据库并在测试完成后清理测试数据。


## 常用命令
```bash
# 安装依赖
go mod tidy

# 构建（Windows）
go build -o webgos.exe cmd/main.go

# 构建（Linux）
GOOS=linux GOARCH=amd64 go build -o webgos cmd/main.go

# 运行（示例）
./webgos -c ./config/config.yaml

# 安全停止（示例）
kill <pid>

# 后台运行（示例）
nohup ./webgos -c ./config/config.yaml > /dev/null 2>&1 &
```

## Swagger 文档生成

如果你需要生成 Swagger 文档并查看接口说明，可以使用 swag 工具（swaggo）。下面是常用的步骤：

```bash
# 安装 swag（仅需一次）
go install github.com/swaggo/swag/cmd/swag@latest

# 在项目根目录生成 Swagger 注释（默认会在 ./docs 目录生成）
swag init -g cmd/main.go

# 生成后启动服务并访问 Swagger UI：
# 访问: http://localhost:8080/swagger/index.html

```

## 日志

- **ACCESS**：访问日志，记录请求处理信息
- **INFO**：常规信息日志
- **ERROR**：错误日志
- **DEBUG**：调试日志
- **WARN**：警告日志
- **SQL**：SQL执行日志

### 日志格式
日志文件按日期和级别分割存储在 `logs` 目录下，格式为：
```
[ACCESS] RequestID=550e8400-e29b-41d4-a716-446655440000 [POST] /users/login 192.168.1.100 200 45ms
```

### 日志系统特性
1. **多级别日志记录**：支持不同级别的日志记录和过滤
2. **日志文件分割**：按日期和级别分割日志文件，便于管理
3. **控制台彩色输出**：在控制台中以不同颜色显示不同级别的日志
4. **异步日志写入**：通过缓冲通道实现异步日志写入，提高性能
5. **请求追踪**：与RequestID中间件配合，实现请求全链路追踪

## 许可证
本项目采用随便玩许可证，详情见 [LICENSE](LICENSE) 文件。

## 开发命令
```bash
swag init -g cmd/main.go
# 访问: http://localhost:8080/swagger/index.html

go run cmd/main.go -c ./config/config.yaml
```