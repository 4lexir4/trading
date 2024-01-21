package main

import (
	"fmt"
	"log"

	"github.com/4lexir4/trading/orderbook"
	"github.com/4lexir4/trading/providers"
)

var symbols = []string{
	"BTCUSDT",
	"ETHUSDT",
	"ATOMUSDT",
	"DOGEUSD",
}

var pairs = map[string]map[string]string{
	"BTCUSD": {
		"Binance":  "BTCUSDT",
		"Kraken":   "XBT/USD",
		"Coinbase": "BTC-USD",
	},
	"ETHUSD": {
		"Binance":  "ETHUSDT",
		"Kraken":   "ETH/USD",
		"Coinbase": "ETH-USD",
	},
}

func mapSymbolsFor(provider string, s string) []string {
	out := make([]string, len(symbols))
	for i, symbol := range symbols {
		out[i] = pairs[symbol][provider]
	}
	return out
}

func main() {
	datach := make(chan orderbook.DataFeed, 1014)
	pvrs := []orderbook.Provider{
		providers.NewKrakenProvider(datach, []string{"XBT/USD", "ETH/USD"}),
		providers.NewCoinbaseProvider(datach, []string{"BTC-USD", "ETH-USD"}),
		providers.NewBinanceOrderbooks(datach, []string{"BTCUSDT", "ETHUSDT"}),
	}

	for _, provider := range pvrs {
		if err := provider.Start(); err != nil {
			log.Fatal(err)
		}
	}

	for data := range datach {
		fmt.Println(data)
	}

	select {}
}
