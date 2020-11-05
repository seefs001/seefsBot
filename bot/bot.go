package bot

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	tb "gopkg.in/tucnak/telebot.v2"
	"seefs-bot/bot/fsm"
	"seefs-bot/pkg/logger"
	"time"
)

var (
	// UserState 用户状态，用于标示当前用户操作所在状态
	UserState map[int64]fsm.UserStatus = make(map[int64]fsm.UserStatus)

	// B telebot
	B *tb.Bot
)

func InitBot() {
	poller := &tb.LongPoller{Timeout: 10 * time.Second}
	logger.Info("bot init successfully")
	// create bot
	var err error

	B, err = tb.NewBot(tb.Settings{
		URL:    viper.GetString("bot.endpoint"),
		Token:  viper.GetString("bot.token"),
		Poller: poller,
	})
	if err != nil {
		logger.Fatal("init err", zap.Error(err))
		return
	}
	logger.Info("bot start successfully")
}

//Start bot
func Start() {
	setCommend()
	setHandle()
	B.Start()
}

func Send(to tb.Recipient, what string, options ...interface{}) (*tb.Message, error) {
	send, err := B.Send(to, what, options)
	return send, err
}
