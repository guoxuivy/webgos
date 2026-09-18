# 定时任务子系统接线并从对外可见面收回

- 日期：2026-09-18
- 负责人：未指定
- 关联分支 / PR：未指定

## 背景与问题

证据：全仓库检索无任何文件 import `webgos/cron`（只有 `cron/demo.go`、`cron/initialize.go` 自身），
Go 只对被 import 的包执行 `init()`，因此 `demo.go` 中的 `Register` 从未被调用；
同时 `bootstrap.Initialize` 未调用 `cron.SetUp()`、`bootstrap.Close` 未调用 `cron.ShutDown()`。

结果：`demo.go` 注释所述「启动后立即聚合一次，避免服务重启后当天汇总长时间为空」从未发生。
调度器实现完整但从未接线，属于「装好了、没通电」的死子系统。

附带问题：`cron/` 位于 `internal/` 之外，使本仓库私有包对外可见（外部模块 `import "webgos/cron"` 是合法的）。

## 范围与非目标

**范围内**：

- `cron/` 移入 `internal/cron/`，不再对外可见（package 名与实现均不变）
- `internal/bootstrap` 接线：`Initialize` 末尾调用 `cron.SetUp()`；`Close` 中于关闭 logger 之前调用 `cron.ShutDown()`
- README 目录树同步

**非目标**：

- 不改动调度语义（单 Job 单 goroutine、上次执行结束才重新计时、panic 后不丢失后续调度）
- 不新增或修改业务定时任务，`demo.go` 保持原样注册，按 5 分钟间隔运行
- 不动 `common/*`；`common/syncx`、`common/time`、`internal/utils/file` 三个零引用包本次**保留**，不删除

## 验收场景

**正常路径**：Given 服务正常启动且 `demo.go` 注册了 `RunOnStart=true` 的 `cron_demo`，
When 进程执行完 `bootstrap.Initialize`，Then 日志出现 `cron setup done, jobs=1`，
并紧随出现一次 `cron job=cron_demo finished`；
此后相邻两轮 `finished` 的间隔 ≥ 5 分钟（验证「执行结束后才重新计时」，而非 Ticker 补跑）。

**拒绝路径**：Given 向进程发送 SIGINT/SIGTERM，When 触发优雅关闭，
Then 日志依次出现：`Shutting down server...` → `Server exiting` → `Closing resources...` → `cron stopped`，
且进程在优雅期内退出，日志中不出现 job 被中断的输出（`ShutDown` 是等待在途任务跑完，不是打断）。
顺序说明：`bootstrap.Close()` 是 `main` 的 defer，需等 `<-idleConnsClosed`（即 `Server exiting` 打印之后）才执行，
因此 `cron stopped` 是进程的最后一行日志，而非出现在 `Server exiting` 之前。

**恢复路径**：Given 某个 Job 的 `Fn` 发生 panic，When 该轮调度执行，
Then `runJob` 的 recover 捕获并记录 `cron job=xxx panic: ...`，且该 Job 后续轮次仍被调度。

## 影响面

- [ ] API 契约（字段 / 错误码 / Swagger）
- [ ] 数据库结构与迁移
- [ ] 权限与鉴权（新增路由即新增 RBAC 权限点）
- [ ] 缓存（增删改后是否需 `DeleteByPrefix`）
- [ ] 前端 / 调用方适配
- [x] 其他：① 进程启动后新增常驻 goroutine，并按 5 分钟周期输出 `cron_demo` 的 Info 日志；
      ② `internal` 之外不再存在 `webgos/cron` 包，原 import 路径 `webgos/cron` 失效（改为 `webgos/internal/cron`）。

## 已知现象

`make run` 用 Ctrl+C 结束后，make 会打印 `make: *** [Makefile:35: run] Error 1`，
但服务本身已完成优雅退出（日志走到最后一行 `cron stopped`，属正常结束而非崩溃）。

原因：`go run` 是包装进程，它收到中断信号时由 cmd/go 内部调用 `base.SetExitStatus(1)`
（`src/cmd/go/internal/work/exec.go` 中 `<-base.Interrupted` 分支），因此以 1 退出；
`go help run` 亦明确说明「The exit status of Run is not the exit status of the compiled binary.」，
即 `go run` 的退出码不等于被运行程序的退出码。本进程自身是以 0 退出的。

规避：`make build` 后直接运行二进制（`./webgos.exe -c ./config/config.yaml`），
中间少了 go 包装层，Ctrl+C 的退出码即程序自身的 0。

## 回滚方案

1. 彻底回滚：回退本次 commit（`internal/bootstrap` 的接线、`internal/cron` 目录、README、本 change 文档），
   行为回到「定时任务从不运行」。注意 `cron/` 原为未跟踪目录（不在 HEAD 中），回退前先确认其跟踪状态。
2. 仅停任务、保留接线：删除 `demo.go` 或移除其中的 `Register` 调用，`SetUp` 将以 `jobs=0` 空转，
   日志变为 `cron setup done, jobs=0`，无需回退其余改动。
   此路径适用于「周期性 demo 日志干扰排障」的场景。

## 验证结果

- 定向测试：2026-09-18 11:29–11:31 实机执行（`make run` 启动后 Ctrl+C），实机日志如下：
  - 正常路径 **passed**：`cron setup done, jobs=1` 后随即出现 `cron job=cron_demo finished cost=517.3µs`（RunOnStart 生效）
  - 拒绝路径 **passed**：`Shutting down server...` → `Server exiting` → `Closing resources...` → `cron stopped`，进程正常退出
  - 未覆盖：① 相邻两轮间隔 ≥ 5 分钟未验证（本次进程仅存活约 2 分钟即退出，未跨第二个周期）；
    ② panic 恢复路径未验证
- `make check`：通过（2026-09-18 执行，fmt + vet + test + build 全绿；`go test` 仅 `webgos/tests/unit` 有用例，`ok 0.046s`）
- 备注：① `go fmt` 顺带重排了 `internal/cron/initialize.go` 缩进（Go 1.25 gofmt 对 label 的处理），属格式化修正、无语义变化；
  ② `make run` 用 Ctrl+C 结束时 make 会打印 `Error 1`，这是 `go run` 包装进程自身的退出码，与定时任务无关，
  详见本文「已知现象」一节。
