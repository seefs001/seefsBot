package bot

import (
	"github.com/phuslu/log"
	tb "gopkg.in/tucnak/telebot.v2"
)

func setCommend() {
	// 设置bot命令提示信息
	commands := []tb.Command{
		{"start", "获取当前用户ID和key"},
		{"price", "查看当前币价"},
	}

	if err := B.SetCommands(commands); err != nil {
		log.Fatal().Err(err)
	}
}
