package event

const (
	SIG_EVENT      int32 = 0
	SIG_HELLO      int32 = 1
	SIG_PING       int32 = 2
	SIG_PONG       int32 = 3
	SIG_RESUME     int32 = 4
	SIG_RECONNECT  int32 = 5
	SIG_RESUME_ACK int32 = 6
	SIG_NACK       int32 = 7
)

// 1:文字消息, 2:图片消息，3:视频消息，4:文件消息， 8:音频消息，9:KMarkdown，10:card 消息，255:系统消息, 其它的暂未开放
const (
	EventTextMsgType   = 1
	EventPicMsgType    = 2
	EventVideoMsgType  = 3
	EventFileMsgType   = 4
	EventVoiceMsgType  = 8
	EventKMDMsgType    = 9
	EVentCardType      = 10
	EventSystemMsgType = 255
)

type EventInterface interface {
	GetType() int
}
type BaseEvent struct {
	ChannelType  string `json:"channel_type"`
	Type         int    `json:"type"`
	TargetId     string `json:"target_id"`
	AuthorId     string `json:"author_id"`
	Content      string `json:"content"`
	MsgId        string `json:"msg_id"`
	MsgTimestamp int64  `json:"msg_timestamp"`
	Nonce        string `json:"nonce"`
	SerialNumber int64  `json:"sn"`
}

func (e *BaseEvent) GetType() int {
	return e.Type
}
