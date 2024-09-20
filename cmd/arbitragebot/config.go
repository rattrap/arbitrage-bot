package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Custom errors for missing configuration values
var (
	ErrMissingAPIKey     = fmt.Errorf("missing KuCoin API keys")
	ErrMissingRPCURL     = fmt.Errorf("missing Ethereum RPC URL")
	ErrMissingPrivateKey = fmt.Errorf("missing Ethereum private key")
)

// Config stores all the configuration values for the arbitrage bot.
type Config struct {
	KucoinAPIKey        string // KuCoin API Key
	KucoinAPISecret     string // KuCoin API Secret
	KucoinAPIPassphrase string // KuCoin API Passphrase
	EthereumRPCURL      string // Ethereum RPC URL (e.g., Infura or Alchemy)
	EthereumPrivateKey  string // Private key to sign transactions on Ethereum
}

// LoadConfig loads the configuration values from environment variables or .env file.
func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	config := &Config{}

	// Load KuCoin API keys from environment variables
	config.KucoinAPIKey = os.Getenv("KUCOIN_API_KEY")
	config.KucoinAPISecret = os.Getenv("KUCOIN_API_SECRET")
	config.KucoinAPIPassphrase = os.Getenv("KUCOIN_API_PASSPHRASE")
	if config.KucoinAPIKey == "" || config.KucoinAPISecret == "" || config.KucoinAPIPassphrase == "" {
		return nil, ErrMissingAPIKey
	}

	// Load Ethereum RPC URL (Infura/Alchemy)
	config.EthereumRPCURL = os.Getenv("ETHEREUM_RPC_URL")
	if config.EthereumRPCURL == "" {
		return nil, ErrMissingRPCURL
	}

	// Load Ethereum private key for signing transactions
	config.EthereumPrivateKey = os.Getenv("ETHEREUM_PRIVATE_KEY")
	if config.EthereumPrivateKey == "" {
		return nil, ErrMissingPrivateKey
	}

	return config, nil
}
