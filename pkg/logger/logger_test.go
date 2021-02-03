package logger_test

import (
	"github.com/seefs001/seefsBot/pkg/logger"
	"testing"

	"go.uber.org/zap"
)

func TestMain(t *testing.M) {
	logger.New(logger.SetEnv("dev"), logger.SetPath("./log"))
	t.Run()
}

func TestGetLogger(t *testing.T) {
	logger.Logger().Info("msg", zap.String("uid", "abc"))
	logger.Logger().Debug("debug", zap.String("uid", "abc"))
	logger.Logger().Error("error", zap.String("uid", "abc"))

	// 多实例日志
	logger.Logger("goim").Info("info", zap.String("uid", "abc"))
	logger.Logger("goim").Error("error", zap.String("uid", "abc"))
	logger.Logger("goim").Debug("debug", zap.String("uid", "abc"))
}
