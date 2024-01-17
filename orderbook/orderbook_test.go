package orderbook

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewOrderbook(t *testing.T) {
	ob := NewOrderbook("BTCUSDT")
	fmt.Printf("%+v\n", ob)
}

func TestLimitFillMultiOrder(t *testing.T) {
	l := NewLimit(10_000)
	askOrderA := NewAskOrder(10)
	askOrderB := NewAskOrder(5)
	l.addOrder(askOrderA)
	l.addOrder(askOrderB)

	marketOrder := NewAskOrder(10)
	l.fillOrder(marketOrder)
	assert.Equal(t, 5.0, l.totalVolume)
	assert.Equal(t, 1, len(l.orders))
	assert.Equal(t, 5.0, l.orders[0].size)
	assert.True(t, marketOrder.isFilled())
}

func TestLimitFillSingleOrder(t *testing.T) {
	l := NewLimit(50_000)
	orderSize := 10.0
	askOrder := NewAskOrder(orderSize)
	l.addOrder(askOrder)

	marketOrderSize := 5.0
	marketOrder := NewAskOrder(marketOrderSize)
	l.fillOrder(marketOrder)
	assert.True(t, marketOrder.isFilled())
	assert.Equal(t, orderSize-marketOrderSize, askOrder.size)
	assert.Equal(t, askOrder.size, l.totalVolume)
}

func TestLimitDeleteOrder(t *testing.T) {
	l := NewLimit(20_000)

	o1 := NewBidOrder(1.0)
	o2 := NewBidOrder(2.0)
	o3 := NewBidOrder(3.0)

	l.addOrder(o1)
	l.addOrder(o2)
	l.addOrder(o3)

	assert.Equal(t, 6.0, l.totalVolume)

	l.deleteOrder(o2)
	assert.Equal(t, 4.0, l.totalVolume)
}

func TestLimitAddOrder(t *testing.T) {
	l := NewLimit(16_000)
	n := 10
	size := 50.0

	for i := 0; i < n; i++ {
		o := NewAskOrder(size)
		l.addOrder(o)
		assert.Equal(t, i, o.limitIndex)
	}

	assert.Equal(t, n, len(l.orders))
	assert.Equal(t, float64(n)*size, l.totalVolume)
}
