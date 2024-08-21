package data_stream_test

import (
	"github.com/lianyun0502/exchange_conn/v1/binance_conn/data_stream"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToNormalTradeData(t *testing.T) {
	rawData := `{"e":"trade","E":1633775480000,"s":"BTCUSDT","t":12345,"T":1633775480000,"p":"60000.00","q":"0.001","m":true}`
	data, err := data_stream.ToNormalTradeData([]byte(rawData))
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, data.Topic, "trade")
	assert.Equal(t, data.Time, int64(1633775480000))
	assert.Equal(t, data.Symbol, "BTCUSDT")
	assert.Equal(t, data.TradeId, "12345")
	assert.Equal(t, data.TradeTime, int64(1633775480000))
	assert.Equal(t, data.Price, "60000.00")
	assert.Equal(t, data.Quantity, "0.001")
	assert.Equal(t, data.Side, "sell")

}

func TestToNormalAggregateTradeData(t *testing.T) {
	rawData := `{"e":"aggTrade","E":1723023583269,"s":"BTCUSDT","a":3106590246,"p":"57564.01000000","q":"0.00014000","f":3731576184,"l":3731576184,"T":1723023583268,"m":false}`
	data, err := data_stream.ToNormalAggregateTradeData([]byte(rawData))
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, data.Topic, "aggTrade")
	assert.Equal(t, data.Time, int64(1723023583269))
	assert.Equal(t, data.Symbol, "BTCUSDT")
	assert.Equal(t, data.TradeId, "3106590246")
	assert.Equal(t, data.TradeTime, int64(1723023583268))
	assert.Equal(t, data.Price, "57564.01000000")
	assert.Equal(t, data.Quantity, "0.00014000")
	assert.Equal(t, data.Side, "buy")
}
