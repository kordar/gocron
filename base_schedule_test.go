package gocron_test

import (
	"encoding/json"
	"github.com/kordar/gocron"
	logger "github.com/kordar/gologger"
	"github.com/robfig/cron/v3"
	"testing"
	"time"
)

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

func (s *TestNameSchedule) Duplicate() int {
	return 2
}

func TestName(t *testing.T) {

	G := gocron.NewGocron(func(job gocron.Schedule) map[string]string {
		return map[string]string{
			"spec": "@every 1s",
		}
	}, func(job gocron.Schedule) bool {
		return true
	})

	schedule := TestNameSchedule{}
	G.AddWithChain(&schedule, func(funcJob cron.Job) cron.Job {
		return cron.NewChain(cron.SkipIfStillRunning(cron.DefaultLogger)).Then(funcJob)
	})
	G.Start()

	time.Sleep(3 * time.Second)
	G.Remove(schedule.GetId())

	time.Sleep(3 * time.Second)
	state := G.State()
	marshal, _ := json.Marshal(&state)
	logger.Infof("------------->%v", string(marshal))
	G.Initializer = func(job gocron.Schedule) map[string]string {
		return map[string]string{
			"spec": "@every 3s",
		}
	}
	G.Reload(schedule.GetId())
	state = G.State()
	marshal, _ = json.Marshal(&state)
	logger.Infof("------------->%v", string(marshal))

	time.Sleep(200 * time.Second)

}
