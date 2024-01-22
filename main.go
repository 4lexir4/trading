package main

import (
	"fmt"
	"log"
	"time"

	"github.com/4lexir4/trading/orderbook"
	"github.com/4lexir4/trading/providers"
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

	//ticker := time.NewTicker(time.Millisecond * 50)
	//go func() {
	//	for {
	//		for _, p := range pvrs {
	//			for _, book := range p.GetOrderbooks() {
	//				var (
	//					spread  = book.Spread()
	//					bestAsk = book.BestAsk()
	//					bestBid = book.BestBid()
	//				)
	//				if bestAsk == nil || bestBid == nil {
	//					continue
	//				}
	//				datach <- orderbook.DataFeed{
	//					Provider: p.Name(),
	//					Symbol:   book.Symbol,
	//					BestAsk:  bestAsk.Price,
	//					BestBid:  bestBid.Price,
	//					Spread:   spread,
	//				}
	//			}
	//		}
	//		<-ticker.C
	//	}
	//}()

	//for data := range datach {
	//	fmt.Printf(
	//		"[%s | %s] ASK %f %f BID [%f] \n",
	//		data.Provider,
	//		data.Symbol,
	//		data.BestAsk,
	//		data.BestBid,
	//		data.Spread,
	//	)
	//}

	type BestSpread struct {
		Symbol  string
		A       string
		B       string
		BestBid float64
		BestAsk float64
		Spread  float64
	}

  func calcBestSpreads(pvrs []orederbook.orderbook.Provider) {
    for i := 0; i < len(pvrs); i++ {
      a := pvrs[i]
      var b orderbook.Provider
      if len(pvrs) -1 == i{
        b = pvrs[0]
      } else {
        b = pvrs[i + 1]
      }

      for _, symbol := range symbols {
        bookA := a.GetOrderbooks()[getSymbolForProvider(a.Name(), symbol)]
        bookB := b.GetOrderbooks()[getSymbolForProvider(b.Name(), symbol)]

        best := BestSpread{
          Symbol: symbol,
        }

        bestBidA := bookA.BestBid()
        bestBidB := bookB.BestBid()
        if bestBidA == nil || bestBidB == nil {
          continue
        }

        if bestBidA.Price < bestBidB.Price {
          best.A = a.Name()
          best.B = b.Name()
          best.BestBid = bestBidA.Price
          best.BestAsk = bookB.BestAsk().Price
        } else {
          best.A = b.Name()
          best.B = a.Name()
          best.BestBid = bestBidB.Price
          best.BestAsk = bookA.BestAsk().Price
        }

        best.Spread = util.Round(best.BestAsk - best.BestBid, 10_000)

        fmt.Println(best)
      }
    }
  }

	select {}
}
