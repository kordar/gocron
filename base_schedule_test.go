package gocron_test

import (
	"encoding/json"
	"github.com/kordar/gocron"
	logger "github.com/kordar/gologger"
	"testing"
	"time"
)

type TestNameSchedule struct {
	gocron.BaseSchedule
}

func (s TestNameSchedule) GetId() string {
	return "test-name"
}

func (s TestNameSchedule) Execute() {
	config := s.Config()
	logger.Infof("--------------AAA--------------%v", config)
}

func (s TestNameSchedule) Duplicate() int {
	return 2
}

//func (s TestNameSchedule) GetSpec() string {
//	return "@every 3s"
//}

//func (s TestNameSchedule) ToCronJob() cron.Job {
//	funcJob := gocron.GenCronJobWithCanRun(&s, func(job gocron.Schedule) bool {
//		return true
//	})
//	return cron.NewChain(cron.SkipIfStillRunning(cron.DefaultLogger)).Then(funcJob)
//}

func TestName(t *testing.T) {

	G := gocron.NewGocron(nil)

	schedule := TestNameSchedule{}
	schedule.SetConfig(map[string]string{
		"spec": "@every 5s",
	})
	G.Add(&schedule)
	//G.AddWithChain(&schedule, func(funcJob cron.Job) cron.Job {
	//	return cron.NewChain(cron.SkipIfStillRunning(cron.DefaultLogger)).Then(funcJob)
	//})
	G.Start()
	G.SetInitFn(func(job gocron.Schedule) map[string]string {
		return map[string]string{
			"spec": "@every 1s",
		}
	})
	G.SetRuntimeFn(func(job gocron.Schedule) bool {
		logger.Warn("Unable to execute current task", job.GetId())
		return false
	})

	//time.Sleep(3 * time.Second)
	//G.Remove(schedule.GetId())

	time.Sleep(3 * time.Second)
	state := G.State()
	marshal, _ := json.Marshal(&state)
	logger.Infof("------------->%v", string(marshal))
	//G.Initializer = func(job gocron.Schedule) map[string]string {
	//	return map[string]string{
	//		"spec": "@every 3s",
	//	}
	//}
	G.Reload(schedule.GetId())
	state = G.State()
	marshal, _ = json.Marshal(&state)
	logger.Infof("------------->%v", string(marshal))

	time.Sleep(200 * time.Second)
}
