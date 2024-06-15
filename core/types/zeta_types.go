package types

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type ZetaCosmosEvmTx struct {
	ChainId *big.Int
	From    common.Address

	Nonce     uint64          // nonce of sender account
	GasFeeCap *big.Int        // wei per gas
	Gas       uint64          // gas limit
	To        *common.Address `rlp:"nil"` // nil means contract creation
	Value     *big.Int        // wei amount
	Data      []byte          // contract invocation input data
}

func (tx *ZetaCosmosEvmTx) txType() byte { return ZetaCosmosEVMTxType }

func (tx *ZetaCosmosEvmTx) copy() TxData {
	cpy := &ZetaCosmosEvmTx{
		ChainId:   new(big.Int),
		Nonce:     tx.Nonce,
		GasFeeCap: new(big.Int),
		Gas:       tx.Gas,
		From:      tx.From,
		To:        nil,
		Value:     new(big.Int),
		Data:      common.CopyBytes(tx.Data),
	}
	if tx.ChainId != nil {
		cpy.ChainId.Set(tx.ChainId)
	}
	if tx.GasFeeCap != nil {
		cpy.GasFeeCap.Set(tx.GasFeeCap)
	}
	if tx.To != nil {
		tmp := *tx.To
		cpy.To = &tmp
	}
	if tx.Value != nil {
		cpy.Value.Set(tx.Value)
	}
	return cpy
}

func (tx *ZetaCosmosEvmTx) chainID() *big.Int         { return tx.ChainId }
func (tx *ZetaCosmosEvmTx) accessList() AccessList    { return nil }
func (tx *ZetaCosmosEvmTx) data() []byte              { return tx.Data }
func (tx *ZetaCosmosEvmTx) gas() uint64               { return tx.Gas }
func (tx *ZetaCosmosEvmTx) gasPrice() *big.Int        { return tx.GasFeeCap }
func (tx *ZetaCosmosEvmTx) gasTipCap() *big.Int       { return bigZero }
func (tx *ZetaCosmosEvmTx) gasFeeCap() *big.Int       { return tx.GasFeeCap }
func (tx *ZetaCosmosEvmTx) value() *big.Int           { return tx.Value }
func (tx *ZetaCosmosEvmTx) nonce() uint64             { return tx.Nonce }
func (tx *ZetaCosmosEvmTx) to() *common.Address       { return tx.To }
func (tx *ZetaCosmosEvmTx) blobGas() uint64           { return 0 }
func (tx *ZetaCosmosEvmTx) blobGasFeeCap() *big.Int   { return nil }
func (tx *ZetaCosmosEvmTx) blobHashes() []common.Hash { return nil }

func (tx *ZetaCosmosEvmTx) rawSignatureValues() (v, r, s *big.Int) {
	return bigZero, bigZero, bigZero
}

func (tx *ZetaCosmosEvmTx) setSignatureValues(chainID, v, r, s *big.Int) {

}

func (tx *ZetaCosmosEvmTx) effectiveGasPrice(dst *big.Int, baseFee *big.Int) *big.Int {
	if baseFee == nil {
		return dst.Set(tx.GasFeeCap)
	}
	return dst.Set(baseFee)
}
