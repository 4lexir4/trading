package providers

import (
	"time"

	"github.com/4lexir4/trading/orderbook"
)

type BinanceProvider struct {
	Orderbooks orderbook.Orderbooks
	symbols    []string
	feedch     chan orderbook.DataFeed
}

func NewBinanceOrderbooks(feedch chan orderbook.DataFeed, symbols []string) *BinanceProvider {
	books := orderbook.Orderbooks{}
	for _, symbol := range symbols {
		books[symbol] = orderbook.NewBook(symbol)
	}
	return &BinanceProvider{
		Orderbooks: books,
		symbols:    symbols,
		feedch:     feedch,
	}
}

func (b *BinanceProvider) feedLoop() {
	time.Sleep(time.Second * 2)
	ticker := time.NewTicker(100 * time.Millisecond)
	for {
		for _, book := range b.Orderbooks {

		}
	}
}
