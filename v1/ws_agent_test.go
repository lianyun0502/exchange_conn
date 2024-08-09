package exchange_conn_test

import (
	"encoding/json"
	"log"
	"testing"
	"time"

	// "github.com/stretchr/testify/assert"
	"github.com/lianyun0502/exchange_conn/v1"
	"github.com/lianyun0502/exchange_conn/v1/binance_conn"
	"github.com/lianyun0502/exchange_conn/v1/bybit_conn"
)

func wsHandler(message []byte) {
	log.Println(string(message))
	j := make(map[string]interface{})
	json.Unmarshal(message, &j)
	log.Printf("%v", j["E"])
	log.Printf("%v", float64(time.Now().UnixNano()/int64(time.Millisecond)))
}
func errorHandler(err error) {
	log.Println(err)
}

func TestBinanceData(t *testing.T) {

	url := "wss://stream.binance.com:443/ws"

	agent := exchange_conn.NewWebSocketAgent(binance_conn.NewWsClient(wsHandler, errorHandler, 10))

	agent.Connect(url)

	go agent.StartLoop()

	agent.SendString(`{"method": "SUBSCRIBE","params": ["btcusdt@trade", "btcusdt@aggTrade", "btcusdt@depth@100ms"],"id": 1}`)

	go func() {
		time.Sleep(10 * time.Second)
		agent.Stop()
	}()

	<-agent.Client.DoneSignal
}

func TestBybitData(t *testing.T) {

	url := "wss://stream.bybit.com/v5/public/spot"

	agent := exchange_conn.NewWebSocketAgent(bybit_conn.NewWsClient(wsHandler, errorHandler, 10))

	agent.Connect(url)

	go agent.StartLoop()

	agent.SendString(`{"req_id": "test","op":"subscribe","args":["orderbook.1.BTCUSDT"]}`)

	go func() {
		time.Sleep(1 * time.Second)
		agent.Stop()
	}()

	<-agent.Client.DoneSignal
}
