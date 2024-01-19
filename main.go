package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/4lexir4/trading/orderbook"
	"github.com/adshao/go-binance/v2"
)

func main() {
	asks := orderbook.NewLimits(false)

	handler := func(event *binance.WsDepthEvent) {
		//fmt.Println(event.Symbol)
		for _, ask := range event.Asks {
			price, _ := strconv.ParseFloat(ask.Price, 64)
			size, _ := strconv.ParseFloat(ask.Quantity, 64)
			asks.Update(price, size)
		}
		fmt.Println("best ask:", asks.Best())
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
