package gocron

import (
	"github.com/robfig/cron/v3"
	"time"
)

type State int

const (
	Ready State = iota
	Running
	Shutdown
	AddFailed
)

func (s State) String() string {
	return [...]string{"Ready", "Running", "Shutdown", "AddFailed"}[s]
}

type Schedule interface {
	GetId() string
	GetSpec() string
	Execute()
	Config() map[string]string
	SetConfig(cfg map[string]string)
	Duplicate() int
	Tag() string
	Description() string
	ToCronJob() cron.Job
}

type JobStateItem struct {
	RegTime time.Time              `json:"reg_time"`
	EntryId cron.EntryID           `json:"entry_id"`
	Params  map[string]interface{} `json:"params"`
}

type JobState struct {
	Id    string
	Job   Schedule
	State State
	Err   error
	Items []JobStateItem
}

type JobStateVO struct {
	JobId string         `json:"job_id"`
	State string         `json:"state"`
	Err   string         `json:"err"`
	Data  []JobStateItem `json:"data"`
}
