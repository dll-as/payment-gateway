package utils

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"
)

func GeneratePrivateKey() (*ecdsa.PrivateKey, error) {
	return crypto.GenerateKey()
}

func GetPublicKey(privateKey *ecdsa.PrivateKey) []byte {
	return crypto.FromECDSAPub(&privateKey.PublicKey)
}
