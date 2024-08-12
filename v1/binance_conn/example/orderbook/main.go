package main

import (
	// "encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/lianyun0502/exchange_conn/v1"
	"github.com/lianyun0502/exchange_conn/v1/binance_conn"
	"github.com/lianyun0502/exchange_conn/v1/binance_conn/orderbook"
	"github.com/lianyun0502/exchange_conn/v1/common"
)

var ob = orderbook.NewOrderBook()

func wsHandler(message []byte) {
	log.Printf("%v", time.Now().UnixNano()/int64(time.Millisecond))
	ob, _ := ob.Update(message)
	if ob != nil {
		// b, _ := json.Marshal(ob)
		log.Println(common.PrettyPrint(ob))
		log.Printf("%v", time.Now().UnixNano()/int64(time.Millisecond))
		return
	}
}

func errorHandler(err error) {
	log.Println(err)
}

func main() {

	ws := binance_conn.NewWsClient(
		wsHandler,
		errorHandler,
		10,
	)

	agent := exchange_conn.NewAgent(binance_conn.NewClient(
		"",
		"",
		"https://api.binance.com",
	))

	err := ws.Connect("wss://stream.binance.com:443/ws")
	if err != nil {
		log.Println(err)
		return
	}
	ws.Send([]byte(`{"method": "SUBSCRIBE","params": ["btcusdt@depth@100ms"],"id": 1}`))

	go ws.StartLoop()

	go func() {
		time.Sleep(1 * time.Second)
		snapshot, _ := agent.Request(http.MethodGet, "/api/v3/depth", false, false).SetQuery("symbol", "BTCUSDT").SetQuery("limit", "50").Send()
		// log.Println(string(snapshot))
		ob.SetSnapshot(snapshot)
	}()

	go func() {
		time.Sleep(2 * time.Second)
		err = ws.Stop()
		if err != nil {
			log.Println(err)
			return
		}
	}()
	<-ws.DoneSignal
}
