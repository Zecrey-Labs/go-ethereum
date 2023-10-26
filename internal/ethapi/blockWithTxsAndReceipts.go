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
		blockReceipts, err := s.BlockReceipts(ctx, rpc.BlockNumberOrHash{BlockNumber: &blockNum})
		if err != nil {
			return nil, err
		}
		blockWithTxs, err := s.GetBlockByNumber(ctx, blockNum, true)
		if err != nil {
			return nil, err
		}
		res = append(res, types.BlockWithTxsAndReceipts{
			Header:       block.Header(),
			Receipts:     blockReceipts,
			Transactions: blockWithTxs["transactions"],
		})
	}

	return res, nil
}
