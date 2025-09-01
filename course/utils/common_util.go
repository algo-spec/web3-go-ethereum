package utils

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func GeneratePrivateKey() (string, error) {
	fmt.Println("GenerateKey: start...")
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		fmt.Println("❌ GenerateKey: error, ", err)
		return "", err
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyHex := hex.EncodeToString(privateKeyBytes)
	fmt.Println("GenerateKey: Private key ", privateKeyHex)
	return privateKeyHex, nil
}

func GetFromAddress(senderPrivateKey string) (common.Address, *ecdsa.PrivateKey, error) {
	privateKey, err := crypto.HexToECDSA(senderPrivateKey)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("❌ Error converting hex to ECDSA: %w", err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return common.Address{}, privateKey, fmt.Errorf("❌ Error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Println("✅ From Address: ", fromAddress.Hex())
	return fromAddress, privateKey, nil
}
