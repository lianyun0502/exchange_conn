package exchange_conn_test

import (
	"encoding/json"
	"testing"
	"time"

	// "github.com/stretchr/testify/assert"
	"github.com/lianyun0502/exchange_conn/v1"
	"github.com/lianyun0502/exchange_conn/v1/binance_conn"
	"github.com/lianyun0502/exchange_conn/v1/bybit_conn"
	"github.com/sirupsen/logrus"
)

// var logger = &logrus.Logger{
// 		Out:          os.Stderr,
// 		Formatter:    &logrus.TextFormatter{
// 			ForceColors: true,
// 			TimestampFormat: "2006-01-02 15:04:05.000",
// 			FullTimestamp: true,
// 		},
// 		Hooks:        make(logrus.LevelHooks),
// 		Level:        logrus.InfoLevel,
// 		ExitFunc:     os.Exit,
// 		ReportCaller: false,
// }

func wsHandler(message []byte) {
	logger.Debug(string(message))
	j := make(map[string]interface{})
	json.Unmarshal(message, &j)
	logger.Debugf("%v", j["E"])
	logger.Debugf("%v", float64(time.Now().UnixNano()/int64(time.Millisecond)))
}
func errorHandler(err error) {
	logger.Error(err)
}

func TestBinanceOrderBookData(t *testing.T) {

	url := "wss://stream.binance.com:9443/ws"

	agent := exchange_conn.NewWebSocketAgent(binance_conn.NewWsClient(wsHandler, errorHandler, 10))
	agent.Client.Logger = logger
	agent.Client.Logger.SetLevel(logrus.DebugLevel)

	resp, err := agent.Connect(url)
	if err != nil {
		logger.Println(resp)
		t.Error(err)
		return
	}

	go func() {
		for {
			<-agent.Client.StartSignal
			agent.SendString(`{"method": "SUBSCRIBE","params": ["btcusdt@trade", "btcusdt@aggTrade", "btcusdt@depth@100ms"],"id": 1}`)
		}
	}()

	go agent.StartLoop()

	go func() {
		time.Sleep(5 * time.Second)
		agent.Stop()
	}()

	<-agent.Client.DoneSignal

}

func TestBybitOrderBookData(t *testing.T) {

	url := "wss://stream.bybit.com/v5/public/spot"

	agent := exchange_conn.NewWebSocketAgent(bybit_conn.NewWsClient(wsHandler, errorHandler, 10))
	agent.Client.Logger = logger
	agent.Client.Logger.SetLevel(logrus.DebugLevel)

	agent.Connect(url)

	go agent.StartLoop()

	agent.SendString(`{"req_id": "test","op":"subscribe","args":["orderbook.50.BTCUSDT"]}`)
	
	go func() {
		time.Sleep(3 * time.Second)
		agent.Stop()
	}()

	<-agent.Client.DoneSignal
}

func TestBybitTradeData(t *testing.T) {
	
	url := "wss://stream.bybit.com/v5/public/linear"

	agent := exchange_conn.NewWebSocketAgent(bybit_conn.NewWsClient(wsHandler, errorHandler, 10))
	agent.Client.Logger = logger
	agent.Client.Logger.SetLevel(logrus.DebugLevel)

	agent.Connect(url)

	go agent.StartLoop()

	agent.Subscribe(`["publicTrade.BTCUSDT"]`)

	go func() {
		time.Sleep(30 * time.Second)
		agent.Stop()
	}()

	<-agent.Client.DoneSignal
}


func TestBybitMarketPriceData(t *testing.T) {

	url := "wss://stream.bybit.com/v5/public/linear"

	agent := exchange_conn.NewWebSocketAgent(bybit_conn.NewWsClient(wsHandler, errorHandler, 10))
	agent.Client.Logger = logger
	agent.Client.Logger.SetLevel(logrus.DebugLevel)

	agent.Connect(url)

	go agent.StartLoop()

	agent.Subscribe(`["tickers.BTCUSDT"]`)

	go func() {
		time.Sleep(30 * time.Second)
		agent.Stop()
	}()

	<-agent.Client.DoneSignal

}