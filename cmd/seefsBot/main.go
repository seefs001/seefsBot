package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/phuslu/log"
	"github.com/seefs001/seefsBot/internal/bot"
	"github.com/seefs001/seefsBot/internal/conf"
	"github.com/seefs001/seefsBot/internal/model"
	"github.com/seefs001/seefsBot/internal/server"
	"github.com/seefs001/seefsBot/pkg/logger"
	"github.com/seefs001/seefsBot/pkg/orm"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
)

func main() {
	run()
}

func run() {
	if err := conf.Init("./config.toml"); err != nil {
		panic(err)
	}

	options := logger.Options{
		LogPath: conf.GetConf().Log.LogPath,
		LogName: conf.GetConf().Log.LogName,
	}
	logger.Init(&options)

	if err := orm.Init(mysql.Open(conf.GetConf().MySQL.DSN)); err != nil {
		log.Info().Msg("MySQL连接失败，自动切换到Sqlite")
		err := orm.Init(sqlite.Open("seefsBot.db"))
		if err != nil {
			panic(err)
		}
	}
	// if err := redis.Init(conf.GetConf().Redis.Addr,
	//	conf.GetConf().Redis.Pass,
	//	conf.GetConf().Redis.DB);err!=nil{
	//	panic(err)
	//}
	orm.DB().AutoMigrate(&model.User{}, &model.Coin{})

	//if err := bot.Init(conf.GetConf().Server.BotToken); err != nil {
	//	panic(err)
	//}
	go func() {
		err := server.Start()
		panic(err)
	}()

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals,
			os.Interrupt,
			os.Kill,
			syscall.SIGQUIT,
			syscall.SIGTERM,
			syscall.SIGINT)
		for {
			sig := <-signals
			switch sig {
			default:
				server.Stop()
				bot.Stop()
				return
			}
		}
	}()

	bot.Start()
}
