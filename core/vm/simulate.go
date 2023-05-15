package vm

import (
	"github.com/ethereum/go-ethereum/crypto"
)

func GetMethodSelector(nameAndParams string) []byte {
	return crypto.Keccak256Hash([]byte(nameAndParams)).Bytes()[:4]
}

type SimulateAssetsChangeResp struct {
	AssetChanges []AssetChange
}

type AssetChange struct {
	AssetAddress          string
	Sender                string
	Receiver              string
	AssetAmount           string
	ReceiverBalanceBefore string
	ReceiverBalanceAfter  string
	Allowance             string
}
