package cron

import (
	"github.com/robfig/cron/v3"
	"seefs-bot/task"
)

var Cron *cron.Cron

func InitTasks() {
	Cron = cron.New(cron.WithSeconds())
	Cron.AddFunc("@every 1s", task.TestTask)
	Cron.Start()
}
