package config

import (
	"os"
	"log"
	"github.com/joho/godotenv"
)

type Config struct {
	EthAlchemyURL string
	EthAlchemyURLWss string
	AlchemyApiKey string
	PrivateKey string
	SenderAddress string
	ReceiverAddress string
}

var AppConfig *Config
func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading config")
		panic(err)
	}

	var cfg Config
	cfg.EthAlchemyURL = os.Getenv("ETH_ALCHEMY_SEPOLIA_URL")
	cfg.EthAlchemyURLWss = os.Getenv("ETH_ALCHEMY_SEPOLIA_URL_WSS")
	cfg.AlchemyApiKey = os.Getenv("ALCHEMY_SEPOLIA_API_KEY")
	cfg.PrivateKey = os.Getenv("PRIVATE_KEY")
	cfg.SenderAddress = os.Getenv("SENDER_ADDRESS")
	cfg.ReceiverAddress = os.Getenv("RECEIVER_ADDRESS")

	AppConfig = &cfg
}