package bybit_test

import (
	"testing"
	"net/http"
	"encoding/json"

	"github.com/lianyun0502/exchange_conn/v2/common"
	"github.com/stretchr/testify/assert"
	"github.com/lianyun0502/exchange_conn/v2/bybit"
	. "github.com/lianyun0502/exchange_conn/v2/http_client"
)

var	apiKey = "L7ksyiOdEgqg0gwIbf"
var	secretKey = "0CVhyQmkwUDKWLcAP6NhtH7jB0P8XqSIVxE1"

func TestMarketTime(t *testing.T) {
	client := bybit.NewSpotClient(apiKey, secretKey)
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
	client := bybit.NewSpotClient(apiKey, secretKey)
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