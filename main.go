package main

import (
	// "math/big"

	"github.com/ethereum/go-ethereum/ethclient"

	"web3-go-ethereum/config"
	// "web3-go-ethereum/course/subscribe"
	// "web3-go-ethereum/course/query"
	// "web3-go-ethereum/course/transfer"
	// "web3-go-ethereum/course/wallet"
	// "web3-go-ethereum/course/contract"
	// util "web3-go-ethereum/course/utils"
	"web3-go-ethereum/task/task01"
)

func main() {

	config.LoadConfig()

	client, err := ethclient.Dial(config.AppConfig.EthAlchemyURL + config.AppConfig.AlchemyApiKey)
	// client, err := ethclient.Dial(config.AppConfig.EthAlchemyURLWss + config.AppConfig.AlchemyApiKey)
	if err != nil {
		panic(err)
	}

	// ===================================================查询接口测试===============================================
	// latestNumber := query.GetLatestBlock(client)
	// query.BlockInfo(client, latestNumber)
	// query.BlockInfo(client, 9045292)

	// query.GetBlockTransactions(client, 9045292)

	// query.GetBlockTransactionCount(client, "0xf40083be740b6d17d88da29ceb95f571fc70b32ae7e3018a3bd427a508f68a11")

	// query.GetTransactionByHash(client, "0x645ccd8283383155a675e42c008c422f69de579270549497a6f28fcc552e8c31")

	// query.GetBlockReceipts(client, "0xf40083be740b6d17d88da29ceb95f571fc70b32ae7e3018a3bd427a508f68a11")

	// query.GetReceiptByHash(client, "0x645ccd8283383155a675e42c008c422f69de579270549497a6f28fcc552e8c31")

	// wallet.CreateWallet()

	// transfer.TransactionEth(client, config.AppConfig.PrivateKey, config.AppConfig.ReceiverAddress, big.NewInt(1000000000000000))

	// transfer.TransactionErc20(client, config.AppConfig.PrivateKey, config.AppConfig.ReceiverAddress, "0x28b149020d2152179873ec60bed6bf7cd705775d", big.NewInt(1))

	// 账户最新区块余额
	// query.GetBalanceOfAccount(client, config.AppConfig.SenderAddress, nil)

	// 账户指定区块余额
	// query.GetBalanceOfAccount(client, config.AppConfig.SenderAddress, big.NewInt(9045292))

	// query.GetERC20Balance(client, "0x25836239F7b632635F815689389C537133248edb", "0xfadea654ea83c00e5003d2ea15c59830b65471c0")

	// subscribe.SubscribeBlock(client)

	// privateKey, _ := util.GeneratePrivateKey()
	// contract.DeployContract(client, config.AppConfig.PrivateKey)
	// contract.ClientDeploy(client, config.AppConfig.PrivateKey)
	// contract.LoadContract(client, "0x4417753278a2513e1d47ca79e9d312f564ec8237")
	// contract.OperateContract(client, config.AppConfig.PrivateKey, "0x0c762f25Aa12b3Bd6e62DEe05ef80F4CA91A062B")

	// contract.OperateContractByAbi(client, "0x4417753278A2513e1D47cA79e9d312f564ec8237", config.AppConfig.PrivateKey, `[{"inputs":[{"internalType":"string","name":"_version","type":"string"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"bytes32","name":"key","type":"bytes32"},{"indexed":false,"internalType":"bytes32","name":"value","type":"bytes32"}],"name":"ItemSet","type":"event"},{"inputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"name":"items","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"key","type":"bytes32"},{"internalType":"bytes32","name":"value","type":"bytes32"}],"name":"setItem","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"version","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"}]`)

// 	contract.GetContractLog(client, "0x0c762f25Aa12b3Bd6e62DEe05ef80F4CA91A062B", `[
// 	{
// 		"inputs": [
// 			{
// 				"internalType": "string",
// 				"name": "_version",
// 				"type": "string"
// 			}
// 		],
// 		"stateMutability": "nonpayable",
// 		"type": "constructor"
// 	},
// 	{
// 		"anonymous": false,
// 		"inputs": [
// 			{
// 				"indexed": true,
// 				"internalType": "bytes32",
// 				"name": "key",
// 				"type": "bytes32"
// 			},
// 			{
// 				"indexed": false,
// 				"internalType": "bytes32",
// 				"name": "value",
// 				"type": "bytes32"
// 			}
// 		],
// 		"name": "ItemSet",
// 		"type": "event"
// 	},
// 	{
// 		"inputs": [
// 			{
// 				"internalType": "bytes32",
// 				"name": "key",
// 				"type": "bytes32"
// 			},
// 			{
// 				"internalType": "bytes32",
// 				"name": "value",
// 				"type": "bytes32"
// 			}
// 		],
// 		"name": "setItem",
// 		"outputs": [],
// 		"stateMutability": "nonpayable",
// 		"type": "function"
// 	},
// 	{
// 		"inputs": [
// 			{
// 				"internalType": "bytes32",
// 				"name": "",
// 				"type": "bytes32"
// 			}
// 		],
// 		"name": "items",
// 		"outputs": [
// 			{
// 				"internalType": "bytes32",
// 				"name": "",
// 				"type": "bytes32"
// 			}
// 		],
// 		"stateMutability": "view",
// 		"type": "function"
// 	},
// 	{
// 		"inputs": [],
// 		"name": "version",
// 		"outputs": [
// 			{
// 				"internalType": "string",
// 				"name": "",
// 				"type": "string"
// 			}
// 		],
// 		"stateMutability": "view",
// 		"type": "function"
// 	}
// ]`)

	// task01.DeployCounterContract(client, config.AppConfig.PrivateKey)
	task01.OperateContract(client, config.AppConfig.PrivateKey, "0x17F285f4dd0a3A185eb428A3bBfdAc31c76861ef")
}
