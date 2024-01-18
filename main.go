package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/4lexir4/trading/orderbook"
	"github.com/adshao/go-binance/v2"
)

func main() {
	asks := orderbook.NewLimits(false)

	ticker := time.NewTicker(time.Second * 3)
	go func() {
		for {
			<-ticker.C
			fmt.Println("best ask:", asks.Best())
		}
	}()

	handler := func(event *binance.WsDepthEvent) {
		//fmt.Println(event.Symbol)
		for _, ask := range event.Asks {
			//fmt.Printf("[%s - %s]\n", ask.Price, ask.Quantity)
			price, _ := strconv.ParseFloat(ask.Price, 64)
			size, _ := strconv.ParseFloat(ask.Quantity, 64)
			asks.Update(price, size)
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
