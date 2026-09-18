# webgos 项目文档

> **AI 与新同学请先读 [`AGENTS.md`](AGENTS.md)（唯一入口）。**

## 项目概述

webgos 是一个基于 Go 的企业级 Web 系统快速开发脚手架，基于 Gin 和 GORM，面向可扩展的业务开发。项目目标是提供一套工程化、可测试、可扩展的模板，包含常见的鉴权、日志、请求追踪、事务、分页、统一响应等能力，帮助团队快速落地业务。

## 项目特点

- **技术先进**：使用 Go 1.25 语言开发，基于 Gin 框架和 GORM ORM 工具
- **架构清晰**：分层设计（Routes → Handlers → Services → Models，DTO 承载出入参）
- **配置灵活**：支持 YAML 格式的配置文件
- **易于扩展**：分层设计，便于新增功能集成
- **维护性强**：统一的响应格式和完善的错误处理机制
- **安全可靠**：基于 RBAC 的权限管理机制，路由自动注册为权限点
- **文档完善**：集成 Swagger API 文档，便于接口调试和使用
- **日志系统**：自定义日志记录系统，支持请求追踪和问题排查
- **优雅关闭**：支持服务的优雅启动和关闭
- **事务支持**：提供便捷的事务处理能力，支持事务嵌套
- **读写分离**：主库写、从库读，Service 层显式选库
- **模型封装**：模型嵌入 `BaseFields` 获得通用基础字段，数据库操作使用 GORM 原生 API

## 技术栈

- **编程语言**：Go 1.25
- **Web 框架**：Gin v1.11
- **ORM 框架**：GORM v1.31
- **数据库**：MySQL / PostgreSQL（`go.mod` 同时引入两个驱动，实际由配置决定）
- **缓存**：github.com/patrickmn/go-cache
- **配置管理**：YAML（gopkg.in/yaml.v3）
- **数据验证**：github.com/go-playground/validator/v10
- **API 文档**：Swaggo（swag + gin-swagger）
- **JWT 认证**：github.com/golang-jwt/jwt/v4
- **UUID 生成**：github.com/google/uuid
- **测试框架**：testify

## 目录结构

> 只列到**目录级**，不逐文件列举——文件清单会随代码快速漂移（历史教训：本节曾列出
> `internal/database/`，而实际路径一直是 `internal/xdb/`）。

```
webgos/
├── AGENTS.md                       # AI / 新人工作入口（唯一入口）
├── Makefile                        # 统一命令入口（make check / make swagger 等）
├── cmd/                            # 程序入口（含 Swagger 通用注解 @title 等）
├── common/                         # 通用工具（bucketx / convert / json / syncx / time）
├── config/                         # 配置（config.yaml 不入库，见 .gitignore）
├── internal/                       # 核心业务代码
│   ├── bootstrap/                  # 项目启动初始化
│   ├── cache/                      # 缓存封装
│   ├── config/                     # 配置加载与校验
│   ├── cron/                       # 定时任务（各文件 init 里 Register，bootstrap 统一启停）
│   ├── dto/                        # 数据传输对象
│   ├── handlers/                   # HTTP 处理器
│   ├── middleware/                 # 中间件
│   ├── models/                     # 数据模型
│   ├── routes/                     # 路由注册（自动同步为 RBAC 权限点）
│   ├── services/                   # 业务逻辑层
│   ├── swagger/                    # Swagger 生成物（make swagger 生成，勿手改）
│   ├── utils/                      # 工具：code / file / param / response
│   ├── xdb/                        # 数据库句柄（主库 GetDB / 从库 GetSlaveDB）
│   │   └── migrate/                # 数据库迁移
│   └── xlog/                       # 日志
├── docs/                           # 开发规范文档（100% 手写，无生成物）
│   ├── backend-conventions.md      # 后端规范唯一权威源
│   ├── dev-workflow.md             # 开发流程与收工清单
│   ├── changes/                    # 一页纸变更记录
│   └── templates/                  # 代码模板（crud-api.md）
├── readme/                         # 补充功能说明（商品/库存/权限/菜单/pprof）
├── tests/
│   └── unit/                       # 单元测试（integration 目录暂为空）
├── go.mod / go.sum
└── README.md
```

## 快速开始

### 环境要求

- Go 1.25 或更高版本
- MySQL 5.7+ 或 PostgreSQL（按所用驱动选择）
- Git

### 配置

在 `config/config.yaml` 中配置数据库连接、服务端口、JWT 密钥等参数。
该文件**不入库**（见 `.gitignore`），首次部署需自行创建。

### 常用命令

全部统一走 `Makefile`：

```bash
make help      # 查看全部命令
make check     # fmt + vet + test + build，一键体检
make test      # go test ./... -count=1
make build     # 编译
make run       # 本地运行（默认 config/config.yaml）
make swagger   # 重新生成 API 文档到 internal/swagger
make fmt       # gofmt
make vet       # go vet
```

也可以直接构建后运行：

```bash
go build -o webgos cmd/main.go
./webgos -c ./config/config.yaml
```

服务启动后，若配置开启 Swagger，可访问 `http://localhost:<port>/swagger/index.html`。

## 开发规范与 AI 协作

**AI 与新同学请先读 [`AGENTS.md`](AGENTS.md)（唯一入口）。**

| 内容 | 位置 |
| --- | --- |
| 后端规范（唯一权威源） | [`docs/backend-conventions.md`](docs/backend-conventions.md) |
| 开发流程与收工清单 | [`docs/dev-workflow.md`](docs/dev-workflow.md) |
| 一页纸变更模板 | [`docs/changes/`](docs/changes/README.md) |
| CRUD 代码生成模板 | [`docs/templates/crud-api.md`](docs/templates/crud-api.md)（**仅为模板，不定义规范**） |

约定：

- `docs/` 是**纯手写**规范区；`internal/swagger/` 是 `make swagger` 的**生成物，勿手改**。两者已物理隔离。
- 历史入口 `readme/backend_rules.md` 已迁移为存根，内容以 `docs/backend-conventions.md` 为准。
- 本仓库**没有 CI 与 git 钩子**，`make check` 靠自觉执行。

> 下文各章节只做**概览**；涉及"怎么写代码"的细节一律以 `docs/backend-conventions.md` 为准，
> 避免产生第二份权威源。

## 主要功能模块

- **用户管理**：用户注册、登录、登出、JWT 认证
- **产品管理**：产品信息的增删改查
- **库存管理**：库存记录的查询与更新、出入库操作
- **权限管理**：基于 RBAC 的权限控制系统，路由自动注册为权限点

## 系统架构

分层架构：

1. **Handlers（表现层）**：处理 HTTP 请求、参数验证与响应（Gin）
2. **Services（业务层）**：组织业务逻辑、事务边界、数据库操作
3. **Models（数据层）**：定义数据结构并嵌入 `BaseFields`
4. **DTO**：入参与出参结构定义与验证
5. **Infrastructure**：配置、数据库连接、日志、middleware 等

各层通过接口解耦，Service 只依赖标准库 `context.Context`（不引用 `*gin.Context`），便于独立测试。
细节见 [`backend-conventions.md` §分层架构](docs/backend-conventions.md)。

## RBAC 权限管理系统

路由在注册时会被收集并同步为权限点，权限标识采用 `路径:HTTP方法`（例如 `/api/products:GET`）。

### 核心概念

- **用户（User）**：系统的使用者，可以被分配一个或多个角色
- **角色（Role）**：一组权限的集合，可以分配给一个或多个用户
- **权限（Permission）**：最小的访问控制单元，由路由节点自动生成

### 权限自动生成机制

1. 路由注册时收集路由信息；
2. 系统启动时将路由信息同步到数据库作为权限点；
3. 权限点已存在则更新描述，不存在则创建。

示例：`/api/products:GET` - 查看商品列表；`/api/products:POST` - 创建商品。

### 权限验证流程

1. 使用 JWT 进行用户身份认证；
2. 通过 RBAC 中间件进行权限检查；
3. 查询当前用户是否拥有该路径与方法的权限；
4. 有权限则继续处理，否则返回 403。

## 数据验证机制

使用 [go-playground/validator](https://github.com/go-playground/validator) 进行入参校验，
统一入口 `param.Validate(c, &dto)`。

- **自定义规则**：已实现 `phone`（中国手机号格式）。
- **`label` 标签**：用于错误消息中的中文字段名。
- **完整验证标签表**见 [`backend-conventions.md` §参数验证标签](docs/backend-conventions.md)。

### 错误消息处理

验证器根据 `label` 生成友好提示：

- required: "为必填项"
- min: "长度不能少于{n}个字符" 或 "不能小于{n}"
- max: "长度不能超过{n}个字符" 或 "不能大于{n}"
- email: "格式不正确"
- gte: "必须大于等于{n}"
- lte: "必须小于等于{n}"
- oneof: "必须是{n}中的一个"

## 中间件系统

全局中间件顺序（`middleware.ApplyMiddlewares`）：
`IPBlacklist → RequestID → Recovery → Logging → CORS → Gzip`

另有安全中间件（敏感路径检测 `CheckSensitivePath`、IP 令牌桶限流 `IPLimiter`）与
路由分组中间件（`JWT` / `Auth` / `Debounce`）。

> 完整说明（含黑名单持久化、限流阈值、执行顺序图）见
> [`backend-conventions.md` §中间件](docs/backend-conventions.md)。

## 数据库连接池配置

使用 GORM 管理连接，连接池参数：

- **最大打开连接数**：10
- **最大空闲连接数**：5
- **连接的最大生命周期**：1 小时

可在 `internal/xdb/` 中调整。

## 统一响应格式

统一 JSON 响应（结构定义见 `internal/utils/response/response.go`）：

```json
{
  "code": 0,
  "message": "success",
  "data": {},
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

- `code`：**业务状态码**，`0` 为成功、`1` 为失败（常量见 `internal/utils/code/bizcode.go`）；
  HTTP 状态码由响应函数另行设置（如 401 / 403）。
- `message`：描述信息
- `data`：返回的数据体（成功时可能为 `null`）
- `request_id`：请求追踪 ID，由 RequestID 中间件生成

调用方统一使用 `response.Success / Error / Unauthorized / Forbidden / ErrorWithCode`，
不自己拼 JSON。细节见 [`backend-conventions.md` §统一响应](docs/backend-conventions.md)。

## 模型层与数据库操作

模型定义在 `internal/models/`，嵌入 `BaseFields` 获得 `ID` / `CreatedAt` / `UpdatedAt` / `DeletedAt`。

数据库操作由 **Service 层直接使用 GORM 原生 API** 完成，通过 `ctxDB(ctx)`（主库/写）与
`ctxSDB(ctx)`（从库/读）获取带上下文的 `*gorm.DB`：

```go
func (s *userService) GetUserInfo(ctx context.Context, id int) (*models.User, error) {
    var u models.User
    if err := ctxSDB(ctx).First(&u, id).Error; err != nil {
        return nil, err
    }
    return &u, nil
}
```

Service 方法首参数统一为 `ctx context.Context`；Handler 调用时直接传 `c *gin.Context`。
细节（选库原则、事务、JSONB 序列化）见 [`backend-conventions.md`](docs/backend-conventions.md)。

## 项目初始化

初始化逻辑集中在 `internal/bootstrap`，`cmd/main.go` 只负责入口与启动。
规范见 [`backend-conventions.md` §项目初始化规范](docs/backend-conventions.md)。

## 测试

```bash
make test                      # 等价于 go test ./... -count=1
go test ./tests/unit/... -v    # 仅单元测试
go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out
```

- 当前只有 `tests/unit` 有用例，`tests/integration` 目录暂为空。
- 注意：直接 `go test ./tests` 会失败（顶层目录没有 Go 包文件），应使用子包路径或 `./...`。
- 涉及数据库的测试请使用独立的测试数据库，并在测试后清理数据。

## 日志

- **ACCESS**：访问日志，记录请求处理信息
- **INFO**：常规信息日志
- **ERROR**：错误日志
- **DEBUG**：调试日志
- **WARN**：警告日志
- **SQL**：SQL 执行日志

### 日志格式

日志文件按日期和级别分割存储在 `logs` 目录下：

```
[ACCESS] RequestID=550e8400-e29b-41d4-a716-446655440000 [POST] /users/login 192.168.1.100 200 45ms
```

### 日志系统特性

1. 多级别日志记录与过滤
2. 日志文件按日期和级别分割
3. 控制台彩色输出
4. 异步日志写入（缓冲通道）
5. 请求追踪：与 RequestID 中间件配合实现全链路追踪

## 性能分析（pprof）

内置基于标准库 `net/http/pprof` 的性能分析，通过独立 debug 端口暴露。
开启方式、接口说明与抓取命令见 [readme/pprof.md](readme/pprof.md)。

## Swagger 文档

修改 Handler 上的 Swagger 注解后，执行：

```bash
make swagger   # 生成到 internal/swagger（不要手敲 swag init，否则会生成回 docs/）
```

启动服务后访问 `http://localhost:<port>/swagger/index.html`。

## 许可证

暂未指定（仓库中暂无 LICENSE 文件）。
