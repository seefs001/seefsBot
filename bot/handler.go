package bot

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

func setHandle() {
	B.Handle(tb.OnText, WelcomeMsg)
	B.Handle(tb.OnDocument, WelcomeMsg)
}
