package bybit

import (
	"fmt"
	"github.com/lianyun0502/exchange_conn/v2/consts"
	"github.com/lianyun0502/exchange_conn/v2/ws_client"
	"github.com/valyala/fastjson"
	
)

func NewWsStreamClient(hostType, apiKey, secretKey string, opts ...func(*WsBybitClient)) (*WsBybitClient, error) {
	var client *WsBybitClient
	switch hostType {
		case consts.Private:
		client = &WsBybitClient{
			WsClient: wsClient.NewWsClient(&wsClient.ExchangeApi{
				Name: consts.Bybit,
				HostType: consts.Private,
				APIKey: apiKey,
				SecretKey: secretKey,
				BaseURL: WEBSOCKET_PRIVATE_MAINNET,
			}, nil),
			maxAliveTime: "",
		}
		default:
			return nil, fmt.Errorf("hostType error")
	}
	for _, opt := range opts {
		opt(client)
	}
	return client, nil
}

func NewWsPrivateClient(apiKey, secretKey string, opts ...func(*WsBybitClient)) (*WsBybitClient, error) {
	client, err := NewWsStreamClient(consts.Private, apiKey, secretKey, opts...)
	if err != nil {
		return nil, err
	}
	client.PingMessage = `{"op":"ping"}`
	opts = append(opts, WithPrivatePingfunction())
	for _, opt := range opts {
		opt(client)
	}
	return client, nil
}


type PrivateQuote[D any] struct {
	ID    string `json:"id,omitempty"`
	Topic string `json:"topic,omitempty"`
	Time  int64  `json:"time,omitempty"`
	Data  []D    `json:"data,omitempty"`
}

func WithPrivatePingfunction() func(client *WsBybitClient) {
	return func(client *WsBybitClient) {
		client.PingFunc = func () (err error) {
			client.ReqMap["pong"] = make(chan []byte, 2)
			client.Send([]byte(`{"op":"ping"}`))
			select {
			// case <- time.After(5 * time.Second):
			// 	err = fmt.Errorf("%s: Ping server timeout", wsc.ExchangeInfo.HostType)
			// 	delete(client.ReqMap, "pong")
			// 	wsc.PingTimeout.Stop()
			// 	wsc.Conn.NetConn().Close()
			// 	// wsc.OnClose(wsc.Conn, err)
			case respData := <-client.ReqMap["pong"]:
				resp := fastjson.MustParseBytes(respData)
				delete(client.ReqMap, "pong")
				if retCode := resp.GetInt("retCode"); retCode != 0 {
					err = fmt.Errorf("ping server failed, retCode=%d, retMsg:%s", retCode, string(resp.GetStringBytes("retMsg")))
					client.OnClose(client.Conn, err)
					// wsc.PingTimeout.Stop()
					// wsc.Conn.NetConn().Close()
				}else{
					client.OnPong(client.Conn, respData)
				}
			case <-client.StopSignal:
				client.Logger.Infof("%s: Stop Ping", client.ExchangeInfo.HostType)
				return
			}
			return 
		}
	}
}
