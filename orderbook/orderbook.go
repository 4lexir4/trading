package orderbook

import (
	"encoding/gob"
	"math/rand"
	"os"
	"time"

	"github.com/VictorLowther/btree"
)

func getBidByPrice(price float64) btree.CompareAgainst[*Limit] {
	return func(l *Limit) int {
		switch {
		case l.price > price:
			return -1
		case l.price < price:
			return 1
		default:
			return 0
		}
	}
}

func getAskByPrice(price float64) btree.CompareAgainst[*Limit] {
	return func(l *Limit) int {
		switch {
		case l.price < price:
			return -1
		case l.price > price:
			return 1
		default:
			return 0
		}
	}
}

func sortByBestBid(a, b *Limit) bool {
	return a.price > b.price
}

func sortByBestAsk(a, b *Limit) bool {
	return a.price < b.price
}

//type LimitMap struct {
//	isBids      bool
//	limits      map[float64]*Limit
//	totalVolume float64
//}

//func NewLimitMap(isBids bool) *LimitMap {
//	return &LimitMap{
//		isBids: isBids,
//		limits: make(map[float64]*Limit),
//	}
//}

func (l *Limits) loadFromFile(src string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	var data map[float64]float64

	if err := gob.NewDecoder(f).Decode(&data); err != nil {
		return err
	}

	for price, size := range data {
		limit := NewLimit(price)
		limit.totalVolume = size

		l.data.Insert(limit)
		l.totalVolume += size
	}

	return nil
}

type Limits struct {
	isBids      bool
	data        *btree.Tree[*Limit]
	totalVolume float64
}

func NewLimits(isBids bool) *Limits {
	f := sortByBestAsk
	if isBids {
		f = sortByBestBid
	}
	return &Limits{
		isBids: isBids,
		data:   btree.New(f),
	}
}

func (l *Limits) addOrder(price float64, o *Order) {
	if o.isBid != l.isBids {
		panic("the side of the limits does not match the side of the order")
	}

	f := getAskByPrice(price)
	if l.isBids {
		f = getBidByPrice(price)
	}

	var (
		limit *Limit
		ok    bool
	)

	limit, ok = l.data.Get(f)
	if !ok {
		limit = NewLimit(price)
		l.data.Insert(limit)
	}

	l.totalVolume += o.size
	limit.addOrder(o)
}

type Orderbook struct {
	pair string
	asks *Limits
	bids *Limits
}

func NewOrderbookFromFile(pair, askSrc, bidSrc string) (*Orderbook, error) {
	asks := NewLimits(false)
	if err := asks.loadFromFile(askSrc); err != nil {
		return nil, err
	}

	bids := NewLimits(true)
	if err := bids.loadFromFile(bidSrc); err != nil {
		return nil, err
	}

	return &Orderbook{
		pair: pair,
		asks: asks,
		bids: bids,
	}, nil
}

func NewOrderbook(pair string) *Orderbook {
	return &Orderbook{
		pair: pair,
		bids: NewLimits(true),
		asks: NewLimits(false),
	}
}

func (ob *Orderbook) placeLimitOrder(price float64, o *Order) {
	if o.isBid {
		ob.bids.addOrder(price, o)
	} else {
		ob.asks.addOrder(price, o)
	}
}

func (ob *Orderbook) bestBid() *Limit {
	iter := ob.bids.data.Iterator(nil, nil)
	iter.Next()
	return iter.Item()
}

func (ob *Orderbook) bestAsk() *Limit {
	iter := ob.asks.data.Iterator(nil, nil)
	iter.Next()
	return iter.Item()
}

func (ob *Orderbook) totalAskVolume() float64 {
	return ob.asks.totalVolume
}

func (ob *Orderbook) totalBidVolume() float64 {
	return ob.bids.totalVolume
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
