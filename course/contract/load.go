package contract

import (
	"fmt"
	store "web3-go-ethereum/course/token/store"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func LoadContract(client *ethclient.Client, contractAddress string) {
	
	sContract, err := store.NewStore(common.HexToAddress(contractAddress), client)
	if err != nil {
		fmt.Println("❌ NewStore error", err)
		return
	}
	version, err := sContract.Version(nil)
	if err != nil {
		fmt.Println("❌ Version error", err)
		return
	}
	fmt.Println("✅ LoadContract success", version)
}
