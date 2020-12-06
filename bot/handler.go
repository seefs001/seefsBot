package bot

import tb "gopkg.in/tucnak/telebot.v2"

func setHandle() {
	B.Handle("/start", start)
	B.Handle("/get_score", getScore)
	B.Handle("/suggest", suggest)
	B.Handle("/invite", inviteMsg)
	B.Handle("/card", recharge)
	B.Handle("/send_notice", sendNotice)
	B.Handle("/gen_card", genCard)
	B.Handle(tb.OnText, queryInfo)
}
