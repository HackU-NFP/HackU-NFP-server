package usecase

// BlockchainUseCase
type IBlockchainUseCase interface {
	CreateNonFungible(userID, contractID, name, meta string) (*TransactionAccepted, error)                 //発行
	MintNonFungible(userID, contractID, tokenType string, name, meta string) (*TransactionAccepted, error) //鋳造
	GetTransaction(txHash string) (*Transaction, error)
	GetNonFungibles(contractID, orderBy, limit, page string) ([]*NonFungible, error)
	GetNonFungibleInfo(contractID, tokenType string) (*NonFungibleInfo, error)
}
