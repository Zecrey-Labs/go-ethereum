package ethapi

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

func (s *BlockChainAPI) BlocksWithTxsAndReceipts(ctx context.Context, blockNums []rpc.BlockNumber) ([]types.BlockWithTxsAndReceipts, error) {
	var res []types.BlockWithTxsAndReceipts
	for _, blockNum := range blockNums {
		block, err := s.b.BlockByNumber(ctx, blockNum)
		if err != nil {
			return nil, err
		}
		hash := block.Hash()
		if hash == (common.Hash{}) {
			hash = block.Header().Hash()
		}
		receipts, err := s.b.GetReceipts(ctx, hash)
		if err != nil {
			return nil, err
		}
		res = append(res, types.BlockWithTxsAndReceipts{
			Header:       block.Header(),
			Receipts:     receipts,
			Transactions: block.Transactions(),
		})
	}

	return res, nil
}
