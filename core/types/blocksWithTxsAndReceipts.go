package types

type BlockWithTxsAndReceipts struct {
	*Block
	Receipts
}

type RpcBlockWithTxsAndReceipts struct {
	*Header
	Receipts     `json:"receipts"`
	Transactions interface{} `json:"transactions"`
}
