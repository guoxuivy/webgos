# AGENTS.md · webgos 仓库工作入口

> **webgos** 是一个基于 Go 的企业级 Web 系统快速开发脚手架，基于 Gin + GORM，
> 面向可扩展的业务开发，内置 RBAC 权限、日志与请求追踪、事务、分页、统一响应等能力。

本文件是 AI / 新同学进入本仓库的**唯一入口**。它只负责导航与红线，不复制规范细节；
细节一律到权威源查阅。

---

## 开工前必读（强制）

开始任何产品、代码、测试、配置或文档改动前：

1. 读本文件。
2. **改后端代码前，必读 [`docs/backend-conventions.md`](docs/backend-conventions.md)**（后端规范唯一权威源）。
3. 按 [文档路由表](#文档路由表) 找到本次影响面对应的材料。
4. 判断是否需要填一页纸 change（见 [§ 什么时候填 change](#什么时候填-change)）。
5. `git status --short` 确认工作树干净，跑一次基线 `make check`。

---

## 文档路由表

| 任务 | 必读权威源 |
| --- | --- |
| 写/改后端代码（模型、DTO、Service、Handler、路由） | [`docs/backend-conventions.md`](docs/backend-conventions.md) |
| 开发流程、分支、提交、收工 | [`docs/dev-workflow.md`](docs/dev-workflow.md) |
| 可观察行为变化的变更记录 | [`docs/changes/README.md`](docs/changes/README.md) |
| 新增 CRUD 模块 | 清单见 [`backend-conventions.md`](docs/backend-conventions.md)；代码骨架见 [`docs/templates/crud-api.md`](docs/templates/crud-api.md) |
| 读写分离 / 主从延迟 | [`backend-conventions.md` §读写分离与选库](docs/backend-conventions.md) |
| 分页返回格式 | [`backend-conventions.md` §分页接口返回格式](docs/backend-conventions.md) |
| Swagger / API 文档生成 | 改注解后 `make swagger`，产物在 `internal/swagger/` |
| 业务功能说明（商品/库存/权限/菜单） | `readme/` 下对应文档 |
| 项目整体结构与模块 | [`README.md`](README.md) |

**单一权威源原则**：同一件事只在一个地方定义。若两处说法冲突，以 `docs/backend-conventions.md`
为准，并修正另一处。[`docs/templates/crud-api.md`](docs/templates/crud-api.md) 只是代码生成模板，**不定义规范**。

---

## 技术栈（不可协商）

| 层 | 选型 |
| --- | --- |
| 语言 | Go 1.25 |
| Web 框架 | Gin v1.11 |
| ORM | GORM v1.31（原生 API，不用泛型 BaseModel 封装） |
| 数据库 | MySQL / PostgreSQL（双驱动，实际由配置决定） |
| 认证 | golang-jwt/jwt/v4 |
| 校验 | go-playground/validator/v10 |
| API 文档 | Swaggo |
| 缓存 | patrickmn/go-cache |
| 测试 | testify |

引入新依赖前先确认现有工具/库无法满足，并在提交信息说明理由。

---

## 目录所有权

| 路径 | 归属 | 说明 |
| --- | --- | --- |
| `internal/**` | 业务代码 | 分层：dto / handlers / services / models / routes / middleware |
| `internal/swagger/**` | **生成物** | 由 `make swagger` 生成，**勿手改** |
| `docs/**` | **手写规范区** | 100% 手写，不含任何生成物；含 `templates/`（CRUD 代码模板） |
| `readme/**` | 业务文档 | 商品/库存/权限/菜单说明 |
| `config/**` | 配置 | `config.yaml` |
| `tests/unit/**` | 单元测试 | |

> 注意：`docs/` 与 `internal/swagger/` 已物理隔离。手写规范**不要**放进 `internal/swagger/`，
> 生成物**不要**放进 `docs/`。

---

## 实现红线

违反以下任一条即为缺陷，不因"能跑通"而放行：

1. **分层固定**：`Routes → Handlers → Services → Models`（DTO 承载入参/出参）。
   Handler 不写业务逻辑、不直接操作 DB；Service 不引用 `*gin.Context`。
2. **Service 方法首参必须是 `ctx context.Context`**；Handler 调用时直接传 `c`。
3. **分页接口必须返回 `items` + `total`**，禁止 `list` / `data` / `rows` / `count` / `sum`。
4. **统一响应**：只用 `response.Success/Error/Unauthorized/Forbidden/ErrorWithCode`，不自己拼 JSON。
5. **选库**：写操作与强一致读用 `ctxDB(ctx)`，普通只读用 `ctxSDB(ctx)`；
   **事务闭包内只用 `tx`**，禁止调用 `ctxDB` / `ctxSDB`。
6. **JSONB 字段手动** `Serialize()` / `Deserialize()`，不依赖 GORM 钩子。
7. **Swagger 注解写在 Handler 上方**，改完必须 `make swagger`。
8. **新增路由即新增 RBAC 权限点**（标识为 `路径:HTTP方法`），需授权后才可访问。
9. **增删改后清缓存**：`cache.GetCache().DeleteByPrefix(...)`。
10. **软删除默认开启**，物理删除显式用 `Unscoped()`。
11. **错误统一用 `response` 返回**，不要直接 `return`。

---

## 什么时候填 change

判断标准是「外部能否观察到行为变化」，不是改动行数：

- **必须填**：新接口、新页面、修改接口字段/错误语义/权限行为、数据库结构变更。
- **不需要**：纯重构、等价性能优化、测试补强、文档/依赖调整。

模板见 [`docs/changes/TEMPLATE.md`](docs/changes/TEMPLATE.md)。
不要用"这只是 bugfix"绕过——若现行行为描述本身错了，那正是要更新规格的场景。

---

## 命令速查

```bash
make help      # 查看全部命令
make check     # 一键体检：fmt + vet + test + build（改完就跑）
make test      # go test ./... -count=1
make build     # 编译
make run       # 本地运行（默认 config/config.yaml）
make swagger   # 重新生成 API 文档到 internal/swagger
make fmt       # gofmt
make vet       # go vet
```

> **本仓库没有 CI，也没有 git 钩子**，`make check` 靠自觉执行。
> 跳过它，规范就只是一堆文字。

---

## 收工检查清单

- [ ] 可观察行为变化已填一页纸 change；
- [ ] 符合 [`backend-conventions.md`](docs/backend-conventions.md) 全部红线（尤其 §3、§10 分页格式）；
- [ ] 改了 Swagger 注解已跑 `make swagger`；
- [ ] DB 结构变更有迁移与回滚方案（不用 AutoMigrate 补生产表）；
- [ ] `make check` 全绿；
- [ ] 暂存区只含本次主题，无密钥/敏感信息；
- [ ] 提交信息格式为 `feat|fix|docs|refactor|chore(模块): 一句话结果`。
