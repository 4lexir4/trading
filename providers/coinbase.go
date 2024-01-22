package providers

import (
	"encoding/json"
	"fmt"
	"strconv"
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
			_, message, err := ws.ReadMessage()
			if err != nil {
				fmt.Println(err)
				break
			}
			msg := Message{}
			if err := json.Unmarshal(message, &msg); err != nil {
				fmt.Println(err)
				break
			}
			if msg.Type == "l2update" {
				c.handleUpdate(msg.ProductID, msg.Changes)
			}
			if msg.Type == "snapshot" {
				c.handleSnapshot(msg.ProductID, msg.Asks, msg.Bids)
			}
		}
	}()

	go c.feedLoop()

	return nil
}

func parseSnapShotChange(change SnapshotChange) (string, float64, float64) {
	// in this case its either "buy" or "sell"
	side := change[0]
	price, _ := strconv.ParseFloat(change[1], 64)
	size, _ := strconv.ParseFloat(change[2], 64)
	return side, price, size
}

func parseSnapShotEntry(entry SnapshotEntry) (float64, float64) {
	price, _ := strconv.ParseFloat(entry[0], 64)
	size, _ := strconv.ParseFloat(entry[1], 64)
	return price, size
}

type CoinbaseMessage struct {
	Type       string   `json:"type"`
	ProductIds []string `json:"product_ids"`
	Channels   []string `json:"channels"`
}

type Message struct {
	Type       string   `json:"type"`
	ProductID  string   `json:"product_id"`
	ProductIds []string `json:"product_ids"`
	//Products      []Product        `json:"products"`
	//Currencies    []Currency       `json:"currencies"`
	TradeID      int    `json:"trade_id,number"`
	OrderID      string `json:"order_id"`
	ClientOID    string `json:"client_oid"`
	Sequence     int64  `json:"sequence,number"`
	MakerOrderID string `json:"maker_order_id"`
	TakerOrderID string `json:"taker_order_id"`
	//Time          Time             `json:"time,string"`
	RemainingSize string           `json:"remaining_size"`
	NewSize       string           `json:"new_size"`
	OldSize       string           `json:"old_size"`
	Size          string           `json:"size"`
	Price         string           `json:"price"`
	Side          string           `json:"side"`
	Reason        string           `json:"reason"`
	OrderType     string           `json:"order_type"`
	Funds         string           `json:"funds"`
	NewFunds      string           `json:"new_funds"`
	OldFunds      string           `json:"old_funds"`
	Message       string           `json:"message"`
	Bids          []SnapshotEntry  `json:"bids,omitempty"`
	Asks          []SnapshotEntry  `json:"asks,omitempty"`
	Changes       []SnapshotChange `json:"changes,omitempty"`
	LastSize      string           `json:"last_size"`
	BestBid       string           `json:"best_bid"`
	BestAsk       string           `json:"best_ask"`
	Channels      []MessageChannel `json:"channels"`
	UserID        string           `json:"user_id"`
	ProfileID     string           `json:"profile_id"`
	LastTradeID   int              `json:"last_trade_id"`
}

type MessageChannel struct {
	Name       string   `json:"name"`
	ProductIds []string `json:"product_ids"`
}

type SnapshotChange [3]string

type SnapshotEntry [2]string
