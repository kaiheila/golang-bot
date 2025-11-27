package main

import (
	"github.com/kaiheila/golang-bot/api/base"
	"github.com/kaiheila/golang-bot/api/helper/compress"
	"github.com/kaiheila/golang-bot/example/conf"
	"github.com/kaiheila/golang-bot/example/handler"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.InfoLevel)
	compress.InitZSTDPool("d:/dev/gowork/src/ws-connector/dict/zstd_v1.zip")
	session := base.NewWebSocketSession(conf.Token, conf.BaseUrl, "./session.pid", "", 1, compress.CompressTypeZstdPerMessage, "2", 1)
	session.On(base.EventReceiveFrame, &handler.ReceiveFrameHandler{})
	session.On("GROUP*", &handler.GroupEventHandler{})
	session.On("GROUP_9", &handler.GroupTextEventHandler{Token: conf.Token, BaseUrl: conf.BaseUrl, Session: session})
	session.Start()
}
