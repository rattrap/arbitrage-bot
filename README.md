# Arbitrage bot

Keep Uniswap V3 pool balanced by trading between it and Kucoin

## Configure

Copy `.env.dist` to `.env` and add your own values.

### Required ENV vars

```
KUCOIN_API_KEY=<API_KEY>
KUCOIN_API_SECRET=<API_SECRET>
KUCOIN_API_PASSPHRASE=<API_PASSPHRASE>
ETHEREUM_RPC_URL=<RPC_URL>
ETHEREUM_PRIVATE_KEY=<PRIVATE_KEY>
UNISWAP_POOL_ADDRESS=<UNISWAP_POOL_ADDRESS>
UNISWAP_TICKLENS_ADDRESS=<UNISWAP_TICKLENS_ADDRESS>
TRADING_PAIR=TOKEN0-TOKEN1
```

### Optional ENV vars

```
TELEGRAM_CHANNEL_ID=<CHANNEL_ID>
TELEGRAM_BOT_TOKEN=<BOT_TOKEN>
```

## Run

```bash
go run ./cmd/arbitragebot
```

## Build

```bash
make build
```

You can now run the binary `./build/arbitragebot`

## Docker

```bash
docker build -t arbitrage-bot:latest
docker run --env-file .env arbitrage-bot:latest
```

## Paper trading mode

If the paper trading mode is enabled, no real transactions will be made on either Uniswap or Kucoin

```bash
./build/arbitragebot --paper
```

## Log Level

Can be one of: debug, info, warn, error, fatal, panic

```bash
./build/arbitragebot --logLevel=info
```