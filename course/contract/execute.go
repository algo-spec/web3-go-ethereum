package contract

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"
	store "web3-go-ethereum/course/token/store"
	util "web3-go-ethereum/course/utils"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types" 
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// 通过go生成合约
func OperateContract(client *ethclient.Client, privateKeyStr string, contractAddress string) {

	// 创建合约
	storeContract, err := store.NewStore(common.HexToAddress(contractAddress), client)
	if err != nil {
		fmt.Println("❌ Load new store error", err)
		return
	}

	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		fmt.Println("❌ Load private key error", err)
		return
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		fmt.Println("❌ Load new chain id error", err)
		return
	}
	fmt.Println("✅ Load new chain id: ", chainID)

	// 调用合约函数
	var key [32]byte
	var value [32]byte

	copy(key[:], []byte("1_test_key"))
	copy(value[:], []byte("1_test_value"))

	opt, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		fmt.Println("❌ New keyed transactor error", err)
		return
	}
	tx, err := storeContract.SetItem(opt, key, value)
	if err != nil {
		fmt.Println("❌ SetItem error", err)
		return
	}
	fmt.Println("✅ SetItem success", tx.Hash().Hex())

	// 查询合约数据验证
	callOpt := &bind.CallOpts{Context: context.Background()}
	valueContract, err := storeContract.Items(callOpt, key)
	if err != nil {
		fmt.Println("❌ Items error", err)
		return
	}
	fmt.Println("✅ Items get success, expect equals to original value : ", valueContract == value)
}

// 通过ABI文件创建合约对象
func OperateContractByAbi(client *ethclient.Client, contractAddress string, privateKeyStr string, abiStr string) {

	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		fmt.Println("❌ Load private key error", err)
		return
	}

	// 获取fromAddress
	fromAddress, _, err := util.GetFromAddress(privateKeyStr)
	if err != nil {
		fmt.Println("❌ Get from address error", err)
		return
	}
	fmt.Println("✅ From Address: ", fromAddress.Hex())
	// 获取nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		fmt.Println("❌ get nonce error", err)
		return
	}
	fmt.Println("✅ Nonce: ", nonce)
	// 获取gasPrice
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		fmt.Println("❌ get gas price error", err)
		return
	}
	fmt.Println("✅ Gas Price: ", gasPrice.String())
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		fmt.Println("❌ Load new chain id error", err)
		return
	}
	fmt.Println("✅ Load new chain id: ", chainID)

	// 准备交易calldata
	contractAbi, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		fmt.Println("❌ read abi error", err)
		return
	}

	methodName := "setItem"
	var key [32]byte
	var value [32]byte

	copy(key[:], []byte("2_test_key"))
	copy(value[:], []byte("2_test_value"))
	input, err := contractAbi.Pack(methodName, key, value)
	if err != nil {
		fmt.Println("❌ pack error", err)
		return
	}

	// 创建交易，签名
	tx := types.NewTransaction(
		nonce,
		common.HexToAddress(contractAddress),
		big.NewInt(0),
		300000,
		gasPrice,
		input,
	)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		fmt.Println("❌ sign error", err)
		return
	}
	// 发送交易
	if err := client.SendTransaction(context.Background(), signedTx); err != nil {
		fmt.Println("❌ send error", err)
		return
	}
	fmt.Printf("✅ Transaction sent: %s\n", signedTx.Hash().Hex())

	time.Sleep(5 * time.Second)
	_, err = waitForReceipt(client, signedTx.Hash())
	if err != nil {
		fmt.Println("❌ wait error", err)
		return
	}
	fmt.Println("✅ Transaction completed")

	// 查询刚刚设置的值
	callInput, err := contractAbi.Pack("items", key)
	if err != nil {
		fmt.Println("❌ pack error", err)
		return
	}
	to := common.HexToAddress(contractAddress)
	callMsg := ethereum.CallMsg{
		To:   &to,
		Data: callInput,
	}

	// 解析返回值
	result, err := client.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		fmt.Println("❌ call error", err)
		return
	}
	var unpacked [32]byte
	contractAbi.UnpackIntoInterface(&unpacked, "items", result)
	fmt.Println("✅ Items unpack, expect equals to original value : ", unpacked == value)
}

func waitForReceipt(client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	for {
		receipt, err := client.TransactionReceipt(context.Background(), txHash)
		if err == nil {
			fmt.Println("✅ Transaction receipt:", receipt)
			return receipt, nil
		}
		if err != ethereum.NotFound {
			return nil, err
		}
		fmt.Println("💦 Transaction not mined yet, waiting...")
		time.Sleep(1 * time.Second)
	}
}
