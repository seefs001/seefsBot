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
		{"sub", "订阅rss源"},
		{"list", "当前订阅的rss源"},
		{"unsub", "退订rss源"},
		{"unsuball", "退订所有rss源"},

		{"set", "设置rss订阅"},
		{"setfeedtag", "设置rss订阅标签"},
		{"setinterval", "设置rss订阅抓取间隔"},

		{"export", "导出订阅为opml文件"},
		{"import", "从opml文件导入订阅"},

		{"check", "检查我的rss订阅状态"},
		{"pauseall", "停止抓取订阅更新"},
		{"activeall", "开启抓取订阅更新"},

		{"help", "使用帮助"},
		{"version", "bot版本"},
	}

	if err := B.SetCommands(commands); err != nil {
		logger.Error("set commend err", zap.Error(err))
	}
}
