package usecase

type UserInfo struct {
	userID        string `json:"userId"`
	WalletAddress string `json:"walletAddress"`
}

type NonFungible struct {
	TokenType   string `json:"tokenType"`
	Name        string `json:"name"`
	Meta        string `json:"meta"`
	CreatedAt   int    `json:"createdAt"`
	TotalSupply string `json:"totalSupply"`
	TotalMint   string `json:"totalMint"`
	TotalBurn   string `json:"totalBurn"`
}

type NonFungibleInfo struct {
	TokenType   string `json:"tokenType"`
	Name        string `json:"name"`
	Meta        string `json:"meta"`
	CreatedAt   int    `json:"createdAt"`
	TotalSupply string `json:"totalSupply"`
	TotalMint   string `json:"totalMint"`
	TotalBurn   string `json:"totalBurn"`
	Token       []Token
}
type Token struct {
	TokenIndex string `json:"tokenIndex"`
	Name       string `json:"name"`
	Meta       string `json:"meta"`
	CreatedAt  int    `json:"createdAt"`
	BurnedAt   *int   `json:"burnedAt"`
}

type Transaction struct {
	Height    uint64 `json:"height"`
	TxHash    string `json:"txhash"`
	Index     int    `json:"index"`
	Code      int    `json:"code"`
	RawLog    string `json:"raw_log"`
	Logs      []Log  `json:"logs"`
	GasWanted uint64 `json:"gasWanted"`
	GasUsed   uint64 `json:"gasUsed"`
	Tx        Tx     `json:"tx"`
	Timestamp int    `json:"timestamp"`
}

type Log struct {
	MsgIndex int     `json:"msg_index"`
	Success  bool    `json:"success"`
	Log      string  `json:"log"`
	Events   []Event `json:"events"`
}

type Event struct {
	Type       string      `json:"type"`
	Attributes []Attribute `json:"attributes"`
}

type Attribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Tx struct {
	Type  string  `json:"type"`
	Value TxValue `json:"value"`
}

type TxValue struct {
	Message    []Message   `json:"msg"`
	Fee        Fee         `json:"fee"`
	Signatures []Signature `json:"signatures"`
	Memo       string      `json:"memo"`
}

type Message struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type TransferBaseCoinMsg struct {
	FromAddress string `json:"fromAddress"`
	ToAddress   string `json:"toAddress"`
	Amount      Amount `json:"amount"`
}

type TransferServiceTokenMsg struct {
	From       string `json:"from"`
	To         string `json:"to"`
	Amount     uint64 `json:"amount"`
	ContractID string `json:"contractId"`
}

type MintNonFungibleMsg struct {
	From       string `json:"from"`
	To         string `json:"to"`
	ContractID string `json:"contractId"`
	Meta       string `json:"meta"`
	Name       string `json:"name"`
	TokenType  string `json:"tokenType"`
}

type Signature struct {
	PubKey    PubKey `json:"pubKey"`
	Signature string `json:"signature"`
}

type PubKey struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Fee struct {
	Amount []Amount `json:"amount"`
	Gas    int      `json:"gas"`
}

type Amount struct {
	Amount int    `json:"amount"`
	Denom  string `json:"denom"`
}

type NonFungibleTxHistory struct {
	PaymentTransaction *Transaction `json:"paymentTransaction"`
	MintTransaction    *Transaction `json:"mintTransaction"`
	PointTransaction   *Transaction `json:"pointTransaction"`
}

type TransactionAccepted struct {
	TxHash string `json:"txHash"`
}

type TransferRequestResult struct {
	RequestSessionToken string `json:"requestSessionToken"`
	RedirectURI         string `json:"redirectUri"`
}
