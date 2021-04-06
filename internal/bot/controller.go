package bot

import (
	"fmt"

	"github.com/phuslu/log"
	"github.com/seefs001/seefsBot/internal/model"
	"github.com/seefs001/seefsBot/pkg/orm"
	"github.com/seefs001/seefsBot/pkg/util"
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
		log.Fatal().Err(err)
	}
}

func neteasencm(m *tb.Message) {
	if !m.Private() {
		return
	}
	_, _ = B.Send(m.Sender, "在写了", &tb.SendOptions{
		ParseMode: tb.ModeMarkdownV2,
	})
}

func updatePrices(m *tb.Message) {
	if !m.Private() {
		return
	}
	price, err := util.UpdatePrice()
	if err != nil {
		log.Info().Msg(err.Error())
		_, _ = B.Send(m.Sender, "出错了")
		return
	}
	for _, result := range price {
		// msg := result.Type + "--------"+strconv.FormatFloat(result.Price,'E',-1,64)+ "--------"+strconv.FormatFloat(result.Increase,'E',-1,64)
		msg := fmt.Sprintf("币种：%s 当前价格：%f 涨幅：%f", result.Type, result.Price, result.Increase)
		log.Info().Msg(msg)
		_, _ = B.Send(m.Sender, msg)
	}
}

func listType(m *tb.Message) {
	if !m.Private() {
		return
	}
	_, _ = B.Send(m.Sender, util.GetCoinType())
}

func addType(m *tb.Message) {
	if !m.Private() {
		return
	}
	coinType := m.Payload
	if coinType == "" {
		_, _ = B.Send(m.Sender, "请使用 /add_type 币种 方法加入监听")
		return
	}
	var coin model.Coin
	result := orm.DB().Model(&model.Coin{}).
		Where("type = ?", coinType).First(&coin).RowsAffected
	if result != 0 {
		_, _ = B.Send(m.Sender, "该币已被监听")
		return
	}
	price, err := util.GetCoinPrice(coinType)
	if err != nil {
		_, _ = B.Send(m.Sender, "系统出错，请联系管理员", err)
		return
	}
	err = orm.DB().Model(&model.Coin{}).
		Create(&model.Coin{
			Type:     coinType,
			Price:    fmt.Sprintf("%f", price),
			Increase: 0,
		}).Error
	if err != nil {
		_, _ = B.Send(m.Sender, "系统出错，请联系管理员", err)
		return
	}
	_, _ = B.Send(m.Sender, "监听成功")
}
