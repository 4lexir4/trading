package main

import (
	"fmt"
	"log"

	"github.com/adshao/go-binance/v2"
)

func main() {
	// asks:= orderbook.NewLimits(false)
	handler := func(event *binance.WsDepthEvent) {
		fmt.Println(event.Symbol)
		for _, ask := range event.Asks {
			fmt.Printf("[%s - %s]\n", ask.Price, ask.Quantity)
		}
		// fmt.Printf("%+v\n", event)
	}
	errHandler := func(err error) {
		fmt.Println(err)
	}
	_, _, err := binance.WsDepthServe("btcusdt", handler, errHandler)
	if err != nil {
		log.Fatal(err)
	}

	select {}
}
