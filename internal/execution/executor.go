package execution

import (
	"rattrap/arbitrage-bot/internal/kucoin"
	"rattrap/arbitrage-bot/internal/logging"
	"rattrap/arbitrage-bot/internal/uniswap"

	"github.com/sirupsen/logrus"
)

// Executor handles trade execution for both KuCoin and Uniswap
type Executor struct {
	kucoinClient  *kucoin.KucoinClient
	uniswapClient *uniswap.UniswapClient
	logger        *logrus.Entry
}

// NewExecutor initializes a new Executor
func NewExecutor(kucoinClient *kucoin.KucoinClient, uniswapClient *uniswap.UniswapClient, logger *logging.Logger) *Executor {
	prefixedLogger := logger.WithField("prefix", "execution")
	prefixedLogger.Debug("Initializing")
	return &Executor{
		kucoinClient:  kucoinClient,
		uniswapClient: uniswapClient,
		logger:        prefixedLogger,
	}
}

// ExecuteArbitrage executes an arbitrage trade
func (e *Executor) ExecuteArbitrage() {
	e.logger.Info("Executing arbitrage trade")
	// Nothing to do here
}

// Close closes the Executor
func (e *Executor) Close() {
	e.logger.Debug("Closing service")
	// Nothing to do here
}
