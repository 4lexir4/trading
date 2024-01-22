package providers

import (
	"fmt"
	"strconv"

	"github.com/4lexir4/trading/orderbook"
	"github.com/adshao/go-binance/v2"
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

func (b *BinanceProvider) Start() error {
	handler := func(event *binance.WsDepthEvent) {
		for _, ask := range event.Asks {
			price, _ := strconv.ParseFloat(ask.Price, 64)
			size, _ := strconv.ParseFloat(ask.Quantity, 64)
			b.Orderbooks[event.Symbol].Asks.Update(price, size)
		}
		for _, bid := range event.Bids {
			price, _ := strconv.ParseFloat(bid.Price, 64)
			size, _ := strconv.ParseFloat(bid.Quantity, 64)
			b.Orderbooks[event.Symbol].Bids.Update(price, size)
		}

		var (
			book    = b.Orderbooks[event.Symbol]
			spread  = book.Spread()
			bestAsk = book.BestAsk()
			bestBid = book.BestBid()
		)

		b.feedch <- orderbook.DataFeed{
			Provider: "Binance",
			Symbol:   book.Symbol,
			BestAsk:  bestAsk.Price,
			BestBid:  bestBid.Price,
			Spread:   spread,
		}

	}
	errHandler := func(err error) {
		fmt.Println(err)
	}

	_, _, err := binance.WsCombinedDepthServe100Ms(b.symbols, handler, errHandler)
	return err
}
