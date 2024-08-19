package binance_conn

import (
	"github.com/lianyun0502/exchange_conn/v1"
)

type WsClient struct {
	*exchange_conn.WsClient
}

func NewWsClient(messageHandle func(message []byte), errHandle func(err error), reconnectTimes int) (client *WsClient) {
	return &WsClient{WsClient: exchange_conn.NewWsClient(messageHandle, errHandle, reconnectTimes)}
}


func (c *WsClient) Subscribe(topics string) {
	c.Send([]byte(`{"method":"SUBSCRIBE","params":` + topics + `,"id":1}`))
}
