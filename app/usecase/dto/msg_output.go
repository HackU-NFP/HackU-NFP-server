package msgdto

// MsgOutput DTO
type MsgOutput struct {
	ReplyToken string
}
type SuccessOutput struct {
	TokenType  string
	ContractId string
	UserId     string
	TxUri      string
	Name       string
}
