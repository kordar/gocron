package gocron

import (
	logger "github.com/kordar/gologger"
	"github.com/robfig/cron/v3"
)

func GenCronJobWithCanRun(job Schedule, canRun func() bool) cron.Job {
	return cron.FuncJob(func() {

		defer func() {
			if err := recover(); err != nil {
				logger.Errorf("[gocron] recover %v", err)
			}
		}()

		if !canRun() {
			return
		}

		job.Execute()
	})
}
