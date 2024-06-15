package types

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

const (
	ZetaCosmosEVMTxType = 0x58
)

// ZetaCosmosEVMTx is the transaction data of the original Ethereum transactions.
type ZetaCosmosEVMTx struct {
	BlockHash common.Hash
	TxHash    common.Hash
	Nonce     uint64          // nonce of sender account
	GasPrice  *big.Int        // wei per gas
	Gas       uint64          // gas limit
	To        *common.Address `rlp:"nil"` // nil means contract creation
	Value     *big.Int        // wei amount
	Data      []byte          // contract invocation input data
	V, R, S   *big.Int        // signature values
}

// NewTransaction creates an unsigned legacy transaction.
// Deprecated: use NewTx instead.
func NewZetaTransaction(nonce uint64, to common.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte) *Transaction {
	return NewTx(&ZetaCosmosEVMTx{
		Nonce:    nonce,
		To:       &to,
		Value:    amount,
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     data,
	})
}

// NewContractCreation creates an unsigned legacy transaction.
// Deprecated: use NewTx instead.
func NewZetaContractCreation(nonce uint64, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte) *Transaction {
	return NewTx(&ZetaCosmosEVMTx{
		Nonce:    nonce,
		Value:    amount,
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     data,
	})
}

// copy creates a deep copy of the transaction data and initializes all fields.
func (tx *ZetaCosmosEVMTx) copy() TxData {
	cpy := &ZetaCosmosEVMTx{
		Nonce: tx.Nonce,
		To:    copyAddressPtr(tx.To),
		Data:  common.CopyBytes(tx.Data),
		Gas:   tx.Gas,
		// These are initialized below.
		Value:    new(big.Int),
		GasPrice: new(big.Int),
		V:        new(big.Int),
		R:        new(big.Int),
		S:        new(big.Int),
	}
	if tx.Value != nil {
		cpy.Value.Set(tx.Value)
	}
	if tx.GasPrice != nil {
		cpy.GasPrice.Set(tx.GasPrice)
	}
	if tx.V != nil {
		cpy.V.Set(tx.V)
	}
	if tx.R != nil {
		cpy.R.Set(tx.R)
	}
	if tx.S != nil {
		cpy.S.Set(tx.S)
	}
	return cpy
}

// accessors for innerTx.
func (tx *ZetaCosmosEVMTx) txType() byte              { return ZetaCosmosEVMTxType }
func (tx *ZetaCosmosEVMTx) chainID() *big.Int         { return deriveChainId(tx.V) }
func (tx *ZetaCosmosEVMTx) accessList() AccessList    { return nil }
func (tx *ZetaCosmosEVMTx) data() []byte              { return tx.Data }
func (tx *ZetaCosmosEVMTx) gas() uint64               { return tx.Gas }
func (tx *ZetaCosmosEVMTx) gasPrice() *big.Int        { return tx.GasPrice }
func (tx *ZetaCosmosEVMTx) gasTipCap() *big.Int       { return tx.GasPrice }
func (tx *ZetaCosmosEVMTx) gasFeeCap() *big.Int       { return tx.GasPrice }
func (tx *ZetaCosmosEVMTx) value() *big.Int           { return tx.Value }
func (tx *ZetaCosmosEVMTx) nonce() uint64             { return tx.Nonce }
func (tx *ZetaCosmosEVMTx) to() *common.Address       { return tx.To }
func (tx *ZetaCosmosEVMTx) blobGas() uint64           { return 0 }
func (tx *ZetaCosmosEVMTx) blobGasFeeCap() *big.Int   { return nil }
func (tx *ZetaCosmosEVMTx) blobHashes() []common.Hash { return nil }

func (tx *ZetaCosmosEVMTx) effectiveGasPrice(dst *big.Int, baseFee *big.Int) *big.Int {
	return dst.Set(tx.GasPrice)
}

func (tx *ZetaCosmosEVMTx) rawSignatureValues() (v, r, s *big.Int) {
	return tx.V, tx.R, tx.S
}

func (tx *ZetaCosmosEVMTx) setSignatureValues(chainID, v, r, s *big.Int) {
	tx.V, tx.R, tx.S = v, r, s
}
