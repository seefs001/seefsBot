package logger

import (
	"github.com/phuslu/log"
)

type Options struct {
	LogPath      string    `json:"log_path" mapstructure:"log_path"`
	LogName      string    `json:"log_name" mapstructure:"log_name"`
	ConsoleLevel log.Level `json:"console_level" mapstructure:"console_level"`
}

func Init(options *Options) {
	if options.LogPath == "" {
		options.LogPath = "./logs"
	}

	if options.LogName == "" {
		options.LogName = "main.log"
	}

	if options.ConsoleLevel == 0 {
		options.ConsoleLevel = log.ErrorLevel
	}

	filename := options.LogPath + options.LogName

	fileWriter := &log.FileWriter{
		Filename:   filename,
		MaxSize:    50 * 1024 * 1024,
		MaxBackups: 7,
		LocalTime:  false,
	}
	consoleWriter := &log.ConsoleWriter{
		ColorOutput:    true,
		QuoteString:    true,
		EndWithMessage: true,
	}
	log.DefaultLogger = log.Logger{
		Level:      log.InfoLevel,
		Caller:     1,
		TimeFormat: "0102 15:04:05.999999",
		Writer: &log.MultiWriter{
			ConsoleWriter: consoleWriter,
			InfoWriter:    fileWriter,
			WarnWriter:    fileWriter,
			ErrorWriter:   fileWriter,
			ConsoleLevel:  log.DebugLevel,
		},
	}
}
