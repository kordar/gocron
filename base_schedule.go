package gocron

type BaseSchedule struct {
	cfg map[string]string
}

func (b *BaseSchedule) SetConfig(cfg map[string]string) {
	b.cfg = cfg
}

func (b *BaseSchedule) GetSpec() string {
	return b.cfg["spec"]
}

func (b *BaseSchedule) Config() map[string]string {
	return b.cfg
}

func (b *BaseSchedule) Execute() {
}

func (b *BaseSchedule) Duplicate() int {
	return 1
}

func (b *BaseSchedule) Tag() string {
	return "main"
}

func (b *BaseSchedule) Description() string {
	return ""
}
