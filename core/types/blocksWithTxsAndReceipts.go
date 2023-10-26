package types

type BlockWithTxsAndReceipts struct {
	*Header
	Receipts     Receipts    `json:"receipts"`
	Transactions interface{} `json:"transactions"`
}
