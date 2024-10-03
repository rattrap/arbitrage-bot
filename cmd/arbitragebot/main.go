package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"rattrap/arbitrage-bot/internal/arbitrage"
	"rattrap/arbitrage-bot/internal/execution"
	"rattrap/arbitrage-bot/internal/kucoin"
	"rattrap/arbitrage-bot/internal/logging"
	"rattrap/arbitrage-bot/internal/pricing"
	"rattrap/arbitrage-bot/internal/telegram"
	"rattrap/arbitrage-bot/internal/uniswap"
	"syscall"
)

var (
	paperTrading bool
	logLevel     string
)

func init() {
	flag.BoolVar(&paperTrading, "paper", false, "Enable paper trading mode")
	flag.StringVar(&logLevel, "logLevel", "debug", "Log level (debug, info, warn, error, fatal, panic)")
	flag.Parse()
}

func main() {
	logger := logging.MakeLogger(logLevel)

	config, err := LoadConfig()
	if err != nil {
		logger.WithError(err).Fatal("Failed to load configuration")
	}

	if paperTrading {
		logger.Info("Bot starting in PAPER TRADING MODE...")
	} else {
		logger.Info("Bot starting in LIVE TRADING MODE...")
	}

	logger.Debugf("Loaded configuration %+v", config)

	ctx, cancel := context.WithCancel(context.Background())

	// Initialize Telegram service
	telegramService := telegram.NewTelegramService(config.TelegramBotToken, config.TelegramChannelID)
	err = telegramService.SendMessage("Arbitrage bot started")
	if err != nil {
		logger.WithError(err).Fatal("Failed to send message to Telegram")
	}

	// Initialize Uniswap client
	err, uniswapClient := uniswap.NewUniswapClient(config.TradingPair, config.EthereumRPCURL, config.UniswapPoolAddress, config.UniswapTickLensAddress, config.EthereumPrivateKey, ctx)
	if err != nil {
		logger.WithError(err).Fatal("Failed to initialize Uniswap client")
	}

	// Initialize KuCoin API client
	err, kucoinClient := kucoin.NewKucoinClient(config.TradingPair, config.KucoinAPIKey, config.KucoinAPISecret, config.KucoinAPIPassphrase, ctx)
	if err != nil {
		logger.WithError(err).Fatal("Failed to initialize KuCoin client")
	}

	// Start fetching prices
	priceService := pricing.NewPricingService(uniswapClient, kucoinClient, logger)
	priceService.Start()

	// Start the execution service
	execution := execution.NewExecutor(paperTrading, config.TradingPair, uniswapClient, kucoinClient, logger)
	execution.Start()

	// Start arbitrage detection and execution loop
	arbitrageService := arbitrage.NewArbitrageService(priceService, execution, telegramService, logger)

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
