package task01

import (
	"context"
	"fmt"
	"math/big"
	"time"

	util "web3-go-ethereum/course/utils"
	counter "web3-go-ethereum/task/task01/contract"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

/**
 * 查询区块信息
 */
func QueryBlockInfo(client *ethclient.Client, blockNumber uint64) {
	block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(blockNumber)))
	if err != nil {
		fmt.Println("❌ BlockByNumber error", err)
		return
	}

	fmt.Println("✅ Block number:", block.Number())
	fmt.Println("✅ Block hash:", block.Hash().Hex())
	fmt.Println("✅ Block transactions:", block.Transactions().Len())
	fmt.Println("✅ Block timestamp:", block.Time())
}

/**
 * 交易以太币
 */
func DoTransfer(client *ethclient.Client, senderPrivateKey string, to string, amount *big.Int) {
	fmt.Println("🔄 Send transaction......")

	// 获取from地址
	fromAddress, privateKey, err := util.GetFromAddress(senderPrivateKey)
	if err != nil {
		fmt.Println("❌ get from address error ", err)
		panic(err)
	}
	// 获取nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		fmt.Println("❌ Error getting nonce", err)
		panic(err)
	}
	fmt.Println("✅ Nonce: ", nonce)

	// 获取gasPrice
	gasPrice, err := client.SuggestGasPrice(context.Background())
	gasLimit := uint64(21000)
	if err != nil {
		fmt.Println("❌ Error getting gas price", err)
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
		fmt.Println("❌ Error getting chain ID", err)
		panic(err)
	}
	fmt.Println("✅ Chain ID: ", chainID.String())

	// 签名交易
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		fmt.Println("❌ Error signing transaction", err)
		panic(err)
	}

	// 广播交易
	if err := client.SendTransaction(context.Background(), signedTx); err != nil {
		fmt.Println("❌ Error sending transaction", err)
		panic(err)
	}
	fmt.Printf("✅ Transaction sent: %s\n", signedTx.Hash().Hex())
}

/**
 * 部署计数器合约
 */
func DeployCounterContract(client *ethclient.Client, privateKeyHex string) {
	fromAddress, privateKey, err := util.GetFromAddress(privateKeyHex)
	if err != nil {
		fmt.Println("❌ DeployCounterContract GetFromAddress error", err)
		return
	}
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		fmt.Println("❌ DeployCounterContract PendingNonceAt error", err)
		return
	}
	fmt.Println("✅ DeployCounterContract Nonce: ", nonce)

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		fmt.Println("❌ DeployCounterContract SuggestGasPrice error", err)
		return
	}
	fmt.Println("✅ DeployCounterContract GasPrice: ", gasPrice)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		fmt.Println("❌ DeployCounterContract NetworkID error", err)
		return
	}
	fmt.Println("✅ DeployCounterContract ChainID: ", chainID.String())

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		fmt.Println("❌ DeployCounterContract NewKeyedTransactorWithChainID error", err)
		return
	}
	fmt.Println("✅ DeployCounterContract auth: ", auth)

	auth.GasPrice = gasPrice
	auth.GasLimit = uint64(3000000)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)

	address, tx, _, err := counter.DeployCounter(auth, client)
	if err != nil {
		fmt.Println("❌ DeployCounterContract DeployCounter error", err)
		return
	}
	fmt.Println("✅ DeployCounterContract address: ", address.Hex())
	fmt.Println("✅ DeployCounterContract tx: ", tx.Hash().Hex())

	// 等待上链，打印回执
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		fmt.Println("❌ DeployCounterContract WaitMined error", err)
		return
	}
	fmt.Println("✅ DeployCounterContract receipt: ", receipt)

	if receipt.BlockNumber != nil {
		fmt.Println("✅ DeployCounterContract block number: ", receipt.BlockNumber.Int64())
	}
	fmt.Println("✅ DeployCounterContract receipt status: ", receipt.Status)
	fmt.Println("✅ DeployCounterContract contract address: ", receipt.ContractAddress.Hex())
}

/**
 * 调用合约方法
 */
func OperateContract(client *ethclient.Client, privateKeyStr string, contractAddress string) {

	// 创建合约
	counterContract, err := counter.NewCounter(common.HexToAddress(contractAddress), client)
	if err != nil {
		fmt.Println("❌ new counter contract error", err)
		return
	}
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		fmt.Println("❌ load private key error", err)
		return
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		fmt.Println("❌ get chain id error", err)
		return
	}
	fmt.Println("✅  chain id: ", chainID)

	opt, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		fmt.Println("❌ New keyed transactor error", err)
		return
	}

	tx, err := counterContract.Increment(opt)
	if err != nil {
		fmt.Println("❌ Increment error", err)
		return
	}
	fmt.Println("✅ Increment success", tx.Hash().Hex())

	callOpt := &bind.CallOpts{Context: context.Background()}
	value, err := counterContract.GetCount(callOpt)
	if err != nil {
		fmt.Println("❌ Count error", err)
		return
	}
	fmt.Println("✅ Count success", value)
}
