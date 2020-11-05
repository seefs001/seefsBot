package bot

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"seefs-bot/model"
	"seefs-bot/pkg/logger"
)

func WelcomeMsg(m *tb.Message) {
	user, _ := model.FindOrCreateUserByTelegramIDAndUserName(int(m.Chat.ID), m.Chat.Username)
	logger.Info(fmt.Sprintf("/start chat_id: %d", user.TelegramID))
	_, _ = B.Send(m.Chat, fmt.Sprintf("欢迎使用pixiv日榜bot!"))
}
