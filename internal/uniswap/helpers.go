package uniswap

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"sort"

	"rattrap/arbitrage-bot/internal/uniswap/contracts"

	coreentities "github.com/daoleno/uniswap-sdk-core/entities"
	"github.com/daoleno/uniswapv3-sdk/constants"
	"github.com/daoleno/uniswapv3-sdk/entities"
	sdkutils "github.com/daoleno/uniswapv3-sdk/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// PriceToSqrtPriceX96 converts a regular price (in float) to SqrtPriceX96.
func PriceToSqrtPriceX96(price float64) *big.Int {
	// Convert the price to big.Float
	priceBigFloat := new(big.Float).SetFloat64(price)
	squaredPrice := new(big.Float).Sqrt(priceBigFloat)

	// Calculate SqrtPriceX96 = (squaredPrice) * (2^96)
	sqrtPriceX96 := new(big.Float).Mul(squaredPrice, new(big.Float).SetInt(constants.Q96))

	// Convert to big.Int
	sqrtPriceX96Int, _ := sqrtPriceX96.Int(nil)

	return sqrtPriceX96Int
}

// ConstructV3Pool constructs a Uniswap V3 pool from the given pool address.
func ConstructV3Pool(client *ethclient.Client, poolAddress common.Address, tickLens *contracts.TickLensCaller, ctx context.Context) (*entities.Pool, error) {
	contractPool, err := contracts.NewUniswapV3PoolCaller(poolAddress, client)
	if err != nil {
		return nil, err
	}

	token0Address, err := contractPool.Token0(nil)
	if err != nil {
		return nil, err
	}

	token1Address, err := contractPool.Token1(nil)
	if err != nil {
		return nil, err
	}

	token0, err := GetTokenEntityFromAddress(client, token0Address, ctx)
	if err != nil {
		return nil, err
	}

	token1, err := GetTokenEntityFromAddress(client, token1Address, ctx)
	if err != nil {
		return nil, err
	}

	liquidity, err := contractPool.Liquidity(nil)
	if err != nil {
		return nil, err
	}

	slot0, err := contractPool.Slot0(nil)
	if err != nil {
		return nil, err
	}

	fee, err := contractPool.Fee(nil)
	if err != nil {
		return nil, err
	}

	ticks := getTicksFromFile(poolAddress.String())
	if len(ticks) == 0 {
		ticks, err := GetPoolTicks(client, fee, poolAddress, tickLens)
		if err != nil {
			return nil, err
		}
		writeTicksToFile(ticks, poolAddress.String())
	}

	fmt.Printf("Pool %s has %d ticks\n", poolAddress.String(), len(ticks))

	// create tick data provider
	p, err := entities.NewTickListDataProvider(ticks, constants.TickSpacings[constants.FeeAmount(fee.Uint64())])
	if err != nil {
		return nil, err
	}

	return entities.NewPool(token0, token1, constants.FeeAmount(fee.Uint64()),
		slot0.SqrtPriceX96, liquidity, int(slot0.Tick.Int64()), p)
}

// GetTokenEntityFromAddress returns a token from the given address.
func GetTokenEntityFromAddress(client *ethclient.Client, tokenAddress common.Address, ctx context.Context) (*coreentities.Token, error) {
	chainId, err := client.ChainID(ctx)
	if err != nil {
		return nil, err
	}

	tokenContract, err := contracts.NewERC20Caller(tokenAddress, client)
	if err != nil {
		return nil, err
	}

	decimals, err := tokenContract.Decimals(nil)
	if err != nil {
		return nil, err
	}

	name, err := tokenContract.Name(nil)
	if err != nil {
		return nil, err
	}

	symbol, err := tokenContract.Symbol(nil)
	if err != nil {
		return nil, err
	}

	token := coreentities.NewToken(uint(chainId.Uint64()), tokenAddress, uint(decimals), name, symbol)

	return token, nil

}

// GetPoolTicks get all ticks of a pool from TickLens smart-contract
func GetPoolTicks(client *ethclient.Client, fee *big.Int, poolAddress common.Address, tickLens *contracts.TickLensCaller) ([]entities.Tick, error) {
	tickSpace := getTickSpacing(float64(fee.Uint64()))

	minWordIndex := sdkutils.MinTick / 256
	poolMinWordIdx := int16(minWordIndex/tickSpace - 1)
	poolMaxWordIdx := -poolMinWordIdx

	// Prepare the list of wordIndexes, the total number of indexes is poolMaxWordIdx-poolMinWordIdx+1
	wordIndexes := make([]int16, 0, poolMaxWordIdx-poolMinWordIdx+1)
	for idx := poolMinWordIdx; idx <= poolMaxWordIdx; idx++ {
		wordIndexes = append(wordIndexes, idx)
	}

	var ticks []entities.Tick

	for _, wordIndex := range wordIndexes {
		populatedTicks, err := tickLens.GetPopulatedTicksInWord(nil, poolAddress, wordIndex)
		if err != nil {
			return nil, err
		}

		for _, pt := range populatedTicks {
			ticks = append(ticks, entities.Tick{
				Index:          int(pt.Tick.Int64()),
				LiquidityGross: pt.LiquidityGross,
				LiquidityNet:   pt.LiquidityNet,
			})
		}
	}

	sort.SliceStable(ticks, func(i, j int) bool {
		return ticks[i].Index < ticks[j].Index
	})

	return ticks, nil
}

func getTickSpacing(swapFee float64) int {
	return constants.TickSpacings[constants.FeeAmount(swapFee)]
}

func writeTicksToFile(ticks []entities.Tick, poolAddress string) {
	ticksJSON, _ := json.Marshal(ticks)

	file := "ticks_" + poolAddress + ".txt"
	f, err := os.Create(file)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	_, err = f.Write(ticksJSON)
	if err != nil {
		panic(err)
	}
}

func getTicksFromFile(poolAddress string) []entities.Tick {
	fmt.Printf("Reading ticks from file %s\n", poolAddress)
	file := "ticks_" + poolAddress + ".txt"
	f, err := os.Open(file)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("File %s does not exist\n", file)
			return []entities.Tick{}
		} else {
			panic(err)
		}
	}

	fi, err := f.Stat()
	if err != nil {
		panic(err)
	}

	// Check if the file is older than 1 hour
	if fi.ModTime().Add(1 * 60 * 60).Before(fi.ModTime()) {
		fmt.Printf("File %s is older than 1 hour\n", file)
		return []entities.Tick{}
	}

	ticks := []entities.Tick{}
	err = json.NewDecoder(f).Decode(&ticks)
	if err != nil {
		panic(err)
	}

	return ticks
}
