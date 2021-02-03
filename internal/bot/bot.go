package bot

import (
	"github.com/seefs001/seefsBot/internal/bot/fsm"
	"github.com/seefs001/seefsBot/pkg/logger"
	"go.uber.org/zap"
	tb "gopkg.in/tucnak/telebot.v2"
	"time"
)

var (
	// UserState 用户状态，用于标示当前用户操作所在状态
	UserState map[int64]fsm.UserStatus = make(map[int64]fsm.UserStatus)

	B *tb.Bot
)

func Init(token string) error {
	poller := &tb.LongPoller{Timeout: 10 * time.Second}
	logger.Logger().Info("bot init successfully")

	var err error

	B, err = tb.NewBot(tb.Settings{
		URL:    "https://api.telegram.org",
		Token:  token,
		Poller: poller,
	})
	if err != nil {
		logger.Logger().Fatal("init err", zap.Error(err))
		return err
	}
	logger.Logger().Info("bot start successfully")
	return nil
}

//Start bot
func Start() {
	setCommend()
	setHandle()
	B.Start()
}
