package usecase

import msgdto "nfp-server/usecase/dto"

// ILineBotUseCase ユースケース
type ILineBotUseCase interface {
	Send(msgdto.MsgInput) msgdto.MsgOutput
	Loading(msgdto.MsgInput) msgdto.MsgOutput
	GetImage(msgdto.MsgInput) msgdto.MsgOutput
	GetTitle(msgdto.MsgInput) msgdto.MsgOutput
	GetDetail(msgdto.MsgInput) msgdto.MsgOutput
	SuccessMint(msgdto.SuccessInput) msgdto.SuccessOutput
	Confirm(msgdto.MsgInput, string, string, string) msgdto.MsgOutput
	HowToUse(msgdto.MsgInput) msgdto.MsgOutput
	StoreImage()
}
