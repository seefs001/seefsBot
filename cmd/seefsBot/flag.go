package main

import "flag"

var (
	fConfig = flag.String("config", "./config.toml", "配置文件路径")
	fHelp   = flag.Bool("h", false, "show help")

	survivalTimeout = int(3e9)
)
