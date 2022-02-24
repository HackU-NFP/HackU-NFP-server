package ipresenter

import msgdto "nfp-server/usecase/dto"

// ILinePresenter LINEBOTプレゼンタ
type ILinePresenter interface {
	Parrot(token, msg string)
	Loading(out msgdto.MsgOutput)
	AskImage(out msgdto.MsgOutput)
	AskTitle(out msgdto.MsgOutput)
	AskDetail(out msgdto.MsgOutput)
	SuccessMint(out msgdto.SuccessOutput)
	Confirm(out msgdto.MsgOutput, image string, title string, meta string)
	HowToUse(out msgdto.MsgOutput)
}
