package tradingview

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetQuote(t *testing.T) {
	c := New()

	v, err := c.GetQuote(UsdTicker)
	assert.Nil(t, err, "GetQuote error")
	assert.NotEqual(t, 0, v)
}
