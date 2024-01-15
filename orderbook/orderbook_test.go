package orderbook

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLimitDeleteOrder(t *testing.T) {
	l := NewLimit(20_000)
	n := 100
	size := 0.5
	for i := 0; i < n; i++ {
		o := NewAskOrder(size)
		l.addOrder(o)
		assert.Equal(t, o.size, l.totalVolume)
		assert.Equal(t, 1, len(l.orders))
		l.deleteOrder(o)
		assert.Equal(t, 0, len(l.orders))
		assert.Equal(t, 0.0, l.totalVolume)
	}
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
