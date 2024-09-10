package bybit_test

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/lianyun0502/exchange_conn/v2/bybit"
	"github.com/lianyun0502/exchange_conn/v2/common"
	. "github.com/lianyun0502/exchange_conn/v2/http_client"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var apiKey = "L7ksyiOdEgqg0gwIbf"
var secretKey = "0CVhyQmkwUDKWLcAP6NhtH7jB0P8XqSIVxE1"

var logger = &logrus.Logger{
	Out: os.Stderr,
	Formatter: &logrus.TextFormatter{
		ForceColors:     true,
		TimestampFormat: "2006-01-02 15:04:05.000",
		FullTimestamp:   true,
	},
	Level: logrus.DebugLevel,
	Hooks: make(logrus.LevelHooks),
}

func TestMarketTime(t *testing.T) {
	client := bybit.NewSpotClient(apiKey, secretKey)
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
	client := bybit.NewSpotClient(apiKey, secretKey)
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
	client := bybit.NewSpotClient(apiKey, secretKey)
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
