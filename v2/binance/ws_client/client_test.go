package binance_test

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/lianyun0502/exchange_conn/v2/binance/ws_client"
	// "github.com/lianyun0502/exchange_conn/v2/common"
	"github.com/lianyun0502/exchange_conn/v2/consts"
	"github.com/sirupsen/logrus"
)

var apiKey = "GNEImPtijyYkCYRBx7WM3bntSPydYYY8bwqFuYA6BVVUNFKQ1cVTi7AHIdCsBQ3Z"
var secretKey = "5lkEjNgcicUQ7YZt1L0CAaTc12UWRuhLrUsahle1iQxIXMk95lwVsLEXRaQfkPQr"

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

func TestBinanceWsQuote(t *testing.T) {
	handle := func(data []byte) {
		logger.Infof(`%s`, string(data))
	}
	client, _ := binance.NewWsQuoteClient(consts.Spot, handle, binance.IsTestNet())
	client.Logger = logger
	client.Logger.SetLevel(logrus.DebugLevel)

	resp, err := client.Connect()
	if err != nil {
		logger.Println(resp)
		t.Error(err)
		return
	}

	go func() {
		for range client.StartSignal {
			client.Subscribe([]string{"btcusdt@trade", "btcusdt@aggTrade", "btcusdt@depth@100ms"})
		}
	}()

	go client.StartLoop()


	time.Sleep(10 * time.Second)
	client.Stop()


	<-client.StopSignal
}

func TestBinanceWsOrder(t *testing.T) {
	client, _ := binance.NewWsAPIClient(consts.Spot, apiKey, secretKey, binance.IsTestNet())
	client.Logger = logger
	client.Logger.SetLevel(logrus.DebugLevel)
	_, err := client.Connect()
	if err != nil {
		t.Error(err)
		return
	}
	go client.StartLoop()

	go func() {
		<-client.StartSignal
		params := binance.ParamMap{
			"symbol":      "BTCUSDT",
			"side":        "SELL",
			"type":        "LIMIT",
			"timeInForce": "GTC",
			"price":       "50000",
			"quantity":    "0.001",
			"timestamp":   strconv.FormatInt(time.Now().UnixMilli(), 10),
		}
		params.Sign(apiKey, secretKey)
		respData, err := client.Request("order.place", params)
		if err != nil {
			t.Error(err)
			return
		}
		logger.Infof(`%s`, string(respData))
	}()


	time.Sleep(10 * time.Second)
	client.Stop()

	<- client.StopSignal
}

func TestBinanceWsOrderFuture(t *testing.T) {
	client, _ := binance.NewWsAPIClient(consts.Future, apiKey, secretKey, binance.IsTestNet())
	client.Logger = logger
	client.Logger.SetLevel(logrus.DebugLevel)
	_, err := client.Connect()

	if err != nil {
		t.Error(err)
		return
	}
	go client.StartLoop()

	go func() {
		<-client.StartSignal
		params := binance.ParamMap{
			"symbol":      "BTCUSDT",
			"side":        "SELL",
			"type":        "LIMIT",
			"timeInForce": "GTC",
			"price":       "50000",
			"quantity":    "0.001",
			"timestamp":   strconv.FormatInt(time.Now().UnixMilli(), 10),
		}
		params.Sign(apiKey, secretKey)
		respData, err := client.Request("order.place", params)
		if err != nil {
			t.Error(err)
			return
		}
		logger.Infof(`%s`, string(respData))
	}()


	time.Sleep(10 * time.Second)
	client.Stop()

	<- client.StopSignal
}

func TestBinanceWsOrderbook(t *testing.T) {
	client, _ := binance.NewWsAPIClient(consts.Spot, apiKey, secretKey, binance.IsTestNet())
	client.Logger = logger
	client.Logger.SetLevel(logrus.DebugLevel)
	resp, err := client.Connect()
	if err != nil {
		t.Log(resp)
		t.Error(err)
		return
	}
	go client.StartLoop()

	go func() {
		<-client.StartSignal
		params := binance.ParamMap{
			"symbol": "BTCUSDT",
			"limit":  "5",
		}
		respData, err := client.Request("depth", params)
		if err != nil {
			t.Error(err)
			return
		}
		logger.Infof(`%s`, string(respData))
	}()

	time.Sleep(10 * time.Second)
	client.Stop()

	<-client.StopSignal

}
