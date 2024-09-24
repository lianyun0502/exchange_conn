package data_stream_test

import (
	"github.com/lianyun0502/exchange_conn/v2/binance/data_stream"
	"github.com/stretchr/testify/assert"
	"testing"
)


func TestUpdateMarketPrice(t *testing.T) {
	rawData := `{"e":"markPriceUpdate","E":1562305380000,"s":"BTCUSDT","T":1562306400000,"p":"11794.15000000","i":"11784.62659091","r":"0.00038167"}`
	data, err := data_stream.UpdateMarketPrice([]byte(rawData))
	if err != nil {
		t.Error(err)
	}
	
	assert.Equal(t, data.Topic, "markPriceUpdate")
	assert.Equal(t, data.Time, int64(1562305380000))
	assert.Equal(t, data.Symbol, "BTCUSDT")
	assert.Equal(t, data.MarketPrice, "11794.15000000")
	assert.Equal(t, data.IndexPrice, "11784.62659091")
	assert.Equal(t, data.FundingRate, "0.00038167")
	assert.Equal(t, data.NextFundingTime, int64(1562306400000))
}