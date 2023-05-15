package vm

import "github.com/ethereum/go-ethereum/crypto"

func GetMethodSelector(nameAndParams string) []byte {
	return crypto.Keccak256Hash([]byte(nameAndParams)).Bytes()[:4]
}
