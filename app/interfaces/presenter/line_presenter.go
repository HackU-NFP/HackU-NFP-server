package presenter

import (
	"fmt"
	msgdto "nfp-server/usecase/dto"
	"os"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/sirupsen/logrus"
)

// NFt作成で使用
const msgFail = "作成に失敗しました。再度はじめからお願いしますm(__)m"
const msgAskImage = "NFTにしたい画像を送信してください"
const msgAskTokenTitle = "NFTのタイトルは何にしますか？"
const msgAskTokenMeta = "NFTの詳細な説明を教えてください"
const msgLoading = "NFTを作成中です..."
const msgCancel = "キャンセルしました"

const maxTextWC = 60

// LinePresenter LINEプレゼンタ
type LinePresenter struct {
	bot *linebot.Client
}

type carouselMsgs struct {
	noResult            string
	altText             string
	postbackActionLabel string
	postbackActionData  string
}

// NewLinePresenter コンストラクタ
func NewLinePresenter() *LinePresenter {
	secret := os.Getenv("LBOT_SECRET")
	token := os.Getenv("LBOT_TOKEN")

	bot, err := linebot.New(secret, token)
	if err != nil {
		logrus.Fatalf("Error creating LINEBOT client: %v", err)
	}

	return &LinePresenter{bot: bot}
}

//おうむ返し
func (presenter *LinePresenter) Parrot(token, msg string) {
	replyToken := token
	presenter.replyMessage(msg, replyToken)
}

// AskIamge NFTにする画像たずねる
func (presenter *LinePresenter) AskIamge(out msgdto.MsgOutput) {
	replyToken := out.ReplyToken

	presenter.replyMessage(msgAskImage, replyToken)
}

// AskTitle NFTのタイトルたずねる
func (presenter *LinePresenter) AskTitle(out msgdto.MsgOutput) {
	replyToken := out.ReplyToken

	presenter.replyMessage(msgAskTokenTitle, replyToken)
}

// AskDetail NFTの説明たずねる
func (presenter *LinePresenter) AskDetail(out msgdto.MsgOutput) {
	replyToken := out.ReplyToken

	presenter.replyMessage(msgAskTokenMeta, replyToken)
}

// Confirm NF作成の確認メッセージ
func (presenter *LinePresenter) Confirm(out msgdto.MsgOutput) {
	replyToken := out.ReplyToken
	//TODO: 値を取得する。
	name := "aaa"
	meta := "詳細いいいい"

	jsonData := []byte(fmt.Sprintf(`{
		"type": "bubble",
		"direction": "ltr",
		"body": {
		  "type": "box",
		  "layout": "vertical",
		  "contents": [
			{
			  "type": "box",
			  "layout": "vertical",
			  "contents": [
				{
				  "type": "text",
				  "text": "このNFTを作成してよろしいですか？",
				  "weight": "bold",
				  "size": "md",
				  "align": "start",
				  "margin": "sm",
				  "wrap": true,
				  "contents": []
				},
				{
				  "type": "spacer"
				}
			  ]
			},
			{
			  "type": "text",
			  "text": "name: %s",
			  "size": "sm",
			  "contents": []
			},
			{
			  "type": "text",
			  "text": "meta: %s",
			  "size": "sm",
			  "wrap": true,
			  "contents": []
			}
		  ]
		},
		"footer": {
		  "type": "box",
		  "layout": "horizontal",
		  "contents": [
			{
			  "type": "button",
			  "action": {
				"type": "postback",
				"label": "作成する",
				"text": "作成する",
				"data": "create"
			  },
			  "style": "primary"
			},
			{
			  "type": "button",
			  "action": {
				"type": "postback",
				"label": "キャンセル",
				"text": "キャンセル",
				"data": "cancel_create"
			  },
			  "style": "secondary"
			}
		  ]
		}
	  }`, name, meta))

	container, err := linebot.UnmarshalFlexMessageJSON(jsonData)
	if err != nil {
		// 正しくUnmarshalできないinvalidなJSONであればerrが返る
		logrus.Errorf("Error invalid JSON")
	}

	message := linebot.NewFlexMessage("NFTを作成しますか？", container)

	presenter.replyFlexMessage(message, replyToken)
}

// func (presenter *LinePresenter) replyCarouselColumn(
// 	msgs carouselMsgs, googleMapOutputs []model.Place, replyToken string) {

// 	if len(googleMapOutputs) == 0 {
// 		presenter.replyMessage(msgs.noResult, replyToken)
// 		return
// 	}

// 	ccs := []*linebot.CarouselColumn{}
// 	for _, gmo := range googleMapOutputs {
// 		addr := gmo.Address
// 		if maxTextWC < utf8.RuneCountInString(addr) {
// 			addr = string([]rune(addr)[:maxTextWC])
// 		}

// 		data := fmt.Sprintf(msgs.postbackActionData, gmo.PlaceID)
// 		cc := linebot.NewCarouselColumn(
// 			gmo.PhotoURL,
// 			gmo.Name,
// 			addr,
// 			linebot.NewURIAction("Open Google Map", gmo.URL),
// 			linebot.NewPostbackAction(msgs.postbackActionLabel, data, "", ""),
// 		).WithImageOptions("#FFFFFF")
// 		ccs = append(ccs, cc)
// 	}

// 	res := linebot.NewTemplateMessage(
// 		msgs.altText,
// 		linebot.NewCarouselTemplate(ccs...).WithImageOptions("rectangle", "cover"),
// 	)

// 	if _, err := presenter.bot.ReplyMessage(replyToken, res).Do(); err != nil {
// 		logrus.Errorf("Error LINEBOT replying message: %v", err)
// 	}
// }
func (presenter *LinePresenter) replyFlexMessage(msg *linebot.FlexMessage, replyToken string) {
	if _, err := presenter.bot.ReplyMessage(replyToken, msg).Do(); err != nil {
		println(err)
		logrus.Errorf("Error LINEBOT replying message: %v", err)
	}
}

func (presenter *LinePresenter) replyMessage(msg string, replyToken string) {
	res := linebot.NewTextMessage(msg)
	logrus.Debug("replying message: %v", res)
	if _, err := presenter.bot.ReplyMessage(replyToken, res).Do(); err != nil {
		println(err)
		logrus.Errorf("Error LINEBOT replying message: %v", err)
	}
}