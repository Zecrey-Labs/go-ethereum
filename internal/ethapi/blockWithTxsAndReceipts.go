package ethapi

import (
	"context"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

func (s *BlockChainAPI) GetBlocksWithTxsAndReceipts(ctx context.Context, blockNums []rpc.BlockNumber) (string, error) {
	var res []*types.BlockWithTxsAndReceipts
	for _, blockNum := range blockNums {
		block, err := s.b.BlockByNumber(ctx, blockNum)
		if err != nil {
			return "", err
		}
		hash := block.Hash()
		if hash == (common.Hash{}) {
			hash = block.Header().Hash()
		}
		receipts, err := s.b.GetReceipts(ctx, hash)
		if err != nil {
			return "", err
		}
		res = append(res, &types.BlockWithTxsAndReceipts{
			Header:   block.Header(),
			Receipts: receipts,
		})
	}
	resStr, err := json.Marshal(res)
	if err != nil {
		return "", err
	}

	return string(resStr), nil
}
