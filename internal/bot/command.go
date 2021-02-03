package bot

import (
	"github.com/seefs001/seefsBot/pkg/logger"
	"go.uber.org/zap"
	tb "gopkg.in/tucnak/telebot.v2"
)

func setCommend() {
	// 设置bot命令提示信息
	commands := []tb.Command{
		{"start", "获取当前用户ID和key"},
	}

	if err := B.SetCommands(commands); err != nil {
		logger.Logger().Error("set commend err", zap.Error(err))
	}
}
