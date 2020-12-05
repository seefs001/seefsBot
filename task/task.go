package task

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	tb "gopkg.in/tucnak/telebot.v2"
	"seefs-bot/bot"
	"seefs-bot/model"
	"seefs-bot/pkg/logger"
)

func UpdateFreeScore() {
	if err := model.DB.Model(&model.User{}).
		Update("free_score", viper.GetInt("score.free_score")).Error; err != nil {
		logger.Error("分配免费额度失败", zap.Error(err))
	}
	var users []model.User
	model.DB.
		Model(&model.User{}).
		Find(&users)
	for _, user := range users {
		logger.Info(fmt.Sprintf("cron 重置次数"))
		bot.B.Send(tb.ChatID(user.ID), "免费查询已重置，谢谢您的支持！", &tb.SendOptions{
			ParseMode: tb.ModeMarkdown,
		})
	}
}
