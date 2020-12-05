package conf

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"path/filepath"
	"seefs-bot/model"
	"seefs-bot/pkg/cron"
	"seefs-bot/pkg/logger"
)

// Init initialize config
func Init(configPath string) error {

	ext := filepath.Ext(configPath)                   // .yml
	filename := configPath[:len(configPath)-len(ext)] // xxx
	viper.SetConfigName(filename)
	viper.SetConfigType(ext[1:])
	viper.AddConfigPath(".")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(fmt.Errorf("no such config file: %s", err))
		} else {
			panic(fmt.Errorf("read config error: %s", err))
		}
	}
	// run mode
	development := false
	if viper.GetString("meta.run-mode") == "debug" {
		development = true
	}
	// initialize logger
	logger.NewLogger(logger.SetAppName(viper.GetString("meta.name")),
		logger.SetDebugFileName("debug"),
		logger.SetErrorFileName("error"),
		logger.SetWarnFileName("warn"),
		logger.SetInfoFileName("info"),
		logger.SetLogFileDir(viper.GetString("log.log-path")),
		logger.SetMaxAge(viper.GetInt("log.max-age")),
		logger.SetMaxBackups(viper.GetInt("log.max-back-up")),
		logger.SetMaxSize(viper.GetInt("log.max-size")),
		logger.SetLevel(zapcore.DebugLevel),
		logger.SetDevelopment(development))

	// initialize database
	dbParams := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		viper.GetString("mysql.user"), viper.GetString("mysql.password"),
		viper.GetString("mysql.host"), viper.GetInt("mysql.port"),
		viper.GetString("mysql.db"), viper.GetString("mysql.charset"),
	)

	err := model.Database(mysql.Open(dbParams))
	// Failed to connect to MySQL, switch to sqlite
	if err != nil {
		logger.Error(fmt.Sprintf("[%s] Failed to connect to MySQL, switch to sqlite", viper.GetString("meta.name")))
		err := model.Database(sqlite.Open(viper.GetString("meta.name") + ".db"))
		if err != nil {
			logger.Error(fmt.Sprintf("[%s] Failed to enable the sqlite", viper.GetString("meta.name")))
			return err
		}
	}
	cron.Init()
	return nil
}
