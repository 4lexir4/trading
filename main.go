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
}

func main() {
	datach := make(chan orderbook.DataFeed, 1014)
	pvrs := []orderbook.Provider{
		providers.NewKrakenProvider(datach, "XBT/USD"),
		providers.NewCoinbaseProvider(datach, "BTC-USD"),
		providers.NewBinanceOrderbooks(datach, "BTCUSDT"),
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
