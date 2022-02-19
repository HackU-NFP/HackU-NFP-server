package msgdto

// MsgInput DTO
type MsgInput struct {
	ReplyToken string
	LineUserID string
	Msg        string
}

type SuccessInput struct {
	TokenType string
	UserId    string
	Tx        string
	Name      string
}
