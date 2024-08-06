package binance_conn_test

import (
	"encoding/json"
	"testing"
	// "fmt"
	"log"
	"time"

	// "github.com/lianyun0502/exchange_conn/v1/binance_conn"
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
	// url := "wss://ws-api.binance.com:443/ws-api/v3"
	// url := "wss://stream.binance.com:9443/stream?streams=btcusdt@trade/btcusdt@aggTrade"

	// client, err := binance_conn.NewWsClient(
	// 	url, 
	// 	&binance_conn.WebSocketEvent{
	// 		Err_Handler: errorHandler,
	// 		Ws_Handler: wsHandler,
	// 	},
	// )
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// client.Conn
	// client.StartLoop()

	// go func() {
	// 	time.Sleep(60*time.Second)
	// 	client.Close()
	// }()
	
	// for {
	// 	select{
	// 		case <- client.DoneSignal:
	// 			fmt.Printf("end")
	// 			return
	// 		case <- client.ReconnectSignal:
	// 			for i :=0; i<10; i++{
	// 				log.Printf("retry connect {%d}", i)
	// 				err = client.Reconnect()
	// 				if err != nil{
	// 					break
	// 				}
	// 				log.Println(err)
	// 			}
				
	// 	}
		
}
	
