package main

import (
	"context"
	"os"
	"os/signal"
	"rattrap/arbitrage-bot/internal/arbitrage"
	"rattrap/arbitrage-bot/internal/execution"
	"rattrap/arbitrage-bot/internal/kucoin"
	"rattrap/arbitrage-bot/internal/logging"
	"rattrap/arbitrage-bot/internal/pricing"
	"rattrap/arbitrage-bot/internal/uniswap"
	"syscall"
)

func main() {
	logger := logging.MakeLogger("debug")
	logger.Debug("Bot starting...")

	config, err := LoadConfig()
	if err != nil {
		logger.WithError(err).Fatal("Failed to load configuration")
	}

	logger.Debugf("Loaded configuration %+v", config)

	ctx, cancel := context.WithCancel(context.Background())

	// Initialize KuCoin API client
	err, kucoinClient := kucoin.NewKucoinClient(config.KucoinAPIKey, config.KucoinAPISecret, config.KucoinAPIPassphrase, ctx)
	if err != nil {
		logger.WithError(err).Fatal("Failed to initialize KuCoin client")
	}

	// Initialize Uniswap client
	err, uniswapClient := uniswap.NewUniswapClient(config.EthereumRPCURL, ctx)
	if err != nil {
		logger.WithError(err).Fatal("Failed to initialize Uniswap client")
	}

	// Start fetching prices
	priceService := pricing.NewPricingService(kucoinClient, uniswapClient, logger)
	priceService.Start()

	// Start arbitrage detection and execution loop
	arbitrageService := arbitrage.NewArbitrageService(priceService, execution.NewExecutor(kucoinClient, uniswapClient, logger), logger)

	// Run the arbitrage loop
	arbitrageService.RunArbitrageLoop()

	// Handle interrupt signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		logger.Debug("Received an interrupt, closing connections...")

		arbitrageService.Close()
		priceService.Close()
		kucoinClient.Close()
		uniswapClient.Close()

		cancel() // Cancel the context to stop any ongoing operations

		logger.Debug("Shutdown complete.")
		os.Exit(0)
	}()

	// Keep the main goroutine running until an interrupt is received
	<-ctx.Done()
}
