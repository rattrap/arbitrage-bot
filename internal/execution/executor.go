package execution

import (
	"rattrap/arbitrage-bot/internal/kucoin"
	"rattrap/arbitrage-bot/internal/logging"
	"rattrap/arbitrage-bot/internal/uniswap"
	"rattrap/arbitrage-bot/internal/utils"
	"strconv"

	"github.com/sirupsen/logrus"
)

// Executor handles trade execution for both KuCoin and Uniswap
type Executor struct {
	paperTrading  bool
	uniswapClient *uniswap.UniswapClient
	kucoinClient  *kucoin.KucoinClient
	logger        *logrus.Entry
	tradingPair   string
	token0        string
	token1        string
	balances      map[string]string
}

// NewExecutor initializes a new Executor
func NewExecutor(paperTrading bool, tradingPair string, uniswapClient *uniswap.UniswapClient, kucoinClient *kucoin.KucoinClient, logger *logging.Logger) *Executor {
	prefixedLogger := logger.WithField("prefix", "execution")
	token0, token1 := utils.GetTokensFromTradingPair(tradingPair)

	return &Executor{
		paperTrading:  paperTrading,
		uniswapClient: uniswapClient,
		kucoinClient:  kucoinClient,
		logger:        prefixedLogger,
		tradingPair:   tradingPair,
		token0:        token0,
		token1:        token1,
		balances:      make(map[string]string),
	}
}

func (e *Executor) Start() {
	e.logger.Debug("Starting service")
	e.GetBalances()
}

// GetBalances
func (e *Executor) GetBalances() {
	ethBalance, err := e.uniswapClient.GetEthBalance()
	if err != nil {
		e.logger.WithError(err).Error("Failed to get ETH balance")
		return
	}

	e.balances["ETH"] = ethBalance.ToExact()

	token0Uniswap, token1Uniswap, err := e.uniswapClient.GetBalances()
	if err != nil {
		e.logger.WithError(err).Error("Failed to get Uniswap balances")
		return
	}

	e.balances["UNISWAP"+e.token0] = token0Uniswap.ToExact()
	e.balances["UNISWAP"+e.token1] = token1Uniswap.ToExact()

	e.logger.Debugf("Uniswap balances: %s %s, %s %s, %s %s", ethBalance.ToExact(), ethBalance.Currency.Symbol(), token0Uniswap.ToExact(), token0Uniswap.Currency.Symbol(), token1Uniswap.ToExact(), token1Uniswap.Currency.Symbol())

	token0Kucoin, err := e.kucoinClient.BalanceOf(e.token0)
	if err != nil {
		e.logger.WithError(err).Error("Failed to get KuCoin balance of token0")
		return
	}

	token1Kucoin, err := e.kucoinClient.BalanceOf(e.token1)
	if err != nil {
		e.logger.WithError(err).Error("Failed to get KuCoin balance of token1")
		return
	}

	e.balances["KUCOIN"+e.token0] = strconv.FormatFloat(token0Kucoin, 'f', -1, 64)
	e.balances["KUCOIN"+e.token1] = strconv.FormatFloat(token1Kucoin, 'f', -1, 64)

	e.logger.Debugf("KuCoin balances: %.18f ELON, %.18f USDT", token0Kucoin, token1Kucoin)
}

// ExecuteArbitrage executes an arbitrage trade
func (e *Executor) ExecuteArbitrage() {
	e.logger.Info("Executing arbitrage trade")
	e.GetBalances()

	kucoinPrice, err := e.kucoinClient.GetPrice()
	if err != nil {
		e.logger.WithError(err).Error("Failed to get KuCoin price")
		return
	}
	uniswapPrice, err := e.uniswapClient.GetPrice()
	if err != nil {
		e.logger.WithError(err).Error("Failed to get Uniswap price")
		return
	}

	// Calculate the average price
	avgPrice := (kucoinPrice + uniswapPrice) / 2

	e.logger.Infof("KuCoin price: %.18f, Uniswap price: %.18f, Average price: %.18f", kucoinPrice, uniswapPrice, avgPrice)

	// Do we buy or sell?
	if uniswapPrice < kucoinPrice {
		// Buy on Uniswap, Sell on KuCoin
		buyAmount, err := e.uniswapClient.GetBuyAmount(avgPrice)
		if err != nil {
			e.logger.WithError(err).Error("Failed to get buy amount")
			return
		}

		e.logger.Infof("Buy %s %s on Uniswap and Sell them on Kucoin", buyAmount.ToExact(), buyAmount.Currency.Symbol())

		err = e.uniswapClient.Trade(buyAmount, e.paperTrading)
		if err != nil {
			e.logger.WithError(err).Error("Failed to trade on Uniswap")
			return
		}

		err = e.kucoinClient.Trade("sell", buyAmount.Currency.Symbol(), buyAmount.ToFixed(2), kucoinPrice, e.paperTrading)
		if err != nil {
			e.logger.WithError(err).Error("Failed to trade on KuCoin")
			return
		}

	} else {
		// Sell on Uniswap, Buy on KuCoin
		sellAmount, err := e.uniswapClient.GetSellAmount(avgPrice)
		if err != nil {
			e.logger.WithError(err).Error("Failed to get sell amount")
			return
		}

		e.logger.Infof("Sell %s %s on Uniswap and Buy them on Kucoin", sellAmount.ToExact(), sellAmount.Currency.Symbol())
		err = e.uniswapClient.Trade(sellAmount, e.paperTrading)
		if err != nil {
			e.logger.WithError(err).Error("Failed to trade on Uniswap")
			return
		}

		err = e.kucoinClient.Trade("buy", sellAmount.Currency.Symbol(), sellAmount.ToFixed(2), kucoinPrice, e.paperTrading)
		if err != nil {
			e.logger.WithError(err).Error("Failed to trade on KuCoin")
			return
		}
	}

	e.logger.Debug("Trade executed successfully")
}

// Close closes the Executor
func (e *Executor) Close() {
	e.logger.Debug("Closing service")
	// Nothing to do here
}
