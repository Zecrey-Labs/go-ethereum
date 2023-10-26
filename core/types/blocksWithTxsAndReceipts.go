package types

type BlockWithTxsAndReceipts struct {
	*Block
	Receipts
}

type RpcBlockWithTxsAndReceipts struct {
	*Header      `json:"header"`
	Receipts     `json:"receipts"`
	Transactions interface{} `json:"transactions"`
}
