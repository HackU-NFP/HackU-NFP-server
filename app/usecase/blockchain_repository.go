package usecase

// BlockchainRepository ブロックチェーンプレゼンタ
type BlockchainRepository interface {
	CreateTokenType()   //新しいtokenTypeを作成
	FindTokenTypeByTx() //transactionHashからtokenTypeを取得
	Mint()              //ミント
}
