package transfer

import (
	"context"
	"fmt"
	"math/big"

	util "web3-go-ethereum/course/utils"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
)

// ETH 交易
func TransactionEth(client *ethclient.Client, senderPrivateKey string, to string, amount *big.Int) {
	fmt.Println("🔄 Transaction ETH......")

	// 获取from地址
	fromAddress, privateKey, err := util.GetFromAddress(senderPrivateKey)
	if err != nil {
		fmt.Println("Error getting from address", err)
		panic(err)
	}

	// 获取nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		fmt.Println("Error getting nonce", err)
		panic(err)
	}
	fmt.Println("✅ Nonce: ", nonce)

	// 获取gasPrice
	gasPrice, err := client.SuggestGasPrice(context.Background())
	gasLimit := uint64(21000)
	if err != nil {
		fmt.Println("Error getting gas price", err)
		panic(err)
	}
	fmt.Println("✅ Gas Price: ", gasPrice.String())

	// 获取to地址
	toAddress := common.HexToAddress(to)
	fmt.Println("✅ To Address: ", toAddress.Hex())

	// 创建交易
	tx := types.NewTransaction(nonce, toAddress, amount, gasLimit, gasPrice, nil)

	// 获取链ID
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		fmt.Println("Error getting chain ID", err)
		panic(err)
	}
	fmt.Println("✅ Chain ID: ", chainID.String())

	// 签名交易
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		fmt.Println("Error signing transaction", err)
		panic(err)
	}

	// 广播交易
	if err := client.SendTransaction(context.Background(), signedTx); err != nil {
		fmt.Println("Error sending transaction", err)
		panic(err)
	}
	fmt.Printf("✅ Transaction sent: %s\n", signedTx.Hash().Hex())
}

/**
 * 创建ERC20代币交易
 */
func TransactionErc20(client *ethclient.Client, senderPrivateKey string, to string, token string, amount *big.Int) {
	fmt.Println("🔄 Transaction ERC20......")

	// 获取from地址
	fromAddress, privateKey, err := util.GetFromAddress(senderPrivateKey)
	if err != nil {
		fmt.Println("Error getting from address", err)
		panic(err)
	}
	// to 地址
	toAddress := common.HexToAddress(to)

	// token 地址
	tokenAddress := common.HexToAddress(token)

	// 手动构造ERC20中transfer函数
	transferFnS := []byte("transfer(address,uint256)")
	// 获取方法ID
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnS)
	methodID := hash.Sum(nil)[:4]
	fmt.Println("✅ Method ID: ", hexutil.Encode(methodID))

	// 获取方法参数
	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	fmt.Println("✅ Padded Address: ", hexutil.Encode(paddedAddress))

	// 转为单位：wei
	tokenAmount := new(big.Int).Mul(amount, big.NewInt(1e18))
	fmt.Println("✅ Convert Amount: ", tokenAmount)
	paddedAmount := common.LeftPadBytes(tokenAmount.Bytes(), 32)
	fmt.Println("✅ Padded Amount: ", hexutil.Encode(paddedAmount))

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		fmt.Println("Error getting gas price", err)
		panic(err)
	}
	fmt.Println("✅ Gas Price: ", gasPrice.String())
	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{From: fromAddress, To: &tokenAddress, Data: data})
	if err != nil {
		fmt.Println("Error getting gas limit", err)
		gasLimit = 60000
	}
	fmt.Println("✅ Gas Limit: ", gasLimit)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		fmt.Println("Error getting nonce", err)
		panic(err)
	}
	fmt.Println("✅ Nonce: ", nonce)

	tx := types.NewTransaction(nonce, tokenAddress, big.NewInt(0), gasLimit, gasPrice, data)
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		fmt.Println("Error getting chain ID", err)
		panic(err)
	}
	fmt.Println("✅ Chain ID: ", chainID.String())
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		fmt.Println("Error signing transaction", err)
		panic(err)
	}
	if err := client.SendTransaction(context.Background(), signedTx); err != nil {
		fmt.Println("Error sending transaction", err)
		panic(err)
	}
	fmt.Printf("✅ Transaction sent: %s\n", signedTx.Hash().Hex())
}
 