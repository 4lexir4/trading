package orderbook

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLimitAddOrder(t *testing.T) {
	l := NewLimit(16_000)
	n := 1_000

	for i := 0; i < n; i++ {
		o := NewAskOrder(50.12)
		l.addOrder(o)
	}

	assert.Equal(t, n, len(l.orders))
}
