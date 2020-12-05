package bot

import (
	"errors"
	"fmt"
	"github.com/seefs001/seefslib-go/xconvertor"
	"github.com/seefs001/seefslib-go/xrandom"
	"github.com/spf13/viper"
	tb "gopkg.in/tucnak/telebot.v2"
	"gorm.io/gorm"
	"seefs-bot/model"
	"seefs-bot/pkg/cache"
	"seefs-bot/pkg/logger"
	"strconv"
	"sync"
)

func start(m *tb.Message) {
	if !m.Private() {
		return
	}
	logger.Info(fmt.Sprintf("%s使用了/start", m.Sender.Username))
	code := xrandom.GenRandomCode(8)
	user := &model.User{
		ID:         m.Sender.ID,
		UserName:   m.Sender.Username,
		Score:      0,
		FreeScore:  viper.GetInt64("score.free_score"),
		InviteCode: code,
	}
	model.DB.Model(&model.User{}).Create(user)
	logger.Info(fmt.Sprintf("%s创建账号", m.Sender.Username))
	userID := strconv.FormatInt(int64(m.Sender.ID), 10)
	_, err := cache.Cache.Get(userID)
	if err == nil {
		_, _ = B.Send(m.Chat, fmt.Sprintf("发送频率过快"))
		logger.Info(fmt.Sprintf("%s发送频率过快", m.Sender.Username))
		return
	}

	err = cache.Cache.Set(userID, []byte("true"))
	if err != nil {
		logger.Info(fmt.Sprintf("%s cache set err", m.Sender.Username))
		_, _ = B.Send(m.Chat, fmt.Sprintf("发送频率过快"))
		return
	}
	defer func() {
		err = cache.Cache.Delete(userID)
		if err != nil {
			logger.Info(fmt.Sprintf("%s delete cache err", m.Sender.Username))
			_, _ = B.Send(m.Chat, fmt.Sprintf("发送频率过快"))
			return
		}
	}()
	_, _ = B.Send(m.Chat, fmt.Sprintf(
		"请发送关键字查询。\n\n"+
			"可查询到：身份户籍、手机机主、开房记录、快递地址、学信网、车牌车主、个人常用密码、顺丰物流、QQ/邮箱/微博/网络账号、就职单位和银行开户等联系方式相关信息。\n\n"+
			"官方交流群: https://t.me/tianyancha123\n"+
			"私聊低价代充的均是骗子，请勿上当！\n\n"+
			"有事找客服: @NanwenBot\n"+
			"您的UID: %d  邀请码为：%s\n"+
			" 当前积分: %d    获取积分：/get_score\n"+
			"免费次数: %d次",
		user.ID,
		user.InviteCode,
		user.Score,
		user.FreeScore))
}

func suggest(m *tb.Message) {
	if !m.Private() {
		return
	}
	logger.Info(fmt.Sprintf("/suggest chat_id: %d  username:%s", m.Sender.ID, m.Sender.Username))
	_, _ = B.Send(m.Chat, fmt.Sprintf("联系客服：@tianyan88bot"))
}

func getScore(m *tb.Message) {
	if !m.Private() {
		return
	}
	userID := strconv.FormatInt(int64(m.Sender.ID), 10)
	_, err := cache.Cache.Get(userID)
	if err == nil {
		logger.Info(fmt.Sprintf("%s cache set err", m.Sender.Username))
		_, _ = B.Send(m.Chat, fmt.Sprintf("发送频率过快"))
		return
	}

	err = cache.Cache.Set(userID, []byte("true"))
	if err != nil {
		logger.Info(fmt.Sprintf("%s cache set err", m.Sender.Username))
		_, _ = B.Send(m.Chat, fmt.Sprintf("发送频率过快"))
		return
	}
	defer func() {
		err = cache.Cache.Delete(userID)
		if err != nil {
			logger.Info(fmt.Sprintf("%s cache delete err", m.Sender.Username))
			_, _ = B.Send(m.Chat, fmt.Sprintf("发送频率过快"))
			return
		}
	}()
	user := model.User{}
	if err := model.DB.First(&user, m.Sender.ID).Error; err != nil {
		logger.Info(fmt.Sprintf("%s没有找到他的", m.Sender.Username))
		_, _ = B.Send(m.Chat, fmt.Sprintf("请先发送/start初始化账号"))
		return
	}
	_, _ = B.Send(m.Chat, fmt.Sprintf(
		"获取积分成功，积分为%d,你今天还有%d次免费次数，可发送你的唯一邀请码推广，双方各奖励5积分\n"+
			"您的邀请码为：%s",
		user.Score,
		user.FreeScore,
		user.InviteCode))
}

func invite(userID int, code string) (bool, error) {
	beInviteCode := code
	if beInviteCode == "" {
		logger.Info(fmt.Sprintf("/invite user id :%d 数值不合法", userID))
		return false, errors.New("请输入合法的数值")
	}
	user := model.User{}
	if err := model.DB.First(&user, userID).Error; err != nil {
		logger.Info(fmt.Sprintf("/invite user id :%d 请先初始化", userID))
		return false, errors.New("没有你这个用户，请先/start初始化")
	}
	if user.BeInvitedCode != nil {
		return false, nil
	}
	if user.InviteCode == beInviteCode {
		logger.Info(fmt.Sprintf("/invite user id :%d 绑定自己", userID))
		return false, errors.New("你不能绑定自己哦")
	}
	inviteUser := model.User{}
	if err := model.DB.Model(&model.User{}).
		Where("invite_code = ?", beInviteCode).
		First(&inviteUser).Error; err != nil {
		return false, nil
	}
	// 绑定
	if err := model.DB.Model(&model.User{}).
		Where("id = ?", user.ID).
		Update("be_invited_code", inviteUser.InviteCode).Error; err != nil {
		return false, errors.New("绑定出错")
	}
	// 为双方添加积分
	if err := model.DB.Model(&model.User{}).
		Where("id = ?", user.ID).
		Update("score",
			gorm.Expr("score + ?", viper.GetInt("score.be_rewarded_score"))).Error; err != nil {
		return false, errors.New("为自己添加积分出错")
	}
	if err := model.DB.Model(&model.User{}).
		Where("id = ?", inviteUser.ID).
		Update("score",
			gorm.Expr("score + ?", viper.GetInt("score.reward_score"))).Error; err != nil {
		return false, errors.New("为邀请者添加积分出错")
	}
	return true, nil
}

func recharge(m *tb.Message) {
	if !m.Private() {
		return
	}
	userID := strconv.FormatInt(int64(m.Sender.ID), 10)
	_, err := cache.Cache.Get(userID)
	if err == nil {
		_, _ = B.Send(m.Chat, fmt.Sprintf("发送频率过快"))
		return
	}

	err = cache.Cache.Set(userID, []byte("true"))
	if err != nil {
		_, _ = B.Send(m.Chat, fmt.Sprintf("发送频率过快"))
		return
	}
	defer func() {
		err = cache.Cache.Delete(userID)
		if err != nil {
			_, _ = B.Send(m.Chat, fmt.Sprintf("发送频率过快"))
			return
		}
	}()
	cardCode := m.Payload
	if cardCode == "" {
		logger.Info(fmt.Sprintf("%s /recharge 不按卡密格式", m.Sender.Username))
		_, _ = B.Send(m.Chat, fmt.Sprintf("请按 /card 卡密 格式充值您的卡密获取积分，如没有，请联系客服"))
		return
	}
	user := model.User{}
	if err := model.DB.First(&user, m.Sender.ID).Error; err != nil {
		logger.Info(fmt.Sprintf("%s /recharge 先发送/start初始化账号", m.Sender.Username))
		_, _ = B.Send(m.Chat, fmt.Sprintf("请先发送/start初始化账号"))
		return
	}
	card := model.Card{}
	if err := model.DB.Model(&model.Card{}).
		Where("content = ?", cardCode).
		First(&card).Error; err != nil {
		logger.Info(fmt.Sprintf("%s /recharge 没有找到该卡密", m.Sender.Username))
		_, _ = B.Send(m.Chat, fmt.Sprintf("没有找到该卡密"))
		return
	}
	if err := model.DB.Model(&model.User{}).
		Where("id = ?", user.ID).
		Update("score",
			gorm.Expr("score + ?", card.Score)).Error; err != nil {
		logger.Info(fmt.Sprintf("%s /recharge 充值积分出错", m.Sender.Username))
		_, _ = B.Send(m.Chat, fmt.Sprintf("充值积分出错，请联系管理员"))
		return
	}
	if err := model.DB.Model(&model.Card{}).
		Where("id = ?", card.ID).
		Delete(&card).Error; err != nil {
		logger.Info(fmt.Sprintf("%s /recharge 卡密处理出错", m.Sender.Username))
		_, _ = B.Send(m.Chat, fmt.Sprintf("卡密处理出错，请联系管理员"))
		return
	}
	if err := model.DB.Model(&model.User{}).
		First(&user, user.ID).Error; err != nil {
		logger.Info(fmt.Sprintf("%s /recharge 查询用户积分出错", m.Sender.Username))
		_, _ = B.Send(m.Chat, fmt.Sprintf("查询用户积分出错，请联系管理员"))
		return
	}
	logger.Info(fmt.Sprintf("%s /recharge 充值成功,当前积分%d", m.Sender.Username, user.Score))
	_, _ = B.Send(m.Chat, fmt.Sprintf("卡密充值成功,你的当前积分为%d", user.Score))
}

func getCardInfo(m *tb.Message) {
	if !m.Private() {
		return
	}
	cardCode := m.Payload
	if cardCode == "" {
		_, _ = B.Send(m.Chat, fmt.Sprintf("请按 /get_card_info 卡密 格式发送"))
		return
	}
	card := model.Card{}
	if err := model.DB.Model(&model.Card{}).
		Where("content = ?", cardCode).
		First(&card).Error; err != nil {
		_, _ = B.Send(m.Chat, fmt.Sprintf("该卡密不存在或已经被使用"))
		return
	}
	_, _ = B.Send(m.Chat, fmt.Sprintf("卡号：%s\n卡内积分：%d", card.Content, card.Score))
}

func sendNotice(m *tb.Message) {
	if !m.Private() {
		return
	}
	if m.Sender.Username != viper.GetString("admin.username") {
		_, _ = B.Send(m.Chat, fmt.Sprintf("您没有权限进行此操作!"))
		return
	}
	if m.Payload == "" {
		_, _ = B.Send(m.Chat, fmt.Sprintf("请按 /send_notice 广播内容 格式输入!"))
		return
	}
	// 广播消息
	var users []model.User
	model.DB.
		Model(&model.User{}).
		Find(&users)
	for _, user := range users {
		B.Send(tb.ChatID(user.ID), m.Payload, &tb.SendOptions{
			ParseMode: tb.ModeMarkdown,
		})
	}
	_, _ = B.Send(m.Chat, fmt.Sprintf("广播成功!"))
}

func genCard(m *tb.Message) {
	if !m.Private() {
		return
	}
	if m.Sender.Username != viper.GetString("admin.username") {
		_, _ = B.Send(m.Chat, fmt.Sprintf("您没有权限进行此操作!"))
		return
	}
	if m.Payload == "" {
		_, _ = B.Send(m.Chat, fmt.Sprintf("请按 /gen_card 卡内积分 格式输入!"))
		return
	}
	score, err := xconvertor.StringToInt64(m.Payload)
	if err != nil {
		_, _ = B.Send(m.Chat, fmt.Sprintf("请输入合法的数字!"))
		return
	}
	cardCode := xrandom.RandStringRunes(10)
	if err = model.DB.Model(&model.Card{}).
		Create(&model.Card{
			Content: cardCode,
			Score:   score,
		}).Error; err != nil {
		_, _ = B.Send(m.Chat, fmt.Sprintf("生成卡密失败，请联系管理员!"))
		return
	}
	_, _ = B.Send(m.Chat, fmt.Sprintf("生成了一张卡密：%s\n内含%d积分", cardCode, score))
}

//func queryAPI(m *tb.Message) {
//	if !m.Private() {
//		return
//	}
//	user := model.User{}
//	if err := model.DB.First(&user, m.Sender.ID).Error; err != nil {
//		_, _ = B.Send(m.Chat, fmt.Sprintf("请先发送/start初始化账号"))
//		return
//	}
//	freeScore := model.FreeScore{}
//	if err := model.DB.Where("user_id = ?", m.Sender.ID).First(&freeScore).Error; err != nil {
//		_, _ = B.Send(m.Chat, fmt.Sprintf("请先发送/start初始化账号"))
//		return
//	}
//	if freeScore.Score <= 0 && user.Score <= 0 {
//		_, _ = B.Send(m.Chat, fmt.Sprintf("您的积分或者免费额度不足，当前积分为%d，免费额度为%d次，可以通过填写邀请码或者卡密充值获得积分哦。", user.Score,
//			freeScore.Score))
//		return
//	}
//	uid := m.Text
//	url := "https://api.data007.org/api/?user=TalkTik&key=CbkSZ4CP1ScptW53VpaqqS22&engine=Weibo&choice=uid&content=" + uid
//	response, err := http.Get(url)
//	if err != nil {
//		_, _ = B.Send(m.Chat, "查询出错，可能API出现了问题")
//		return
//	}
//	defer response.Body.Close()
//	result, err := ioutil.ReadAll(response.Body)
//	if err != nil {
//		_, _ = B.Send(m.Chat, "查询出错，可能API出现了问题")
//		return
//	}
//	if freeScore.Score > 0 {
//		if err := model.DB.Model(&model.FreeScore{}).
//			Where("user_id = ?", user.ID).
//			Update("score",
//				gorm.Expr("score - ?", 1)).Error; err != nil {
//			_, _ = B.Send(m.Chat, "免费额度扣除失败，退出查询")
//			return
//		}
//	} else {
//		if err := model.DB.Model(&model.User{}).
//			Where("id = ?", user.ID).
//			Update("score",
//				gorm.Expr("score - ?", 1)).Error; err != nil {
//			_, _ = B.Send(m.Chat, "积分扣除失败，退出查询")
//			return
//		}
//	}
//	status := gjson.Get(string(result), "status")
//	if status.String() == "false" {
//		_, _ = B.Send(m.Chat, fmt.Sprintf("查询完毕，手机号没有泄露，当前剩余积分%d，当前剩余免费额度%d次", user.Score-1,
//			freeScore.Score))
//		return
//	}
//	phoneArray := gjson.Get(string(result), "result")
//	for _, res := range phoneArray.Array() {
//		phone := res.Get("phone").String()
//		logger.Info(fmt.Sprintf("uid:%s---->phone:%s，当前剩余积分%d，当前剩余免费额度%d次", uid, phone, user.Score-1, freeScore.Score))
//		_, _ = B.Send(m.Chat, fmt.Sprintf("查询成功，uid:%s---->phone:%s，当前剩余积分%d，当前剩余免费额度%d次", uid, phone, user.Score-1,
//			freeScore.Score))
//	}
//}

func queryInfo(m *tb.Message) {
	if !m.Private() {
		return
	}
	userID := strconv.FormatInt(int64(m.Sender.ID), 10)
	_, err := cache.Cache.Get(userID)
	if err == nil {
		_, _ = B.Send(m.Chat, fmt.Sprintf("发送频率过快"))
		return
	}

	err = cache.Cache.Set(userID, []byte("true"))
	if err != nil {
		_, _ = B.Send(m.Chat, fmt.Sprintf("发送频率过快"))
		return
	}
	defer func() {
		err = cache.Cache.Delete(userID)
		if err != nil {
			_, _ = B.Send(m.Chat, fmt.Sprintf("发送频率过快"))
			return
		}
	}()
	user := model.User{}
	if err := model.DB.First(&user, m.Sender.ID).Error; err != nil {
		logger.Info(fmt.Sprintf("%s /query 请先初始化", m.Sender.Username))
		_, _ = B.Send(m.Chat, fmt.Sprintf("请先发送/start初始化账号"))
		return
	}
	success, err := invite(m.Sender.ID, m.Text)
	if err != nil {
		logger.Info(fmt.Sprintf("%s /query invite err %s", m.Sender.Username, err.Error()))
		_, _ = B.Send(m.Chat, fmt.Sprintf(err.Error()))
		return
	}
	if success {
		logger.Info(fmt.Sprintf("%s /query 邀请码绑定成功", m.Sender.Username))
		_, _ = B.Send(m.Chat, fmt.Sprintf("邀请码绑定成功，双方账户自动奖励5积分，感谢对天眼的支持。"))
		return
	}
	if user.FreeScore <= 0 && user.Score <= 0 {
		logger.Info(fmt.Sprintf("%s /query 积分不足", m.Sender.Username))
		_, _ = B.Send(m.Chat, fmt.Sprintf("您的积分或者免费额度不足，当前积分为%d，免费额度为%d次，可以通过填写邀请码或者卡密充值获得积分哦。", user.Score,
			user.FreeScore))
		return
	}
	msg := m.Text
	exist := false
	wg := &sync.WaitGroup{}
	wg.Add(8)
	go func() {
		info, success := queryCar(msg)
		if success {
			_, _ = B.Send(m.Chat, info)
			exist = true
		}
		wg.Done()
	}()
	go func() {
		info, success := queryPingAn(msg)
		if success {
			_, _ = B.Send(m.Chat, info)
			exist = true
		}
		wg.Done()
	}()
	go func() {
		info, success := queryQQ(msg)
		if success {
			_, _ = B.Send(m.Chat, info)
			exist = true
		}
		wg.Done()
	}()
	go func() {
		info, success := querySf(msg)
		if success {
			_, _ = B.Send(m.Chat, info)
			exist = true
		}
		wg.Done()
	}()
	go func() {
		info, success := queryWebQQAPI(msg)
		if success {
			_, _ = B.Send(m.Chat, info)
			exist = true
		}
		wg.Done()
	}()
	go func() {
		info, success := queryWebQQAPIByPhone(msg)
		if success {
			_, _ = B.Send(m.Chat, info)
			exist = true
		}
		wg.Done()
	}()
	go func() {
		info, success := queryWebWBAPIByPhone(msg)
		if success {
			_, _ = B.Send(m.Chat, info)
			exist = true
		}
		wg.Done()
	}()
	go func() {
		info, success := queryWebWBAPI(msg)
		if success {
			_, _ = B.Send(m.Chat, info)
			exist = true
		}
		wg.Done()
	}()
	wg.Wait()
	if exist {
		logger.Info(strconv.FormatInt(user.Score, 10))
		logger.Info(strconv.FormatInt(user.FreeScore, 10))
		if user.FreeScore > 0 {
			if err := model.DB.Model(&model.User{}).
				Where("id = ?", user.ID).
				Update("free_score",
					gorm.Expr("free_score - ?", 1)).Error; err != nil {
				logger.Info(fmt.Sprintf("%s /query 免费额度扣除失败", m.Sender.Username))
				_, _ = B.Send(m.Chat, "免费额度扣除失败，退出查询")
				return
			}
		} else {
			if err := model.DB.Model(&model.User{}).
				Where("id = ?", user.ID).
				Update("score",
					gorm.Expr("score - ?", 1)).Error; err != nil {
				logger.Info(fmt.Sprintf("%s /query 积分扣除失败", m.Sender.Username))
				_, _ = B.Send(m.Chat, "积分扣除失败，退出查询")
				return
			}
		}
		return
	}
	logger.Info(fmt.Sprintf("%s /query 没有查到你要的信息", m.Sender.Username))
	_, _ = B.Send(m.Chat, "没有查到你要的信息")
}
