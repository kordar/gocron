# gocron

## 使用

```go
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
    logger.Infof("----------------------------%v", config)
}
```