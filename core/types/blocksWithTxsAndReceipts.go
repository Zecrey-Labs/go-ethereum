package types

type BlockWithTxsAndReceipts struct {
	*Header
	Receipts     Receipts     `json:"receipts"`
	Transactions Transactions `json:"transactions"`
}
