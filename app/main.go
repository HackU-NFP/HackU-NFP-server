package main

import (
	"fmt"
	"nfp-server/infrastructure"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")

	//もし err がnilではないなら、"読み込み出来ませんでした"が出力されます。
	if err != nil {
		fmt.Printf("環境変数を読み込み出来ませんでした: %v", err)
	}
	infrastructure.Init()
}
