package query

import (
	"fmt"
	"math"

	"context"
	"math/big"

	erc20 "web3-go-ethereum/course/token/erc20"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// =====================================================查询区块===================================================
/**
 * 获取最新块高
 */
func GetLatestBlock(client *ethclient.Client) uint64 {
	block, err := client.BlockByNumber(context.Background(), nil)
	if err != nil {
		fmt.Println("BlockByNumber error", err)
		return 0
	}
	return block.Number().Uint64()
}

/**
 * 获取块信息
 */
func BlockInfo(client *ethclient.Client, blockNumber uint64) {
	fmt.Println("Block Info: block number is ", blockNumber)
	header, err := client.HeaderByNumber(context.Background(), big.NewInt(int64(blockNumber)))

	if err != nil {
		fmt.Println("Block Header get by Number error...", err)
		return
	}

	fmt.Println("Header: block number is ", header.Number.Uint64())
	fmt.Println("Header: block hash is ", header.Hash().Hex())
	fmt.Println("Header: block time is ", header.Time)
	fmt.Println("Header: block difficulty is ", header.Difficulty.Uint64())

}

// =====================================================查询交易===================================================

/**
 * 获取块中所有交易信息
 */
func GetBlockTransactions(client *ethclient.Client, blockNumber uint64) {
	block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(blockNumber)))
	if err != nil {
		fmt.Println("BlockByNumber error", err)
		return
	}

	chainId, _ := client.ChainID(context.Background())

	for _, tx := range block.Transactions() {
		fmt.Println("Transaction hash: ", tx.Hash().Hex())
		fmt.Println("Transaction to: ", tx.To().Hex())
		fmt.Println("Transaction value: ", tx.Value().String())
		fmt.Println("Transaction gas: ", tx.Gas())
		fmt.Println("Transaction gas price: ", tx.GasPrice().String())
		fmt.Println("Transaction nonce: ", tx.Nonce())
		// fmt.Println("Transaction data: ", string(tx.Data()))

		if from, err := types.Sender(types.NewEIP155Signer(chainId), tx); err == nil {
			fmt.Println("Transaction from: ", from.Hex())
		}

		if receipt, err := client.TransactionReceipt(context.Background(), tx.Hash()); err == nil {
			fmt.Println("Transaction receipt status: ", receipt.Status)
			fmt.Println("Transaction receipt logs: ", receipt.Logs)
		}
		fmt.Println("==========================================================================================")
	}
}

/**
 * 通过交易hash获取块中交易信息
 */
func GetBlockTransactionCount(client *ethclient.Client, bHash string) {
	blockHash := common.HexToHash(bHash)
	count, err := client.TransactionCount(context.Background(), blockHash)
	if err != nil {
		fmt.Println("TransactionCount error", err)
		return
	}
	fmt.Println("TransactionCount: ", count)
	for i := range count {
		tx, err := client.TransactionInBlock(context.Background(), blockHash, i)
		if err != nil {
			fmt.Println("TransactionInBlock error", err)
			return
		}
		fmt.Println("TransactionInBlock: ", tx.Hash().Hex())
		fmt.Println("============================================================================================")
	}
}

/**
 * 通过交易hash获取交易信息
 */
func GetTransactionByHash(client *ethclient.Client, txHash string) {
	tx, isPending, err := client.TransactionByHash(context.Background(), common.HexToHash(txHash))
	if err != nil {
		fmt.Println("TransactionByHash error", err)
		return
	}
	fmt.Println("TransactionByHash: ", tx.Hash().Hex())
	fmt.Println("TransactionByHash isPending: ", isPending)
}

// =====================================================查询收据===================================================
/**
 * 通过区块hash获取所有收据
 */
func GetBlockReceipts(client *ethclient.Client, blockHash string) {
	receipts, err := client.BlockReceipts(context.Background(), rpc.BlockNumberOrHashWithHash(common.HexToHash(blockHash), false))
	if err != nil {
		fmt.Println("BlockReceipts error", err)
		return
	}
	for _, receipt := range receipts {
		fmt.Println("BlockReceipt tx hash: ", receipt.TxHash.Hex())
		fmt.Println("BlockReceipt status: ", receipt.Status)
		fmt.Println("BlockReceipt logs: ", receipt.Logs)
		fmt.Println("BlockReceipt transaction index: ", receipt.TransactionIndex)
		fmt.Println("BlockReceipt contract address: ", receipt.ContractAddress.Hex())
		fmt.Println("==========================================================================================")
	}
}

/**
 * 通过交易hash获取收据
 */
func GetReceiptByHash(client *ethclient.Client, txHash string) {
	receipt, err := client.TransactionReceipt(context.Background(), common.HexToHash(txHash))
	if err != nil {
		fmt.Println("TransactionReceipt error", err)
		return
	}
	fmt.Println("BlockReceipt tx hash: ", receipt.TxHash.Hex())
	fmt.Println("BlockReceipt status: ", receipt.Status)
	fmt.Println("BlockReceipt logs: ", receipt.Logs)
	fmt.Println("BlockReceipt transaction index: ", receipt.TransactionIndex)
	fmt.Println("BlockReceipt contract address: ", receipt.ContractAddress.Hex())
}

/**
 *获取账户余额
 */
func GetBalanceOfAccount(client *ethclient.Client, account string, blockNumber *big.Int) {
	balance, err := client.BalanceAt(context.Background(), common.HexToAddress(account), blockNumber)
	if err != nil {
		fmt.Println("BalanceAt error", err)
		return
	}
	fmt.Println("BalanceOfAccount: ", balance)

	// 转换成ETH
	balanceInEth := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(math.Pow10(18)))
	fmt.Println("BalanceOfAccount in ETH: ", balanceInEth)
}

/**
 * 获取待处理的余额
 */
func GetPendingBalance(client *ethclient.Client, account string) {
	balance, err := client.PendingBalanceAt(context.Background(), common.HexToAddress(account))
	if err != nil {
		fmt.Println("PendingBalanceAt error", err)
		return
	}
	fmt.Println("PendingBalanceOfAccount: ", balance)
	balanceInEth := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(math.Pow10(18)))
	fmt.Println("PendingBalanceOfAccount in ETH: ", balanceInEth)
}

func GetERC20Balance(client *ethclient.Client, account string, tokenAddr string) {
	tokenAddress := common.HexToAddress(tokenAddr)

	instance, err := erc20.NewToken(tokenAddress, client)
	if err != nil {
		fmt.Println("❌ Error creating token instance:", err)
		return
	}
	address := common.HexToAddress(account)

	balance, err := instance.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		fmt.Println("❌ Error getting balance:", err)
		return
	}

	name, err := instance.Name(&bind.CallOpts{})
	if err != nil {
		fmt.Println("❌ Error getting name:", err)
		return
	}

	symbol, err := instance.Symbol(&bind.CallOpts{})
	if err != nil {
		fmt.Println("❌ Error getting symbol:", err)
		return
	}

	decimals, err := instance.Decimals(&bind.CallOpts{})
	if err != nil {
		fmt.Println("❌ Error getting decimals:", err)
		return
	}

	fmt.Printf("ERC20 Token Name: %s\n", name)
	fmt.Printf("ERC20 Token Symbol: %s\n", symbol)
	fmt.Printf("ERC20 Token Decimals: %d\n", decimals)
	fmt.Println("ERC20 balance: ", balance)

	ethBalance := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(math.Pow10(int(decimals))))
	fmt.Println("ERC20 balance in ETH: ", ethBalance)

}
