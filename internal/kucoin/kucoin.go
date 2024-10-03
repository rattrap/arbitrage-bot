package kucoin

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	kucoin "github.com/Kucoin/kucoin-go-sdk"
)

// KucoinClient represents a client to interact with KuCoin
type KucoinClient struct {
	client  *kucoin.ApiService
	context context.Context
}

// NewKucoinClient initializes a new KuCoin API client
func NewKucoinClient(apiKey, apiSecret, apiPassphrase string, context context.Context) (error, *KucoinClient) {
	client := kucoin.NewApiService(
		// kucoin.ApiBaseURIOption("https://api.kucoin.com"),
		kucoin.ApiKeyOption(apiKey),
		kucoin.ApiSecretOption(apiSecret),
		kucoin.ApiPassPhraseOption(apiPassphrase),
		kucoin.ApiKeyVersionOption(kucoin.ApiKeyVersionV2),
	)

	status, err := client.ServiceStatus(context)
	if err != nil {
		return fmt.Errorf("Failed to connect to the KuCoin API"), nil
	}

	var s struct {
		Status string `json:"status"`
	}
	err = json.Unmarshal(status.RawData, &s)
	if err != nil {
		return fmt.Errorf("Failed to parse KuCoin API status"), nil
	}

	if s.Status != "open" {
		return fmt.Errorf("KuCoin API is not open: %s", s.Status), nil
	}
	return nil, &KucoinClient{
		client:  client,
		context: context,
	}
}

// BalanceOf returns the balance of a currency
func (c *KucoinClient) BalanceOf(currency string) (float64, error) {
	account, err := c.client.Accounts(c.context, "", "")
	if err != nil {
		return 0, fmt.Errorf("Failed to get account list: %s", err)
	}

	accounts := &kucoin.AccountsModel{}
	if err := account.ReadData(accounts); err != nil {
		return 0, fmt.Errorf("Failed to read account data: %s", err)
	}

	for _, a := range *accounts {
		if a.Currency == currency {
			balance, err := strconv.ParseFloat(a.Available, 64)
			if err != nil {
				return 0, fmt.Errorf("Failed to parse balance: %s", err)
			}
			return balance, nil
		}
	}

	return 0, fmt.Errorf("Currency %s not found", currency)
}

// GetPrice returns the current price of a trading pair
func (c *KucoinClient) GetPrice(tradingPair string) (float64, error) {
	ticker, err := c.client.TickerLevel1(c.context, tradingPair)
	if err != nil {
		return 0, fmt.Errorf("Failed to get ticker for %s: %s", tradingPair, err)
	}

	t := &kucoin.TickerLevel1Model{}
	if err := ticker.ReadData(t); err != nil {
		return 0, fmt.Errorf("Failed to read ticker data for %s: %s", tradingPair, err)
	}

	price, err := strconv.ParseFloat(t.Price, 64)
	if err != nil {
		return 0, fmt.Errorf("Failed to parse price for %s: %s", tradingPair, err)
	}

	return price, nil
}

// Close closes the KuCoin client
func (c *KucoinClient) Close() {
}
