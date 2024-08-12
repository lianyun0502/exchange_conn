package binance_conn_test

import (
	"encoding/json"
	"testing"
	// "fmt"
	"log"
	"time"

	"github.com/lianyun0502/exchange_conn/v1/binance_conn"
)


func wsHandler(message []byte) {
	log.Println(string(message))
	j := make(map[string]interface{})
	json.Unmarshal(message, &j)
	log.Printf("%v", j["E"])
	log.Printf("%v", time.Now().UnixNano()/int64(time.Millisecond))
}
func errorHandler(err error) {
	log.Println(err)
}


func TestWsClient(t *testing.T) {
	// url := "wss://stream.binance.com:9443/ws/btcusdt@depth@100ms"
	// url := "wss://stream.binance.com:9443/ws/btcusdt@aggTrade"
	url := "wss://ws-api.binance.com:443/ws-api/v3"
	// url := "wss://stream.binance.com:9443/stream?streams=btcusdt@trade/btcusdt@aggTrade"

	client := binance_conn.NewWsClient(
		wsHandler,
		errorHandler,
		10,
	)
	err := client.Connect(url)
	if err != nil {
		log.Println(err)
		return
	}

	go client.StartLoop()
	go func() {
		time.Sleep(30*time.Second)
		err = client.Stop()
		if err != nil {
			log.Println(err)
			return
		} 
	}()
	<- client.DoneSignal
		
}
	
