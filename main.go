package main

import (
	//"fmt"
	"log"
	"time"

	"github.com/4lexir4/trading/orderbook"
	"github.com/4lexir4/trading/providers"
	"github.com/4lexir4/trading/socket"
	//"github.com/4lexir4/trading/util"
)

var symbols = []string{
	"BTCUSD",
	"ETHUSD",
	"ADAUSD",
}

var pairs = map[string]map[string]string{
	"ADAUSD": {
		"Binance":  "ADAUSDT",
		"Kraken":   "ADA/USD",
		"Coinbase": "ADA-USD",
	},
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

func getSymbolForProvider(p string, symbol string) string {
	return pairs[symbol][p]
}

func mapSymbolsFor(provider string) []string {
	out := make([]string, len(symbols))
	for i, symbol := range symbols {
		out[i] = pairs[symbol][provider]
	}
	return out
}

func main() {
	datach := make(chan orderbook.DataFeed, 1024)
	pvrs := []orderbook.Provider{
		providers.NewKrakenProvider(datach, mapSymbolsFor("Kraken")),
		providers.NewCoinbaseProvider(datach, mapSymbolsFor("Coinbase")),
		providers.NewBinanceOrderbooks(datach, mapSymbolsFor("Binance")),
	}

	for _, provider := range pvrs {
		if err := provider.Start(); err != nil {
			log.Fatal(err)
		}
	}

	crossSpreadch := make(chan map[string][]orderbook.CrossSpread, 1024)
	go func() {
		ticker := time.NewTicker(time.Microsecond * 100)
		for {
			calcCrossSpreads(crossSpreadch, pvrs)
			<-ticker.C
		}
	}()

	socketServer := socket.NewServer(crossSpreadch)
	socketServer.Start()
}

func calcCrossSpreads(datach chan map[string][]orderbook.CrossSpread, pvrs []orderbook.Provider) {
	data := map[string][]orderbook.CrossSpread{}

	for _, symbol := range symbols {
		crossSpreads := []orderbook.CrossSpread{}
		for i := 0; i < len(pvrs); i++ {
			a := pvrs[i]
			var b orderbook.Provider
			if len(pvrs)-1 == i {
				b = pvrs[0]
			} else {
				b = pvrs[i+1]
			}

			bookA := a.GetOrderbooks()[getSymbolForProvider(a.Name(), symbol)]
			bookB := b.GetOrderbooks()[getSymbolForProvider(b.Name(), symbol)]

			best := orderbook.CrossSpread{
				Symbol: symbol,
			}

			bestBidA := bookA.BestBid()
			bestBidB := bookB.BestBid()
			if bestBidA == nil || bestBidB == nil {
				continue
			}

			bestAsk := orderbook.BestPrice{}
			bestBid := orderbook.BestPrice{}
			if bestBidA.Price < bestBidB.Price {
				bestAsk.Provider = a.Name()
				bestBid.Provider = b.Name()
				bestBid.Price = bestBidA.Price
				bestAsk.Price = bookB.BestAsk().Price
			} else {
				best.A = b.Name()
				best.B = a.Name()
				best.BestBid = bestBidB.Price
				best.BestAsk = bookA.BestAsk().Price
			}

			best.Spread = best.BestAsk - best.BestBid //util.Round(best.BestAsk-best.BestBid, 10000)
			crossSpreads = append(crossSpreads, best)

		}
		data[symbol] = crossSpreads
	}
	datach <- data
}

func minMaxLimit(a, b *orderbook.Limit) (*orderbook.Limit, *orderbook.Limit) {
	if a.Price <= b.Price {
		return a, b
	}
	return b, a
}
