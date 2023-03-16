package vm

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// TODO used for zecrey os
var (
	addrType, _  = abi.NewType("address", "", nil)
	bytesType, _ = abi.NewType("bytes", "", nil)
)

func PadAddressIntoInput(caller common.Address, contract common.Address, input []byte) ([]byte, error) {
	if contract.Hex() == common.HexToAddress("0x3f").Hex() {
		args := abi.Arguments{
			{
				Type: addrType,
			},
			{
				Type: bytesType,
			},
		}
		input, err := args.Pack(
			caller,
			input,
		)
		if err != nil {
			return nil, err
		}
		return input, nil
	}
	return input, nil
}
