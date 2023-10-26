package ethapi

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

func (s *BlockChainAPI) GetBlocksWithTxsAndReceipts(ctx context.Context, blockNums []rpc.BlockNumber) ([]*types.RpcBlockWithTxsAndReceipts, error) {
	var res []*types.RpcBlockWithTxsAndReceipts
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
		blockWithTxs, err := s.GetBlockByNumber(ctx, blockNum, true)
		if err != nil {
			return nil, err
		}
		res = append(res, &types.RpcBlockWithTxsAndReceipts{
			Header:       block.Header(),
			Receipts:     receipts,
			Transactions: blockWithTxs["transactions"],
		})
	}

	return res, nil
}
