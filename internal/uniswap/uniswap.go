package uniswap

import (
	"context"
	"fmt"
	"math/big"

	"rattrap/arbitrage-bot/internal/uniswap/contracts"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// UniswapClient represents a client to interact with Uniswap
type UniswapClient struct {
	client        *ethclient.Client
	context       context.Context
	pool          *contracts.UniswapV3PoolCaller
	token0Address common.Address
	token1Address common.Address
	token0        *contracts.IERC20MinimalCaller
	token1        *contracts.IERC20MinimalCaller
}

// NewUniswapClient initializes a new Uniswap client
func NewUniswapClient(ethereumRPCUrl string, ctx context.Context) (error, *UniswapClient) {
	client, err := ethclient.Dial(ethereumRPCUrl)
	if err != nil {
		return fmt.Errorf("Failed to connect to the Ethereum client"), nil
	}
	defer client.Close()

	pool, err := contracts.NewUniswapV3PoolCaller(common.HexToAddress("0x543842CBfef3B3F5614B2153c28936967218A0E6"), client)
	if err != nil {
		return fmt.Errorf("Failed to connect to the Uniswap V3 pool"), nil
	}

	token0Address, err := pool.Token0(nil)
	if err != nil {
		return fmt.Errorf("Failed to get token0 address"), nil
	}

	token1Address, err := pool.Token1(nil)
	if err != nil {
		return fmt.Errorf("Failed to get token1 address"), nil
	}

	token0, err := contracts.NewIERC20MinimalCaller(token0Address, client)
	if err != nil {
		return fmt.Errorf("Failed to connect to token0"), nil
	}

	token1, err := contracts.NewIERC20MinimalCaller(token1Address, client)
	if err != nil {
		return fmt.Errorf("Failed to connect to token1"), nil
	}

	return nil, &UniswapClient{
		client:        client,
		context:       ctx,
		pool:          pool,
		token0Address: token0Address,
		token1Address: token1Address,
		token0:        token0,
		token1:        token1,
	}
}

// GetPrice returns the current price of a token on Uniswap
func (c *UniswapClient) GetPrice(token string) (float64, error) {
	var slot0 struct {
		SqrtPriceX96               *big.Int
		Tick                       *big.Int
		ObservationIndex           uint16
		ObservationCardinality     uint16
		ObservationCardinalityNext uint16
		FeeProtocol                uint8
		Unlocked                   bool
	}
	slot0, err := c.pool.Slot0(nil)
	if err != nil {
		return 0, fmt.Errorf("Failed to get slot0: %s", err)
	}

	// price := new(big.Float).SetInt(slot0.SqrtPriceX96)
	// price.Mul(price, price)
	// price.Quo(price, new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(2), big.NewInt(192), nil)))

	sqrtPriceX96 := new(big.Float).SetInt(slot0.SqrtPriceX96)

	scaleFactor := new(big.Int).Exp(big.NewInt(2), big.NewInt(96), nil)
	price := new(big.Float).Quo(sqrtPriceX96, new(big.Float).SetInt(scaleFactor)) // sqrtPriceX96 / 2^96
	price.Mul(price, price)

	token0Decimals, err := c.token0.Decimals(nil)
	if err != nil {
		return 0, fmt.Errorf("Failed to get token0 decimals: %s", err)
	}

	token1Decimals, err := c.token1.Decimals(nil)
	if err != nil {
		return 0, fmt.Errorf("Failed to get token1 decimals: %s", err)
	}

	price.Mul(price, new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(token0Decimals)-int64(token1Decimals)), nil)))

	// Convert the price to a float64
	priceFloat, _ := price.Float64()

	return priceFloat, nil
}

// Get liquidity returns the current liquidity of a token on Uniswap
func (c *UniswapClient) GetLiquidity(token string) (float64, error) {
	liquidity, err := c.pool.Liquidity(nil)
	if err != nil {
		return 0, fmt.Errorf("Failed to get liquidity: %s", err)
	}

	liquidityToken0, err := c.token0.BalanceOf(nil, c.token0Address)
	if err != nil {
		return 0, fmt.Errorf("Failed to get token0 balance: %s", err)
	}

	liquidityToken1, err := c.token1.BalanceOf(nil, c.token1Address)
	if err != nil {
		return 0, fmt.Errorf("Failed to get token1 balance: %s", err)
	}

	fmt.Printf("Liquidity: %s\n", liquidity.String())
	fmt.Printf("Token0: %s\n", liquidityToken0.String())
	fmt.Printf("Token1: %s\n", liquidityToken1.String())

	return 0, nil
}

// Close closes the Uniswap client
func (c *UniswapClient) Close() {
	c.client.Close()
}
