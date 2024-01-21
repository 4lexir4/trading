package providers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/4lexir4/trading/orderbook"
	"github.com/gorilla/websocket"
)

type CoinbaseProvider struct {
	Orderbooks orderbook.Orderbooks
	symbols    []string
	feedch     chan orderbook.DataFeed
}

func NewCoinbaseProvider(feedch chan orderbook.DataFeed, symbols []string) *CoinbaseProvider {
	books := orderbook.Orderbooks{}
	for _, symbol := range symbols {
		books[symbol] = orderbook.NewBook(symbol)
	}
	return &CoinbaseProvider{
		Orderbooks: books,
		symbols:    symbols,
		feedch:     feedch,
	}
}

func (c *CoinbaseProvider) handleUpdate(symbol string, changes []SnapshotChange) error {
	for _, change := range changes {
		side, price, size := parseSnapShotChange(change)
		if side == "sell" {
			c.Orderbooks[symbol].Asks.Update(price, size)
		} else {
			c.Orderbooks[symbol].Bids.Update(price, size)
		}
	}
	return nil
}

func (c *CoinbaseProvider) handleSnapshot(symbol string, asks []SnapshotEntry, bids []SnapshotEntry) error {
	for _, entry := range asks {
		price, size := parseSnapShotEntry(entry)
		c.Orderbooks[symbol].Asks.Update(price, size)
	}
	for _, entry := range bids {
		price, size := parseSnapShotEntry(entry)
		c.Orderbooks[symbol].Bids.Update(price, size)
	}
	return nil
}

func (c *CoinbaseProvider) feedLoop() {
	time.Sleep(time.Second * 2)
	ticker := time.NewTicker(100 * time.Millisecond)
	for {
		for _, book := range c.Orderbooks {
			spread := book.Spread()
			bestAsk := book.BestAsk().Price
			bestBid := book.BestBid().Price
			c.feedch <- orderbook.DataFeed{
				Provider: "Coinbase",
				Symbol:   book.Symbol,
				BestAsk:  bestAsk,
				BestBid:  bestBid,
				Spread:   spread,
			}
		}
		<-ticker.C
	}
}

func (c *CoinbaseProvider) Start() error {
	ws, _, err := websocket.DefaultDialer.Dial("wss://ws-feed.exchange.coingase.com", nil)
	if err != nil {
		return err
	}

	ws.WriteJSON(CoinbaseMessage{
		Type:       "subscribe",
		ProductIds: c.symbols,
		Channels:   []string{"level2"},
	})

	go func() {
		for {
			_, message, err := ws.ReadMeassage()
			if err != nil {
				fmt.Println(err)
				break
			}
			msg := CoinbaseMessageResponse{}
			if err := json.Unmarshal(message, &msg); err != nil {
				fmt.Println(err)
				break
			}
			if msg.Type == "l2update" {
				continue
			}
			fmt.Println(msg)
		}
	}()

	return nil
}

type CoinabaseChannelInfo struct {
	Name       string   `json:"name"`
	ProductIds []string `json:"prodcuct_ids"`
}

type CoinbaseMessage struct {
	Type       string   `json:"type"`
	ProductIds []string `json:"product_ids"`
	Channels   []string `json:"channels"`
}

type CoinbaseMessageResponse struct {
	Type      string `json:"type"`
	ProductID string `json:"product_id"`
}
