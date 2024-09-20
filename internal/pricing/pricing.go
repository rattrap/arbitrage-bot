package pricing

import (
	"rattrap/arbitrage-bot/internal/kucoin"
	"rattrap/arbitrage-bot/internal/logging"
	"rattrap/arbitrage-bot/internal/uniswap"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// PricingService is a struct to manage pricing from multiple sources
type PricingService struct {
	kucoinClient  *kucoin.KucoinClient
	uniswapClient *uniswap.UniswapClient
	logger        *logrus.Entry
	stopChan      chan struct{}
	lock          sync.RWMutex
	kucoinPrice   float64
	uniswapPrice  float64
}

// NewPricingService initializes a new PricingService
func NewPricingService(kucoinClient *kucoin.KucoinClient, uniswapClient *uniswap.UniswapClient, logger *logging.Logger) *PricingService {
	prefixedLogger := logger.WithField("prefix", "pricing")
	prefixedLogger.Info("Initializing service")

	return &PricingService{
		kucoinClient:  kucoinClient,
		uniswapClient: uniswapClient,
		logger:        prefixedLogger,
		stopChan:      make(chan struct{}),
		lock:          sync.RWMutex{},
		kucoinPrice:   0,
		uniswapPrice:  0,
	}
}

// FetchPrices fetches prices from KuCoin and Uniswap
func (ps *PricingService) FetchPrices() {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	// Fetch prices from KuCoin
	kucoinPrice, err := ps.kucoinClient.GetPrice("ELON-USDT")
	if err != nil {
		ps.logger.WithError(err).Error("Failed to get KuCoin price")
	}

	// Fetch prices from Uniswap
	uniswapPrice, err := ps.uniswapClient.GetPrice("ELON")
	if err != nil {
		ps.logger.WithError(err).Error("Failed to get Uniswap price")
	}

	// Store the prices
	ps.kucoinPrice = kucoinPrice
	ps.uniswapPrice = uniswapPrice
}

// Start starts the PricingService
func (ps *PricingService) Start() {
	ps.logger.Info("Starting service")
	go (func() {
		for {
			select {
			case <-ps.stopChan:
				ps.logger.Info("Stopping service")
				return
			default:
				ps.FetchPrices()
				time.Sleep(10 * time.Second)
			}
		}
	})()
}

// GetKucoinPrice returns the current price of a trading pair from KuCoin
func (ps *PricingService) GetKucoinPrice(tradingPair string) float64 {
	ps.lock.Lock()
	defer ps.lock.Unlock()
	return ps.kucoinPrice
}

// GetUniswapPrice returns the current price of a token on Uniswap
func (ps *PricingService) GetUniswapPrice(token string) float64 {
	ps.lock.Lock()
	defer ps.lock.Unlock()
	return ps.uniswapPrice
}

// Close closes the PricingService
func (ps *PricingService) Close() {
	ps.logger.Info("Closing service")
	close(ps.stopChan)
}
