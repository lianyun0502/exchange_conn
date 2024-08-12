package binance_conn_test

import (
	"log"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/lianyun0502/exchange_conn/v1/common"
	"github.com/lianyun0502/exchange_conn/v1/binance_conn"
)

var errHandler = func(err error) {
	if err != nil {
		log.Println(err)
	}
}

// type WsApiPingResponse struct {
// 	ID string `json:"id"`
// 	Method string `json:"method"`
// }


func TestWsApiPing(t *testing.T) {
	assert := assert.New(t)
	client, err := binance_conn.NewWebSocketAPI(
		apiKey,
		secretKey,
		"wss://testnet.binance.vision/ws-api/v3",
		errHandler,
	)
	if err != nil {
		t.Error(err)
		return
	}
	client.StartLoop()
	for i := 0; i < 20; i++ {
		resp, err := client.PingServer()
		if err != nil {
			t.Error(err)
			return
		}
		data :=  common.PrettyPrint(resp)
		log.Printf("response:\n%s", data)
		assert.NotEqual("{}", data)
	}
	client.StopLoop()
}