package tasks

import (
	"time"

	"github.com/phuslu/log"
	"github.com/robfig/cron/v3"
	"github.com/seefs001/seefsBot/internal/bot"
	"github.com/seefs001/seefsBot/internal/conf"
	"github.com/seefs001/seefsBot/pkg/util"
)

var Task = &cron.Cron{}

func Start() error {
	Task = cron.New(cron.WithSeconds())
	spec := conf.GetConf().Task.Cron
	_, err := Task.AddFunc(spec, func() {
		results, err := util.UpdatePrice()
		if err != nil {
			log.Warn().Time("time", time.Now()).
				Msg("任务执行失败 ")
		}
		for _, result := range results {
			if result.Increase >= conf.GetConf().Task.UpperLimit {
				bot.Broadcast(result.Type + "上涨幅度过大")
			}
			if result.Increase <= conf.GetConf().Task.LowerLimit {
				bot.Broadcast(result.Type + "下跌幅度过大")
			}
		}
	})
	if err != nil {
		return err
	}
	Task.Start()
	return nil
}

func Stop() {
	Task.Stop()
}
