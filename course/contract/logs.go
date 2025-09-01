package contract

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// var StoreABI = `[{"inputs":[{"internalType":"string","name":"_version","type":"string"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"bytes32","name":"key","type":"bytes32"},{"indexed":false,"internalType":"bytes32","name":"value","type":"bytes32"}],"name":"ItemSet","type":"event"},{"inputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"name":"items","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"key","type":"bytes32"},{"internalType":"bytes32","name":"value","type":"bytes32"}],"name":"setItem","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"version","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"}]`

func GetContractLog(client *ethclient.Client, contractAddressStr string, storeAbi string) {

	contractAddress := common.HexToAddress(contractAddressStr)

	// latestBlock, err := client.BlockNumber(context.Background())
	// if err != nil {
	// 	fmt.Println("获取最新区块号错误", err)
	// 	return
	// }

	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(9107893),
		ToBlock: big.NewInt(9107894),
		Addresses: []common.Address{
			contractAddress,
		},
	}
	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		fmt.Println("FilterLogs error", err)
		return
	}

	contractAbi, err := abi.JSON(strings.NewReader(storeAbi))
	if err != nil {
		fmt.Println("❌ get abi error", err)
		return
	}

	for _, vLog := range logs {
		fmt.Println("vLog block hash:", vLog.BlockHash.Hex())
		fmt.Println("vLog block number:", vLog.BlockNumber)
		fmt.Println("vLog tx hash:", vLog.TxHash.Hex())
		event := struct {
			Key   [32]byte
			Value [32]byte
		}{}
		err := contractAbi.UnpackIntoInterface(&event, "ItemSet", vLog.Data)
		if err != nil {
			fmt.Println("❌ unpack error", err)
			return
		}
		fmt.Println("✅ ItemSet unpack, key:", common.Bytes2Hex(event.Key[:]) , "value:", common.Bytes2Hex(event.Value[:]))

		var topics []string
		for i := range vLog.Topics {
			topics = append(topics, vLog.Topics[i].Hex())
		}
		fmt.Println("✅ ItemSet topics:", topics)
		fmt.Println("✅ ItemSet topics[0]:", topics[1:])
	}
	eventSignature := []byte("ItemSet(bytes32,string)")
	hash := crypto.Keccak256Hash(eventSignature)
	fmt.Println("✅ ItemSet hash:", hash.Hex())
}
