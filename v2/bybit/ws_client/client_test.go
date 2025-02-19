package bybit_test

import (
	"os"
	"testing"
	"time"

	"github.com/lianyun0502/exchange_conn/v2/bybit/ws_client"
	"github.com/lianyun0502/exchange_conn/v2/common"
	"github.com/lianyun0502/exchange_conn/v2/consts"
	"github.com/sirupsen/logrus"
)

var apiKey = "L7ksyiOdEgqg0gwIbf"
var secretKey = "0CVhyQmkwUDKWLcAP6NhtH7jB0P8XqSIVxE1"

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

func TestBybitWsApiOrder(t *testing.T) {
	client, _ := bybit.NewWsTradeClient(apiKey, secretKey, bybit.IsTestNet())
	client.Logger = logger
	client.PostStartFunc = func() error{
		resp, err := client.Auth(15)
		if err != nil {
			t.Log(string(resp))
			t.Error(err)
			return err
		}
		return nil
	}
	client.Connect()

	type ParamMap map[string]string
	param := ParamMap{
		"category":  "spot",
		"symbol":    "BTCUSDT",
		"side":      "Buy",
		"orderType": "Limit",
		"qty":       "0.001",
		"price":     "50000",
	}
	resp2, err := client.Order("order.create", []ParamMap{param})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(resp2.OrderID))

	go func() {
		time.Sleep(50 * time.Second)
		client.Stop()
	}()

	<-client.StopSignal
}

func TestBybitQuote(t *testing.T) {

	handle := func(rawData []byte) {
		logger.Infof(`%s`, string(rawData))
	}

	client, _ := bybit.NewWsQuoteClient(consts.Future, handle, bybit.IsTestNet())
	client.Logger = logger
	logger.SetLevel(logrus.DebugLevel)
	client.Connect()
	// client.Subscribe([]string{"orderbook.1.BTCUSDT", "publicTrade.BTCUSDT"})
	go func() {
		time.Sleep(40 * time.Second)
		client.Stop()
	}()
	common.WaitForClose(logger, client.StopSignal)
}
func BenchmarkBybitWsApiOrder(b *testing.B) {
	client, _ := bybit.NewWsTradeClient(apiKey, secretKey, bybit.IsTestNet())
	client.Logger = logger
	// logger.SetLevel(logrus.ErrorLevel)
	client.Connect()
	client.Auth(15)
	type ParamMap map[string]string
	param := ParamMap{
		"category":  "linear",
		"symbol":    "ADAUSDT",
		"side":      "Buy",
		"orderType": "Market",
		"qty":       "5.52",
		// "price":     "50000",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.Order("order.create", []ParamMap{param})
	}

}

func TestPrivateWsQuote(t *testing.T) {
	handle := func(rawData []byte) {
		logger.Infof(`%s`, string(rawData))
	}

	client, _ := bybit.NewWsPrivateClient(apiKey, secretKey, bybit.IsTestNet(), bybit.WithWsHandle(handle))
	client.Logger = logger
	logger.SetLevel(logrus.DebugLevel)
	client.Connect()
	res, err := client.Auth(15)
	if err != nil {
		t.Log(string(res))
		t.Error(err)
		return
	}
	res, err = client.Subscribe([]string{"position"})
	if err != nil {
		t.Log(string(res))
		t.Error(err)
		return
	}

	go func() {
		time.Sleep(20 * time.Second)
		client.Stop()
	}()

	common.WaitForClose(logger, client.StopSignal)
}

func TestPrivateWsTrade(t *testing.T) {
	handle := func(rawData []byte) {
		logger.Infof(`%s`, string(rawData))
	}

	// client, _ := bybit.NewWsPrivateClient(apiKey, secretKey, bybit.IsTestNet(), bybit.WithWsHandle(handle))
	client, _ := bybit.NewWsPrivateClient("L9THu88IZnbyg21MIt", "f4ViW6m6wDbTuTdeNlaKnOWFkvKH4DWJkD1k", bybit.WithWsHandle(handle))
	client.Logger = logger
	logger.SetLevel(logrus.DebugLevel)
	client.Connect()

	go func() {
		for range client.StartSignal {
			client.Auth(15)
			resp, err := client.Subscribe([]string{"position", "execution"})
			if err != nil {
				t.Log(string(resp))
				t.Error(err)
				client.Stop()
				return
			}
		}
	}()

	go client.StartLoop()

	go func() {
		time.Sleep(60 * time.Second)
		client.Stop()
	}()

	common.WaitForClose(logger, client.StopSignal)
}