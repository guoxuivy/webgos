# webgos 开发流程

> 精简版流程。配套规范：[`backend-conventions.md`](backend-conventions.md)（后端权威源）、
> [`changes/`](changes/README.md)（一页纸变更模板）。
>
> **本仓库没有 CI，也没有 git 钩子**，`make check` 靠自觉执行。这不是可以跳过的理由——
> 跳过它，`docs/` 下的规范就只是一堆文字。

---

## 1. 开工清单

开始任何代码改动前：

1. `git status --short` —— 确认工作树干净，先识别并处理他人/其他任务的改动。
2. 从最新 `dev` 建分支，禁止直接在 `dev` / `master` 上开发。
3. 读影响面对应的规范：
   - 改后端代码 → [`backend-conventions.md`](backend-conventions.md)
   - 不确定 → 先读根 `AGENTS.md` 的文档路由表
4. 跑一次基线：`make check`（确认改动前就是绿的）。
5. 判断是否需要填一页纸 change（见 §3）。

### 分支命名

```
feature/<简短描述>     # 新功能
fix/<简短描述>         # 缺陷修复
docs/<简短描述>        # 仅文档
refactor/<简短描述>    # 重构
```

小写短横线，不放密钥、客户数据或敏感内容。

---

## 2. 实现内环

小步闭环，**失败 → 最小实现 → 验证 → 再改**：

1. 先让问题可复现（Bug）或先写能失败的测试（新功能）。
2. 做最小实现，不顺便重构无关代码。
3. 定向验证：`go test ./internal/xxx/... -count=1`。
4. 改完跑 `make check`。

### 命令速查

```bash
make help      # 查看全部命令
make fmt       # gofmt 格式化
make vet       # go vet 静态检查
make test      # go test ./... -count=1
make build     # 编译
make run       # 本地运行
make swagger   # 重新生成 API 文档（改了 Swagger 注解后必须跑）
make check     # fmt + vet + test + build，一键体检
```

---

## 3. 什么时候填一页纸 change

判断标准不是改动行数，而是**外部能否观察到行为变化**。

| 变更类型 | 是否填 change |
| --- | --- |
| 新接口、新页面、可观察行为变化 | **必须填** |
| 修改接口字段、错误语义、权限行为 | **必须填** |
| 数据库结构变更 | **必须填**（含回滚方案） |
| 纯重构、等价性能优化、测试补强 | 不需要（需说明行为不变） |
| 文档、依赖、构建调整 | 通常不需要 |

模板：[`changes/TEMPLATE.md`](changes/TEMPLATE.md)，复制到 `changes/<yyyy-mm-dd>-<动词>-<能力>.md` 后填写。

**不要用"这只是 bugfix"绕过**：如果现行行为描述本身就错了，那正是要更新规格的场景。

---

## 4. 提交约定

```
feat(模块): 一句话结果
fix(模块): 一句话结果
docs(模块): 一句话结果
refactor(模块): 一句话结果
chore(模块): 一句话结果
```

- 标题说**结果**，正文说**原因**与**影响**。
- 禁止 `fix bug`、`update` 这类无信息标题。
- 只暂存本次任务的文件，提交前 `git diff --cached` 自查是否混入密钥、签名 URL 或无关改动。

---

## 5. 收工检查清单

- [ ] 需求/问题明确，范围与非目标清楚；
- [ ] 可观察行为变化已填写一页纸 change（见 §3）；
- [ ] 代码符合 [`backend-conventions.md`](backend-conventions.md)：分层、ctx 首参、选库、`items`+`total` 分页、统一响应；
- [ ] 改了 Swagger 注解已跑 `make swagger`；
- [ ] 数据库结构变更有迁移与回滚方案（不用 AutoMigrate 补生产表）；
- [ ] `make check` 全绿；
- [ ] 暂存区只含本次主题，无敏感信息；
- [ ] 提交信息可追溯，分支已推送。

任何一项未满足，明确记录为阻断项或后续任务，不要用"基本完成"掩盖。
