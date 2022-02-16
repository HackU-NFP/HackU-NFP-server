package controllers

import (
	"nfp-server/usecase"
	msgdto "nfp-server/usecase/dto"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/sirupsen/logrus"
)

// LinebotController LINEBOTコントローラ
type LinebotController struct {
	linebotInteractor    usecase.LineBotInteractor
	blockchainInteractor usecase.BlockchainInteractor
	bot                  *linebot.Client
}

// NewLinebotController コンストラクタ
func NewLinebotController() *LinebotController {

	secret := os.Getenv("LBOT_SECRET")
	token := os.Getenv("LBOT_TOKEN")

	bot, err := linebot.New(secret, token)
	if err != nil {
		logrus.Fatalf("Error creating LINEBOT client: %v", err)
	}

	return &LinebotController{
		linebotInteractor:    usecase.LineBotInteractor{},
		blockchainInteractor: usecase.BlockchainInteractor{},
		bot:                  bot,
	}
}

// CatchEvents LINEBOTに関する処理
func (controller *LinebotController) CatchEvents() echo.HandlerFunc {
	return func(c echo.Context) error {
		events, err := controller.bot.ParseRequest(c.Request())
		if err != nil {
			logrus.Fatalf("Error LINEBOT parsing request: %v", err)
		}

		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch event.Message.(type) {
				case *linebot.TextMessage:
					controller.replyToTextMessage(event)
				}
			} else if event.Type == linebot.EventTypePostback {
				// controller.replyToEventTypePostback(event)
			}
		}

		return nil
	}
}

func (controller *LinebotController) replyToTextMessage(e *linebot.Event) {
	msg := e.Message.(*linebot.TextMessage).Text

	input := msgdto.MsgInput{
		ReplyToken: e.ReplyToken,
		Msg:        msg,
	}
	controller.linebotInteractor.Send(input)
}

// func (controller *LinebotController) replyToEventTypePostback(e *linebot.Event) {
// 	dataMap := createDataMap(e.Postback.Data)

// 	if dataMap["action"] == "addFavorite" {
// 		favoriteAddInput := favoritedto.AddInput{
// 			ReplyToken: e.ReplyToken,
// 			LineUserID: e.Source.UserID,
// 			PlaceID:    dataMap["placeId"],
// 		}
// 		controller.favoriteInteractor.Add(favoriteAddInput)
// 	} else if dataMap["action"] == "removeFavorite" {
// 		favoriteRemoveInput := favoritedto.RemoveInput{
// 			ReplyToken: e.ReplyToken,
// 			LineUserID: e.Source.UserID,
// 			PlaceID:    dataMap["placeId"],
// 		}
// 		controller.favoriteInteractor.Remove(favoriteRemoveInput)
// 	}
// }

// createDataMap Postbackで受け取ったデータをパースしてマップ形式で保存する
// e.g.
// input : "action=favorite&placeId=xxxxxx"
// output: dataMap["action"] = "favorite", dataMap["placeId"] = "xxxxx"
func createDataMap(q string) map[string]string {
	dataMap := make(map[string]string)

	dataArr := strings.Split(q, "&")
	for _, data := range dataArr {
		splitedData := strings.Split(data, "=")
		dataMap[splitedData[0]] = splitedData[1]
	}

	return dataMap
}
