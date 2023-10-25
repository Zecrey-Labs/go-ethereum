package types

type BlockWithTxsAndReceipts struct {
	*Header
	Receipts
	Transactions
}
