package gocron

import (
	"github.com/robfig/cron/v3"
	"time"
)

type Schedule interface {
	GetId() string
	GetSpec() string
	Execute()
	Config() map[string]string
	SetConfig(cfg map[string]string)
	Duplicate() int
	Tag() string
	Description() string
}

type CachedJob struct {
	Job     Schedule
	CronJob cron.Job
}

type EntryItem struct {
	Id      string                 `json:"id"`
	RegTime time.Time              `json:"reg_time"`
	EntryId cron.EntryID           `json:"entry_id"`
	Params  map[string]interface{} `json:"params"`
}

type StateEntryItem struct {
	Id    string       `json:"id"`
	State string       `json:"state"`
	Data  []*EntryItem `json:"data"`
}
