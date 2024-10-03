package arbitrage

import (
	"fmt"
	"math"
	"rattrap/arbitrage-bot/internal/execution"
	"rattrap/arbitrage-bot/internal/logging"
	"rattrap/arbitrage-bot/internal/pricing"
	"rattrap/arbitrage-bot/internal/telegram"
	"time"

	"github.com/sirupsen/logrus"
)

// ArbitrageService is a struct to manage arbitrage opportunities
type ArbitrageService struct {
	pricingService *pricing.PricingService
	executor       *execution.Executor
	telegram       *telegram.TelegramService
	logger         *logrus.Entry
	stopChan       chan struct{}
}

// NewArbitrageService initializes a new ArbitrageService
func NewArbitrageService(pricingService *pricing.PricingService, executor *execution.Executor, telegramService *telegram.TelegramService, logger *logging.Logger) *ArbitrageService {
	prefixedLogger := logger.WithField("prefix", "arbitrage")
	prefixedLogger.Debug("Starting service")
	return &ArbitrageService{
		pricingService: pricingService,
		executor:       executor,
		telegram:       telegramService,
		logger:         prefixedLogger,
		stopChan:       make(chan struct{}),
	}
}

func (a *ArbitrageService) RunArbitrageLoop() {
	go func() {
		for {
			select {
			case <-a.stopChan:
				a.logger.Debug("Stopping arbitrage loop")
				return
			default:
				a.pricingService.FetchPrices()
				a.logger.Debug("Checking for arbitrage opportunities...")
				uniswapPrice, kucoinPrice := a.pricingService.GetPrices()

				// Calculate the difference between the two prices
				priceDifference := kucoinPrice - uniswapPrice
				priceDifferencePercentage := (priceDifference / uniswapPrice) * 100

				stat := fmt.Sprintf("KuCoin price: %.18f, Uniswap price: %.18f, Price difference: %.18f (%.2f%%)", kucoinPrice, uniswapPrice, priceDifference, priceDifferencePercentage)

				a.logger.Info(stat)
				err := a.telegram.SendMessage(telegram.FormatMessage(stat))
				if err != nil {
					a.logger.WithError(err).Error("Failed to send message to Telegram")
				}

				if math.Abs(priceDifferencePercentage) > 1 {
					a.logger.Info("Arbitrage opportunity found")
					a.executor.ExecuteArbitrage()
				}

				time.Sleep(1 * time.Minute)
			}
		}
	}()
}

// Close closes the ArbitrageService
func (a *ArbitrageService) Close() {
	a.logger.Debug("Closing service")
	close(a.stopChan)
	a.executor.Close()
}
