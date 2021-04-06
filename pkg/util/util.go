package util

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/phuslu/log"
	"github.com/seefs001/seefsBot/internal/model"
	"github.com/seefs001/seefsBot/pkg/orm"
)

// GenRandomString 生成随机字符串
// length 生成长度
// specialChar 是否生成特殊字符
func GenRandomString(length int, specialChar bool) string {

	letterBytes := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	special := "!@#%$*.="

	if specialChar {
		letterBytes = letterBytes + special
	}

	chars := []byte(letterBytes)

	if length == 0 {
		return ""
	}

	clen := len(chars)
	maxrb := 255 - (256 % clen)
	b := make([]byte, length)
	r := make([]byte, length+(length/4)) // storage for random bytes.
	i := 0
	for {
		if _, err := rand.Read(r); err != nil {
			return ""
		}
		for _, rb := range r {
			c := int(rb)
			if c > maxrb {
				continue // Skip this number to avoid modulo bias.
			}
			b[i] = chars[c%clen]
			i++
			if i == length {
				return string(b)
			}
		}
	}
}

type Result struct {
	Type     string  `json:"type"`
	Price    float64 `json:"price"`
	Increase float64 `json:"increase"`
}

func GetCoinType() string {
	var coins []model.Coin
	orm.DB().Model(&model.Coin{}).Find(&coins)
	if len(coins) == 0 {
		log.Info().Msg("数据库没有要监控的币")
		return ""
	}
	var coinType []string
	for _, coin := range coins {
		coinType = append(coinType, coin.Type)
	}

	return strings.Join(coinType, ",")
}

func UpdatePrice() ([]Result, error) {
	log.Info().Msg("正在查询" + GetCoinType())

	url := fmt.Sprintf("https://min-api.cryptocompare.com/data/pricemulti?fsyms=%s&tsyms=USD", GetCoinType())
	resp, err := http.Get(url) //nolint:gosec
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	prices := make(map[string]interface{})
	err = json.Unmarshal(body, &prices)
	if err != nil {
		return nil, err
	}

	var results []Result
	for k, v := range prices {
		var tmp model.Coin
		orm.DB().Model(&model.Coin{}).
			Where("type = ?", k).First(&tmp)
		priceDBFloat64, err := strconv.ParseFloat(tmp.Price, 64)
		if err != nil {
			return nil, err
		}
		priceFloat64 := v.(map[string]interface{})["USD"]
		increase := (priceFloat64.(float64) - priceDBFloat64) / priceDBFloat64

		log.Info().Float64(k, priceDBFloat64).Msg("当前货币价格")
		log.Info().Float64(fmt.Sprintf("%s--涨幅", k), increase).Msg("涨幅")

		orm.DB().Model(&model.Coin{}).
			Where("type = ?", k).
			Update("increase", increase).
			Update("price", priceFloat64)
		results = append(results, Result{
			Type:     k,
			Price:    priceFloat64.(float64),
			Increase: increase,
		})
	}
	return results, nil
}

func GetCoinPrice(coinType string) (float64, error) {
	url := fmt.Sprintf("https://min-api.cryptocompare.com/data/price?fsym=%s&tsyms=USD", coinType)
	resp, err := http.Get(url) //nolint:gosec
	defer resp.Body.Close()
	if err != nil {
		return 0, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	prices := make(map[string]float64)
	err = json.Unmarshal(body, &prices)
	if err != nil {
		return 0, err
	}
	return prices["USD"], nil
}
