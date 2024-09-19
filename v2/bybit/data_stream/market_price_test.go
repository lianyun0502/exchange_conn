package data_stream_test

import (
	"encoding/json"
	"github.com/lianyun0502/exchange_conn/v2/bybit/data_stream"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUpdateMarketPrice(t *testing.T) {
	rawData := `{
    "topic": "tickers.BTCUSDT",
    "type": "snapshot",
    "data": {
        "symbol": "BTCUSDT",
        "tickDirection": "PlusTick",
        "price24hPcnt": "0.017103",
        "lastPrice": "17216.00",
        "prevPrice24h": "16926.50",
        "highPrice24h": "17281.50",
        "lowPrice24h": "16915.00",
        "prevPrice1h": "17238.00",
        "markPrice": "17217.33",
        "indexPrice": "17227.36",
        "openInterest": "68744.761",
        "openInterestValue": "1183601235.91",
        "turnover24h": "1570383121.943499",
        "volume24h": "91705.276",
        "nextFundingTime": "1673280000000",
        "fundingRate": "-0.000212",
        "bid1Price": "17215.50",
        "bid1Size": "84.489",
        "ask1Price": "17216.00",
        "ask1Size": "83.020"
    },
    "cs": 24987956059,
    "ts": 1673272861686
}`
	// closure return a function
	// UpdateMarketPrice := data_stream.MarketPrice()

	market_data := data_stream.NewMarketData()

	data, err := market_data.Update([]byte(rawData))
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "tickers.BTCUSDT", data.Topic)
	assert.Equal(t, int64(1673272861686), data.Time)
	assert.Equal(t, "BTCUSDT", data.Symbol)
	assert.Equal(t, "17217.33", data.MarketPrice)
	assert.Equal(t, "17227.36", data.IndexPrice)
	assert.Equal(t, "-0.000212", data.FundingRate)
	assert.Equal(t, int64(1673280000000), data.NextFundingTime)

	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Error(err)
		t.Log(string(jsonData))
	}
	rawData = `{"topic": "tickers.BTCUSDT","type": "delta","data": {"symbol": "BTCUSDT","markPrice": "17217.00","indexPrice": "17227.00"},"cs": 24987956059,"ts": 1673272861700}`
	data, err = market_data.Update([]byte(rawData))
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "tickers.BTCUSDT", data.Topic)
	assert.Equal(t, int64(1673272861700), data.Time)
	assert.Equal(t, "BTCUSDT", data.Symbol)
	assert.Equal(t, "17217.00", data.MarketPrice)
	assert.Equal(t, "17227.00", data.IndexPrice)
	assert.Equal(t, "-0.000212", data.FundingRate)
	assert.Equal(t, int64(1673280000000), data.NextFundingTime)
}
