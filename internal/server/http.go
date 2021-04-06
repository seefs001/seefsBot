package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/phuslu/log"
	"github.com/seefs001/seefsBot/internal/bot"
	"github.com/seefs001/seefsBot/internal/conf"
	"github.com/seefs001/seefsBot/internal/model"
	"github.com/seefs001/seefsBot/pkg/util"
	tb "gopkg.in/tucnak/telebot.v2"
)

var Server = &fiber.App{}

func Start() error {
	Server = fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	Server.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "pong",
		})
	})
	Server.Get("/push", func(c *fiber.Ctx) error {
		secretKey := c.Query("key")
		text := c.Query("text")
		if secretKey == "" || text == "" {
			return c.JSON(fiber.Map{
				"message": "请检查key或text是否传入",
			})
		}
		return handlePush(c, secretKey, text)
	})

	Server.Get("/update_prices", func(ctx *fiber.Ctx) error {
		results, err := util.UpdatePrice()
		if err != nil {
			return err
		}
		return ctx.JSON(fiber.Map{
			"message": "success",
			"data":    results,
		})
	})

	Server.Post("/push", func(c *fiber.Ctx) error {
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

	return Server.Listen(conf.GetConf().Server.Addr)
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

func Stop() error {
	return Server.Shutdown()
}
