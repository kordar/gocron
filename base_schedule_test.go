package gocron_test

import (
	"github.com/kordar/gocron"
	logger "github.com/kordar/gologger"
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

func TestName(t *testing.T) {

	g := gocron.NewGocron(func(job gocron.Schedule) map[string]string {
		return map[string]string{
			"spec": "@every 3s",
		}
	}, func(job gocron.Schedule) bool {
		return true
	})

	g.Add(&TestNameSchedule{})
	g.Start()
	time.Sleep(200 * time.Second)

}
