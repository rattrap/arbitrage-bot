package main

import (
	"fmt"
	"os"
	"strconv"

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
	TelegramChannelID   int64  // Telegram Channel ID
	TelegramBotToken    string // Telegram Bot Token
}

// LoadConfig loads the configuration values from environment variables or .env file.
func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

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

	// Load Telegram configuration
	telegramChannelID := os.Getenv("TELEGRAM_CHANNEL_ID")
	if telegramChannelID != "" {
		// convert from string to int
		tgID, err := strconv.ParseInt(telegramChannelID, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid Telegram Channel ID: %w", err)
		}
		config.TelegramChannelID = tgID
	}
	config.TelegramBotToken = os.Getenv("TELEGRAM_BOT_TOKEN")

	return config, nil
}
