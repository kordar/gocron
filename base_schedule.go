package gocron

type BaseSchedule struct {
	cfg map[string]string
}

func (b *BaseSchedule) SetConfig(cfg map[string]string) {
	b.cfg = cfg
}

func (b *BaseSchedule) GetSpec() string {
	if b.cfg == nil || b.cfg["spec"] == "" {
		return "@every 10m"
	} else {
		return b.cfg["spec"]
	}
}

func (b *BaseSchedule) Config() map[string]string {
	if b.cfg == nil {
		return map[string]string{}
	} else {
		return b.cfg
	}
}

func (b *BaseSchedule) Execute() {
}

func (b *BaseSchedule) Duplicate() int {
	return 1
}

func (b *BaseSchedule) Tag() string {
	return ""
}

func (b *BaseSchedule) Description() string {
	return ""
}
