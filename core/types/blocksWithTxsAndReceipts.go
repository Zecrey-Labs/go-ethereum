package types

type BlockWithTxsAndReceipts struct {
	*Header               `json:"header"`
	Receipts              `json:"receipts"`
	Transactions          `json:"transactions"`
	FormattedTransactions map[string]interface{} `json:"formattedTransactions"`
}
