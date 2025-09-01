package wallet

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
)

/**
 * 创建钱包
 */
func CreateWallet() {
	fmt.Println("CreateWallet")

	// 生成私钥
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		fmt.Println("Error generating key", err)
		panic(err)
	}

	// privateKey, err2 := crypto.HexToECDSA("242dc20274e4b06f46562f04fdac86ea97a7c26ec25611e624193ab701b6e19b")
	// if err2 != nil {
	// 	log.Println("Error converting hex to ECDSA", err2)
	// 	panic(err2)
	// }

	// 转为字节，私钥需要剥离前缀0x
	privateKeyBytes := crypto.FromECDSA(privateKey)
	fmt.Println("Private key: ", hexutil.Encode(privateKeyBytes)[2:])

	// 通过私钥获取公钥
	publicKey := privateKey.Public()

	// 转为字节
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("Error casting public key to ECDSA")
	}
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	fmt.Println("Public key: ", hexutil.Encode(publicKeyBytes)[4:])

	// 通过公钥获取公共地址【公钥的keccak-256】
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println("Address: ", address)
	hash := sha3.NewLegacyKeccak512()
	hash.Write(publicKeyBytes[1:])
	fmt.Println("Full Address: ", hexutil.Encode(hash.Sum(nil)[:]))

	// 截取前12位，保留后20位
	fmt.Println("View Address: ", hexutil.Encode(hash.Sum(nil)[12:]))
}
