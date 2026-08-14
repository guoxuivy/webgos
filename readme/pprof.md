# 性能分析（pprof）使用说明

项目内置了基于标准库 `net/http/pprof` 的性能分析能力，通过**独立 debug 端口**暴露，不影响业务路由，也不受 JWT/RBAC 中间件拦截。

## 开启配置

在 `config/config.yaml` 的 `server` 节点下：

```yaml
server:
  pprof: true        # 是否启用 pprof（建议仅本地/排查问题时开启，生产环境设为 false）
                     # pprof 监听端口固定为 6060，无需配置
```

启用后服务会额外监听 `:6060`，提供以下接口：

| 接口 | 说明 |
|---|---|
| `/debug/pprof/` | 分析首页（含各 profile 链接） |
| `/debug/pprof/heap` | 堆内存快照 |
| `/debug/pprof/profile` | CPU profile（默认采样 30s，`?seconds=N` 可调） |
| `/debug/pprof/goroutine` | goroutine 数量与栈 |
| `/debug/pprof/block` | 阻塞分析 |
| `/debug/pprof/mutex` | 锁竞争分析 |
| `/debug/pprof/threadcreate` | 线程创建分析 |
| `/debug/pprof/trace` | 执行追踪（`?seconds=N`） |

> 排查完成后请将 `pprof` 设为 `false` 并重启，无需改代码即可关闭。注意：**生产环境务必关闭 `pprof` 与 `swag`**，二者仅在启动时把 Swagger UI 静态资源内嵌进内存就会占用约 8MB 常驻堆。

## 抓取与分析

```bash
# 1) 启动服务（已开启 pprof）
go run cmd/main.go -c ./config/config.yaml

# 2) 在有业务流量时抓取（先压一波请求，再抓才有意义）
# CPU profile（采样 30 秒）
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# 堆内存快照（建议跑过业务后再抓，避免 profile 自身压缩噪声）
go tool pprof http://localhost:6060/debug/pprof/heap

# goroutine / 阻塞 / 锁
go tool pprof http://localhost:6060/debug/pprof/goroutine
go tool pprof http://localhost:6060/debug/pprof/block
go tool pprof http://localhost:6060/debug/pprof/mutex

# 执行追踪（5 秒）
curl "http://localhost:6060/debug/pprof/trace?seconds=5" -o trace.out
go tool trace trace.out
```

进入 `go tool pprof` 交互后可用 `top`、`web`、`png`、`list 函数名` 等命令定位热点。

### 查看可视化图（火焰图 / 调用图）

- `web`：生成 SVG 调用图并在浏览器打开（需要安装 graphviz）。
- `png`：导出 PNG 调用图（需要安装 graphviz）。
- 火焰图可借助 `go tool pprof -http=:8081` 直接起本地 Web UI，在浏览器查看火焰图、拓扑图、源码热行。

```bash
# 起本地 Web UI（推荐，无需 graphviz）
curl http://localhost:6060/debug/pprof/profile?seconds=30 -o cpu.prof
go tool pprof cpu.prof && top
# 或者开启ui展示
# http://localhost:8081/ui/flamegraph # 火焰图
# http://localhost:8081/ui/top
# http://localhost:8081/ui/source
# http://localhost:8081/ui/peek
go tool pprof -http=:8081 cpu.prof 

# 或者直接查看
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
top10   # 看耗时前10函数
# list 函数名   # 看某函数逐行耗时
```

## 实测记录

对启动后空闲状态的堆做了一次快照，结论：**无内存泄漏**，堆约 15MB、进程 RSS 约 40MB，属正常水平。大头来自两部分：

- Swagger UI 内嵌静态资源（`swaggo/files` init）常驻约 8MB —— 生产关 `swag` 即可省下；
- 抓取 heap 时 pprof 自身 gzip 压缩产生的约 7.5MB 临时缓冲 —— 测量噪声，非业务常驻。

另对空载（无效路由压测）做了一次 CPU profile：业务代码未进入热点，CPU 几乎全耗在运行时系统调用与网络 I/O 上（Windows 下 `runtime.cgocall` 占比最高），说明网关层开销极小、无异常自旋。要看真实业务热点，需用有效 token 压测真实业务接口（如 `/api/menu/list`）。

## 排查建议

1. **内存泄漏**：对比两次 `heap` 快照（`go tool pprof -base=old.pb.gz new.pb.gz`），看哪些对象只增不减。
2. **CPU 热点**：务必在有真实业务流量时抓取 `profile`，再用 `top`/`web` 定位耗时函数。
3. **goroutine 暴涨**：抓 `goroutine` 看是否大量阻塞在 channel / 锁 / 网络读，排查泄漏的 goroutine。
4. **锁竞争**：抓 `mutex`；阻塞：`block`。

## 压测命令

```bash
netstat -ano | grep -E ":8080|:6060" | head -3; echo "---login---"; cd /d/Goroot/webgos && curl -s -X POST http://localhost:8080/auth/login -H "Content-Type: application/json" -d "{\"username\":\"admin\",\"password\":\"123456\"}"

cd /d/Goroot/webgos && TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3ODY3MDA2MzAsInVzZXJfaWQiOjEsInVzZXJuYW1lIjoiYWRtaW4ifQ.IkOtywkSDz7cQGjwIHcuQeeKmwppuDp0M20SdduGIyI" && hey -n 200000 -c 300 -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/menu/list 2>&1 | grep -E "Requests/sec|Slowest|Fastest|Average|99%|Status code|200"
```