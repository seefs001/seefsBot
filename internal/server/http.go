package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/phuslu/log"
	"github.com/seefs001/seefsBot/internal/bot"
	"github.com/seefs001/seefsBot/internal/conf"
	"github.com/seefs001/seefsBot/internal/model"
	"github.com/seefs001/seefsBot/pkg/orm"
	tb "gopkg.in/tucnak/telebot.v2"
)

func Start() error {
	r := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	r.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "pong",
		})
	})
	r.Get("/push", func(c *fiber.Ctx) error {
		secretKey := c.Query("key")
		text := c.Query("text")
		if secretKey == "" || text == "" {
			return c.JSON(fiber.Map{
				"message": "请检查key或text是否传入",
			})
		}
		return handlePush(c, secretKey, text)
	})

	r.Get("/prices", func(ctx *fiber.Ctx) error {
		err := UpdatePrice()
		if err != nil {
			return err
		}
		return ctx.JSON(fiber.Map{
			"message": "success",
		})
	})

	r.Post("/push", func(c *fiber.Ctx) error {
		secretKey := c.FormValue("key")
		text := c.FormValue("text")
		if secretKey == "" || text == "" {
			return c.JSON(fiber.Map{
				"message": "请检查key或text是否传入",
			})
		}
		return handlePush(c, secretKey, text)
	})
	log.Info().Msgf("API init on host %s", conf.GetConf().Server.Addr)
	return r.Listen(conf.GetConf().Server.Addr)
}

func handlePush(c *fiber.Ctx, secretKey, text string) error {
	userID, err := model.GetUserIDBySecretKey(secretKey)
	if err != nil {
		return c.JSON(fiber.Map{
			"message": "key不合法",
		})
	}
	_, err = bot.B.Send(tb.ChatID(userID), text, &tb.SendOptions{
		ParseMode: tb.ModeMarkdown,
	})
	if err != nil {
		return c.JSON(fiber.Map{
			"message": "发送失败",
		})
	}
	return c.JSON(fiber.Map{
		"message": "success",
	})
}

func UpdatePrice() error {
	var coins []model.Coin
	orm.DB().Model(&model.Coin{}).Find(&coins)
	var coinType []string
	for _, coin := range coins {
		coinType = append(coinType, coin.Type)
	}
	fmt.Println(coinType)
	join := strings.Join(coinType, ",")
	fmt.Println(join)
	url := fmt.Sprintf("https://min-api.cryptocompare.com/data/pricemulti?fsyms=%s&tsyms=USD", join)
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	prices := make(map[string]interface{})
	err = json.Unmarshal(body, &prices)
	if err != nil {
		return err
	}
	fmt.Println(prices)
	for k, v := range prices {
		var tmp string
		fmt.Println(k, v)
		orm.DB().Model(&model.Coin{}).
			Where("type = ?", k).Find("price", &tmp)
		priceDbFloat64, err := strconv.ParseFloat(tmp, 64)
		if err != nil {
			return err
		}
		priceFloat64, err := strconv.ParseFloat(v.(map[string]string)["USD"], 64)
		if err != nil {
			return err
		}
		fmt.Println(priceFloat64)
		fmt.Printf("%f ------ %f", priceFloat64, priceDbFloat64)
		orm.DB().Model(&model.Coin{}).
			Where("type = ?", k).Update("price", v)
	}
	return nil
}
