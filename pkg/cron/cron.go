package cron

import (
	"github.com/robfig/cron/v3"
	"seefs-bot/task"
)

var Cron *cron.Cron

func Init() {
	Cron = cron.New(cron.WithSeconds())
	//Cron.AddFunc("*/5 * * * * ?", task.UpdateFreeScore)
	Cron.AddFunc("CRON_TZ=Asia/Shanghai 0 0 0 */1 * ?", task.UpdateFreeScore)
	Cron.Start()
}
