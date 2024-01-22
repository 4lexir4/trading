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

	if err := kraken.SubscribeBook(p.symbols, 1000); err != nil {
		return err
	}

	go func() {
		for {
			update := <-kraken.Listen()
			switch data := update.Data.(type) {
			case ws.OrderBookUpdate:
				book := p.Orderbooks[update.Pair]
				for _, ask := range data.Asks {
					if !ask.Republish {
						price, _ := ask.Price.Float64()
						size, _ := ask.Volume.Float64()
						book.Asks.Update(price, size)
					}
				}
				for _, bid := range data.Bids {
					if !bid.Republish {
						price, _ := bid.Price.Float64()
						size, _ := bid.Volume.Float64()
						book.Bids.Update(price, size)
					}
				}

				spread := book.Spread()
				bestAsk := book.BestAsk()
				bestBid := book.BestBid()

				if bestAsk == nil || bestBid == nil {
					continue
				}

				p.feedch <- orderbook.DataFeed{
					Provider: "Kraken",
					Symbol:   book.Symbol,
					BestAsk:  bestAsk.Price,
					BestBid:  bestBid.Price,
					Spread:   spread,
				}
			}
		}
	}()
	return nil
}
