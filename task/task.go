package task

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"seefs-bot/bot"
	"seefs-bot/model"
)

func TestTask() {
	var users []model.User
	model.DB.
		Model(&model.User{}).
		Find(&users)
	for _, user := range users {
		bot.Send(&tb.User{
			ID: user.ID,
		}, "xxx", &tb.SendOptions{
			ReplyTo: &tb.Message{
				ID: user.ID,
			},
			ParseMode: tb.ModeMarkdown,
		})
	}
}
