package tasks

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/phuslu/log"
	"github.com/robfig/cron/v3"
	"github.com/seefs001/seefsBot/internal/conf"
	"github.com/seefs001/seefsBot/internal/model"
	"github.com/seefs001/seefsBot/pkg/orm"
)

var Task = &cron.Cron{}

type Result struct {
	Type     string  `json:"type"`
	Price    float64 `json:"price"`
	Increase float64 `json:"increase"`
}

func UpdatePrice() ([]Result, error) {
	var coins []model.Coin
	orm.DB().Model(&model.Coin{}).Find(&coins)
	if len(coins) == 0 {
		log.Info().Msg("数据库没有要监控的币")
		return nil, nil
	}
	var coinType []string
	for _, coin := range coins {
		coinType = append(coinType, coin.Type)
	}

	join := strings.Join(coinType, ",")

	log.Info().Msg("正在查询" + join)

	url := fmt.Sprintf("https://min-api.cryptocompare.com/data/pricemulti?fsyms=%s&tsyms=USD", join)
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

func Start() error {
	Task = cron.New(cron.WithSeconds())
	spec := conf.GetConf().Task.Cron
	_, err := Task.AddFunc(spec, func() {
		_, err := UpdatePrice()
		if err != nil {
			log.Warn().Time("time", time.Now()).
				Msg("任务执行失败 ")
		}
	})
	if err != nil {
		return err
	}
	Task.Start()
	return nil
}

func Stop() {
	Task.Stop()
}
