# seefsBot
自个人用的telegram bot

## 交叉编译
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build main.go