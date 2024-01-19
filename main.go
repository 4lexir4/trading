package main

import (
	"log"

	"github.com/4lexir4/trading/providers"
	"github.com/gorilla/websocket"
)

var symbols = []string{
	"BTCUSDT",
	"ETHUSDT",
	"ATOMUSDT",
}

func main() {
	cb := providers.NewCoinbaseProvider("BTC-USD")
	cb.Start()

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
