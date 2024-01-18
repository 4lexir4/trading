package orderbook

import (
	"encoding/gob"
	"math/rand"
	"os"
	"time"
)

type LimitMap struct {
	limits      map[float64]*Limit
	totalVolume float64
}

func NewLimitMap() *LimitMap {
	return &LimitMap{
		limits: make(map[float64]*Limit),
	}
}

func (m *LimitMap) loadFromFile(src string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	var data map[float64]float64

	if err := gob.NewDecoder(f).Decode(&data); err != nil {
		return err
	}

	for price, size := range data {
		l := NewLimit(price)
		l.totalVolume = size
		m.limits[price] = l
		m.totalVolume = size
	}

	return nil
}

type Orderbook struct {
	ticker string
	asks   *LimitMap
}

func NewOrderbookFromFile(ticker, askSrc, bidSrc string) (*Orderbook, error) {
	askMap := NewLimitMap()
	if err := askMap.loadFromFile(askSrc); err != nil {
		return nil, err
	}

	return &Orderbook{
		ticker: ticker,
	}, nil
}

func NewOrderbook(ticker string) *Orderbook {
	return &Orderbook{
		asks:   NewLimitMap(),
		ticker: ticker,
	}
}

type Limit struct {
	price       float64
	orders      []*Order
	totalVolume float64
}

func NewLimit(price float64) *Limit {
	return &Limit{
		price:  price,
		orders: []*Order{},
	}

}

func (l *Limit) fillOrder(marketOrder *Order) {
	ordersToDelete := []*Order{}
	for _, limitOrder := range l.orders {
		maxo, mino := maxMinOrder(limitOrder, marketOrder)
		sizeFilled := mino.size
		maxo.size -= sizeFilled
		l.totalVolume -= sizeFilled
		mino.size = 0.0

		if limitOrder.isFilled() {
			ordersToDelete = append(ordersToDelete, limitOrder)
		}

		if marketOrder.isFilled() {
			break
		}
	}

	for _, order := range ordersToDelete {
		l.deleteOrder(order)
	}
}

func (l *Limit) addOrder(o *Order) {
	o.limitIndex = len(l.orders)
	l.orders = append(l.orders, o)
	l.totalVolume += o.size
}

func (l *Limit) deleteOrder(o *Order) {
	l.orders[o.limitIndex] = l.orders[len(l.orders)-1]
	l.orders = l.orders[:len(l.orders)-1]
	if !o.isFilled() {
		l.totalVolume -= o.size
	}
}

type Order struct {
	id         int64
	size       float64
	timestamp  int64
	isBid      bool
	limitIndex int
}

func NewOrder(isBid bool, size float64) *Order {
	return &Order{
		id:        rand.Int63n(100_000), // TODO fix later...
		size:      size,
		timestamp: time.Now().UnixNano(),
		isBid:     isBid,
	}
}

func NewBidOrder(size float64) *Order {
	return NewOrder(true, size)
}

func NewAskOrder(size float64) *Order {
	return NewOrder(false, size)
}

func (o *Order) isFilled() bool {
	return o.size == 0
}

func maxMinOrder(a, b *Order) (*Order, *Order) {
	if a.size >= b.size {
		return a, b
	}
	return b, a
}
