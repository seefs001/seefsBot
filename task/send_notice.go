package task

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"seefs-bot/bot"
	"seefs-bot/model"
)

func SendNotice(msg string) {
	bot.Init()
	var users []model.User
	model.DB.
		Model(&model.User{}).
		Find(&users)
	for _, user := range users {
		bot.B.Send(tb.ChatID(user.ID), msg, &tb.SendOptions{
			ParseMode: tb.ModeMarkdown,
		})
	}
}
