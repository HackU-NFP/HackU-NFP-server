package presenter

import (
	"os"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/sirupsen/logrus"
)

// NFt作成で使用
const msgFail = "作成に失敗しました。再度はじめからお願いしますm(__)m"
const msgAskImage = "NFTにしたい画像を送信してください"
const msgAskTokenTitle = "NFTをタイトルは何にしますか？"
const msgAskTokenMeta = "NFTの詳細説明を文で教えてください"
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

func (presenter *LinePresenter) replyMessage(msg string, replyToken string) {
	res := linebot.NewTextMessage(msg)
	logrus.Debug("replying message: %v", res)
	if _, err := presenter.bot.ReplyMessage(replyToken, res).Do(); err != nil {
		println(err)
		logrus.Errorf("Error LINEBOT replying message: %v", err)
	}
}
