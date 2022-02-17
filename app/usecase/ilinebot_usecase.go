package usecase

import msgdto "nfp-server/usecase/dto"

// ILineBotUseCase ユースケース
type ILineBotUseCase interface {
	Send(msgdto.MsgInput) msgdto.MsgOutput
	GetImage(msgdto.MsgInput) msgdto.MsgOutput
	GetTitle(msgdto.MsgInput) msgdto.MsgOutput
	GetDetail(msgdto.MsgInput) msgdto.MsgOutput
	Confirm(msgdto.MsgInput) msgdto.MsgOutput
	StoreImage()
}
