package gocron

import "github.com/robfig/cron/v3"

type CachedJob struct {
	Job     Schedule
	CronJob cron.Job
}
