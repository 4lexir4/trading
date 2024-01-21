package providers

import "github.com/4lexir4/trading/orderbook"

type KrakenProvider struct {
	Orderbooks orderbook.Orderbooks
	symbols    []string
	feedch     chan orderbook.DataFeed
}

func NewKrakenProvider(feedch chan orderbook.DataFeed, symbols []string) *KrakenProvider {
	books := orderbook.Orderbooks{}
	for _, symbol := range symbols {
		books[symbol] = orderbook.NewBook(symbol)
	}
	return &KrakenProvider{
		Orderbooks: books,
		symbols:    symbols,
		feedch:     feedch,
	}
}
