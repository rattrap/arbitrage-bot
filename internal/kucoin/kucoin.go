package kucoin

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	kucoin "github.com/Kucoin/kucoin-go-sdk"

	"rattrap/arbitrage-bot/internal/utils"
)

// KucoinClient represents a client to interact with KuCoin
type KucoinClient struct {
	client      *kucoin.ApiService
	context     context.Context
	tradingPair string
	token0      string
	token1      string
}

// NewKucoinClient initializes a new KuCoin API client
func NewKucoinClient(tradingPair, apiKey, apiSecret, apiPassphrase string, context context.Context) (error, *KucoinClient) {
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

	token0, token1 := utils.GetTokensFromTradingPair(tradingPair)

	return nil, &KucoinClient{
		client:      client,
		context:     context,
		tradingPair: tradingPair,
		token0:      token0,
		token1:      token1,
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
func (c *KucoinClient) GetPrice() (float64, error) {
	ticker, err := c.client.TickerLevel1(c.context, c.tradingPair)
	if err != nil {
		return 0, fmt.Errorf("Failed to get ticker for %s: %s", c.tradingPair, err)
	}

	t := &kucoin.TickerLevel1Model{}
	if err := ticker.ReadData(t); err != nil {
		return 0, fmt.Errorf("Failed to read ticker data for %s: %s", c.tradingPair, err)
	}

	price, err := strconv.ParseFloat(t.Price, 64)
	if err != nil {
		return 0, fmt.Errorf("Failed to parse price for %s: %s", c.tradingPair, err)
	}

	return price, nil
}

// Trade executes a trade
func (c *KucoinClient) Trade(side, symbol, size string, priceLimit float64) error {
	priceLimitStr := fmt.Sprintf("%.18f", priceLimit)
	order, err := c.client.CreateOrder(c.context, &kucoin.CreateOrderModel{
		ClientOid:   kucoin.IntToString(time.Now().UnixNano()),
		Symbol:      c.tradingPair,
		Side:        side,
		Type:        "limit",
		Price:       priceLimitStr,
		Size:        size,
		TimeInForce: "GTC",
	})
	fmt.Printf("order=%+v\n", order)
	return err
}

// Close closes the KuCoin client
func (c *KucoinClient) Close() {
}
