package vm

import (
	"bytes"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"math/big"
	"time"
)

func GetMethodSelector(nameAndParams string) []byte {
	return crypto.Keccak256Hash([]byte(nameAndParams)).Bytes()[:4]
}

const initGas = 80000000

var (
	transferFromSelector = GetMethodSelector("transferFrom(address,address,uint256)")
	approveSelector      = GetMethodSelector("approve(address,uint256)")
	allowanceSelector    = GetMethodSelector("allowance(address,address)")
	balanceOfSelector    = GetMethodSelector("balanceOf(address)")
	transferSelector     = GetMethodSelector("transfer(address,uint256)")
)

type SimulateResponse struct {
	AssetChanges         []AssetChange
	GasCost              uint64
	SuccessWithPrePay    bool
	SuccessWithoutPrePay bool
	ErrInfo              string
}

type AssetChange struct {
	AssetAddress  string
	Sender        string
	Receiver      string
	AssetAmount   string
	Spender       string
	Allowance     string
	SenderBalance string
	ActionType    string
}

func (evm *EVM) simulateCall(caller ContractRef, addr common.Address, input []byte, gas uint64, value *big.Int) (ret []byte, leftOverGas uint64, err error) {
	// Fail if we're trying to execute above the call depth limit
	if evm.depth > int(params.CallCreateDepth) {
		return nil, gas, ErrDepth
	}
	// Fail if we're trying to transfer more than the available balance
	gas = initGas
	evm.simulateNativeAsset(caller.Address(), addr, value)
	gas -= 21000
	snapshot := evm.StateDB.Snapshot()
	p, isPrecompile := evm.precompile(addr)

	if !evm.StateDB.Exist(addr) {
		if !isPrecompile && evm.chainRules.IsEIP158 && value.Sign() == 0 {
			// Calling a non existing account, don't do anything, but ping the tracer
			if evm.Config.Debug {
				if evm.depth == 0 {
					evm.Config.Tracer.CaptureStart(evm, caller.Address(), addr, false, input, gas, value)
					evm.Config.Tracer.CaptureEnd(ret, 0, 0, nil)
				} else {
					evm.Config.Tracer.CaptureEnter(CALL, caller.Address(), addr, input, gas, value)
					evm.Config.Tracer.CaptureExit(ret, 0, nil)
				}
			}
			return nil, gas, nil
		}
		evm.StateDB.CreateAccount(addr)
	}
	evm.Context.Transfer(evm.StateDB, caller.Address(), addr, value)

	// Capture the tracer start/end events in debug mode
	if evm.Config.Debug {
		if evm.depth == 0 {
			evm.Config.Tracer.CaptureStart(evm, caller.Address(), addr, false, input, gas, value)
			defer func(startGas uint64, startTime time.Time) { // Lazy evaluation of the parameters
				evm.Config.Tracer.CaptureEnd(ret, startGas-gas, time.Since(startTime), err)
			}(gas, time.Now())
		} else {
			// Handle tracer events for entering and exiting a call frame
			evm.Config.Tracer.CaptureEnter(CALL, caller.Address(), addr, input, gas, value)
			defer func(startGas uint64) {
				evm.Config.Tracer.CaptureExit(ret, startGas-gas, err)
			}(gas)
		}
	}

	if isPrecompile {
		ret, gas, err = RunPrecompiledContract(p, input, gas)
	} else {
		// Initialise a new contract and set the code that is to be used by the EVM.
		// The contract is a scoped environment for this execution context only.
		code := evm.StateDB.GetCode(addr)
		if len(code) == 0 {
			ret, err = nil, nil // gas is unchanged
		} else {
			addrCopy := addr
			// If the account has no code, we can abort here
			// The depth-check is already done, and precompiles handled above
			contract := NewContract(caller, AccountRef(addrCopy), value, gas)
			contract.SetCallCode(&addrCopy, evm.StateDB.GetCodeHash(addrCopy), code)
			// if that's simulate, do assets change check
			ret, err = evm.simulateAction(contract, caller, addr, input)
			gas = contract.Gas
		}
	}
	// When an error was returned by the EVM or when setting the creation code
	// above we revert to the snapshot and consume any gas remaining. Additionally
	// when we're in homestead this also counts for code storage gas errors.
	if err != nil {
		evm.StateDB.RevertToSnapshot(snapshot)
		if err != ErrExecutionReverted {
			evm.SimulateResp.GasCost += initGas - gas
			gas = 0
		}
		// TODO: consider clearing up unused snapshots:
		//} else {
		//	evm.StateDB.DiscardSnapshot(snapshot)
		evm.SimulateResp.SuccessWithPrePay = false
	}
	evm.SimulateResp.SuccessWithPrePay = true
	if gas != 0 {
		evm.SimulateResp.GasCost += initGas - gas
	}
	return ret, gas, err
}

func (evm *EVM) erc20Allowance(contract *Contract, from, to common.Address) *big.Int {
	// get allowance
	var buf bytes.Buffer
	buf.Write(allowanceSelector)
	buf.Write(new(big.Int).SetBytes(from.Bytes()).FillBytes(make([]byte, 32)))
	buf.Write(new(big.Int).SetBytes(to.Bytes()).FillBytes(make([]byte, 32)))
	var (
		allowanceRet []byte
		err          error
	)
	allowanceRet, err = evm.interpreter.Run(contract, buf.Bytes(), false)
	if err != nil {
		log.Warn("simulate: cannot get allowance:", err)
	}
	return new(big.Int).SetBytes(allowanceRet)
}

func (evm *EVM) erc20Balance(contract *Contract, from common.Address) *big.Int {
	// get balance
	var buf bytes.Buffer
	buf.Write(balanceOfSelector)
	buf.Write(new(big.Int).SetBytes(from.Bytes()).FillBytes(make([]byte, 32)))
	var (
		balanceRet []byte
		err        error
	)
	balanceRet, err = evm.interpreter.Run(contract, buf.Bytes(), true)
	if err != nil {
		log.Warn("simulate: cannot get balance for sender:", err)
	}
	return new(big.Int).SetBytes(balanceRet)
}

func (evm *EVM) erc20Approve(caller ContractRef, fromAddr common.Address, addr common.Address, amount *big.Int) {
	// force approve
	var buf bytes.Buffer
	buf.Write(approveSelector)
	buf.Write(new(big.Int).SetBytes(caller.Address().Bytes()).FillBytes(make([]byte, 32)))
	buf.Write(amount.FillBytes(make([]byte, 32)))
	var err error
	_, _, err = evm.Call(AccountRef(fromAddr), addr, buf.Bytes(), 80000, big.NewInt(0))
	if err != nil {
		log.Warn("simulate: cannot approve for sender:", err)
	}
}

func (evm *EVM) simulateAction(contract *Contract, caller ContractRef, addr common.Address, input []byte) (ret []byte, err error) {
	// catch transferFrom call
	// if that's transferFrom call, decode inputs
	var assetChange AssetChange
	if bytes.Equal(transferFromSelector, input[:4]) && len(input) == 100 {
		info := input[4:]
		fromAddr := common.BytesToAddress(info[:32])
		toAddr := common.BytesToAddress(info[32:64])
		amount := new(big.Int).SetBytes(info[64:])
		// get allowance
		allowance := evm.erc20Allowance(contract, fromAddr, caller.Address())
		assetChange.Allowance = allowance.String()
		// force approve
		if allowance.Cmp(amount) < 0 {
			evm.erc20Approve(caller, fromAddr, addr, amount)
		}
		// fill asset change info
		assetChange.AssetAddress = addr.Hex()
		assetChange.AssetAmount = amount.String()
		assetChange.Sender = fromAddr.Hex()
		balance := evm.erc20Balance(contract, fromAddr)
		assetChange.SenderBalance = balance.String()
		if balance.Cmp(amount) < 0 {
			evm.SimulateResp.SuccessWithoutPrePay = false
		}
		assetChange.Receiver = toAddr.Hex()
		assetChange.Spender = caller.Address().Hex()
		assetChange.ActionType = "transferFrom"
		evm.SimulateResp.AssetChanges = append(evm.SimulateResp.AssetChanges, assetChange)
	} else if bytes.Equal(transferSelector, input[:4]) && len(input) == 68 {
		info := input[4:]
		toAddr := common.BytesToAddress(info[:32])
		amount := new(big.Int).SetBytes(info[32:])
		// fill asset change info
		assetChange.AssetAddress = addr.Hex()
		assetChange.AssetAmount = amount.String()
		assetChange.Sender = caller.Address().Hex()
		assetChange.SenderBalance = evm.erc20Balance(contract, caller.Address()).String()
		assetChange.Receiver = toAddr.Hex()
		assetChange.Spender = common.Address{}.Hex()
		assetChange.Allowance = "0"
		assetChange.ActionType = "transfer"
		evm.SimulateResp.AssetChanges = append(evm.SimulateResp.AssetChanges, assetChange)
	}
	ret, err = evm.interpreter.Run(contract, input, false)
	if err != nil {
		evm.SimulateResp.ErrInfo = err.Error()
		log.Warn("simulate: unable to run contract:", err)
	}
	return ret, nil
}

func (evm *EVM) simulateNativeAsset(from, to common.Address, value *big.Int) {
	if value.Cmp(big.NewInt(0)) == 0 {
		return
	}
	// catch transferFrom call
	// if that's transferFrom call, decode inputs
	var assetChange AssetChange
	// fill asset change info
	assetChange.AssetAddress = common.Address{}.Hex()
	assetChange.AssetAmount = value.String()
	assetChange.Sender = from.Hex()
	assetChange.SenderBalance = evm.StateDB.GetBalance(from).String()
	assetChange.Receiver = to.Hex()
	assetChange.Spender = common.Address{}.Hex()
	assetChange.Allowance = "0"
	assetChange.ActionType = "native"
	evm.SimulateResp.AssetChanges = append(evm.SimulateResp.AssetChanges, assetChange)
}
