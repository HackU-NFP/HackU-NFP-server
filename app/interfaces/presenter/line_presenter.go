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

//loading
func (presenter *LinePresenter) Loading(out msgdto.MsgOutput) {
	replyToken := out.ReplyToken
	presenter.replyMessage(msgLoading, replyToken)
}

// AskImage NFTにする画像たずねる
func (presenter *LinePresenter) AskImage(out msgdto.MsgOutput) {
	replyToken := out.ReplyToken

	res := linebot.NewTextMessage(msgAskImage)

	cameraQuickReplyButton := linebot.NewQuickReplyButton("", linebot.NewCameraAction("カメラ"))
	cameraRollQuickReplyButton := &linebot.QuickReplyButton{
		ImageURL: "",
		Action: linebot.NewCameraRollAction("カメラロール"),
	}
	cancelQuickReplyButton := linebot.NewQuickReplyButton("", linebot.NewMessageAction("キャンセル", "キャンセル"))

	quickReply := linebot.NewQuickReplyItems(cameraQuickReplyButton, cameraRollQuickReplyButton, cancelQuickReplyButton)
	res.WithQuickReplies(quickReply)
	logrus.Debug("replying message: %v", res)
	if _, err := presenter.bot.ReplyMessage(replyToken, res).Do(); err != nil {
		println(err)
		logrus.Errorf("Error LINEBOT replying message: %v", err)
	}
}

// AskTitle NFTのタイトルたずねる
func (presenter *LinePresenter) AskTitle(out msgdto.MsgOutput) {
	replyToken := out.ReplyToken

	res := linebot.NewTextMessage(msgAskTokenTitle)
	cancelQuickReplyButton := linebot.NewQuickReplyButton("", linebot.NewMessageAction("キャンセル", "キャンセル"))
	quickReply := linebot.NewQuickReplyItems(cancelQuickReplyButton)
	res.WithQuickReplies(quickReply)
	logrus.Debug("replying message: %v", res)

	if _, err := presenter.bot.ReplyMessage(replyToken, res).Do(); err != nil {
		println(err)
		logrus.Errorf("Error LINEBOT replying message: %v", err)
	}
}

// AskDetail NFTの説明たずねる
func (presenter *LinePresenter) AskDetail(out msgdto.MsgOutput) {
	replyToken := out.ReplyToken

	res := linebot.NewTextMessage(msgAskTokenMeta)
	cancelQuickReplyButton := linebot.NewQuickReplyButton("", linebot.NewMessageAction("キャンセル", "キャンセル"))
	quickReply := linebot.NewQuickReplyItems(cancelQuickReplyButton)
	res.WithQuickReplies(quickReply)
	logrus.Debug("replying message: %v", res)

	if _, err := presenter.bot.ReplyMessage(replyToken, res).Do(); err != nil {
		println(err)
		logrus.Errorf("Error LINEBOT replying message: %v", err)
	}
}

// Confirm NF作成の確認メッセージ
func (presenter *LinePresenter) Confirm(out msgdto.MsgOutput, image string, title string, meta string) {
	replyToken := out.ReplyToken

	jsonData := []byte(fmt.Sprintf(`{
  "type": "bubble",
  "direction": "ltr",
  "body": {
    "type": "box",
    "layout": "vertical",
    "contents": [
      {
        "type": "text",
        "text": "このNFTを作成しますか？",
        "weight": "bold",
        "size": "md",
        "align": "start"
      },
      {
        "type": "image",
        "url": "%s",
        "size": "full",
        "aspectRatio": "5:4",
        "aspectMode": "cover",
        "margin": "md"
      },
      {
        "type": "text",
        "text": "name: %s",
        "size": "md",
        "contents": [],
        "margin": "md"
      },
      {
        "type": "text",
        "text": "meta: %s",
        "size": "md",
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
        "style": "secondary",
        "margin": "md"
      }
    ]
  }
}`, image, title, meta))

	container, err := linebot.UnmarshalFlexMessageJSON(jsonData)
	if err != nil {
		// 正しくUnmarshalできないinvalidなJSONであればerrが返る
		logrus.Errorf("Error invalid JSON")
	}

	message := linebot.NewFlexMessage("NFTを作成しますか？", container)

	presenter.replyFlexMessage(message, replyToken)
}

//NFT作成成功時のメッセージ
func (presenter *LinePresenter) SuccessMint(out msgdto.SuccessOutput) {
	name := out.Name
	tokenType := out.TokenType
	contractId := out.ContractId
	userId := out.UserId
	txUri := out.TxUri
	image := out.Image

	jsonData := []byte(fmt.Sprintf(`{
		"type": "bubble",
		"hero": {
		  "type": "image",
		  "url": "%s",
		  "size": "full",
		  "aspectRatio": "20:13",
		  "aspectMode": "cover",
		  "action": {
			"type": "uri",
			"label": "Line",
			"uri": "https://linecorp.com/"
		  }
		},
		"body": {
		  "type": "box",
		  "layout": "vertical",
		  "contents": [
			{
			  "type": "text",
			  "text": "%s",
			  "weight": "bold",
			  "size": "xl",
			  "contents": []
			},
			{
			  "type": "box",
			  "layout": "vertical",
			  "spacing": "sm",
			  "margin": "lg",
			  "contents": [
				{
				  "type": "box",
				  "layout": "baseline",
				  "spacing": "sm",
				  "contents": [
					{
					  "type": "text",
					  "text": "tokenType",
					  "size": "sm",
					  "color": "#AAAAAA",
					  "contents": []
					},
					{
					  "type": "text",
					  "text": "%s",
					  "flex": 2,
					  "contents": []
					}
				  ]
				},
				{
				  "type": "box",
				  "layout": "baseline",
				  "spacing": "sm",
				  "contents": [
					{
					  "type": "text",
					  "text": "contractId",
					  "size": "sm",
					  "color": "#AAAAAA",
					  "flex": 1,
					  "wrap": true,
					  "contents": []
					},
					{
					  "type": "text",
					  "text": "%s",
					  "size": "sm",
					  "color": "#666666",
					  "flex": 2,
					  "wrap": true,
					  "contents": []
					}
				  ]
				},
				{
				  "type": "box",
				  "layout": "baseline",
				  "spacing": "sm",
				  "contents": [
					{
					  "type": "text",
					  "text": "owner",
					  "size": "sm",
					  "color": "#AAAAAA",
					  "wrap": true,
					  "contents": []
					},
					{
					  "type": "text",
					  "text": "%s",
					  "size": "sm",
					  "color": "#666666",
					  "flex": 2,
					  "wrap": true,
					  "contents": []
					}
				  ]
				}
			  ]
			}
		  ]
		},
		"footer": {
		  "type": "box",
		  "layout": "vertical",
		  "flex": 0,
		  "spacing": "sm",
		  "contents": [
			{
			  "type": "button",
			  "action": {
				"type": "uri",
				"label": "トランザクション",
				"uri": "%s"
			  },
			  "height": "sm",
			  "style": "link"
			},
			{
			  "type": "spacer",
			  "size": "sm"
			}
		  ]
		}
	  }`, image, name, tokenType, contractId, userId, txUri))

	container, err := linebot.UnmarshalFlexMessageJSON(jsonData)
	if err != nil {
		// 正しくUnmarshalできないinvalidなJSONであればerrが返る
		logrus.Errorf("Error invalid JSON")
	}

	message := linebot.NewFlexMessage("NFTを作成しました", container)

	presenter.pushFlexMessage(message, userId)
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
func (presenter *LinePresenter) pushFlexMessage(msg *linebot.FlexMessage, userId string) {
	if _, err := presenter.bot.PushMessage(userId, msg).Do(); err != nil {
		println(err)
		logrus.Errorf("Error LINEBOT push message: %v", err)
	}
}

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
