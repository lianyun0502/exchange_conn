package exchange_conn_test

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lianyun0502/exchange_conn/v1"
	"github.com/lianyun0502/exchange_conn/v1/binance_conn"
	"github.com/lianyun0502/exchange_conn/v1/binance_conn/enums"
	"github.com/lianyun0502/exchange_conn/v1/bybit_conn"
	"github.com/lianyun0502/exchange_conn/v1/common"

	"github.com/sirupsen/logrus"
)

var logger = &logrus.Logger{
	Out: os.Stderr,
	Formatter: &logrus.TextFormatter{
		ForceColors:     true,
		TimestampFormat: "2006-01-02 15:04:05.000",
		FullTimestamp:   true,
	},
	Hooks:        make(logrus.LevelHooks),
	Level:        logrus.InfoLevel,
	ExitFunc:     os.Exit,
	ReportCaller: false,
}

func TestBinancePing(t *testing.T) {
	agent := exchange_conn.NewAgent(binance_conn.NewClient("YourAPIKey", "YourSecretKey", "https://api.binance.com"))
	agent.Client.Logger = logger
	agent.Client.Logger.SetLevel(logrus.DebugLevel)

	data, err := agent.Request(http.MethodGet, "/api/v3/ping", false, false).Send()
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, string(data), "{}")
}

func TestBinanceGetInfo(t *testing.T) {
	agent := exchange_conn.NewAgent(binance_conn.NewClient("YourAPIKey", "YourSecretKey", "https://api.binance.com"))
	agent.Client.Logger = logger

	data, err := agent.Request(http.MethodGet, "/api/v3/exchangeInfo", false, false).Send()
	if err != nil {
		t.Error(err)
		return
	}

	assert.NotEqual(t, string(data), "{}")
	j := new(interface{})
	json.Unmarshal(data, &j)
	logger.Println(common.PrettyPrint(j))
}

func TestBinanceOrderBook(t *testing.T) {
	agent := exchange_conn.NewAgent(binance_conn.NewClient("YourAPIKey", "YourSecretKey", "https://api.binance.com"))
	agent.Client.Logger = logger

	data, err := agent.Request(http.MethodGet, "/api/v3/depth", false, false).SetQuery("symbol", "BTCUSDT").SetQuery("limit", "10").Send()
	if err != nil {
		t.Error(err)
		return
	}

	assert.NotEqual(t, string(data), "{}")
	j := new(interface{})
	json.Unmarshal(data, &j)
	logger.Println(common.PrettyPrint(j))

}

func TestBinanceOrder(t *testing.T) {
	apiKey := "xTz5GK0rSyPANKeQTR5o1zohOdF7DmPRGR2ReAFKNLj0GjoIacB2Ld5Sjzd2p8Wk"
	secretKey := "Hvsqtth66iAyXw7lnbzQGdw0ZCLPru5MWZPllLbcAuHpGMPNiuWoxXAE6LjpKqNg"
	agent := exchange_conn.NewAgent(binance_conn.NewClient(apiKey, secretKey, "https://api.binance.com"))
	agent.Client.Logger = logger

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
		logger.Println(data)
		t.Error(err)
		return
	}

	assert.NotEqual(t, string(data), "{}")
	j := new(interface{})
	json.Unmarshal(data, &j)
	logger.Println(common.PrettyPrint(j))

}

func TestBinanceTestNewOrder(t *testing.T) {
	apiKey := "xTz5GK0rSyPANKeQTR5o1zohOdF7DmPRGR2ReAFKNLj0GjoIacB2Ld5Sjzd2p8Wk"
	secretKey := "Hvsqtth66iAyXw7lnbzQGdw0ZCLPru5MWZPllLbcAuHpGMPNiuWoxXAE6LjpKqNg"
	agent := exchange_conn.NewAgent(binance_conn.NewClient(apiKey, secretKey, "https://api.binance.com"))
	agent.Client.Logger = logger

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
		logger.Println(data)
		t.Error(err)
		return
	}

	assert.NotEqual(t, string(data), "{}")
	j := new(interface{})
	json.Unmarshal(data, &j)
	logger.Println(common.PrettyPrint(j))

}

func TestBinanceAccountInfo(t *testing.T) {
	apiKey := "xTz5GK0rSyPANKeQTR5o1zohOdF7DmPRGR2ReAFKNLj0GjoIacB2Ld5Sjzd2p8Wk"
	secretKey := "Hvsqtth66iAyXw7lnbzQGdw0ZCLPru5MWZPllLbcAuHpGMPNiuWoxXAE6LjpKqNg"
	agent := exchange_conn.NewAgent(binance_conn.NewClient(apiKey, secretKey, "https://api.binance.com"))
	agent.Client.Logger = logger

	data, err := agent.Request(http.MethodGet, "/api/v3/account", true, true).Send()
	if err != nil {
		t.Error(err)
		return
	}

	assert.NotEqual(t, string(data), "{}")
	j := new(interface{})
	json.Unmarshal(data, &j)
	logger.Println(common.PrettyPrint(j))
}

func TestBybitMarketTime(t *testing.T) {
	apiKey := "L7ksyiOdEgqg0gwIbf"
	secretKey := "0CVhyQmkwUDKWLcAP6NhtH7jB0P8XqSIVxE1"
	agent := exchange_conn.NewAgent(bybit_conn.NewClient(apiKey, secretKey, "https://api-testnet.bybit.com"))
	agent.Client.Logger = logger

	data, err := agent.Request(http.MethodGet, "/v5/market/time", false, false).Send()
	if err != nil {
		t.Error(err)
		return
	}

	assert.NotEqual(t, string(data), "{}")
	j := new(interface{})
	json.Unmarshal(data, &j)
	logger.Println(common.PrettyPrint(j))

}

func TestBybitVolatility(t *testing.T) {
	apiKey := "L7ksyiOdEgqg0gwIbf"
	secretKey := "0CVhyQmkwUDKWLcAP6NhtH7jB0P8XqSIVxE1"
	agent := exchange_conn.NewAgent(bybit_conn.NewClient(apiKey, secretKey, "https://api-testnet.bybit.com"))
	agent.Client.Logger = logger

	req := agent.Request(http.MethodGet, "/v5/market/historical-volatility", false, false)
	req.SetQueries(map[string]any{
		"category": "option",
		"baseCoin": "BTC",
	},
	)
	data, err := req.Send()
	if err != nil {
		t.Error(err)
		return
	}

	assert.NotEqual(t, string(data), "{}")
	j := new(interface{})
	json.Unmarshal(data, &j)
	logger.Println(common.PrettyPrint(j))

}

func TestBybitOrder(t *testing.T) {
	apiKey := "L7ksyiOdEgqg0gwIbf"
	secretKey := "0CVhyQmkwUDKWLcAP6NhtH7jB0P8XqSIVxE1"
	agent := exchange_conn.NewAgent(bybit_conn.NewClient(apiKey, secretKey, "https://api-testnet.bybit.com"))
	agent.Client.Logger = logger

	req := agent.Request(http.MethodPost, "/v5/order/create", true, true)
	req.SetParams(map[string]any{
		"category":  "spot",
		"symbol":    "BTCUSDT",
		"side":      "Buy",
		"orderType": "Limit",
		"qty":       "0.001",
		"price":     "50000",
	})

	data, err := req.Send()
	if err != nil {
		t.Error(err)
		return
	}

	assert.NotEqual(t, string(data), "{}")
	j := new(interface{})
	json.Unmarshal(data, &j)
	logger.Println(common.PrettyPrint(j))
}
