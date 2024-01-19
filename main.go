package main

import (
	//"fmt"
	"log"
	//"time"

	//"github.com/4lexir4/trading/orderbook"
	"github.com/gorilla/websocket"
)

var symbols = []string{
	"BTCUSDT",
	"ETHUSDT",
	"ATOMUSDT",
}

type ChannelInfo struct {
	Name       string   `json:"name"`
	ProductIds []string `json:"prodcuct_ids"`
}

type CoinbaseMessage struct {
	Type       string        `json:"type"`
	ProductIds []string      `json:"product_ids"`
	Channels   []ChannelInfo `json:"channels"`
}

func main() {
	c, _, err := websocket.DefaultDialer.Dial("wss://ws-feed.exchange.coingase.com", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	productIds := []string{"ETF-USD"}
	c.WriteJSON(CoinbaseMessage{
		Type:       "subscribe",
		ProductIds: productIds,
		Channels: []ChannelInfo{
			{
				Name:       "ticker",
				ProductIds: productIds,
			},
		},
	})
	for {
		_, message, err := c.ReadMeassage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Printf("recv: %s", message)
	}

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

//{
//  "type": "subscribe",
//  "product_ids": [
//    "ETH-USD",
//    "ETH-EUR"
//  ],
//  "channels": [
//    "level2",
//    "heartbeat",
//    {
//      "name": "ticker",
//      "product_ids": [
//        "ETH-BTC",
//        "ETH-USD"
//      ]
//    }
//  ]
//}

//func main() {
//	asks := orderbook.NewLimits(false)
//	bids := orderbook.NewLimits(true)
//
//	handler := func(event *binance.WsDepthEvent) {
//		for _, ask := range event.Asks {
//			price, _ := strconv.ParseFloat(ask.Price, 64)
//			size, _ := strconv.ParseFloat(ask.Quantity, 64)
//			asks.Update(price, size)
//		}
//		for _, bid := range event.Bids {
//			price, _ := strconv.ParseFloat(bid.Price, 64)
//			size, _ := strconv.ParseFloat(bid.Quantity, 64)
//			bids.Update(price, size)
//		}
//		fmt.Printf("ask [%.1f] [%.1f] bid\n", asks.Best().Price, bids.Best().Price)
//	}
//	errHandler := func(err error) {
//		fmt.Println(err)
//	}
//	_, _, err := binance.WsDepthServe100Ms("btcusdt", handler, errHandler)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	select {}
//}
