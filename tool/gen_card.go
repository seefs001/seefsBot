package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/seefs001/seefslib-go/xfile"
	"github.com/seefs001/seefslib-go/xrandom"
	"os"
	"seefs-bot/model"
	"seefs-bot/pkg/conf"
	"seefs-bot/pkg/logger"
)

var (
	fConfig = flag.String("config", "./config.toml", "配置文件路径")
	fCount  = flag.Int("count", 10, "制卡数量")
	fPath   = flag.String("path", "card.txt", "生成位置")
	fScore  = flag.Int64("score", 10, "每张卡的积分")
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
	if err := conf.Init(*fConfig); err != nil {
		panic(err)
	}
	count := *fCount
	score := *fScore
	path := *fPath
	file, err := xfile.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	write := bufio.NewWriter(file)
	for i := 0; i < count; i++ {
		cardCode := xrandom.RandStringRunes(10)
		model.DB.Model(&model.Card{}).
			Create(&model.Card{
				Content: cardCode,
				Score:   score,
			})
		write.WriteString(fmt.Sprintf("%s\n", cardCode))
	}
	write.Flush()
	logger.Info(fmt.Sprintf("本次共生成了%d张卡密，每张积分为%d分", count, score))
}
