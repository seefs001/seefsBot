package main

import (
	"flag"
	"github.com/seefs001/seefsBot/internal/bot"
	"github.com/seefs001/seefsBot/internal/conf"
	"github.com/seefs001/seefsBot/internal/http"
	"github.com/seefs001/seefsBot/pkg/logger"
	"github.com/seefs001/seefsBot/pkg/orm"
	"gorm.io/driver/mysql"
	"os"
)

func main() {
	flag.Parse()
	if *fHelp {
		flag.Usage()
	} else {
		run()
	}
}

func run() {
	if err := conf.Init(*fConfig); err != nil {
		panic(err)
	}

	logger.New(logger.SetEnv("dev"), logger.SetPath(conf.GetConf().Log.Path))

	if err := orm.Init(mysql.Open(conf.GetConf().MySQL.DSN)); err != nil {
		panic(err)
	}
	//if err := redis.Init(conf.GetConf().Redis.Addr,
	//	conf.GetConf().Redis.Pass,
	//	conf.GetConf().Redis.DB);err!=nil{
	//	panic(err)
	//}

	if err := bot.Init(os.Getenv("BOT_TOKEN")); err != nil {
		panic(err)
	}

	go func() {
		err := http.Start()
		panic(err)
	}()

	bot.Start()
}
