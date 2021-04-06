package bot

import (
	"github.com/phuslu/log"
	tb "gopkg.in/tucnak/telebot.v2"
)

func setCommend() {
	// 设置bot命令提示信息
	commands := []tb.Command{
		{"start", "获取当前用户ID和key"},
		{"list_type", "查看当前正在监听的币"},
		{"add_type", "增加监听币类型"},
		{"price", "查看当前币价"},
	}

	if err := B.SetCommands(commands); err != nil {
		log.Fatal().Err(err)
	}
}
