package subscribe

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func SubscribeBlock(client *ethclient.Client) {
	headers := make(chan *types.Header)
	subscription, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		fmt.Println("SubscribeBlock: subscribe new head error ", err)
		panic(err)
	}
	for {
		select {
		case header := <-headers:
			fmt.Println("Block: block hash is ", header.Hash().Hex())
			fmt.Println("Block: block number is ", header.Number.Uint64())
			fmt.Println("Block: block time is ", header.Time)
			fmt.Println("Block: block nonce is ", header.Nonce)
			fmt.Println("Block: block difficulty is ", header.Difficulty.Uint64())
		case err := <-subscription.Err():
			fmt.Println("SubscribeBlock: subscribe error ", err)
		}
	}
}
