package bybit_test

import (
	"os"
	"testing"
	"time"

	"github.com/lianyun0502/exchange_conn/v2/bybit/ws_client"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)


var	apiKey = "L7ksyiOdEgqg0gwIbf"
var	secretKey = "0CVhyQmkwUDKWLcAP6NhtH7jB0P8XqSIVxE1"

var logger = &logrus.Logger{
	Out: os.Stderr,
	Formatter: &logrus.TextFormatter{
		ForceColors:     true,
		TimestampFormat: "2006-01-02 15:04:05.000000",
		FullTimestamp:   true,
	},
	Level: logrus.DebugLevel,
	Hooks: make(logrus.LevelHooks),
}
func TestBybitOrder(t *testing.T) {
	client := bybit.NewWsTradeClient(apiKey, secretKey, bybit.IsTestNet())
	client.Logger = logger
	client.Connect()
	go client.Conn.ReadLoop()
	resp, err := client.GetSignature()
	if err != nil {
		t.Log(string(resp))
		t.Error(err)
		return
	}
	type ParamMap map[string]string
	param := ParamMap{
		"category":  "spot",
		"symbol":    "BTCUSDT",
		"side":      "Buy",
		"orderType": "Limit",
		"qty":       "0.001",
		"price":     "50000",
	}
	resp, err = client.Order("order.create", []ParamMap{param})
	if err != nil {
		t.Error(err)
		return
	}
	assert.NotEqual(t, string(resp), "{}")

	<- client.StopSignal
}

func TestBybitQuote(t *testing.T) {
	client := bybit.NewWsSpotClient(bybit.IsTestNet())
	client.Logger = logger
	client.Connect()
	go client.Conn.ReadLoop()
	resp, err := client.Subscribe([]string{"orderbook.1.BTCUSDT", "publicTrade.BTCUSDT"})
	if err != nil {
		t.Log(string(resp))
		t.Error(err)
		client.Stop()
		return
	}

	go func() {
		time.Sleep(10 * time.Second)
		client.Stop()
	}()

	<- client.StopSignal
}