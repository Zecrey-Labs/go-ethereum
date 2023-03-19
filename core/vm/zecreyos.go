package vm

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

var (
	customPrecompiledContracts map[common.Address]StatefulPrecompiledContract
	tmpCtx                     context.Context
	tmpCommit                  bool
)

type StatefulPrecompiledContract interface {
	PrecompiledContract
	RunStateful(ctx context.Context, evm *EVM, addr common.Address, input []byte, value *big.Int) (ret []byte, err error)
}

func SetTmpConfig(ctx context.Context, commit bool) {
	tmpCtx = ctx
	tmpCommit = commit
}

func SetCustomPrecompiledContracts(customContracts map[common.Address]StatefulPrecompiledContract) {
	customPrecompiledContracts = customContracts
}

// RunStatefulPrecompiledContract runs a stateful precompiled contract and ignores the address and
// value arguments. It uses the RunPrecompiledContract function from the geth vm package
func (e *EVM) RunStatefulPrecompiledContract(
	ctx context.Context,
	p StatefulPrecompiledContract,
	caller common.Address, // address arg is unused
	input []byte,
	suppliedGas uint64,
	value *big.Int,
) (ret []byte, remainingGas uint64, err error) {
	gasCost := p.RequiredGas(input)
	if suppliedGas < gasCost {
		return nil, 0, ErrOutOfGas
	}
	suppliedGas -= gasCost
	output, err := p.RunStateful(ctx, e, caller, input, value)
	return output, suppliedGas, err
}

// PreRunStatefulPrecompiledContract runs a stateful precompiled contract and ignores the address and
// value arguments. It uses the RunPrecompiledContract function from the geth vm package
func (e *EVM) PreRunStatefulPrecompiledContract(
	ctx context.Context,
	p StatefulPrecompiledContract,
	caller common.Address, // address arg is unused
	input []byte,
	suppliedGas uint64,
	value *big.Int,
) (ret []byte, remainingGas uint64, err error) {
	gasCost := p.RequiredGas(input)
	if suppliedGas < gasCost {
		return nil, 0, ErrOutOfGas
	}
	suppliedGas -= gasCost
	//output, err := p.RunStateful(e, caller, input, value)
	return nil, suppliedGas, err
}
