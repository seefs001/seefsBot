package bot

import (
	"go.uber.org/zap"
	tb "gopkg.in/tucnak/telebot.v2"
	"seefs-bot/pkg/logger"
)

func setCommend() {
	// 设置bot命令提示信息
	commands := []tb.Command{
		{"start", "开始使用"},
		{"get_score", "获取当前积分"},
		{"invite", "发送邀请码奖励5积分"},
		{"card", "卡密充值"},
		{"suggest", "给作者提建议"},
	}

	if err := B.SetCommands(commands); err != nil {
		logger.Error("set commend err", zap.Error(err))
	}
}
