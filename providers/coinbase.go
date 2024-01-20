package providers

import (
	"encoding/json"
	"fmt"

	"github.com/4lexir4/trading/orderbook"
	"github.com/gorilla/websocket"
)

type CoinbaseProvider struct {
	Orderbooks orderbook.Orderbooks
	symbols    []string
}

func NewCoinbaseProvider(symbols ...string) *CoinbaseProvider {
	books := orderbook.Orderbooks{}
	for _, symbol := range symbols {
		books[symbol] = orderbook.NewBook(symbol)
	}
	return &CoinbaseProvider{
		Orderbooks: books,
		symbols:    symbols,
	}
}

func (c *CoinbaseProvider) handleUpdate() error {
	return nil
}

func (c *CoinbaseProvider) handleSnapshot() error {
	return nil
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
