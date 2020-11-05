package main

import (
	"flag"
	"seefs-bot/bot"
	"seefs-bot/pkg/conf"
)

var (
	fConfig = flag.String("config", "./config.toml", "配置文件路径")
	fHelp   = flag.Bool("h", false, "show help")
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
	// initialize config
	if err := conf.InitConfig(*fConfig); err != nil {
		panic(err)
	}
	bot.InitBot()
	bot.Start()
}
