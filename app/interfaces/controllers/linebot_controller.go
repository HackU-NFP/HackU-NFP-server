package controllers

import (
	"context"
	"fmt"
	"io"
	"nfp-server/usecase"
	msgdto "nfp-server/usecase/dto"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"cloud.google.com/go/storage"
	"github.com/labstack/echo/v4"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
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
		LineUserID: uid,
		Msg:        msg,
	}

	if msg == "使い方" {
		controller.linebotInteractor.HowToUse(input)
	} else if msg == "キャンセル" {
		controller.cancel(input)
	} else if msg == "NFTを作る" {
		controller.linebotInteractor.GetImage(input)
		sessions[key{uid, "state"}] = "image"
	} else if msg == "NFTテスト" {
		logrus.Debug("NFTテスト")
		userId := e.Source.UserID
		contractId := os.Getenv("CONTRACT_ID")
		name := "HelloWorld" //TODO: stateの値にする
		meta := "HelloWorld" //TODO:stateの値にする
		controller.mint(e, userId, contractId, name, meta)
	} else {
		state := sessions[key{uid, "state"}]

		switch state {
		case "image":
			controller.linebotInteractor.GetImage(input)
		case "title":
			if !checkTitle(msg) {
				input.Msg = "タイトルは英数字のみ(スペースなし)で入力してください"
				controller.linebotInteractor.Send(input)
				break
			} else if utf8.RuneCountInString(msg) < 3 || utf8.RuneCountInString(msg) > 20 {
				input.Msg = "タイトルは3~20文字で入力してください"
				controller.linebotInteractor.Send(input)
				break
			}
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
				sessions[key{userId, "state"}] = ""
				// mint
				controller.mint(e, userId, contractId, name, meta)
			} else {
				controller.linebotInteractor.Confirm(input, sessions[key{uid, "image"}], sessions[key{uid, "title"}], sessions[key{uid, "meta"}])
			}
		default:
			// 使い方
			controller.linebotInteractor.HowToUse(input)
		}
	}
}

func checkTitle(title string) bool {
	matched, err := regexp.MatchString("^[0-9a-zA-Z]+$", title)
	if err != nil {
		return false
	}
	if len(title) == 0 || !matched {
		return false
	}
	return true //true
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
			controller.cancel(input)
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
	state := sessions[key{uid, "state"}]
	input := msgdto.MsgInput{
		ReplyToken: e.ReplyToken,
		LineUserID: uid,
		Msg:        "NFT画像",
	}

	switch state {
	case "image":
		// 受け取った画像の処理
		content, err := controller.bot.GetMessageContent(msgId).Do()
		if err != nil {
			// 画像を取得できない場合errが返る
			logrus.Errorf("ERROR: Image not found", err)
		}
		defer content.Content.Close()

		// 画像アップロード
		err = putGCS(content.Content, msgId)
		if err != nil {
			logrus.Errorf("ERROR: putGCS failure", err)
		}
		imageUrl := os.Getenv("STORAGE_BASE_URI") + msgId

		sessions[key{uid, "image"}] = imageUrl

		controller.linebotInteractor.GetTitle(input)
		sessions[key{uid, "state"}] = "title"
	case "":
		// 使い方送信
		controller.linebotInteractor.HowToUse(input)
	default:
		// エラー初期化
		controller.error(input)
	}
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

func putGCS(content io.ReadCloser, object string) error {
	bucket := os.Getenv("BUCKET")
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsJSON([]byte(os.Getenv("GCP_CREDENTIALS"))))
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := client.Bucket(bucket).Object(object).NewWriter(ctx)
	wc.ChunkSize = 0 // note retries are not supported for chunk size 0.

	if _, err = io.Copy(wc, content); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	// Data can continue to be added to the file until the writer is closed.
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	fmt.Println("putGCS")

	return nil
}

func changeObjectName(preName string, afterName string) error {
	bucket := os.Getenv("BUCKET")
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsJSON([]byte(os.Getenv("GCP_CREDENTIALS"))))
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	src := client.Bucket(bucket).Object(preName)
	dst := client.Bucket(bucket).Object(afterName)

	// Copy content.
	_, err = dst.CopierFrom(src).Run(ctx)
	if err != nil {
		return fmt.Errorf("Copy content: %v", err)
	}

	// Delete src
	if err = src.Delete(ctx); err != nil {
		return fmt.Errorf("Delete src: %v", err)
	}
	fmt.Println("changeObjectName")

	return nil
}

func deleteGCS(object string) error {
	bucket := os.Getenv("BUCKET")
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsJSON([]byte(os.Getenv("GCP_CREDENTIALS"))))
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Delete object
	if err = client.Bucket(bucket).Object(object).Delete(ctx); err != nil {
		return fmt.Errorf("Delete object: %v", err)
	}
	fmt.Println("deleteGCS")

	return nil
}

func (controller *LinebotController) cancel(input msgdto.MsgInput) {
	uid := input.LineUserID
	sessions[key{uid, "state"}] = ""
	input.Msg = "キャンセルしました"
	controller.linebotInteractor.Send(input)
	if sessions[key{uid, "image"}] != "" {
		object := strings.Replace(sessions[key{uid, "image"}], os.Getenv("STORAGE_BASE_URI"), "", -1)
		deleteGCS(object)
	}
	fmt.Println("cancel")
}

func (controller *LinebotController) error(input msgdto.MsgInput) {
	uid := input.LineUserID
	sessions[key{uid, "state"}] = ""
	input.Msg = "エラーが発生しました、お手数ですが最初からやり直してください"
	controller.linebotInteractor.Send(input)
	if sessions[key{uid, "image"}] != "" {
		object := strings.Replace(sessions[key{uid, "image"}], os.Getenv("STORAGE_BASE_URI"), "", -1)
		deleteGCS(object)
	}
	fmt.Println("error")
}

func (controller *LinebotController) mint(e *linebot.Event, userId, contractId, name, meta string) {
	loadingInput := msgdto.MsgInput{
		ReplyToken: e.ReplyToken,
		Msg:        "作成中です...",
	}
	//作成中です...メッセージ送信
	controller.linebotInteractor.Loading(loadingInput)

	txhash, err := controller.blockchainInteractor.CreateNonFungible(userId, contractId, name, meta)
	fmt.Println("txhash: ", txhash)
	time.Sleep(time.Second * 10) //sleep
	if err != nil {
		sessions[key{userId, "state"}] = ""
		logrus.Debug("NFTの作成に失敗しました: ", err)
		return
	}
	tx, err := controller.blockchainInteractor.GetTransaction(txhash.TxHash)
	fmt.Println("tx: ", tx)
	time.Sleep(time.Second * 10) //sleep
	if err != nil {
		sessions[key{userId, "state"}] = ""
		logrus.Debug("NFTの作成に失敗しました: ", err)
		return
	}
	tokenType := *&tx.Logs[0].Events[0].Attributes[1].Value

	// 画像をstorageにアップロードする.tokenTypeをファイル名にする
	preObjectName := strings.Replace(sessions[key{userId, "image"}], os.Getenv("STORAGE_BASE_URI"), "", -1)
	// object名前変更
	if err = changeObjectName(preObjectName, tokenType); err != nil {
		logrus.Debug("オブジェクトの名前変更に失敗しました。: ", err)
	}
	fmt.Println("Change object name")
	sessions[key{userId, "image"}] = os.Getenv("STORAGE_BASE_URI") + tokenType

	// ミント
	mintTx, err := controller.blockchainInteractor.MintNonFungible(userId, contractId, tokenType, name, meta)
	time.Sleep(time.Second * 5) //sleep
	if err != nil {
		logrus.Debug("mintに失敗しました: ", err)
		input := msgdto.SuccessInput{
			TokenType: tokenType,
			UserId:    userId,
			Tx:        txhash.TxHash,
			Name:      name,
			Image:     sessions[key{userId, "image"}],
		}
		//ミント成功メッセージ送信
		sessions[key{userId, "state"}] = ""
		controller.linebotInteractor.SuccessMint(input)
		return
	}

	input := msgdto.SuccessInput{
		TokenType: tokenType,
		UserId:    userId,
		Tx:        mintTx.TxHash,
		Name:      name,
		Image:     sessions[key{userId, "image"}],
	}
	//ミント成功メッセージ送信
	controller.linebotInteractor.SuccessMint(input)
	sessions[key{userId, "state"}] = ""
	return
}
