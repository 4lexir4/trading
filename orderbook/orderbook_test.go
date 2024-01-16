package orderbook

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLimitFillOrder(t *testing.T) {
	l := NewLimit(50_000)
	askOrder := NewAskOrder(10)
	l.addOrder(askOrder)

	marketOrder := NewAskOrder(5)
	l.fillOrder(marketOrder)
	assert.True(t, marketOrder.isFilled())
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
