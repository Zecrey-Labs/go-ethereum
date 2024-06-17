// Copyright 2021 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/rlp"
	"io"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/holiman/uint256"
)

// txJSON is the JSON representation of transactions.
type txJSON struct {
	Type hexutil.Uint64 `json:"type"`

	ChainID              *hexutil.Big    `json:"chainId,omitempty"`
	Nonce                *hexutil.Uint64 `json:"nonce"`
	To                   *common.Address `json:"to"`
	Gas                  *hexutil.Uint64 `json:"gas"`
	GasPrice             *hexutil.Big    `json:"gasPrice"`
	MaxPriorityFeePerGas *hexutil.Big    `json:"maxPriorityFeePerGas"`
	MaxFeePerGas         *hexutil.Big    `json:"maxFeePerGas"`
	MaxFeePerDataGas     *hexutil.Big    `json:"maxFeePerDataGas,omitempty"`

	// eip 4844 upgrade
	MaxFeePerBlobGas    *hexutil.Big   `json:"maxFeePerBlobGas,omitempty"`
	Value               *hexutil.Big   `json:"value"`
	Input               *hexutil.Bytes `json:"input"`
	AccessList          *AccessList    `json:"accessList,omitempty"`
	BlobVersionedHashes []common.Hash  `json:"blobVersionedHashes,omitempty"`
	V                   *hexutil.Big   `json:"v"`
	R                   *hexutil.Big   `json:"r"`
	S                   *hexutil.Big   `json:"s"`

	// Deposit transaction fields
	SourceHash *common.Hash    `json:"sourceHash,omitempty"`
	From       *common.Address `json:"from,omitempty"`
	Mint       *hexutil.Big    `json:"mint,omitempty"`
	EthValue   *hexutil.Big    `json:"ethValue,omitempty"`
	EthTxValue *hexutil.Big    `json:"ethTxValue,omitempty"`
	IsSystemTx *bool           `json:"isSystemTx,omitempty"`

	// Arbitrum fields:
	RequestId           *common.Hash    `json:"requestId,omitempty"`           // Contract SubmitRetryable Deposit
	TicketId            *common.Hash    `json:"ticketId,omitempty"`            // Retry
	MaxRefund           *hexutil.Big    `json:"maxRefund,omitempty"`           // Retry
	SubmissionFeeRefund *hexutil.Big    `json:"submissionFeeRefund,omitempty"` // Retry
	RefundTo            *common.Address `json:"refundTo,omitempty"`            // SubmitRetryable Retry
	L1BaseFee           *hexutil.Big    `json:"l1BaseFee,omitempty"`           // SubmitRetryable
	DepositValue        *hexutil.Big    `json:"depositValue,omitempty"`        // SubmitRetryable
	RetryTo             *common.Address `json:"retryTo,omitempty"`             // SubmitRetryable
	RetryValue          *hexutil.Big    `json:"retryValue,omitempty"`          // SubmitRetryable
	RetryData           *hexutil.Bytes  `json:"retryData,omitempty"`           // SubmitRetryable
	Beneficiary         *common.Address `json:"beneficiary,omitempty"`         // SubmitRetryable
	MaxSubmissionFee    *hexutil.Big    `json:"maxSubmissionFee,omitempty"`    // SubmitRetryable
	EffectiveGasPrice   *hexutil.Uint64 `json:"effectiveGasPrice,omitempty"`   // ArbLegacy
	L1BlockNumber       *hexutil.Uint64 `json:"l1BlockNumber,omitempty"`       // ArbLegacy

	// Only used for encoding:
	Hash      common.Hash `json:"hash"`
	BlockHash common.Hash `json:"blockHash"`

	// L1 message transaction fields: [Scroll]
	Sender     *common.Address `json:"sender,omitempty"`
	QueueIndex *hexutil.Uint64 `json:"queueIndex,omitempty"`
}

// MarshalJSON marshals as JSON with a hash.
func (tx *Transaction) MarshalJSON() ([]byte, error) {
	var enc txJSON
	// These are set for all tx types.
	enc.Hash = tx.Hash()
	enc.Type = hexutil.Uint64(tx.Type())

	// Other fields are set conditionally depending on tx type.
	switch itx := tx.inner.(type) {
	case *LegacyTx:
		enc.Nonce = (*hexutil.Uint64)(&itx.Nonce)
		enc.To = tx.To()
		enc.Gas = (*hexutil.Uint64)(&itx.Gas)
		enc.GasPrice = (*hexutil.Big)(itx.GasPrice)
		enc.Value = (*hexutil.Big)(itx.Value)
		enc.Input = (*hexutil.Bytes)(&itx.Data)
		enc.V = (*hexutil.Big)(itx.V)
		enc.R = (*hexutil.Big)(itx.R)
		enc.S = (*hexutil.Big)(itx.S)
	case *ZetaCosmosEVMTx:
		enc.Nonce = (*hexutil.Uint64)(&itx.Nonce)
		enc.To = tx.To()
		enc.Gas = (*hexutil.Uint64)(&itx.Gas)
		enc.GasPrice = (*hexutil.Big)(itx.GasPrice)
		enc.Value = (*hexutil.Big)(itx.Value)
		enc.Input = (*hexutil.Bytes)(&itx.Data)
		enc.V = (*hexutil.Big)(itx.V)
		enc.R = (*hexutil.Big)(itx.R)
		enc.S = (*hexutil.Big)(itx.S)
	case *AccessListTx:
		enc.ChainID = (*hexutil.Big)(itx.ChainID)
		enc.Nonce = (*hexutil.Uint64)(&itx.Nonce)
		enc.To = tx.To()
		enc.Gas = (*hexutil.Uint64)(&itx.Gas)
		enc.GasPrice = (*hexutil.Big)(itx.GasPrice)
		enc.Value = (*hexutil.Big)(itx.Value)
		enc.Input = (*hexutil.Bytes)(&itx.Data)
		enc.AccessList = &itx.AccessList
		enc.V = (*hexutil.Big)(itx.V)
		enc.R = (*hexutil.Big)(itx.R)
		enc.S = (*hexutil.Big)(itx.S)

	case *DynamicFeeTx:
		enc.ChainID = (*hexutil.Big)(itx.ChainID)
		enc.Nonce = (*hexutil.Uint64)(&itx.Nonce)
		enc.To = tx.To()
		enc.Gas = (*hexutil.Uint64)(&itx.Gas)
		enc.MaxFeePerGas = (*hexutil.Big)(itx.GasFeeCap)
		enc.MaxPriorityFeePerGas = (*hexutil.Big)(itx.GasTipCap)
		enc.Value = (*hexutil.Big)(itx.Value)
		enc.Input = (*hexutil.Bytes)(&itx.Data)
		enc.AccessList = &itx.AccessList
		enc.V = (*hexutil.Big)(itx.V)
		enc.R = (*hexutil.Big)(itx.R)
		enc.S = (*hexutil.Big)(itx.S)

	case *BlobTx:
		enc.ChainID = (*hexutil.Big)(itx.ChainID.ToBig())
		enc.Nonce = (*hexutil.Uint64)(&itx.Nonce)
		enc.Gas = (*hexutil.Uint64)(&itx.Gas)
		enc.MaxFeePerGas = (*hexutil.Big)(itx.GasFeeCap.ToBig())
		enc.MaxPriorityFeePerGas = (*hexutil.Big)(itx.GasTipCap.ToBig())
		enc.MaxFeePerDataGas = (*hexutil.Big)(itx.BlobFeeCap.ToBig())
		enc.Value = (*hexutil.Big)(itx.Value.ToBig())
		enc.Input = (*hexutil.Bytes)(&itx.Data)
		enc.AccessList = &itx.AccessList
		enc.BlobVersionedHashes = itx.BlobHashes
		enc.To = tx.To()
		enc.V = (*hexutil.Big)(itx.V.ToBig())
		enc.R = (*hexutil.Big)(itx.R.ToBig())
		enc.S = (*hexutil.Big)(itx.S.ToBig())
	}
	return json.Marshal(&enc)
}

// UnmarshalJSON unmarshals from JSON.
func (tx *Transaction) UnmarshalJSON(input []byte) error {
	var dec txJSON
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}

	// Decode / verify fields according to transaction type.
	var inner TxData
	switch dec.Type {
	case LegacyTxType:
		var itx LegacyTx
		inner = &itx
		if dec.Nonce == nil {
			return errors.New("missing required field 'nonce' in transaction")
		}
		itx.Nonce = uint64(*dec.Nonce)
		if dec.To != nil {
			itx.To = dec.To
		}
		if dec.Gas == nil {
			return errors.New("missing required field 'gas' in transaction")
		}
		itx.Gas = uint64(*dec.Gas)
		if dec.GasPrice == nil {
			return errors.New("missing required field 'gasPrice' in transaction")
		}
		itx.GasPrice = (*big.Int)(dec.GasPrice)
		if dec.Value == nil {
			return errors.New("missing required field 'value' in transaction")
		}
		itx.Value = (*big.Int)(dec.Value)
		if dec.Input == nil {
			return errors.New("missing required field 'input' in transaction")
		}
		itx.Data = *dec.Input
		if dec.V == nil {
			return errors.New("missing required field 'v' in transaction")
		}
		itx.V = (*big.Int)(dec.V)
		if dec.R == nil {
			return errors.New("missing required field 'r' in transaction")
		}
		itx.R = (*big.Int)(dec.R)
		if dec.S == nil {
			return errors.New("missing required field 's' in transaction")
		}
		itx.S = (*big.Int)(dec.S)
		withSignature := itx.V.Sign() != 0 || itx.R.Sign() != 0 || itx.S.Sign() != 0
		if withSignature {
			if err := sanityCheckSignature(itx.V, itx.R, itx.S, true); err != nil {
				return err
			}
		}
	case ZetaCosmosEVMTxType:
		var itx ZetaCosmosEVMTx
		itx.TxHash = dec.Hash
		itx.BlockHash = dec.BlockHash
		inner = &itx
		if dec.Nonce == nil {
			return errors.New("missing required field 'nonce' in transaction")
		}
		itx.Nonce = uint64(*dec.Nonce)
		if dec.To != nil {
			itx.To = dec.To
		}
		if dec.Gas == nil {
			return errors.New("missing required field 'gas' in transaction")
		}
		itx.Gas = uint64(*dec.Gas)
		if dec.GasPrice == nil {
			return errors.New("missing required field 'gasPrice' in transaction")
		}
		itx.GasPrice = (*big.Int)(dec.GasPrice)
		if dec.Value == nil {
			return errors.New("missing required field 'value' in transaction")
		}
		itx.Value = (*big.Int)(dec.Value)
		if dec.Input == nil {
			return errors.New("missing required field 'input' in transaction")
		}
		itx.Data = *dec.Input
		itx.V = big.NewInt(0)
		itx.R = big.NewInt(0)
		itx.S = big.NewInt(0)
	case AccessListTxType:
		var itx AccessListTx
		inner = &itx
		if dec.ChainID == nil {
			return errors.New("missing required field 'chainId' in transaction")
		}
		itx.ChainID = (*big.Int)(dec.ChainID)
		if dec.Nonce == nil {
			return errors.New("missing required field 'nonce' in transaction")
		}
		itx.Nonce = uint64(*dec.Nonce)
		if dec.To != nil {
			itx.To = dec.To
		}
		if dec.Gas == nil {
			return errors.New("missing required field 'gas' in transaction")
		}
		itx.Gas = uint64(*dec.Gas)
		if dec.GasPrice == nil {
			return errors.New("missing required field 'gasPrice' in transaction")
		}
		itx.GasPrice = (*big.Int)(dec.GasPrice)
		if dec.Value == nil {
			return errors.New("missing required field 'value' in transaction")
		}
		itx.Value = (*big.Int)(dec.Value)
		if dec.Input == nil {
			return errors.New("missing required field 'input' in transaction")
		}
		itx.Data = *dec.Input
		if dec.V == nil {
			return errors.New("missing required field 'v' in transaction")
		}
		if dec.AccessList != nil {
			itx.AccessList = *dec.AccessList
		}
		itx.V = (*big.Int)(dec.V)
		if dec.R == nil {
			return errors.New("missing required field 'r' in transaction")
		}
		itx.R = (*big.Int)(dec.R)
		if dec.S == nil {
			return errors.New("missing required field 's' in transaction")
		}
		itx.S = (*big.Int)(dec.S)
		withSignature := itx.V.Sign() != 0 || itx.R.Sign() != 0 || itx.S.Sign() != 0
		if withSignature {
			if err := sanityCheckSignature(itx.V, itx.R, itx.S, false); err != nil {
				return err
			}
		}

	case DynamicFeeTxType:
		var itx DynamicFeeTx
		inner = &itx
		if dec.ChainID == nil {
			return errors.New("missing required field 'chainId' in transaction")
		}
		itx.ChainID = (*big.Int)(dec.ChainID)
		if dec.Nonce == nil {
			return errors.New("missing required field 'nonce' in transaction")
		}
		itx.Nonce = uint64(*dec.Nonce)
		if dec.To != nil {
			itx.To = dec.To
		}
		if dec.Gas == nil {
			return errors.New("missing required field 'gas' for txdata")
		}
		itx.Gas = uint64(*dec.Gas)
		if dec.MaxPriorityFeePerGas == nil {
			return errors.New("missing required field 'maxPriorityFeePerGas' for txdata")
		}
		itx.GasTipCap = (*big.Int)(dec.MaxPriorityFeePerGas)
		if dec.MaxFeePerGas == nil {
			return errors.New("missing required field 'maxFeePerGas' for txdata")
		}
		itx.GasFeeCap = (*big.Int)(dec.MaxFeePerGas)
		if dec.Value == nil {
			return errors.New("missing required field 'value' in transaction")
		}
		itx.Value = (*big.Int)(dec.Value)
		if dec.Input == nil {
			return errors.New("missing required field 'input' in transaction")
		}
		itx.Data = *dec.Input
		if dec.V == nil {
			return errors.New("missing required field 'v' in transaction")
		}
		if dec.AccessList != nil {
			itx.AccessList = *dec.AccessList
		}
		itx.V = (*big.Int)(dec.V)
		if dec.R == nil {
			return errors.New("missing required field 'r' in transaction")
		}
		itx.R = (*big.Int)(dec.R)
		if dec.S == nil {
			return errors.New("missing required field 's' in transaction")
		}
		itx.S = (*big.Int)(dec.S)
		withSignature := itx.V.Sign() != 0 || itx.R.Sign() != 0 || itx.S.Sign() != 0
		if withSignature {
			if err := sanityCheckSignature(itx.V, itx.R, itx.S, false); err != nil {
				return err
			}
		}

	case BlobTxType:
		var itx BlobTx
		inner = &itx
		if dec.ChainID == nil {
			return errors.New("missing required field 'chainId' in transaction")
		}
		itx.ChainID = uint256.MustFromBig((*big.Int)(dec.ChainID))
		if dec.Nonce == nil {
			return errors.New("missing required field 'nonce' in transaction")
		}
		itx.Nonce = uint64(*dec.Nonce)
		if dec.To != nil {
			itx.To = dec.To
		}
		if dec.Gas == nil {
			return errors.New("missing required field 'gas' for txdata")
		}
		itx.Gas = uint64(*dec.Gas)
		if dec.MaxPriorityFeePerGas == nil {
			return errors.New("missing required field 'maxPriorityFeePerGas' for txdata")
		}
		itx.GasTipCap = uint256.MustFromBig((*big.Int)(dec.MaxPriorityFeePerGas))
		if dec.MaxFeePerGas == nil {
			return errors.New("missing required field 'maxFeePerGas' for txdata")
		}
		itx.GasFeeCap = uint256.MustFromBig((*big.Int)(dec.MaxFeePerGas))
		// eip 4844 upgrade

		//if dec.MaxFeePerDataGas == nil {
		//	return errors.New("missing required field 'maxFeePerDataGas' for txdata")
		//}
		//itx.BlobFeeCap = uint256.MustFromBig((*big.Int)(dec.MaxFeePerDataGas))

		if dec.MaxFeePerBlobGas == nil {
			return errors.New("missing required field 'maxFeePerBlobGas' for txdata")
		}
		itx.BlobFeeCap = uint256.MustFromBig((*big.Int)(dec.MaxFeePerBlobGas))
		if dec.Value == nil {
			return errors.New("missing required field 'value' in transaction")
		}
		itx.Value = uint256.MustFromBig((*big.Int)(dec.Value))
		if dec.Input == nil {
			return errors.New("missing required field 'input' in transaction")
		}
		itx.Data = *dec.Input
		if dec.V == nil {
			return errors.New("missing required field 'v' in transaction")
		}
		if dec.AccessList != nil {
			itx.AccessList = *dec.AccessList
		}
		if dec.BlobVersionedHashes == nil {
			return errors.New("missing required field 'blobVersionedHashes' in transaction")
		}
		itx.BlobHashes = dec.BlobVersionedHashes
		itx.V = uint256.MustFromBig((*big.Int)(dec.V))
		if dec.R == nil {
			return errors.New("missing required field 'r' in transaction")
		}
		itx.R = uint256.MustFromBig((*big.Int)(dec.R))
		if dec.S == nil {
			return errors.New("missing required field 's' in transaction")
		}
		itx.S = uint256.MustFromBig((*big.Int)(dec.S))
		withSignature := itx.V.Sign() != 0 || itx.R.Sign() != 0 || itx.S.Sign() != 0
		if withSignature {
			if err := sanityCheckSignature(itx.V.ToBig(), itx.R.ToBig(), itx.S.ToBig(), false); err != nil {
				return err
			}
		}
	case DepositTxType:
		if dec.AccessList != nil || dec.MaxFeePerGas != nil ||
			dec.MaxPriorityFeePerGas != nil {
			return errors.New("unexpected field(s) in deposit transaction")
		}
		if dec.GasPrice != nil && dec.GasPrice.ToInt().Cmp(common.Big0) != 0 {
			return errors.New("deposit transaction GasPrice must be 0")
		}
		if (dec.V != nil && dec.V.ToInt().Cmp(common.Big0) != 0) ||
			(dec.R != nil && dec.R.ToInt().Cmp(common.Big0) != 0) ||
			(dec.S != nil && dec.S.ToInt().Cmp(common.Big0) != 0) {
			return errors.New("deposit transaction signature must be 0 or unset")
		}

		if dec.EthValue != nil {
			var itx DepositTxMantle
			inner = &itx
			if dec.To != nil {
				itx.To = dec.To
			}
			if dec.Gas == nil {
				return errors.New("missing required field 'gas' for txdata")
			}
			itx.Gas = uint64(*dec.Gas)
			if dec.Value == nil {
				return errors.New("missing required field 'value' in transaction")
			}
			itx.Value = (*big.Int)(dec.Value)
			// mint may be omitted or nil if there is nothing to mint.
			itx.Mint = (*big.Int)(dec.Mint)

			// ethValue may be omitted or nil if there is nothing to mint.
			if dec.EthValue != nil {
				itx.EthValue = (*big.Int)(dec.EthValue)
			}
			//
			//// ethValue may be omitted or nil if there is nothing to transfer to msg.To.
			if dec.EthTxValue != nil {
				itx.EthTxValue = (*big.Int)(dec.EthTxValue)
			}

			if dec.Input == nil {
				return errors.New("missing required field 'input' in transaction")
			}
			itx.Data = *dec.Input
			if dec.From == nil {
				return errors.New("missing required field 'from' in transaction")
			}
			itx.From = *dec.From
			if dec.SourceHash != nil {
				itx.SourceHash = *dec.SourceHash
			}
			// IsSystemTx may be omitted. Defaults to false.
			if dec.IsSystemTx != nil {
				itx.IsSystemTransaction = *dec.IsSystemTx
			}

			if dec.Nonce != nil {
				inner = &depositMantleTxWithNonce{DepositTxMantle: itx, EffectiveNonce: uint64(*dec.Nonce)}
			}
		} else if dec.Sender != nil {
			var itx L1MessageTx
			inner = &itx
			if dec.QueueIndex == nil {
				return errors.New("missing required field 'queueIndex' in transaction")
			}
			itx.QueueIndex = uint64(*dec.QueueIndex)
			if dec.Gas == nil {
				return errors.New("missing required field 'gas' in transaction")
			}
			itx.Gas = uint64(*dec.Gas)
			if dec.To != nil {
				itx.To = dec.To
			}
			if dec.Value == nil {
				return errors.New("missing required field 'value' in transaction")
			}
			itx.Value = (*big.Int)(dec.Value)
			if dec.Input == nil {
				return errors.New("missing required field 'input' in transaction")
			}
			itx.Data = *dec.Input
			if dec.Sender == nil {
				return errors.New("missing required field 'sender' in transaction")
			}
			itx.Sender = *dec.Sender
		} else {
			var itx DepositTx
			inner = &itx
			if dec.To != nil {
				itx.To = dec.To
			}
			if dec.Gas == nil {
				return errors.New("missing required field 'gas' for txdata")
			}
			itx.Gas = uint64(*dec.Gas)
			if dec.Value == nil {
				return errors.New("missing required field 'value' in transaction")
			}
			itx.Value = (*big.Int)(dec.Value)
			// mint may be omitted or nil if there is nothing to mint.
			itx.Mint = (*big.Int)(dec.Mint)

			if dec.Input == nil {
				return errors.New("missing required field 'input' in transaction")
			}
			itx.Data = *dec.Input
			if dec.From == nil {
				return errors.New("missing required field 'from' in transaction")
			}
			itx.From = *dec.From
			if dec.SourceHash == nil {
				return errors.New("missing required field 'sourceHash' in transaction")
			}
			itx.SourceHash = *dec.SourceHash
			// IsSystemTx may be omitted. Defaults to false.
			if dec.IsSystemTx != nil {
				itx.IsSystemTransaction = *dec.IsSystemTx
			}

			if dec.Nonce != nil {
				inner = &depositTxWithNonce{DepositTx: itx, EffectiveNonce: uint64(*dec.Nonce)}
			}
		}
	//case ZetaCosmosEVMTxType:
	//	fmt.Println(dec.Hash.Hex())
	//	marshal, err := json.Marshal(dec)
	//	if err != nil {
	//		return err
	//	}
	//	fmt.Println(string(marshal))
	//	var itx ZetaCosmosEvmTx
	//	inner = &itx
	//	itx.Data = dec.Hash.Bytes()

	case ArbitrumLegacyTxType:
		var itx LegacyTx
		if dec.To != nil {
			itx.To = dec.To
		}
		if dec.Nonce == nil {
			return errors.New("missing required field 'nonce' in transaction")
		}
		itx.Nonce = uint64(*dec.Nonce)
		if dec.GasPrice == nil {
			return errors.New("missing required field 'gasPrice' in transaction")
		}
		itx.GasPrice = (*big.Int)(dec.GasPrice)
		if dec.Gas == nil {
			return errors.New("missing required field 'gas' in transaction")
		}
		itx.Gas = uint64(*dec.Gas)
		if dec.Value == nil {
			return errors.New("missing required field 'value' in transaction")
		}
		itx.Value = (*big.Int)(dec.Value)
		if dec.Input == nil {
			return errors.New("missing required field 'input' in transaction")
		}
		itx.Data = *dec.Input
		if dec.V == nil {
			return errors.New("missing required field 'v' in transaction")
		}
		itx.V = (*big.Int)(dec.V)
		if dec.R == nil {
			return errors.New("missing required field 'r' in transaction")
		}
		itx.R = (*big.Int)(dec.R)
		if dec.S == nil {
			return errors.New("missing required field 's' in transaction")
		}
		itx.S = (*big.Int)(dec.S)
		withSignature := itx.V.Sign() != 0 || itx.R.Sign() != 0 || itx.S.Sign() != 0
		if withSignature {
			if err := sanityCheckSignature(itx.V, itx.R, itx.S, true); err != nil {
				return err
			}
		}
		if dec.EffectiveGasPrice == nil {
			return errors.New("missing required field 'EffectiveGasPrice' in transaction")
		}
		if dec.L1BlockNumber == nil {
			return errors.New("missing required field 'L1BlockNumber' in transaction")
		}
		inner = &ArbitrumLegacyTxData{
			LegacyTx:          itx,
			HashOverride:      dec.Hash,
			EffectiveGasPrice: uint64(*dec.EffectiveGasPrice),
			L1BlockNumber:     uint64(*dec.L1BlockNumber),
			Sender:            dec.From,
		}

	case ArbitrumInternalTxType:
		if dec.ChainID == nil {
			return errors.New("missing required field 'chainId' in transaction")
		}
		if dec.Input == nil {
			return errors.New("missing required field 'input' in transaction")
		}
		inner = &ArbitrumInternalTx{
			ChainId: (*big.Int)(dec.ChainID),
			Data:    *dec.Input,
		}

	case ArbitrumDepositTxType:
		if dec.ChainID == nil {
			return errors.New("missing required field 'chainId' in transaction")
		}
		if dec.RequestId == nil {
			return errors.New("missing required field 'requestId' in transaction")
		}
		if dec.To == nil {
			return errors.New("missing required field 'to' in transaction")
		}
		if dec.From == nil {
			return errors.New("missing required field 'from' in transaction")
		}
		if dec.Value == nil {
			return errors.New("missing required field 'value' in transaction")
		}
		inner = &ArbitrumDepositTx{
			ChainId:     (*big.Int)(dec.ChainID),
			L1RequestId: *dec.RequestId,
			To:          *dec.To,
			From:        *dec.From,
			Value:       (*big.Int)(dec.Value),
		}

	case ArbitrumUnsignedTxType:
		if dec.ChainID == nil {
			return errors.New("missing required field 'chainId' in transaction")
		}
		if dec.From == nil {
			return errors.New("missing required field 'from' in transaction")
		}
		if dec.Nonce == nil {
			return errors.New("missing required field 'nonce' in transaction")
		}
		if dec.MaxFeePerGas == nil {
			return errors.New("missing required field 'maxFeePerGas' for txdata")
		}
		if dec.Gas == nil {
			return errors.New("missing required field 'gas' in txdata")
		}
		if dec.Value == nil {
			return errors.New("missing required field 'value' in transaction")
		}
		if dec.Input == nil {
			return errors.New("missing required field 'input' in transaction")
		}
		inner = &ArbitrumUnsignedTx{
			ChainId:   (*big.Int)(dec.ChainID),
			From:      *dec.From,
			Nonce:     uint64(*dec.Nonce),
			GasFeeCap: (*big.Int)(dec.MaxFeePerGas),
			Gas:       uint64(*dec.Gas),
			To:        dec.To,
			Value:     (*big.Int)(dec.Value),
			Data:      *dec.Input,
		}

	case ArbitrumContractTxType:
		if dec.ChainID == nil {
			return errors.New("missing required field 'chainId' in transaction")
		}
		if dec.RequestId == nil {
			return errors.New("missing required field 'requestId' in transaction")
		}
		if dec.From == nil {
			return errors.New("missing required field 'from' in transaction")
		}
		if dec.MaxFeePerGas == nil {
			return errors.New("missing required field 'maxFeePerGas' for txdata")
		}
		if dec.Gas == nil {
			return errors.New("missing required field 'gas' in txdata")
		}
		if dec.Value == nil {
			return errors.New("missing required field 'value' in transaction")
		}
		if dec.Input == nil {
			return errors.New("missing required field 'input' in transaction")
		}
		inner = &ArbitrumContractTx{
			ChainId:   (*big.Int)(dec.ChainID),
			RequestId: *dec.RequestId,
			From:      *dec.From,
			GasFeeCap: (*big.Int)(dec.MaxFeePerGas),
			Gas:       uint64(*dec.Gas),
			To:        dec.To,
			Value:     (*big.Int)(dec.Value),
			Data:      *dec.Input,
		}

	case ArbitrumRetryTxType:
		if dec.ChainID == nil {
			return errors.New("missing required field 'chainId' in transaction")
		}
		if dec.Nonce == nil {
			return errors.New("missing required field 'nonce' in transaction")
		}
		if dec.From == nil {
			return errors.New("missing required field 'from' in transaction")
		}
		if dec.MaxFeePerGas == nil {
			return errors.New("missing required field 'maxFeePerGas' for txdata")
		}
		if dec.Gas == nil {
			return errors.New("missing required field 'gas' in txdata")
		}
		if dec.Value == nil {
			return errors.New("missing required field 'value' in transaction")
		}
		if dec.Input == nil {
			return errors.New("missing required field 'input' in transaction")
		}
		if dec.TicketId == nil {
			return errors.New("missing required field 'ticketId' in transaction")
		}
		if dec.RefundTo == nil {
			return errors.New("missing required field 'refundTo' in transaction")
		}
		if dec.MaxRefund == nil {
			return errors.New("missing required field 'maxRefund' in transaction")
		}
		if dec.SubmissionFeeRefund == nil {
			return errors.New("missing required field 'submissionFeeRefund' in transaction")
		}
		inner = &ArbitrumRetryTx{
			ChainId:             (*big.Int)(dec.ChainID),
			Nonce:               uint64(*dec.Nonce),
			From:                *dec.From,
			GasFeeCap:           (*big.Int)(dec.MaxFeePerGas),
			Gas:                 uint64(*dec.Gas),
			To:                  dec.To,
			Value:               (*big.Int)(dec.Value),
			Data:                *dec.Input,
			TicketId:            *dec.TicketId,
			RefundTo:            *dec.RefundTo,
			MaxRefund:           (*big.Int)(dec.MaxRefund),
			SubmissionFeeRefund: (*big.Int)(dec.SubmissionFeeRefund),
		}

	case ArbitrumSubmitRetryableTxType:
		if dec.ChainID == nil {
			return errors.New("missing required field 'chainId' in transaction")
		}
		if dec.RequestId == nil {
			return errors.New("missing required field 'requestId' in transaction")
		}
		if dec.From == nil {
			return errors.New("missing required field 'from' in transaction")
		}
		if dec.L1BaseFee == nil {
			return errors.New("missing required field 'l1BaseFee' in transaction")
		}
		if dec.DepositValue == nil {
			return errors.New("missing required field 'depositValue' in transaction")
		}
		if dec.MaxFeePerGas == nil {
			return errors.New("missing required field 'maxFeePerGas' for txdata")
		}
		if dec.Gas == nil {
			return errors.New("missing required field 'gas' in txdata")
		}
		if dec.Beneficiary == nil {
			return errors.New("missing required field 'beneficiary' in transaction")
		}
		if dec.MaxSubmissionFee == nil {
			return errors.New("missing required field 'maxSubmissionFee' in transaction")
		}
		if dec.RefundTo == nil {
			return errors.New("missing required field 'refundTo' in transaction")
		}
		if dec.RetryValue == nil {
			return errors.New("missing required field 'retryValue' in transaction")
		}
		if dec.RetryData == nil {
			return errors.New("missing required field 'retryData' in transaction")
		}
		inner = &ArbitrumSubmitRetryableTx{
			ChainId:          (*big.Int)(dec.ChainID),
			RequestId:        *dec.RequestId,
			From:             *dec.From,
			L1BaseFee:        (*big.Int)(dec.L1BaseFee),
			DepositValue:     (*big.Int)(dec.DepositValue),
			GasFeeCap:        (*big.Int)(dec.MaxFeePerGas),
			Gas:              uint64(*dec.Gas),
			RetryTo:          dec.RetryTo,
			RetryValue:       (*big.Int)(dec.RetryValue),
			Beneficiary:      *dec.Beneficiary,
			MaxSubmissionFee: (*big.Int)(dec.MaxSubmissionFee),
			FeeRefundAddr:    *dec.RefundTo,
			RetryData:        *dec.RetryData,
		}
	default:
		fmt.Println(dec.Type)
		return ErrTxTypeNotSupported
	}

	// Now set the inner transaction.
	tx.setDecoded(inner, 0)

	// TODO: check hash here?
	return nil
}

type depositTxWithNonce struct {
	DepositTx
	EffectiveNonce uint64
}

func (tx *depositTxWithNonce) skipAccountChecks() bool {
	//TODO implement me
	panic("implement me")
}

// EncodeRLP ensures that RLP encoding this transaction excludes the nonce. Otherwise, the tx Hash would change
func (tx *depositTxWithNonce) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, tx.DepositTx)
}

func (tx *depositTxWithNonce) effectiveNonce() *uint64 { return &tx.EffectiveNonce }

type depositMantleTxWithNonce struct {
	DepositTxMantle
	EffectiveNonce uint64
}

func (tx *depositMantleTxWithNonce) skipAccountChecks() bool {
	//TODO implement me
	panic("implement me")
}

// EncodeRLP ensures that RLP encoding this transaction excludes the nonce. Otherwise, the tx Hash would change
func (tx *depositMantleTxWithNonce) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, tx.DepositTxMantle)
}

func (tx *depositMantleTxWithNonce) effectiveNonce() *uint64 { return &tx.EffectiveNonce }
