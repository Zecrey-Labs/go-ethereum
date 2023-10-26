package types

type BlockWithTxsAndReceipts struct {
	Header       *Header     `json:"header"`
	Receipts     Receipts    `json:"receipts"`
	Transactions interface{} `json:"transactions"`
}
