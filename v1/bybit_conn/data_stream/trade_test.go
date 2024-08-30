package data_stream_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lianyun0502/exchange_conn/v1/bybit_conn/data_stream"
	"github.com/lianyun0502/exchange_conn/v1/common"
)

func TestTradeData(t *testing.T) {
	rawData := []byte(`
	{"topic":"publicTrade.BTCUSDT","type":"snapshot","ts":1724048334483,
	"data":[
	{"T":1724048334477,"s":"BTCUSDT","S":"Sell","v":"0.018","p":"58450.00","L":"ZeroMinusTick","i":"21b325fa-234e-5d46-8663-df511eb76e8e","BT":false},
	{"T":1724048334477,"s":"BTCUSDT","S":"Sell","v":"0.004","p":"58450.00","L":"ZeroMinusTick","i":"830b035b-4a8c-5589-af40-b613112b4a73","BT":false},
	{"T":1724048334477,"s":"BTCUSDT","S":"Sell","v":"0.010","p":"58450.00","L":"ZeroMinusTick","i":"753e0939-ac3d-5470-8255-da5ed1082493","BT":false},
	{"T":1724048334478,"s":"BTCUSDT","S":"Buy","v":"0.067","p":"58450.00","L":"ZeroMinusTick","i":"989b6607-46ad-564f-9276-af24f2dfd6d5","BT":false}
	]}`)
	trade := data_stream.NewTrade()
	data, _ := trade.Update(rawData)
	assert.Equal(t, data.Trades[2].Topic, "publicTrade.BTCUSDT")
	assert.Equal(t, data.Trades[2].Time, int64(1724048334483))
	assert.Equal(t, data.Trades[2].Symbol, "BTCUSDT")
	assert.Equal(t, data.Trades[2].TradeId, "753e0939-ac3d-5470-8255-da5ed1082493")
	assert.Equal(t, data.Trades[2].TradeTime, int64(1724048334477))
	assert.Equal(t, data.Trades[2].Price, "58450.00")
	assert.Equal(t, data.Trades[2].Quantity, "0.010")
	assert.Equal(t, data.Trades[2].Side, "SELL")

	// s, err := json.Marshal(data)
	// if err != nil {
	// 	t.Error(err)
	// }
	t.Log(common.PrettyPrint(data))

}
