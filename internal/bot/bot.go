package bot

import (
	"time"

	"github.com/phuslu/log"
	"github.com/seefs001/seefsBot/internal/bot/fsm"
	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	// UserState 用户状态，用于标示当前用户操作所在状态
	UserState map[int64]fsm.UserStatus = make(map[int64]fsm.UserStatus)

	B *tb.Bot
)

func Init(token string) error {
	poller := &tb.LongPoller{Timeout: 10 * time.Second}
	log.Info().Msg("bot init")

	var err error

	B, err = tb.NewBot(tb.Settings{
		URL:    "https://api.telegram.org",
		Token:  token,
		Poller: poller,
	})
	if err != nil {
		log.Fatal().Err(err)
		return err
	}
	log.Info().Msg("bot start successfully")
	return nil
}

// Start bot
func Start() {
	setCommend()
	setHandle()
	B.Start()
}

func Stop() {
	B.Stop()
}
