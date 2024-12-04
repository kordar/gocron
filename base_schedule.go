package gocron

type BaseSchedule struct {
	config map[string]string
}

func (b *BaseSchedule) SetConfig(cfg map[string]string) {
	b.config = cfg
}

func (b *BaseSchedule) GetSpec() string {
	if b.config == nil || b.config["spec"] == "" {
		return "@every 10m"
	} else {
		return b.config["spec"]
	}
}

func (b *BaseSchedule) Config() map[string]string {
	if b.config == nil {
		return map[string]string{}
	} else {
		return b.config
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
