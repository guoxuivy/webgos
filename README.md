# hserp 项目文档

## 项目概述
hserp 是一个基于 Go 语言开发的企业资源计划（ERP）系统，主要用于企业内部的库存、产品和用户管理。该项目采用现代化的开发技术和架构，提供高效、可维护的企业级解决方案。

## 项目特点
- **技术先进**：使用 Go 语言开发，基于 Gin 框架和 GORM ORM 工具
- **架构清晰**：采用 MVC 架构模式，分层设计清晰
- **配置灵活**：支持 YAML 格式的配置文件
- **易于扩展**：分层设计，便于新增功能集成
- **维护性强**：统一的响应格式和完善的错误处理机制
- **安全可靠**：基于 RBAC 的权限管理机制，自动化注册路由节点为权限点
- **文档完善**：集成 Swagger API 文档，便于接口调试和使用
- **日志系统**：完善的自定义日志记录系统，支持请求追踪和问题排查
- **优雅关闭**：支持服务的优雅启动和关闭

## 技术栈
- **编程语言**：Go 1.24.4
- **Web 框架**：Gin v1.10.1
- **ORM 框架**：GORM v1.30.0
- **数据库**：MySQL 驱动（gorm.io/driver/mysql v1.5.2）
- **配置管理**：YAML 格式（gopkg.in/yaml.v3 v3.0.1）
- **数据验证**：go-playground/validator v10
- **API 文档**：Swagger (github.com/swaggo/gin-swagger)
- **测试框架**：testify (github.com/stretchr/testify)

## 目录结构
```
hserp/
├── cmd/                            # 可执行文件相关代码
│   └── main.go                     # 程序入口文件
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
│   │   ├── rbac.go                 # RBAC相关DTO
│   │   └── user.go                 # 用户相关DTO
│   ├── handlers/                   # HTTP请求处理器
│   │   ├── inventory.go            # 库存相关请求处理
│   │   ├── product.go              # 产品相关请求处理
│   │   ├── rbac.go                 # RBAC相关请求处理
│   │   ├── test.go                 # 测试相关请求处理
│   │   └── user.go                 # 用户相关请求处理
│   ├── middleware/                 # Gin框架中间件
│   │   ├── auth.go                 # rbac权限认证中间件
│   │   ├── cors.go                 # 跨域中间件
│   │   ├── debounce.go             # 防抖中间件
│   │   ├── jwt.go                  # JWT登录认证中间件
│   │   ├── logging.go              # 日志记录中间件
│   │   ├── middleware.go           # 中间件接口定义
│   │   ├── recovery.go             # 恢复中间件
│   │   └── requestid.go            # 请求ID中间件
│   ├── models/                     # 数据访问层
│   │   ├── base_model.go           # 基础模型
│   │   ├── inventory_record.go     # 库存记录数据模型
│   │   ├── product.go              # 产品数据模型
│   │   ├── rbac.go                 # RBAC权限数据模型
│   │   └── user.go                 # 用户数据模型
│   ├── routes/                     # 路由注册
│   │   ├── router_wrapper.go       # 路由注册rbac包装器
│   │   └── routes.go               # 路由注册和管理
│   ├── services/                   # 业务逻辑层
│   │   ├── inventory.go            # 库存业务逻辑
│   │   ├── product.go              # 产品业务逻辑
│   │   ├── rbac.go                 # RBAC业务逻辑
│   │   └── user.go                 # 用户业务逻辑
│   ├── utils/                      # 公共工具函数
│   │   ├── response/               # 响应处理
│   │   │   └── response.go         # 统一响应格式
│   │   ├── utils.go                # 通用工具函数
│   │   └── validator.go            # 通用dto参数验证器
│   └── xlog/                       # 日志工具
│       ├── glog.go                 # gorm日志接入实现
│       └── xlog.go                 # 日志系统实现
├── readme/                         # 详细功能说明文档
│   ├── 商品管理.md                  # 商品管理功能说明
│   ├── 库存管理.md                  # 库存管理功能说明
│   └── 权限管理.md                  # 权限管理功能说明
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
- **用户管理**：用户注册、登录、登出、JWT认证
- **产品管理**：产品信息的增删改查
- **库存管理**：库存记录的查询与更新、出入库操作
- **权限管理**：基于RBAC的权限控制系统，角色和权限管理

## 系统架构

本系统采用分层架构设计，遵循标准的MVC模式：

1. **表现层（Handlers）**：处理HTTP请求和响应，参数验证
2. **业务逻辑层（Services）**：实现核心业务逻辑
3. **数据访问层（Models）**：封装数据库操作
4. **数据传输层（DTO）**：定义接口数据结构
5. **基础设施层**：包括配置管理、数据库连接、日志系统等

各层之间通过接口和依赖注入进行解耦，保证系统的可维护性和可扩展性。

## RBAC权限管理系统

本系统采用基于角色的访问控制（RBAC）模型实现权限管理，通过用户、角色和权限的多对多关系，实现灵活的权限控制。

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
1. **RequestID中间件**：为每个请求生成唯一标识，用于日志追踪
2. **Recovery中间件**：捕获系统panic，防止服务崩溃
3. **Logging中间件**：记录请求日志，便于问题追踪
4. **CORS中间件**：处理跨域请求
5. **JWT中间件**：处理用户身份认证
6. **Auth中间件**：处理RBAC权限验证
7. **Debounce中间件**：防止重复提交

### 中间件执行顺序
```
RequestID -> Recovery -> Logging -> CORS -> JWT -> Auth -> 业务处理 -> Auth -> JWT -> CORS -> Logging -> Recovery -> RequestID
```

## 数据库连接池配置

项目中使用 GORM 管理数据库连接，并配置了连接池参数：

### 连接池配置
- **最大打开连接数**：10
- **最大空闲连接数**：5
- **连接的最大生命周期**：1小时

这些参数可以根据实际需求在 [database/db.go](file:///d:/Goroot/hserp/internal/database/db.go) 文件中进行调整。

## 统一响应格式

系统采用统一的JSON响应格式：
```json
{
  "code": 200,
  "msg": "success",
  "data": {}
}
```

其中：
- `code`：HTTP状态码或自定义业务状态码
- `msg`：响应消息，成功时通常为"success"，失败时为错误描述
- `data`：响应数据，可以是对象、数组或null

## BaseModel 核心功能

BaseModel 是系统中所有数据模型的基础类，提供了通用的数据库操作方法。

### 核心特性
1. **泛型支持**：使用 Go 的泛型特性支持不同类型的模型
2. **CRUD 操作**：提供基本的增删改查操作
3. **链式查询**：支持链式调用构建复杂查询
4. **事务支持**：提供事务处理方法
5. **并发安全**：通过创建新实例实现并发安全的链式调用

### 主要方法
- `Create(item *T) error` - 创建记录
- `Read(id int) (*T, error)` - 根据ID读取记录
- `Update(item *T) error` - 更新记录
- `Delete(id int) error` - 删除记录（软删除）
- `More() ([]T, error)` - 查询多条记录
- `One() (*T, error)` - 查询单条记录
- `Count() (int64, error)` - 统计记录数
- `Page(page, pageSize int) ([]T, error)` - 分页查询

### 链式查询方法
- `Where(query any, args ...any) IActiveRecode[T]` - 添加WHERE条件
- `Select(query any, args ...any) IActiveRecode[T]` - 指定查询字段
- `Order(value any) IActiveRecode[T]` - 添加排序条件
- `Limit(limit int) IActiveRecode[T]` - 限制返回记录数
- `Group(query string) IActiveRecode[T]` - 添加分组条件
- `Joins(query string, args ...any) IActiveRecode[T]` - 添加JOIN连接查询

### 事务方法
- `Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) (err error)` - 在事务中执行数据库操作
- `WithTx(tx *gorm.DB) *BaseModel[T]` - 将查询器与事务对象绑定

## 开发与部署

### 开发环境要求
- Go 1.24.4 或更高版本
- MySQL 5.7 或更高版本
- Git 版本管理工具

### 配置文件
在 [config/config.yaml](file:///d:/Goroot/hserp/config/config.yaml) 中配置数据库连接、服务端口、JWT密钥等参数。

### 构建和运行
```bash
# 克隆项目
git clone <项目地址>

# 进入项目目录
cd hserp

# 安装依赖
go mod tidy

# 构建项目
go build -o hserp cmd/main.go

# 运行项目
./hserp
```

### 命令行参数
```bash
# 指定配置文件路径
./hserp -c ./config/config.yaml
```

### 优雅关闭
项目支持优雅关闭，使用 `kill <pid>` 或 `Ctrl+C` 可以安全关闭服务。

## 测试

项目包含单元测试和集成测试，位于 [tests](file:///d:/Goroot/hserp/tests/) 目录下。

### 单元测试

单元测试位于 [tests/unit](file:///d:/Goroot/hserp/tests/unit/) 目录下，主要测试各个模块的功能：

1. **Handlers测试**：测试HTTP请求处理器函数
2. **Services测试**：测试业务逻辑层函数
3. **Models测试**：测试数据访问层函数
4. **Utils测试**：测试工具函数
5. **BaseModel测试**：测试基础模型的核心功能

单元测试不依赖外部服务，运行速度快，主要用于测试业务逻辑的正确性。

### 集成测试

集成测试位于 [tests/integration](file:///d:/Goroot/hserp/tests/integration/) 目录下，主要测试需要依赖外部服务的功能：

1. **数据库集成测试**：测试与数据库的交互
2. **BaseModel集成测试**：测试 BaseModel 的实际数据库操作，包括：
   - CRUD 操作测试（创建、读取、更新、删除）
   - 查询操作测试（Where 条件查询、分页查询等）
   - 链式调用测试（验证链式查询方法的正确性）
   - 事务操作测试（验证事务的提交和回滚功能）

集成测试需要连接真实的数据库服务，用于验证系统各组件之间的协作是否正常。

### 运行测试

```bash
# 运行所有测试
go test ./tests/... -v

# 运行单元测试
go test ./tests/unit/... -v

# 运行集成测试
go test ./tests/integration/... -v

# 运行特定测试
go test -v ./tests/integration/base_model_integration_test.go -run TestBaseModelCRUDIntegration
```

### 测试覆盖率

```bash
# 生成测试覆盖率报告
go test ./tests/... -coverprofile=coverage.out

# 在浏览器中查看覆盖率报告
go tool cover -html=coverage.out
```

### 测试环境配置

集成测试需要配置数据库连接，当前测试使用与主应用相同的数据库配置进行测试。在实际项目中，应该使用独立的测试数据库以避免影响生产数据。

测试环境会自动初始化日志系统和数据库连接，测试完成后会清理测试数据和关闭连接。

## Swagger API 文档

本项目已集成 Swagger，启动服务后可访问：

- [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

### 使用方法

1. 安装 swag 工具：`go install github.com/swaggo/swag/cmd/swag@latest`
2. 在项目根目录执行：`swag init -g cmd/main.go`
3. 启动服务后访问 `http://localhost:8080/swagger/index.html` 查看接口文档

接口注释示例见 `internal/handlers/*.go` 文件。

## 日志系统

系统采用自定义日志系统实现结构化日志记录：

### 日志级别
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

私有项目

## 常用命令
```bash
# 安装依赖
go mod tidy
# 构建项目
go build -o hserp.exe cmd/main.go
#linux 构建
GOOS=linux GOARCH=amd64 go build -o hserp cmd/main.go
# 安全停止服务
kill <pid> 
# 启动服务
nohup ./hserp -c ./config.yaml > /dev/null 2>&1 &

```