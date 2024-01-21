package providers

import (
	"github.com/4lexir4/trading/orderbook"

	ws "github.com/aopoltorzhicky/go_kraken/websocket"
)

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

func (p *KrakenProvider) GetOrderbooks() orderbook.Orderbooks {
	return p.Orderbooks
}

func (p *KrakenProvider) Start() error {
	kraken := ws.NewKraken(ws.ProdBaseURL)
	if err := kraken.Connect(); err != nil {
		return err
	}
	return nil
}
