package bot

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"seefs-bot/model"
)

func queryPingAn(msg string) (infoStr string, success bool) {
	info := model.PingAn{}
	if err := model.DB.Model(&model.PingAn{}).
		Where("name = ?", msg).
		First(&info).Error; err != nil {
		if err := model.DB.Model(&model.PingAn{}).
			Where("phone = ?", msg).
			First(&info).Error; err != nil {
			if err := model.DB.Model(&model.PingAn{}).
				Where("email = ?", msg).
				First(&info).Error; err != nil {
				if err := model.DB.Model(&model.PingAn{}).
					Where("id_card = ?", msg).
					First(&info).Error; err != nil {
					success = false
					return
				}
			}
		}
	}
	success = true
	infoStr = fmt.Sprintf("姓名:%s\n手机号：%s\n身份证：%s\n性别：%s\nEmail:%s\n地址：%s\n月收入:%s\n是否已婚:%s\n\n-----来源：平安保险",
		info.Name, info.Phone, info.IDCard, info.Gender, info.Email, info.Province+info.City, info.MonthInCome,
		info.IsMarried)
	return
}

func queryCar(msg string) (infoStr string, success bool) {
	info := model.Car{}
	if err := model.DB.Model(&model.Car{}).
		Where("name = ?", msg).
		First(&info).Error; err != nil {
		if err := model.DB.Model(&model.Car{}).
			Where("phone = ?", msg).
			First(&info).Error; err != nil {
			if err := model.DB.Model(&model.Car{}).
				Where("email = ?", msg).
				First(&info).Error; err != nil {
				if err := model.DB.Model(&model.Car{}).
					Where("id_card = ?", msg).
					First(&info).Error; err != nil {
					success = false
					return
				}
			}
		}
	}
	success = true
	infoStr = fmt.Sprintf("姓名:%s\n手机号：%s\n身份证：%s\n性别：%s\nEmail:%s\n地址：%s\n月收入:%s\n是否已婚:%s\n生日：%s\n行业：%s\n"+
		"教育程度：%s\n\n-----来源：全国购车数据",
		info.Name, info.Phone, info.IDCard, info.Gender, info.Email, info.Address, info.MonthInCome,
		info.IsMarried, info.Birthday, info.Industry, info.Education)
	return
}

func queryQQ(msg string) (infoStr string, success bool) {
	info := model.QQ{}
	if err := model.DB.Model(&model.QQ{}).
		Where("qq = ?", msg).
		First(&info).Error; err != nil {
		if err := model.DB.Model(&model.QQ{}).
			Where("phone = ?", msg).
			First(&info).Error; err != nil {
			success = false
			return
		}
	}
	success = true
	infoStr = fmt.Sprintf("手机号：%s\nQQ：%s\n\n-----来源：泄露的QQ数据库", info.Phone, info.QQ)
	return
}

func querySf(msg string) (infoStr string, success bool) {
	info := model.SF{}
	if err := model.DB.Model(&model.SF{}).
		Where("name = ?", msg).
		First(&info).Error; err != nil {
		if err := model.DB.Model(&model.SF{}).
			Where("phone1 = ?", msg).
			First(&info).Error; err != nil {
			if err := model.DB.Model(&model.SF{}).
				Where("phone2 = ?", msg).
				First(&info).Error; err != nil {
				success = false
				return
			}
		}
	}
	success = true
	infoStr = fmt.Sprintf("姓名:%s\n手机号1：%s\n手机号2：%s\n地址：%s\n省份：%s\n城市:%s\n\n-----来源：顺丰快递",
		info.Name, info.Phone1, info.Phone2, info.Address, info.Province, info.City)
	return
}

func queryWebQQAPI(msg string) (string, bool) {
	uid := msg
	url := "https://api.data007.org/api/?user=TalkTik&key=CbkSZ4CP1ScptW53VpaqqS22&engine=QQ&choice=qq&content=" + uid
	response, err := http.Get(url)
	if err != nil {
		return "", false
	}
	defer response.Body.Close()
	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", false
	}
	status := gjson.Get(string(result), "status")
	if status.String() == "false" {
		return "", false
	}
	phoneArray := gjson.Get(string(result), "result")
	for _, res := range phoneArray.Array() {
		phone := res.Get("phone").String()
		qq := res.Get("qq").String()
		return fmt.Sprintf("QQ:%s\n手机号:%s\n\n-----来源：QQ数据", qq, phone), true
	}
	return "", false
}

func queryWebWBAPI(msg string) (string, bool) {
	uid := msg
	url := "https://api.data007.org/api/?user=TalkTik&key=CbkSZ4CP1ScptW53VpaqqS22&engine=Weibo&choice=uid&content" +
		"=" + uid
	response, err := http.Get(url)
	if err != nil {
		return "", false
	}
	defer response.Body.Close()
	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", false
	}
	status := gjson.Get(string(result), "status")
	if status.String() == "false" {
		return "", false
	}
	phoneArray := gjson.Get(string(result), "result")
	for _, res := range phoneArray.Array() {
		phone := res.Get("phone").String()
		return fmt.Sprintf("UID:%s\n手机号:%s\n\n-----来源：微博", uid, phone), true
	}
	return "", false
}

func queryWebWBAPIByPhone(msg string) (string, bool) {
	phone := msg
	url := "https://api.data007.org/api/?user=TalkTik&key=CbkSZ4CP1ScptW53VpaqqS22&engine=Weibo&choice=phone&content" +
		"=" + phone
	response, err := http.Get(url)
	if err != nil {
		return "", false
	}
	defer response.Body.Close()
	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", false
	}
	status := gjson.Get(string(result), "status")
	if status.String() == "false" {
		return "", false
	}
	phoneArray := gjson.Get(string(result), "result")
	for _, res := range phoneArray.Array() {
		uid := res.Get("uid").String()
		return fmt.Sprintf("UID:%s\n手机号:%s\n\n-----来源：微博", uid, phone), true
	}
	return "", false
}

func queryWebQQAPIByPhone(msg string) (string, bool) {
	phone := msg
	url := "https://api.data007.org/api/?user=TalkTik&key=CbkSZ4CP1ScptW53VpaqqS22&engine=QQ&choice=phone&content" +
		"=" + phone
	response, err := http.Get(url)
	if err != nil {
		return "", false
	}
	defer response.Body.Close()
	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", false
	}
	status := gjson.Get(string(result), "status")
	if status.String() == "false" {
		return "", false
	}
	phoneArray := gjson.Get(string(result), "result")
	for _, res := range phoneArray.Array() {
		qq := res.Get("qq").String()
		return fmt.Sprintf("QQ:%s\n手机号:%s\n\n-----来源：QQ", qq, phone), true
	}
	return "", false
}
