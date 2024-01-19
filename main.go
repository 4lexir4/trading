package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/4lexir4/trading/orderbook"
	"github.com/adshao/go-binance/v2"
)

var symbols = []string{
	"BTCUSDT",
	"ETHUSDT",
	"ATOMUSDT",
}

func main() {
	asks := orderbook.NewLimits(false)
	bids := orderbook.NewLimits(true)

	handler := func(event *binance.WsDepthEvent) {
		for _, ask := range event.Asks {
			price, _ := strconv.ParseFloat(ask.Price, 64)
			size, _ := strconv.ParseFloat(ask.Quantity, 64)
			asks.Update(price, size)
		}
		for _, bid := range event.Bids {
			price, _ := strconv.ParseFloat(bid.Price, 64)
			size, _ := strconv.ParseFloat(bid.Quantity, 64)
			bids.Update(price, size)
		}
		fmt.Printf("ask [%.1f] [%.1f] bid\n", asks.Best().Price, bids.Best().Price)
	}
	errHandler := func(err error) {
		fmt.Println(err)
	}
	_, _, err := binance.WsDepthServe100Ms("btcusdt", handler, errHandler)
	if err != nil {
		log.Fatal(err)
	}

	select {}
}
