package usecase

import (
	msgdto "nfp-server/usecase/dto"
	"nfp-server/usecase/ipresenter"
	"os"
)

// MEMO: 参考資料 https://zenn.dev/lilpacy/articles/0d91349742db1c

// LineBotInteractor LINE botインタラクタ
type LineBotInteractor struct {
	linePresenter ipresenter.ILinePresenter
}

// NewFavoriteInteractor コンストラクタ
func NewLineBotInteractor(
	linePresenter ipresenter.ILinePresenter) *LineBotInteractor {

	return &LineBotInteractor{
		linePresenter: linePresenter,
	}
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

//画像を送信してくださいメッセージ
func (interactor *LineBotInteractor) GetImage(in msgdto.MsgInput) msgdto.MsgOutput {
	out := msgdto.MsgOutput{
		ReplyToken: in.ReplyToken,
	}
	if out.ReplyToken != "" {
		interactor.linePresenter.AskImage(out)
	}

	return out
}

//画像を受けっとってタイトル送信メッセージ
func (interactor *LineBotInteractor) GetTitle(in msgdto.MsgInput) msgdto.MsgOutput {
	//TODO: 画像保存処理&state変更

	out := msgdto.MsgOutput{
		ReplyToken: in.ReplyToken,
	}
	if out.ReplyToken != "" {
		interactor.linePresenter.AskTitle(out)
	}

	return out
}

//タイトルを受け取って詳細送信してくださいメッセージ
func (interactor *LineBotInteractor) GetDetail(in msgdto.MsgInput) msgdto.MsgOutput {
	//TODO: タイトル保存処理&state変更

	out := msgdto.MsgOutput{
		ReplyToken: in.ReplyToken,
	}
	if out.ReplyToken != "" {
		interactor.linePresenter.AskDetail(out)
	}

	return out
}

//詳細を受け取ってNFT作成確認メッセージ
func (interactor *LineBotInteractor) Confirm(in msgdto.MsgInput, title string, image string, meta string) msgdto.MsgOutput {
	//TODO: 詳細保存処理&state変更

	out := msgdto.MsgOutput{
		ReplyToken: in.ReplyToken,
	}
	if out.ReplyToken != "" {
		interactor.linePresenter.Confirm(out, title, image, meta)
	}

	return out
}

// storeImage NFT化する画像を受け取りstorageに保存する tokenIdに対応するurlにする必要あり
func (interactor *LineBotInteractor) StoreImage() {

}

//ミント成功時のメッセージ送信
func (interactor *LineBotInteractor) SuccessMint(in msgdto.SuccessInput) msgdto.SuccessOutput {
	out := msgdto.SuccessOutput{
		UserId:     in.UserId,
		ContractId: os.Getenv("CONTRACT_ID"),
		TokenType:  in.TokenType,
		TxUri:      "https://explorer.blockchain.line.me/cashew/transaction/" + in.Tx,
		Name:       in.Name,
		Image:      in.Image,
	}
	if out.UserId != "" {
		interactor.linePresenter.SuccessMint(out)
	}

	return out
}

//作成中...
func (interactor *LineBotInteractor) Loading(in msgdto.MsgInput) msgdto.MsgOutput {
	out := msgdto.MsgOutput{
		ReplyToken: in.ReplyToken,
	}
	if out.ReplyToken != "" {
		interactor.linePresenter.Loading(out)
	}

	return out
}

// 使い方
func (interactor *LineBotInteractor) HowToUse(in msgdto.MsgInput) msgdto.MsgOutput {
	out := msgdto.MsgOutput{
		ReplyToken: in.ReplyToken,
	}
	if out.ReplyToken != "" {
		interactor.linePresenter.HowToUse(out)
	}

	return out
}
