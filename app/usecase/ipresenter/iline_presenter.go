package ipresenter

import msgdto "nfp-server/usecase/dto"

// ILinePresenter LINEBOTプレゼンタ
type ILinePresenter interface {
	Parrot(token, msg string)
	AskIamge(out msgdto.MsgOutput)
	AskTitle(out msgdto.MsgOutput)
	AskDetail(out msgdto.MsgOutput)
	Confirm(out msgdto.MsgOutput)
}
