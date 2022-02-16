package usecase

import msgdto "nfp-server/usecase/dto"

// ILineBotUseCase ユースケース
type ILineBotUseCase interface {
	Send(msgdto.MsgInput) msgdto.MsgOutput
	StoreImage()
}
