package main

import (
	"fmt"

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

	kraken := providers.NewKrakenProvider(datach, "XBT/USD")
	kraken.Start()

	for data := range datach {
		fmt.Println(data)
	}

	return

	coinbase := providers.NewCoinbaseProvider(datach, "BTC-USD", "ETH-USD", "DOGE-USD", "ADA-USD")
	coinbase.Start()

	binance := providers.NewBinanceOrderbooks(datach, "BTCUSDT", "ETHUSDT", "DOGEUSDT", "ADAUSDT")
	binance.Start()

	//b := orderbook.NewBinanceOrderbooks(symbols...)
	//b.Start()

	//go func() {
	//	for {
	//		time.Sleep(1 * time.Second)
	//		for _, book := range b.Orderbooks {
	//			fmt.Println(book.Asks.Best())
	//		}
	//	}
	//}()

	select {}
}
