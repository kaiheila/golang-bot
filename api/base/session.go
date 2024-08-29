package base

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/gookit/event"
	event2 "github.com/kaiheila/golang-bot/api/base/event"
	log "github.com/sirupsen/logrus"
	"io"
)

const EventReceiveFrame = "EVENT-GLOBAL-RECEIVE_FRAME"
const EventDataFrameKey = "frame"
const EventDataSessionKey = "session"

type Session struct {
	Compressed          int
	ReceiveFrameHandler func(frame *event2.FrameMap) (error, []byte)
	ProcessDataHandler  func(data []byte) (error, []byte)
	EventSyncHandle     bool
}

func (s *Session) On(message string, handler event.Listener) {
	event.On(message, handler)
}
func (s *Session) Trigger(eventName string, params event.M) {
	if s.EventSyncHandle {
		event.Trigger(eventName, params)
	} else {
		event.AsyncFire(event.NewBasic(eventName, params))
	}
}

func (s *Session) ReceiveData(data []byte) (error, []byte) {
	if s.Compressed == 1 {
		b := bytes.NewReader(data)
		r, err := zlib.NewReader(b)
		if err != nil {
			return err, nil
		}

		data, err = io.ReadAll(r)
		if err != nil {
			log.Error(err)
			return err, nil
		}

	}
	_, err := sonic.Get(data)
	if err != nil {
		log.Error("Json Unmarshal err.", err)
		return err, nil
	}
	if s.ProcessDataHandler != nil {
		err, data = s.ProcessDataHandler(data)
		if err != nil {
			log.WithError(err).Error("ProcessDataHandler")
			return err, nil
		}
	}
	frame := event2.ParseFrameMapByData(data)
	log.WithField("frame", frame).Info("Receive frame from server")
	if frame != nil {
		if s.ReceiveFrameHandler != nil {
			return s.ReceiveFrameHandler(frame)
		} else {
			return s.ReceiveFrame(frame)
		}
	} else {
		log.Warnf("数据不是合法的frame:%s", string(data))
	}
	return nil, nil

}

func (s *Session) ReceiveFrame(frame *event2.FrameMap) (error, []byte) {
	event.Trigger(EventReceiveFrame, map[string]interface{}{"frame": frame})
	if frame.SignalType == event2.SIG_EVENT {
		eventType := frame.Data["type"]
		if _, ok := frame.Data["channel_type"]; !ok {
			log.Errorf("frame data not contain channel_type,%+v", frame.Data)
			return nil, nil
		}
		channelType := frame.Data["channel_type"].(string)
		if eventType != "" {
			name := fmt.Sprintf("%s_%d", channelType, int64(eventType.(float64)))
			fireEvent := event.NewBasic(name, map[string]interface{}{EventDataFrameKey: frame, EventDataSessionKey: s})
			if s.EventSyncHandle {
				event.Trigger(fireEvent.Name(), fireEvent.Data())
			} else {
				event.AsyncFire(fireEvent)
			}
		}
	}
	return nil, nil

}
