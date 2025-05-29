package handler

import (
	"errors"
	"github.com/bytedance/sonic"
	"github.com/gookit/event"
	"github.com/kaiheila/golang-bot/api/base"
	event2 "github.com/kaiheila/golang-bot/api/base/event"
	"github.com/kaiheila/golang-bot/api/helper"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"sync/atomic"
)

type ReceiveFrameHandler struct {
}

func (rf *ReceiveFrameHandler) Handle(e event.Event) error {
	//log.WithField("event", e).WithField("data", e.Data()).Info("ReceiveFrameHandler receive frame.")

	return nil
}

type GroupEventHandler struct {
	ReceivedNum atomic.Int32
}

func (ge *GroupEventHandler) Handle(e event.Event) error {
	ge.ReceivedNum.Add(1)
	log.WithField("receivedNum", ge.ReceivedNum).WithField("e", e).Warn("GroupEventHandler receive event.%+v", e)
	return nil
}

type GroupTextEventHandler struct {
	Token   string
	BaseUrl string
	MsgNum  atomic.Int64
	Session *base.WebSocketSession
}

func (gteh *GroupTextEventHandler) Handle(e event.Event) error {
	//log.WithField("event", fmt.Sprintf("%+v", e.Data())).Info("收到频道内的文字消息.")
	err := func() error {
		if _, ok := e.Data()[base.EventDataFrameKey]; !ok {
			return errors.New("data has no frame field")
		}
		frame := e.Data()[base.EventDataFrameKey].(*event2.FrameMap)
		data, err := sonic.Marshal(frame.Data)
		if err != nil {
			return err
		}
		msgEvent := &event2.MessageKMarkdownEvent{}
		err = sonic.Unmarshal(data, msgEvent)
		gteh.MsgNum.Add(1)
		log.Infof("MsgNum:%d, Received json event:%+v", gteh.MsgNum.Load(), msgEvent)
		if err != nil {
			return err
		}
		if strings.Contains(msgEvent.Content, "nack") {
			items := strings.Split(msgEvent.Content, ":")
			sns := strings.Split(items[1], ",")
			isns := make([]int64, 0)
			for _, sn := range sns {
				isn, err2 := strconv.ParseInt(sn, 10, 64)
				if err2 != nil {
					log.Error(err2)
				}
				isns = append(isns, isn)
			}
			gteh.Session.NAck(isns)
			return nil
		}
		client := helper.NewApiHelper("/v3/message/create", gteh.Token, gteh.BaseUrl, "", "")
		if msgEvent.Author.Bot {
			log.Info("bot message")
			return nil
		}
		echoData := map[string]string{
			"channel_id": msgEvent.TargetId,
			"content":    "echo:" + msgEvent.KMarkdown.RawContent,
		}
		echoDataByte, err := sonic.Marshal(echoData)
		if err != nil {
			return err
		}
		resp, err := client.SetBody(echoDataByte).Post()
		log.Info("sent post:%s", client.String())
		if err != nil {
			log.WithError(err).Errorf("resp:%s", string(resp))
			return err
		}
		log.Infof("resp:%s", string(resp))
		return nil
	}()
	if err != nil {
		log.WithError(err).Error("GroupTextEventHandler err")
	}

	return nil
}
