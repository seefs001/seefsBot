# seefsBot
个人用的telegram bot

## 交叉编译Linux
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build main.go

## 功能
- [x] 类似于Server酱的推送功能
- [] RSS订阅
- [] 网易云ncm转flac
- [] B站任务
- [] 网易云打卡听歌
