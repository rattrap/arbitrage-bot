package execution

import (
	"rattrap/arbitrage-bot/internal/kucoin"
	"rattrap/arbitrage-bot/internal/logging"
	"rattrap/arbitrage-bot/internal/uniswap"

	"github.com/sirupsen/logrus"
)

// Executor handles trade execution for both KuCoin and Uniswap
type Executor struct {
	uniswapClient *uniswap.UniswapClient
	kucoinClient  *kucoin.KucoinClient
	logger        *logrus.Entry
}

// NewExecutor initializes a new Executor
func NewExecutor(uniswapClient *uniswap.UniswapClient, kucoinClient *kucoin.KucoinClient, logger *logging.Logger) *Executor {
	prefixedLogger := logger.WithField("prefix", "execution")
	prefixedLogger.Debug("Initializing")
	return &Executor{
		uniswapClient: uniswapClient,
		kucoinClient:  kucoinClient,
		logger:        prefixedLogger,
	}
}

// ExecuteArbitrage executes an arbitrage trade
func (e *Executor) ExecuteArbitrage() {
	e.logger.Info("Executing arbitrage trade")

	kucoinPrice, err := e.kucoinClient.GetPrice("ELON-USDT")
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
	} else {
		// Sell on Uniswap, Buy on KuCoin
		sellAmount, err := e.uniswapClient.GetSellAmount(avgPrice)
		if err != nil {
			e.logger.WithError(err).Error("Failed to get sell amount")
			return
		}

		e.logger.Infof("Sell %s %s on Uniswap and Buy them on Kucoin", sellAmount.ToExact(), sellAmount.Currency.Symbol())

	}

	token0Balance, token1Balance, err := e.uniswapClient.GetBalances()
	if err != nil {
		e.logger.WithError(err).Error("Failed to get balances")
		return
	}

	ethBalance, err := e.uniswapClient.GetEthBalance()
	if err != nil {
		e.logger.WithError(err).Error("Failed to get ETH balance")
		return
	}

	e.logger.Infof("Uniswap balances: %s %s, %s %s, %s %s", ethBalance.ToExact(), ethBalance.Currency.Symbol(), token0Balance.ToExact(), token0Balance.Currency.Symbol(), token1Balance.ToExact(), token1Balance.Currency.Symbol())

	token0BalanceKucoin, err := e.kucoinClient.BalanceOf("ELON")
	if err != nil {
		e.logger.WithError(err).Error("Failed to get KuCoin balance")
		return
	}

	token1BalanceKucoin, err := e.kucoinClient.BalanceOf("USDT")
	if err != nil {
		e.logger.WithError(err).Error("Failed to get KuCoin balance")
		return
	}

	e.logger.Infof("KuCoin balances: %.18f ELON, %.18f USDT", token0BalanceKucoin, token1BalanceKucoin)

}

// Close closes the Executor
func (e *Executor) Close() {
	e.logger.Debug("Closing service")
	// Nothing to do here
}
