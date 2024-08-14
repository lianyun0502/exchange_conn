package main


import (
	"encoding/json"
	"log"
	"time"
	"github.com/lianyun0502/exchange_conn/v1/binance_conn"
	"github.com/lianyun0502/exchange_conn/v1/binance_conn/trade_stream"
)


func wsHandler(message []byte) {
	data, err := trade_stream.ToNormalTradeData(message)
	if err != nil {
		log.Println(err)
	}
	j, err  :=json.Marshal(data)
	if err != nil {
		log.Println(err)
	}

	log.Println(string(j))
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

	_, err := ws.Connect("wss://stream.binance.com:9443/ws")
	if err != nil {
		log.Println(err)
		return
	}
	ws.Send([]byte(`{"method": "SUBSCRIBE","params": ["btcusdt@trade"],"id": 1}`))

	go ws.StartLoop()

	go func() {
		time.Sleep(30 * time.Second)
		err = ws.Stop()
		if err != nil {
			log.Println(err)
			return
		}
	}()
	<-ws.DoneSignal


}