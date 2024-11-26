package bybit_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/lianyun0502/exchange_conn/v2/bybit/http_client"
	"github.com/lianyun0502/exchange_conn/v2/common"
	. "github.com/lianyun0502/exchange_conn/v2/http_client"
	"github.com/stretchr/testify/assert"
)

func TestMarketTime(t *testing.T) {
	client := bybit.NewSpotClient(apiKey, secretKey, bybit.IsTestNet())
	client.Log = logger
	data, err := client.Request(http.MethodGet, "/v5/market/time").Send()
	if err != nil {
		t.Error(err)
		return
	}

	assert.NotEqual(t, string(data), "{}")
	j := new(interface{})
	json.Unmarshal(data, &j)
	t.Log(common.PrettyPrint(j))
}

func TestBybitOrder(t *testing.T) {
	client := bybit.NewSpotClient(apiKey, secretKey, bybit.IsTestNet())
	client.Log = logger
	param := ParamMap{
		"category":  "spot",
		"symbol":    "BTCUSDT",
		"side":      "Buy",
		"orderType": "Limit",
		"qty":       "0.001",
		"price":     "50000",
	}
	data, err := client.Request(http.MethodPost, "/v5/order/create", bybit.SetSercurityType(true, true)).SetParam(param).Send()
	if err != nil {
		t.Error(err)
		return
	}

	assert.NotEqual(t, string(data), "{}")
	j := new(interface{})
	json.Unmarshal(data, &j)
	t.Log(common.PrettyPrint(j))
}

func TestBybitVolatility(t *testing.T) {
	client := bybit.NewSpotClient(apiKey, secretKey, bybit.IsTestNet())
	client.Log = logger
	req := client.Request(http.MethodGet, "/v5/market/historical-volatility")
	query := QueryMap{
		"category": "option",
		"baseCoin": "BTC",
	}
	req.SetQuery(query)
	data, err := req.Send()
	if err != nil {
		t.Error(err)
		return
	}

	assert.NotEqual(t, string(data), "{}")
	j := new(interface{})
	json.Unmarshal(data, &j)
	t.Log(common.PrettyPrint(j))
}


func BenchmarkOrder(b *testing.B) {
	client := bybit.NewSpotClient(apiKey, secretKey, bybit.IsTestNet())
	client.Log = logger
	// logger.SetLevel(logrus.ErrorLevel)
	param := ParamMap{
		"category":  "spot",
		"symbol":    "BTCUSDT",
		"side":      "Buy",
		"orderType": "Limit",
		"qty":       "0.001",
		"price":     "50000",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.Request(http.MethodPost, "/v5/order/create", bybit.SetSercurityType(true, true)).SetParam(param).Send()
	}
	
}