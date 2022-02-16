package usecase

import (
	"nfp-server/interfaces/presenter"
	msgdto "nfp-server/usecase/dto"
)

// MEMO: 参考資料 https://zenn.dev/lilpacy/articles/0d91349742db1c

// LineBotInteractor LINE botインタラクタ
type LineBotInteractor struct {
	lineBotRepository LineBotRepository
	linePresenter     presenter.LinePresenter
}

// メッセージ返すだけ
func (interactor *LineBotInteractor) Send(in msgdto.MsgInput) msgdto.MsgOutput {
	out := msgdto.MsgOutput{
		ReplyToken: in.ReplyToken,
	}
	if out.ReplyToken != "" {
		interactor.linePresenter.Parrot(in.ReplyToken, in.Msg)
	}

	return out
}

// storeImage NFT化する画像を受け取りstorageに保存する tokenIdに対応するurlにする必要あり
func (interactor *LineBotInteractor) storeImage() {

}
