package orderbook

import (
	"math/rand"
	"time"
)

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

type Order struct {
	id        int64
	size      float64
	timestamp int64
	isBid     bool
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
