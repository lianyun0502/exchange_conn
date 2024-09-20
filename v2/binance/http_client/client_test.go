package binance_test

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/lianyun0502/exchange_conn/v2/binance/enums"
	"github.com/lianyun0502/exchange_conn/v2/binance/http_client"
	"github.com/lianyun0502/exchange_conn/v2/common"
	. "github.com/lianyun0502/exchange_conn/v2/http_client"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var apiKey = "GNEImPtijyYkCYRBx7WM3bntSPydYYY8bwqFuYA6BVVUNFKQ1cVTi7AHIdCsBQ3Z"
var secretKey = "5lkEjNgcicUQ7YZt1L0CAaTc12UWRuhLrUsahle1iQxIXMk95lwVsLEXRaQfkPQr"

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

func TestBinancePing(t *testing.T) {
	client, _ := binance.NewAPISpotClient(apiKey, secretKey, binance.IsTestNet())
	client.Log = logger
	data, err := client.Request(http.MethodGet, "/api/v3/ping").Send()
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, string(data), "{}")
	j := new(interface{})
	json.Unmarshal(data, &j)
	t.Log(common.PrettyPrint(j))
}

func TestBinanceGetInfo(t *testing.T) {
	client, _ := binance.NewAPISpotClient(apiKey, secretKey, binance.IsTestNet())
	client.Log = logger
	data, err := client.Request(http.MethodPost, "/api/v3/exchangeInfo").Send()
	if err != nil {
		t.Error(err)
		return
	}

	assert.NotEqual(t, string(data), "{}")
	j := new(interface{})
	json.Unmarshal(data, &j)
	t.Log(common.PrettyPrint(j))
}

func TestBinanceOrderBook(t *testing.T) {
	client, _ := binance.NewAPISpotClient(apiKey, secretKey, binance.IsTestNet())
	client.Log = logger
	req := client.Request(http.MethodGet, "/api/v3/depth")
	query := QueryMap{"symbol": "BTCUSDT", "limit": "10"}
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
func TestBinanceOrder(t *testing.T) {
	client, _ := binance.NewAPISpotClient(apiKey, secretKey, binance.IsTestNet())
	client.Log = logger
	req := client.Request(http.MethodPost, "/api/v3/order", binance.SetSercurityType(true, true))
	param := ParamMap{
		"symbol":      "BTCUSDT",
		"side":        enums.Buy,
		"type":        enums.Limit,
		"timeInForce": enums.GTC,
		"quantity":    0.0001,
		"price":       "50000",
	}
	req.SetParam(param)
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

func TestBinanceAccountInfo(t *testing.T) {
	client, _ := binance.NewAPISpotClient(apiKey, secretKey, binance.IsTestNet())
	client.Log = logger
	req := client.Request(http.MethodGet, "/api/v3/account", binance.SetSercurityType(true, true))
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
