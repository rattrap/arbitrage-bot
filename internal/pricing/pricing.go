package pricing

import (
	"rattrap/arbitrage-bot/internal/kucoin"
	"rattrap/arbitrage-bot/internal/logging"
	"rattrap/arbitrage-bot/internal/uniswap"
	"sync"

	"github.com/sirupsen/logrus"
)

// PricingService is a struct to manage pricing from multiple sources
type PricingService struct {
	uniswapClient *uniswap.UniswapClient
	kucoinClient  *kucoin.KucoinClient
	logger        *logrus.Entry
	stopChan      chan struct{}
	lock          sync.RWMutex
	uniswapPrice  float64
	kucoinPrice   float64
}

// NewPricingService initializes a new PricingService
func NewPricingService(uniswapClient *uniswap.UniswapClient, kucoinClient *kucoin.KucoinClient, logger *logging.Logger) *PricingService {
	prefixedLogger := logger.WithField("prefix", "pricing")
	prefixedLogger.Debug("Initializing service")

	return &PricingService{
		uniswapClient: uniswapClient,
		kucoinClient:  kucoinClient,
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

	ps.logger.Debug("Fetching prices")

	// Fetch prices from Uniswap
	uniswapPrice, err := ps.uniswapClient.GetPrice()
	if err != nil {
		ps.logger.WithError(err).Error("Failed to get Uniswap price")
	}

	// Fetch prices from KuCoin
	kucoinPrice, err := ps.kucoinClient.GetPrice("ELON-USDT")
	if err != nil {
		ps.logger.WithError(err).Error("Failed to get KuCoin price")
	}

	// Store the prices
	ps.uniswapPrice = uniswapPrice
	ps.kucoinPrice = kucoinPrice
}

// Start starts the PricingService
func (ps *PricingService) Start() {
	ps.logger.Debug("Starting service")
}

// GetUniswapPrice returns the current price of a token on Uniswap
func (ps *PricingService) GetUniswapPrice() float64 {
	ps.lock.Lock()
	defer ps.lock.Unlock()
	return ps.uniswapPrice
}

// GetKucoinPrice returns the current price of a trading pair from KuCoin
func (ps *PricingService) GetKucoinPrice(tradingPair string) float64 {
	ps.lock.Lock()
	defer ps.lock.Unlock()
	return ps.kucoinPrice
}

// Close closes the PricingService
func (ps *PricingService) Close() {
	ps.logger.Debug("Closing service")
	close(ps.stopChan)
}
