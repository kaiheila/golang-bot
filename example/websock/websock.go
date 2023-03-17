package main

import (
	"github.com/kaiheila/golang-bot/api/base"
	"github.com/kaiheila/golang-bot/example/conf"
	"github.com/kaiheila/golang-bot/example/handler"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.InfoLevel)

	session := base.NewWebSocketSession(conf.Token, conf.BaseUrl, "./session.pid", "", 1)
	session.On(base.EventReceiveFrame, &handler.ReceiveFrameHandler{})
	session.On("GROUP*", &handler.GroupEventHandler{})
	session.On("GROUP_9", &handler.GroupTextEventHandler{Token: conf.Token, BaseUrl: conf.BaseUrl})
	session.Start()
}
