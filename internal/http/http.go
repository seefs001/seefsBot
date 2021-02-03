package http

import (
	"github.com/gin-gonic/gin"
	"github.com/seefs001/seefsBot/internal/bot"
	"github.com/seefs001/seefsBot/internal/conf"
	"github.com/seefs001/seefsBot/internal/model"
	tb "gopkg.in/tucnak/telebot.v2"
)

func Start() error {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/push", func(c *gin.Context) {
		secretKey, contain := c.GetQuery("key")
		text, contain := c.GetQuery("text")
		if !contain {
			c.JSON(200, gin.H{
				"message": "请检查key或text是否传入",
			})
			return
		}
		handlePush(c, secretKey, text)
	})

	r.POST("/push", func(c *gin.Context) {
		secretKey, contain := c.GetPostForm("key")
		text, contain := c.GetPostForm("text")
		if !contain {
			c.JSON(200, gin.H{
				"message": "请检查key或text是否传入",
			})
			return
		}
		handlePush(c, secretKey, text)
	})
	return r.Run(conf.GetConf().Server.Addr)
}

func handlePush(c *gin.Context, secretKey, text string) {
	userID, err := model.GetUserIDBySecretKey(secretKey)
	if err != nil {
		c.JSON(200, gin.H{
			"message": "key不合法",
		})
		return
	}
	_, err = bot.B.Send(tb.ChatID(userID), text, &tb.SendOptions{
		ParseMode: tb.ModeMarkdown,
	})
	if err != nil {
		c.JSON(500, gin.H{
			"message": "发送失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "success",
	})
	return
}
