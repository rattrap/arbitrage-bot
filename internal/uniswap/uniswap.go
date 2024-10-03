package uniswap

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"time"

	"rattrap/arbitrage-bot/internal/uniswap/contracts"
	"rattrap/arbitrage-bot/internal/utils"

	coreentities "github.com/daoleno/uniswap-sdk-core/entities"
	"github.com/daoleno/uniswapv3-sdk/entities"
	"github.com/daoleno/uniswapv3-sdk/examples/helper"
	"github.com/daoleno/uniswapv3-sdk/periphery"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// UniswapClient represents a client to interact with Uniswap
type UniswapClient struct {
	client      *ethclient.Client
	wallet      *Wallet
	context     context.Context
	pool        *entities.Pool
	tradingPair string
	token0      string
	token1      string
}

// NewUniswapClient initializes a new Uniswap client
func NewUniswapClient(tradingPair, ethereumRPCUrl string, uniswapPoolAddress, uniswapTickLensAddress common.Address, ethereumPrivateKey string, ctx context.Context) (error, *UniswapClient) {
	client, err := ethclient.Dial(ethereumRPCUrl)
	if err != nil {
		return fmt.Errorf("Failed to connect to the Ethereum client"), nil
	}
	defer client.Close()

	wallet := InitWallet(ethereumPrivateKey)
	if wallet == nil {
		return fmt.Errorf("Failed to initialize the wallet"), nil
	}

	ticklens, err := contracts.NewTickLensCaller(uniswapTickLensAddress, client)
	if err != nil {
		return fmt.Errorf("Failed to connect to the TickLens"), nil
	}

	pool, err := ConstructV3Pool(client, uniswapPoolAddress, ticklens, ctx)
	if err != nil {
		return fmt.Errorf("Failed to connect to the Uniswap V3 pool"), nil
	}

	token0, token1 := utils.GetTokensFromTradingPair(tradingPair)

	return nil, &UniswapClient{
		client:      client,
		wallet:      wallet,
		context:     ctx,
		pool:        pool,
		tradingPair: tradingPair,
		token0:      token0,
		token1:      token1,
	}
}

// GetPrice returns the current price of a token on Uniswap
func (c *UniswapClient) GetPrice() (float64, error) {
	price, err := c.pool.PriceOf(c.pool.Token0)
	if err != nil {
		return 0, err
	}

	priceFloat, err := strconv.ParseFloat(price.ToFixed(18), 64)
	if err != nil {
		return 0, err
	}

	return priceFloat, nil
}

func (c *UniswapClient) GetEthBalance() (*coreentities.CurrencyAmount, error) {
	balance, err := c.client.BalanceAt(c.context, c.wallet.PublicKey, nil)
	if err != nil {
		return nil, err
	}
	return coreentities.FromRawAmount(coreentities.EtherOnChain(1), balance), nil
}

// BalanceOf returns the balance of a token in the wallet
func (c *UniswapClient) BalanceOf(token *coreentities.Token) (*coreentities.CurrencyAmount, error) {
	tokenContract, err := contracts.NewERC20Caller(token.Address, c.client)
	if err != nil {
		return nil, err
	}

	balance, err := tokenContract.BalanceOf(nil, c.wallet.PublicKey)
	if err != nil {
		return nil, err
	}

	return coreentities.FromRawAmount(token, balance), nil
}

// GetBalances returns the balances of the wallet
func (c *UniswapClient) GetBalances() (*coreentities.CurrencyAmount, *coreentities.CurrencyAmount, error) {

	token0Balance, err := c.BalanceOf(c.pool.Token0)
	if err != nil {
		return coreentities.FromRawAmount(c.pool.Token0, big.NewInt(0)), coreentities.FromRawAmount(c.pool.Token1, big.NewInt(0)), err
	}

	token1Balance, err := c.BalanceOf(c.pool.Token1)
	if err != nil {
		return coreentities.FromRawAmount(c.pool.Token0, big.NewInt(0)), coreentities.FromRawAmount(c.pool.Token1, big.NewInt(0)), err
	}

	return token0Balance, token1Balance, nil
}

// TargetPriceToSqrtPriceX96 converts a target price to a square root price
func (c *UniswapClient) TargetPriceToSqrtPriceX96(targetPrice float64) *big.Int {
	targetSqrtPrice := new(big.Float).SetFloat64(targetPrice)
	targetSqrtPrice.Quo(targetSqrtPrice, big.NewFloat(math.Pow(10, float64(c.pool.Token0.Decimals()-c.pool.Token1.Decimals()))))
	targetSqrtPriceFloat, _ := targetSqrtPrice.Float64()
	price := PriceToSqrtPriceX96(targetSqrtPriceFloat)
	return price

}

// GetBuyAmount returns the amount of token0 needed to buy token1
func (c *UniswapClient) GetBuyAmount(targetPrice float64) (*coreentities.CurrencyAmount, error) {
	outputAmount := coreentities.FromRawAmount(c.pool.Token0, coreentities.MaxUint256)
	inputAmount, _, err := c.pool.GetInputAmount(outputAmount, c.TargetPriceToSqrtPriceX96(targetPrice))
	if err != nil {
		return coreentities.FromRawAmount(c.pool.Token0, big.NewInt(0)), err
	}

	return inputAmount, nil
}

// GetSellAmount returns the amount of token1 needed to sell token0
func (c *UniswapClient) GetSellAmount(targetPrice float64) (*coreentities.CurrencyAmount, error) {
	outputAmount := coreentities.FromRawAmount(c.pool.Token1, coreentities.MaxUint256)
	inputAmount, _, err := c.pool.GetInputAmount(outputAmount, c.TargetPriceToSqrtPriceX96(targetPrice))
	if err != nil {
		return coreentities.FromRawAmount(c.pool.Token0, big.NewInt(0)), err
	}

	return inputAmount, nil
}

// Trade trades tokens on Uniswap
func (c *UniswapClient) Trade(amount *coreentities.CurrencyAmount) error {
	//0.01%
	slippageTolerance := coreentities.NewPercent(big.NewInt(1), big.NewInt(1000))
	//after 5 minutes
	d := time.Now().Add(time.Minute * time.Duration(15)).Unix()
	deadline := big.NewInt(d)

	var output *coreentities.Token
	if amount.Currency.Equal(c.pool.Token0) {
		output = c.pool.Token1
	} else {
		output = c.pool.Token0
	}

	// single trade input
	// single-hop exact input
	r, err := entities.NewRoute([]*entities.Pool{c.pool}, amount.Currency, output)
	if err != nil {
		return err
	}

	trade, err := entities.FromRoute(r, amount, coreentities.ExactInput)
	if err != nil {
		return err
	}

	fmt.Printf("trade=%+v\n", trade)

	fmt.Printf("%v %v\n", trade.Swaps[0].InputAmount.Quotient(), trade.Swaps[0].OutputAmount.Wrapped().Quotient())
	params, err := periphery.SwapCallParameters([]*entities.Trade{trade}, &periphery.SwapOptions{
		SlippageTolerance: slippageTolerance,
		Recipient:         c.wallet.PublicKey,
		Deadline:          deadline,
	})
	if err != nil {
		return err
	}

	tx, err := SendTX(c.client, common.HexToAddress(helper.ContractV3SwapRouterV1), big.NewInt(0), params.Calldata, c.wallet)
	if err != nil {
		return err
	}
	fmt.Println(tx.Hash().String())

	return nil

}

// Close closes the Uniswap client
func (c *UniswapClient) Close() {
	c.client.Close()
}
