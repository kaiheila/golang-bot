package event

type ChannelAddUserEventFrame struct {
	Frame
	ChannelAddUserEvent
}

type ChannelAddUserEvent struct {
	BaseEvent
	extra *ChannelAddUserExtra `json:"extra"`
}
type Emoji struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
type ChannelAddUserExtra struct {
	Type string `json:"type"`
	Body struct {
		ChannelId string `json:"channel_id"`
		Emoji     *Emoji
		UserId    string `json:"user_id"`
		MsgId     string `json:"msg_id"`
	} `json"body"`
}
