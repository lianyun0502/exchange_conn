package main

import (
	// "encoding/json"
	// "log"
	"net/http"
	"os"
	"time"

	"github.com/lianyun0502/exchange_conn/v1"
	"github.com/lianyun0502/exchange_conn/v1/binance_conn"
	"github.com/lianyun0502/exchange_conn/v1/binance_conn/data_stream"
	"github.com/lianyun0502/exchange_conn/v1/common"
	"github.com/sirupsen/logrus"
)

var ob = data_stream.NewOrderBook()

var logger = &logrus.Logger{
	Out:          os.Stdout,
	Formatter:    &logrus.TextFormatter{
		ForceColors: true,
		TimestampFormat: "2006-01-02 15:04:05.000",
		FullTimestamp: true,
	},
	Hooks:        make(logrus.LevelHooks),
	Level:        logrus.InfoLevel,
}

func wsHandler(message []byte) {
	logger.Infof("%v", time.Now().UnixNano()/int64(time.Millisecond))
	ob, _ := ob.Update(message)
	if ob != nil {
		logger.Info(common.PrettyPrint(ob))
		logger.Infof("%v", time.Now().UnixNano()/int64(time.Millisecond))
		return
	}
}

func errorHandler(err error) {
	logger.Error(err)
}

func main() {

	ws := binance_conn.NewWsClient(
		wsHandler,
		errorHandler,
		10,
	)
	ws.Logger = logger
	logger.SetLevel(logrus.DebugLevel)

	agent := exchange_conn.NewAgent(binance_conn.NewClient(
		"",
		"",
		"https://api.binance.com",
	))

	_, err := ws.Connect("wss://stream.binance.com:443/ws")
	if err != nil {
		logger.Error(err)
		return
	}
	ws.Send([]byte(`{"method": "SUBSCRIBE","params": ["btcusdt@depth@100ms"],"id": 1}`))

	go ws.StartLoop()

	go func() {
		time.Sleep(1 * time.Second)
		snapshot, _ := agent.Request(http.MethodGet, "/api/v3/depth", false, false).SetQuery("symbol", "BTCUSDT").SetQuery("limit", "50").Send()
		logger.Debug(string(snapshot))
		ob.SetSnapshot(snapshot)
	}()

	go func() {
		time.Sleep(10 * time.Second)
		err = ws.Stop()
		if err != nil {
			logger.Error(err)
			return
		}
	}()
	<-ws.DoneSignal
}
