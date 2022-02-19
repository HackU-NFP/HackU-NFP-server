package controllers

import (
	"fmt"
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

// info_key = "state" | "image" | "title" | "meta"
type key struct {
	uid, info_key string
}

// state, image, title, meta 一時保持
var sessions = make(map[key]string)

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
			}
			//	テキスト送信も同時に行われるため、Postback分岐省略
			//  else if event.Type == linebot.EventTypePostback {
			// 	controller.replyToEventTypePostback(event)
			// }
		}

		return nil
	}
}

func (controller *LinebotController) replyToTextMessage(e *linebot.Event) {
	msg := e.Message.(*linebot.TextMessage).Text
	uid := e.Source.UserID

	input := msgdto.MsgInput{
		ReplyToken: e.ReplyToken,
		Msg:        msg,
	}

	if msg == "キャンセル" {
		sessions[key{uid, "state"}] = ""
		input.Msg = "キャンセルしました"
		controller.linebotInteractor.Send(input)
	} else if msg == "NFTを作る" {
		controller.linebotInteractor.GetImage(input)
		sessions[key{uid, "state"}] = "title"
	} else if msg == "NFTテスト" {
		logrus.Debug("NFTテスト")
		userId := e.Source.UserID
		contractId := os.Getenv("CONTRACT_ID")
		name := "HelloWorld" //TODO: stateの値にする
		meta := "HelloWorld" //TODO:stateの値にする
		controller.mint(e, userId, contractId, name, meta)
	} else {
		state := sessions[key{uid, "state"}] //TODO: state管理

		switch state {
		case "title":
			sessions[key{uid, "title"}] = msg
			controller.linebotInteractor.GetDetail(input)
			sessions[key{uid, "state"}] = "detail"
		case "detail":
			sessions[key{uid, "meta"}] = msg
			controller.linebotInteractor.Confirm(input, sessions[key{uid, "image"}], sessions[key{uid, "title"}], sessions[key{uid, "meta"}])
			sessions[key{uid, "state"}] = "confirm"
		case "confirm":
			if msg == "作成する" {
				userId := e.Source.UserID
				contractId := os.Getenv("CONTRACT_ID")
				name := sessions[key{uid, "title"}]
				meta := sessions[key{uid, "meta"}]
				// mint
				controller.mint(e, userId, contractId, name, meta)
			} else {
				controller.linebotInteractor.Confirm(input, sessions[key{uid, "image"}], sessions[key{uid, "title"}], sessions[key{uid, "meta"}])
			}
		default:
			// 使い方送信
			input.Msg = "使い方"
			controller.linebotInteractor.Send(input)
		}
	}
}

func (controller *LinebotController) replyToEventTypePostback(e *linebot.Event) {
	fmt.Println(e.Postback.Data)
	// dataMap := createDataMap(e.Postback.Data)
	// msg := e.Message.(*linebot.TextMessage).Text
	uid := e.Source.UserID
	state := sessions[key{uid, "state"}]

	input := msgdto.MsgInput{
		ReplyToken: e.ReplyToken,
		Msg:        "",
	}

	if state == "confirm" {
		if e.Postback.Data == "create" {
			//NFT作成するボタン押された時
			input.Msg = "作成中..."
			controller.linebotInteractor.Send(input)
			// mint
			//NFT作成するボタン押された時
			// userId := e.Source.UserID
			// contractId := os.Getenv("CONTRACT_ID")
			// name := "gmmm"   //TODO: stateの値にする
			// meta := "gmmmmm" //TODO:stateの値にする
			// // ミント
			// controller.mint(e, userId, contractId, name, meta)
		} else if e.Postback.Data == "cancel_create" {
			//キャンセルボタン押された時
			sessions[key{uid, "state"}] = ""
			input.Msg = "キャンセルしました"
			controller.linebotInteractor.Send(input)
		}
	} else {
		// エラー初期化
		sessions[key{uid, "state"}] = ""
		input.Msg = "エラーが発生しました、最初からやり直してください"
		controller.linebotInteractor.Send(input)
	}
}

func (controller *LinebotController) replyToImageMessage(e *linebot.Event) {
	uid := e.Source.UserID
	msgId := e.Message.(*linebot.ImageMessage).ID
	fmt.Println((msgId))

	input := msgdto.MsgInput{
		ReplyToken: e.ReplyToken,
		Msg:        "NFT画像",
	}

	//TODO: 受け取った画像の処理
	content, err := controller.bot.GetMessageContent(msgId).Do()
	if err != nil {
		// 画像を取得できない場合errが返る
		logrus.Errorf("ERROR: Image not found")
	}
	defer content.Content.Close()
	// 画像アップロード
	imageUrl := "https://1.bp.blogspot.com/-DgQkaAeOGgc/X9lJVi_Yv9I/AAAAAAABc34/S867MFYTC30KImIFJWIMYgg29mGgyPj0gCNcBGAsYHQ/s659/food_yamunyomu_chiken.png"

	sessions[key{uid, "image"}] = imageUrl

	controller.linebotInteractor.GetTitle(input)
	sessions[key{uid, "state"}] = "title"
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
		sessions[key{userId, "state"}] = ""
		logrus.Debug("NFTの作成に失敗しました: ", err)
		return
	}
	tx, err := controller.blockchainInteractor.GetTransaction(txhash.TxHash)
	time.Sleep(time.Second * 3) //sleep
	if err != nil {
		sessions[key{userId, "state"}] = ""
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
		sessions[key{userId, "state"}] = ""
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
	sessions[key{userId, "state"}] = ""
	return
}
