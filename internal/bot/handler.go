package bot

import tb "gopkg.in/tucnak/telebot.v2"

func setHandle() {
	B.Handle("/start", start)
	B.Handle(tb.OnAudio, neteasencm)
}
