package conf

import (
	"os"

	"github.com/BurntSushi/toml"
)

var conf Config

func GetConf() *Config {
	return &conf
}

type Config struct {
	Server *Server `toml:"server" json:"server,omitempty"`
	MySQL  *MySQL  `toml:"mysql" json:"mysql,omitempty"`
	Redis  *Redis  `toml:"redis" json:"redis,omitempty"`
	Log    *Log    `toml:"log" json:"log,omitempty"`
}

type Server struct {
	Addr     string `toml:"addr"`
	BotToken string `toml:"bot_token"`
}

type MySQL struct {
	DSN string `toml:"dsn"`
}

type Log struct {
	LogPath string `json:"log_path" toml:"log_path"`
	LogName string `json:"log_name" toml:"log_name"`
}

type Redis struct {
	Addr string `toml:"addr"`
	DB   int    `toml:"db"`
	Pass string `toml:"pass"`
}

func Init(filepath string) error {
	file, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	_, err = toml.Decode(string(file), &conf)
	if err != nil {
		return err
	}
	return nil
}
