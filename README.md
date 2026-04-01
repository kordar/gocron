# gocron

轻量、直接可用的 Go 定时任务调度库，基于 robfig/cron v3 封装。现已全面切换到标准库 slog 进行日志输出。

## 安装

```bash
go get github.com/kordar/gocron
```

## 使用

```go
import (
    "log/slog"
    "os"
)

// 配置 slog（示例为 JSON 输出）
handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
slog.SetDefault(slog.New(handler))

// 初始化cron处理器，需要设置配置函数，过滤执行函数
g := gocron.NewGocron(nil)

// 添加schedule
g.Add(&TestNameSchedule{})

// 启动服务
g.Start()
```

## 实现

```go
type TestNameSchedule struct {
    gocron.BaseSchedule
}

func (s *TestNameSchedule) GetId() string {
    return "test-name"
}

func (s *TestNameSchedule) Execute() {
    config := s.Config()
    slog.Info("execute", "config", config)
}
```

## 说明

- 已移除第三方 gologger，统一使用 Go 标准库 slog。
- 需要 Go 1.21+。
- 日志为结构化输出，推荐使用键值对，如：`slog.Info("add job", "id", id, "duplicate", i)`.

## 运行时控制

支持初始化配置与运行期过滤：

```go
g := gocron.NewGocron(nil)

// 初始化时为不同任务下发配置（例如动态 spec）
g.SetInitFn(func(job gocron.Schedule) map[string]string {
    return map[string]string{
        "spec": "@every 5s",
    }
})

// 运行前过滤：返回 false 则本次不执行
g.SetRuntimeFn(func(job gocron.Schedule) bool {
    slog.Info("runtime check", "id", job.GetId())
    return true
})
```

## 重载、移除与状态

```go
// 重新加载某个任务（会先移除再添加）
g.Reload("test-name")

// 移除任务
g.Remove("test-name")

// 查看当前状态
state := g.State() // []JobStateVO
slog.Info("state", "data", state)
```

## 并发与自定义 Cron Job

- Duplicate：同一任务可注册多个副本。
- ToCronJob：可自定义 cron.Job（例如加上 SkipIfStillRunning）。

```go
type ChainSchedule struct{ gocron.BaseSchedule }
func (s ChainSchedule) GetId() string { return "chain" }
func (s ChainSchedule) Execute()      { slog.Info("do work") }
func (s ChainSchedule) ToCronJob() cron.Job {
    // 自定义链：仅示例，如无需可不实现
    funcJob := gocron.GenCronJobWithCanRun(&s, nil)
    return cron.NewChain(cron.SkipIfStillRunning(cron.DefaultLogger)).Then(funcJob)
}
```

## 迁移指南（从 gologger 到 slog）

- 将 `logger.Infof(\"...%v\", v)` 替换为 `slog.Info(\"...\", \"v\", v)`。
- 将 `logger.Warnf/Warn` 替换为 `slog.Warn`；`logger.Errorf/Error` 替换为 `slog.Error`。
- 如需格式化字符串，可直接使用 `slog.Info(\"msg\", \"detail\", fmt.Sprintf(\"...\", v))`，更推荐结构化字段。
