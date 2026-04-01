package gocron

import (
	"log/slog"

	"github.com/robfig/cron/v3"
)

func GenCronJobWithCanRun(job Schedule, canRun RuntimeFunction) cron.Job {
	return cron.FuncJob(func() {

		defer func() {
			if err := recover(); err != nil {
				slog.Error("[gocron] recover", "err", err)
			}
		}()

		if canRun != nil && !canRun(job) {
			return
		}

		job.Execute()
	})
}
