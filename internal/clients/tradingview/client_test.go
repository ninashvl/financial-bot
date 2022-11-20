package tradingview

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetQuote(t *testing.T) {
	c := New()

	v, err := c.GetQuote(context.TODO(), UsdTicker)
	assert.Nil(t, err, "GetQuote error")
	assert.NotEqual(t, 0, v)
}
