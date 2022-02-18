package controllers

import (
	"nfp-server/usecase"
	msgdto "nfp-server/usecase/dto"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/sirupsen/logrus"
)

// LinebotController LINEBOTコントローラ
type LinebotController struct {
	linebotInteractor    usecase.ILineBotUseCase
	blockchainInteractor usecase.IBlockchainUseCase
	bot                  *linebot.Client
}

// NewLinebotController コンストラクタ
func NewLinebotController(linebotInteractor usecase.ILineBotUseCase, blockchainInteractor usecase.IBlockchainUseCase) *LinebotController {

	secret := os.Getenv("LBOT_SECRET")
	token := os.Getenv("LBOT_TOKEN")

	bot, err := linebot.New(secret, token)
	if err != nil {
		logrus.Fatalf("Error creating LINEBOT client: %v", err)
	}

	return &LinebotController{
		linebotInteractor:    linebotInteractor,
		blockchainInteractor: blockchainInteractor,
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
				case *linebot.ImageMessage:
					controller.replyToImageMessage(event)
				}
			} else if event.Type == linebot.EventTypePostback {
				controller.replyToEventTypePostback(event)
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

	if msg == "NFTを作成する" {
		controller.linebotInteractor.GetImage(input)
	} else if msg == "NFTテスト" {
		logrus.Debug("NFTテスト")
		userId := e.Source.UserID
		contractId := os.Getenv("CONTRACT_ID")
		name := "HelloWorld" //TODO: stateの値にする
		meta := "HelloWorld" //TODO:stateの値にする
		controller.mint(e, userId, contractId, name, meta)
	} else {
		state := "" //TODO: state管理

		switch state {
		case "detail":
			controller.linebotInteractor.GetDetail(input)
		case "confirm":
			controller.linebotInteractor.Confirm(input)
		default:
			controller.linebotInteractor.Send(input)
		}
	}
}

func (controller *LinebotController) replyToEventTypePostback(e *linebot.Event) {
	dataMap := createDataMap(e.Postback.Data)

	if dataMap["action"] == "create" {
		//NFT作成するボタン押された時
		// userId := e.Source.UserID
		// contractId := os.Getenv("CONTRACT_ID")
		// name := "gmmm"   //TODO: stateの値にする
		// meta := "gmmmmm" //TODO:stateの値にする
		// // ミント
		// controller.mint(e, userId, contractId, name, meta)

	} else if dataMap["action"] == "cancel" {
		//キャンセルボタン押された時

	}
}

func (controller *LinebotController) replyToImageMessage(e *linebot.Event) {
	//TODO: 受け取った画像の処理

	input := msgdto.MsgInput{
		ReplyToken: e.ReplyToken,
		Msg:        "NFT画像",
	}

	controller.linebotInteractor.GetTitle(input)
}

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

func (controller *LinebotController) mint(e *linebot.Event, userId, contractId, name, meta string) {
	loadingInput := msgdto.MsgInput{
		ReplyToken: e.ReplyToken,
		Msg:        "作成中です...",
	}
	//作成中です...メッセージ送信
	controller.linebotInteractor.Loading(loadingInput)

	txhash, err := controller.blockchainInteractor.CreateNonFungible(userId, contractId, name, meta)
	time.Sleep(time.Second * 3) //sleep
	if err != nil {
		logrus.Debug("NFTの作成に失敗しました: ", err)
		return
	}
	tx, err := controller.blockchainInteractor.GetTransaction(txhash.TxHash)
	time.Sleep(time.Second * 3) //sleep
	if err != nil {
		logrus.Debug("NFTの作成に失敗しました: ", err)
		return
	}
	tokenType := *&tx.Logs[0].Events[0].Attributes[1].Value
	// ミント
	mintTx, err := controller.blockchainInteractor.MintNonFungible(userId, contractId, tokenType, name, meta)
	time.Sleep(time.Second * 5) //sleep
	if err != nil {
		input := msgdto.SuccessInput{
			TokenType: tokenType,
			UserId:    userId,
			Tx:        txhash.TxHash,
			Name:      name,
		}
		//ミント成功メッセージ送信
		controller.linebotInteractor.SuccessMint(input)
		return
	}
	// TODO:画像をstorageにアップロードする.tokenTypeをファイル名にする

	input := msgdto.SuccessInput{
		TokenType: tokenType,
		UserId:    userId,
		Tx:        mintTx.TxHash,
		Name:      name,
	}
	//ミント成功メッセージ送信
	controller.linebotInteractor.SuccessMint(input)
	return
}
