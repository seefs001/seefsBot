package bot

import tb "gopkg.in/tucnak/telebot.v2"

func setHandle() {
	B.Handle("/start", start)
	B.Handle("/list_type", listType)
	B.Handle("/add_type", addType)
	B.Handle("/price", updatePrices)
	B.Handle(tb.OnAudio, neteasencm)
}
