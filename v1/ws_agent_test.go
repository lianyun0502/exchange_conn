package exchange_conn_test

import (
	"testing"
	"log"
	"encoding/json"
	"time"

	// "github.com/stretchr/testify/assert"
	"github.com/lianyun0502/exchange_conn/v1/binance_conn"
	"github.com/lianyun0502/exchange_conn/v1"
	// "github.com/lianyun0502/exchange_conn/v1/common"
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
	// url := "wss://stream.binance.com:9443/ws/btcusdt@depth@100ms"
	// url := "wss://stream.binance.com:9443/ws/btcusdt@aggTrade"
	url := "wss://ws-api.binance.com:443/ws-api/v3"
	// url := "wss://stream.binance.com:9443/stream?streams=btcusdt@trade/btcusdt@aggTrade"

	wsEvent := &binance_conn.WebSocketEvent{
		Err_Handler: errorHandler,
		Ws_Handler: wsHandler,
	}

	client, err := binance_conn.NewWsClient(
		url, 
		wsEvent,
	)
	if err != nil {
		log.Println(err)
		return
	}
	
	agent := exchange_conn.NewWebSocketAgent(client, 5)
	wsEvent.Close_Handler = func() {
		agent.AddEvent(&exchange_conn.Event{
			Name: "reconnect",
			IsBlock: true,
			Handler: agent.Reconnect,
		})
	}
	defer agent.Stop()
	agent.Start()

	
	<-agent.DoneSignal
}
	
