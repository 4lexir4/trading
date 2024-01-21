package providers

import "github.com/4lexir4/trading/orderbook"

type BinanceProvider struct {
	Orderbooks orderbook.Orderbooks
	symbols    []string
	feedch     chan orderbook.DataFeed
}
