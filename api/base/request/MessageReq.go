package request

type SendMessageReq struct {
	Type         int    `json:"type"`
	TargetId     string `json:"target_id"`
	Content      string `json:"content"`
	Quote        string `json:"quote"`
	Nonce        string `json:"nonce"`
	TempTargetId string `json:"temp_target_id"`
}
