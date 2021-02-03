package bot

import (
	"fmt"
	"github.com/seefs001/seefsBot/internal/model"
	"github.com/seefs001/seefsBot/pkg/logger"
	"github.com/seefs001/seefsBot/pkg/orm"
	"github.com/seefs001/seefsBot/pkg/util"
	"go.uber.org/zap"
	tb "gopkg.in/tucnak/telebot.v2"
)

func start(m *tb.Message) {
	if !m.Private() {
		return
	}
	user := model.User{
		ID:        int64(m.Sender.ID),
		Role:      model.NormalRole,
		SecretKey: util.GenRandomString(10, false),
	}
	orm.DB().Where(&model.User{ID: int64(m.Sender.ID)}).
		FirstOrCreate(&user)
	msg := fmt.Sprintf("欢迎使用seefsBot\n"+
		"您的*userID*为*%d*\n"+
		"您的*secretKey*为*%s*",
		user.ID, user.SecretKey)
	_, err := B.Send(tb.ChatID(user.ID), msg, &tb.SendOptions{
		ParseMode: tb.ModeMarkdownV2,
	})
	if err != nil {
		logger.Logger().Error("bot send err", zap.Error(err))
	}
}

func neteasencm(m *tb.Message) {
	_, _ = B.Send(m.Sender, "在写了", &tb.SendOptions{
		ParseMode: tb.ModeMarkdownV2,
	})
}
