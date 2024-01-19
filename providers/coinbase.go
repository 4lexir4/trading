package providers

import (
	"log"

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

func (c *CoinbaseProvider) Start() error {
	ws, _, err := websocket.DefaultDialer.Dial("wss://ws-feed.exchange.coingase.com", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer ws.Close()

	ws.WriteJSON(CoinbaseMessage{
		Type:       "subscribe",
		ProductIds: c.symbols,
		Channels:   []string{"level2"},
	})

	go func() {
		for {
			_, message, err := ws.ReadMeassage()
			if err != nil {
				break
			}
			log.Printf("recv: %s", message)
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
