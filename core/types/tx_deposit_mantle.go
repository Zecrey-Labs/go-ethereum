package types

import (
	"bytes"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
)

type DepositTxMantle struct {
	// SourceHash uniquely identifies the source of the deposit
	SourceHash common.Hash
	// From is exposed through the types.Signer, not through TxData
	From common.Address
	// nil means contract creation
	To *common.Address `rlp:"nil"`
	// Mint is minted on L2, locked on L1, nil if no minting.
	Mint *big.Int `rlp:"nil"`
	// Value is transferred from L2 balance, executed after Mint (if any)
	Value *big.Int
	// gas limit
	Gas uint64
	// Field indicating if this transaction is exempt from the L2 gas limit.
	IsSystemTransaction bool
	// EthValue means L2 BVM_ETH mint tag, nil means that there is no need to mint BVM_ETH.
	EthValue *big.Int `rlp:"nil"`
	// Normal Tx data
	Data []byte
	// EthTxValue means L2 BVM_ETH tx tag, nil means that there is no need to transfer BVM_ETH to msg.To.
	EthTxValue *big.Int `rlp:"optional"`
}

func (tx *DepositTxMantle) skipAccountChecks() bool {
	//TODO implement me
	panic("implement me")
}

func (tx *DepositTxMantle) blobGas() uint64 {
	//TODO implement me
	panic("implement me")
}

func (tx *DepositTxMantle) blobGasFeeCap() *big.Int {
	//TODO implement me
	panic("implement me")
}

func (tx *DepositTxMantle) blobHashes() []common.Hash {
	//TODO implement me
	panic("implement me")
}

// copy creates a deep copy of the transaction data and initializes all fields.
func (tx *DepositTxMantle) copy() TxData {
	cpy := &DepositTxMantle{
		SourceHash:          tx.SourceHash,
		From:                tx.From,
		To:                  copyAddressPtr(tx.To),
		Mint:                nil,
		Value:               new(big.Int),
		Gas:                 tx.Gas,
		IsSystemTransaction: tx.IsSystemTransaction,
		Data:                common.CopyBytes(tx.Data),
	}
	if tx.Mint != nil {
		cpy.Mint = new(big.Int).Set(tx.Mint)
	}
	if tx.Value != nil {
		cpy.Value.Set(tx.Value)
	}
	return cpy
}

// accessors for innerTx.
func (tx *DepositTxMantle) txType() byte           { return DepositTxType }
func (tx *DepositTxMantle) chainID() *big.Int      { return common.Big0 }
func (tx *DepositTxMantle) accessList() AccessList { return nil }
func (tx *DepositTxMantle) data() []byte           { return tx.Data }
func (tx *DepositTxMantle) gas() uint64            { return tx.Gas }
func (tx *DepositTxMantle) gasFeeCap() *big.Int    { return new(big.Int) }
func (tx *DepositTxMantle) gasTipCap() *big.Int    { return new(big.Int) }
func (tx *DepositTxMantle) gasPrice() *big.Int     { return new(big.Int) }
func (tx *DepositTxMantle) value() *big.Int        { return tx.Value }
func (tx *DepositTxMantle) nonce() uint64          { return 0 }
func (tx *DepositTxMantle) to() *common.Address    { return tx.To }
func (tx *DepositTxMantle) isSystemTx() bool       { return tx.IsSystemTransaction }

func (tx *DepositTxMantle) effectiveGasPrice(dst *big.Int, baseFee *big.Int) *big.Int {
	return dst.Set(new(big.Int))
}

func (tx *DepositTxMantle) effectiveNonce() *uint64 { return nil }

func (tx *DepositTxMantle) rawSignatureValues() (v, r, s *big.Int) {
	return common.Big0, common.Big0, common.Big0
}

func (tx *DepositTxMantle) setSignatureValues(chainID, v, r, s *big.Int) {
	// this is a noop for deposit transactions
}

func (tx *DepositTxMantle) encode(b *bytes.Buffer) error {
	return rlp.Encode(b, tx)
}

func (tx *DepositTxMantle) decode(input []byte) error {
	return rlp.DecodeBytes(input, tx)
}
