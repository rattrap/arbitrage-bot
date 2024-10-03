package uniswap

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"strconv"

	"rattrap/arbitrage-bot/internal/uniswap/contracts"

	coreentities "github.com/daoleno/uniswap-sdk-core/entities"
	"github.com/daoleno/uniswapv3-sdk/entities"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// UniswapClient represents a client to interact with Uniswap
type UniswapClient struct {
	client  *ethclient.Client
	wallet  *Wallet
	context context.Context
	pool    *entities.Pool
}

// NewUniswapClient initializes a new Uniswap client
func NewUniswapClient(ethereumRPCUrl string, uniswapPoolAddress, uniswapTickLensAddress common.Address, ethereumPrivateKey string, ctx context.Context) (error, *UniswapClient) {
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

	return nil, &UniswapClient{
		client:  client,
		wallet:  wallet,
		context: ctx,
		pool:    pool,
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

func (c *UniswapClient) TargetPriceToSqrtPriceX96(targetPrice float64) *big.Int {
	targetSqrtPrice := new(big.Float).SetFloat64(targetPrice)
	targetSqrtPrice.Quo(targetSqrtPrice, big.NewFloat(math.Pow(10, float64(c.pool.Token0.Decimals()-c.pool.Token1.Decimals()))))
	targetSqrtPriceFloat, _ := targetSqrtPrice.Float64()
	price := PriceToSqrtPriceX96(targetSqrtPriceFloat)
	return price

}

func (c *UniswapClient) GetBuyAmount(targetPrice float64) (*coreentities.CurrencyAmount, error) {
	outputAmount := coreentities.FromRawAmount(c.pool.Token0, coreentities.MaxUint256)
	inputAmount, _, err := c.pool.GetInputAmount(outputAmount, c.TargetPriceToSqrtPriceX96(targetPrice))
	if err != nil {
		return coreentities.FromRawAmount(c.pool.Token0, big.NewInt(0)), err
	}

	return inputAmount, nil
}

func (c *UniswapClient) GetSellAmount(targetPrice float64) (*coreentities.CurrencyAmount, error) {
	outputAmount := coreentities.FromRawAmount(c.pool.Token1, coreentities.MaxUint256)
	inputAmount, _, err := c.pool.GetInputAmount(outputAmount, c.TargetPriceToSqrtPriceX96(targetPrice))
	if err != nil {
		return coreentities.FromRawAmount(c.pool.Token0, big.NewInt(0)), err
	}

	return inputAmount, nil
}

// Close closes the Uniswap client
func (c *UniswapClient) Close() {
	c.client.Close()
}
