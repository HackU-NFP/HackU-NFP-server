package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"nfp-server/api"
	"os"
	"regexp"
)

// BlockchainInteractor ブロックチェーン処理
type BlockchainInteractor struct{}

// NewFavoriteInteractor コンストラクタ
func NewBlockchainInteractor() *BlockchainInteractor {
	return &BlockchainInteractor{}
}

var (
	errInvalidParam = errors.New("invalid URL params")
)

func checkUrlParam(params ...string) bool {
	for _, param := range params {
		matched, err := regexp.MatchString("^[a-zA-Z0-9-_=]*$", param)
		if err != nil {
			return false
		}
		if len(param) == 0 || !matched {
			return false
		}
	}
	return false //true
}

// //
// func (interactor *LineBotInteractor) createNft() {

// }

// //
// func (interactor *LineBotInteractor) getTransactionInfo() {

// }
func (interactor *BlockchainInteractor) CreateNonFungible(userID, contractID, name, meta string) (*TransactionAccepted, error) {
	if checkUrlParam(contractID) {
		return nil, errInvalidParam
	}
	path := fmt.Sprintf("/v1/item-tokens/%s/non-fungibles/", contractID)

	// marshaledMeta, err := json.Marshal(meta)
	// if err != nil {
	// 	return nil, err
	// }

	params := map[string]interface{}{
		"name":         name,
		"meta":         meta,
		"ownerAddress": os.Getenv("WALLET_ADRESS"),
		"ownerSecret":  os.Getenv("WALLET_SECRET"),
	}

	apiResult, err := api.CallAPI(path, "POST", nil, params)
	if err != nil {
		return nil, err
	}

	txAccepted := &TransactionAccepted{}

	if err := json.Unmarshal(apiResult, txAccepted); err != nil {
		return nil, err
	}

	return txAccepted, nil
}

func (interactor *BlockchainInteractor) GetTransaction(txHash string) (*Transaction, error) {
	// if checkUrlParam(txHash) {
	// 	return nil, errInvalidParam
	// }
	path := fmt.Sprintf("/v1/transactions/%s", txHash)

	apiResult, err := api.CallAPI(path, "GET", nil, nil)

	if err != nil {
		return nil, err
	}

	tx := &Transaction{}
	if err := json.Unmarshal(apiResult, tx); err != nil {
		return nil, err
	}

	return tx, nil
}

func (interactor *BlockchainInteractor) MintNonFungible(userID, contractID, tokenType string, name, meta string) (*TransactionAccepted, error) {
	if checkUrlParam(contractID, tokenType) {
		return nil, errInvalidParam
	}
	path := fmt.Sprintf("/v1/item-tokens/%s/non-fungibles/%s/mint", contractID, tokenType)

	// marshaledMeta, err := json.Marshal(meta)
	// if err != nil {
	// 	return nil, err
	// }

	params := map[string]interface{}{
		"toUserId":     userID,
		"name":         name,
		"meta":         meta,
		"ownerAddress": os.Getenv("WALLET_ADRESS"),
		"ownerSecret":  os.Getenv("WALLET_SECRET"),
	}

	apiResult, err := api.CallAPI(path, "POST", nil, params)
	if err != nil {
		return nil, err
	}

	txAccepted := &TransactionAccepted{}

	if err := json.Unmarshal(apiResult, txAccepted); err != nil {
		return nil, err
	}

	return txAccepted, nil
}

func (interactor *BlockchainInteractor) GetNonFungibles(contractID, orderBy, limit, page string) ([]*NonFungible, error) {
	// if checkUrlParam(contractID) {
	// 	return nil, errInvalidParam
	// }
	query := map[string]string{
		"orderBy": orderBy,
		"limit":   limit,
		"page":    page,
	}

	path := fmt.Sprintf("/v1/item-tokens/%s/non-fungibles", contractID)

	apiResult, err := api.CallAPI(path, "GET", query, nil)

	if err != nil {
		return nil, err
	}

	nonFungible := make([]*NonFungible, 0)

	if err := json.Unmarshal(apiResult, &nonFungible); err != nil {
		return nil, err
	}

	return nonFungible, nil
}

func (interactor *BlockchainInteractor) GetNonFungibleInfo(contractID, tokenType string) (*NonFungibleInfo, error) {
	if checkUrlParam(contractID, tokenType) {
		return nil, errInvalidParam
	}
	path := fmt.Sprintf("/v1/item-tokens/%s/non-fungibles/%s", contractID, tokenType)

	apiResult, err := api.CallAPI(path, "GET", nil, nil)

	if err != nil {
		return nil, err
	}

	nonFungibleInfo := &NonFungibleInfo{}

	if err := json.Unmarshal(apiResult, nonFungibleInfo); err != nil {
		return nil, err
	}

	return nonFungibleInfo, nil
}

func (interactor *BlockchainInteractor) GetUserNonFungibles(contractID, userId, orderBy, limit, page string) ([]*NonFungible, error) {
	// if checkUrlParam(contractID) {
	// 	return nil, errInvalidParam
	// }
	query := map[string]string{
		"orderBy": orderBy,
		"limit":   limit,
		"page":    page,
	}

	path := fmt.Sprintf("/v1/users/%s/item-tokens/%s/non-fungibles", userId, contractID)

	apiResult, err := api.CallAPI(path, "GET", query, nil)

	if err != nil {
		return nil, err
	}

	nonFungible := make([]*NonFungible, 0)

	if err := json.Unmarshal(apiResult, &nonFungible); err != nil {
		return nil, err
	}

	return nonFungible, nil
}
