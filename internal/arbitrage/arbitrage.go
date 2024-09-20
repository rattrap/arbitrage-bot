package arbitrage

import (
	"math"
	"rattrap/arbitrage-bot/internal/execution"
	"rattrap/arbitrage-bot/internal/logging"
	"rattrap/arbitrage-bot/internal/pricing"
	"time"

	"github.com/sirupsen/logrus"
)

// ArbitrageService is a struct to manage arbitrage opportunities
type ArbitrageService struct {
	pricingService *pricing.PricingService
	executor       *execution.Executor
	logger         *logrus.Entry
	stopChan       chan struct{}
}

// NewArbitrageService initializes a new ArbitrageService
func NewArbitrageService(pricingService *pricing.PricingService, executor *execution.Executor, logger *logging.Logger) *ArbitrageService {
	contextLogger := logger.WithField("service", "arbitrage")
	contextLogger.Info("Initializing service")
	return &ArbitrageService{
		pricingService: pricingService,
		executor:       executor,
		logger:         contextLogger,
		stopChan:       make(chan struct{}),
	}
}

func (a *ArbitrageService) RunArbitrageLoop() {
	go func() {
		for {
			select {
			case <-a.stopChan:
				a.logger.Info("Stopping arbitrage loop")
				return
			default:
				a.logger.Info("Checking for arbitrage opportunities...")
				kucoinPrice := a.pricingService.GetKucoinPrice("ELON-USDT")
				uniswapPrice := a.pricingService.GetUniswapPrice("ELON")

				// Calculate the difference between the two prices
				priceDifference := kucoinPrice - uniswapPrice
				priceDifferencePercentage := (priceDifference / uniswapPrice) * 100

				a.logger.Infof("KuCoin price: %.18f, Uniswap price: %.18f, Price difference: %.18f (%.2f%%)", kucoinPrice, uniswapPrice, priceDifference, priceDifferencePercentage)

				if math.Abs(priceDifferencePercentage) > 1 {
					a.logger.Info("Arbitrage opportunity found")
					a.executor.ExecuteArbitrage()
				}

				time.Sleep(10 * time.Second)
			}
		}
	}()
}

// Close closes the ArbitrageService
func (a *ArbitrageService) Close() {
	a.logger.Info("Closing service")
	close(a.stopChan)
	a.executor.Close()
}
