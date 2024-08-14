package exchange_conn_test

import (
	"encoding/json"
	// "log"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lianyun0502/exchange_conn/v1"
	"github.com/lianyun0502/exchange_conn/v1/binance_conn"
	"github.com/lianyun0502/exchange_conn/v1/binance_conn/enums"
	"github.com/lianyun0502/exchange_conn/v1/common"

	log "github.com/sirupsen/logrus"
)

func TestBinancePing(t *testing.T) {
	agent := exchange_conn.NewAgent(binance_conn.NewClient("YourAPIKey", "YourSecretKey", "https://api.binance.com"))

	data, err := agent.Request(http.MethodGet, "/api/v3/ping", false, false).Send()
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, string(data), "{}")
}

func TestBinanceGetInfo(t *testing.T) {
	agent := exchange_conn.NewAgent(binance_conn.NewClient("YourAPIKey", "YourSecretKey", "https://api.binance.com"))

	data, err := agent.Request(http.MethodGet, "/api/v3/exchangeInfo", false, false).Send()
	if err != nil {
		t.Error(err)
		return
	}

	assert.NotEqual(t, string(data), "{}")
	j := new(interface{})
	json.Unmarshal(data, &j)
	fmt.Println(common.PrettyPrint(j))
	log.SetFormatter(&log.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05.000"})
	log.Info("Binance exchange info")
}

func TestBinanceOrderBook(t *testing.T) {
	agent := exchange_conn.NewAgent(binance_conn.NewClient("YourAPIKey", "YourSecretKey", "https://api.binance.com"))

	data, err := agent.Request(http.MethodGet, "/api/v3/depth", false, false).SetQuery("symbol", "BTCUSDT").SetQuery("limit", "10").Send()
	if err != nil {
		t.Error(err)
		return
	}

	assert.NotEqual(t, string(data), "{}")
	j := new(interface{})
	json.Unmarshal(data, &j)
	fmt.Println(common.PrettyPrint(j))

}

var (
	apiKey    = "xTz5GK0rSyPANKeQTR5o1zohOdF7DmPRGR2ReAFKNLj0GjoIacB2Ld5Sjzd2p8Wk"
	secretKey = "Hvsqtth66iAyXw7lnbzQGdw0ZCLPru5MWZPllLbcAuHpGMPNiuWoxXAE6LjpKqNg"
)

func TestBinanceOrder(t *testing.T) {
	agent := exchange_conn.NewAgent(binance_conn.NewClient(apiKey, secretKey, "https://api.binance.com"))

	req := agent.Request(http.MethodPost, "/api/v3/order", true, true)
	req.SetQueries(map[string]any{
		"symbol":      "BTCUSDT",
		"side":        enums.Buy,
		"type":        enums.Limit,
		"timeInForce": enums.GTC,
		"quantity":    0.0001,
		"price":       "50000",
	})
	data, err := req.Send()
	if err != nil {
		log.Println(data)
		t.Error(err)
		return
	}

	assert.NotEqual(t, string(data), "{}")
	j := new(interface{})
	json.Unmarshal(data, &j)
	fmt.Println(common.PrettyPrint(j))

}

func TestBinanceTestNewOrder(t *testing.T) {
	agent := exchange_conn.NewAgent(binance_conn.NewClient(apiKey, secretKey, "https://api.binance.com"))

	req := agent.Request(http.MethodPost, "/api/v3/order/test", true, true)
	req.SetParams(map[string]any{
		"symbol":                 "BTCUSDT",
		"side":                   "BUY",
		"type":                   "LIMIT",
		"timeInForce":            enums.GTC,
		"quantity":               0.0001,
		"price":                  "50000",
		"computeCommissionRates": true,
	})
	data, err := req.Send()
	if err != nil {
		log.Println(data)
		t.Error(err)
		return
	}

	assert.NotEqual(t, string(data), "{}")
	j := new(interface{})
	json.Unmarshal(data, &j)
	fmt.Println(common.PrettyPrint(j))

}

func TestBinanceAccountInfo(t *testing.T) {
	agent := exchange_conn.NewAgent(binance_conn.NewClient(apiKey, secretKey, "https://api.binance.com"))

	data, err := agent.Request(http.MethodGet, "/api/v3/account", true, true).Send()
	if err != nil {
		t.Error(err)
		return
	}

	assert.NotEqual(t, string(data), "{}")
	j := new(interface{})
	json.Unmarshal(data, &j)
	fmt.Println(string(data))
}
