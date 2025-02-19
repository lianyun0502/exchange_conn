package bybit

import (
	"fmt"
	"strconv"

	"github.com/lianyun0502/exchange_conn/v2/consts"
	"github.com/lianyun0502/exchange_conn/v2/ws_client"
	"github.com/valyala/fastjson"
)

func NewWsQuoteClient(category string, quoteHandle func([]byte), opts ...func(*WsBybitClient)) (*WsBybitClient, error) {
	var exchInfo *wsClient.ExchangeApi
	switch category {
		case consts.Spot:
		exchInfo = &wsClient.ExchangeApi{
			Name: consts.Bybit,
			HostType: consts.Spot,
			BaseURL: SPOT_MAINNET,
		}
		case consts.Future:
		exchInfo = &wsClient.ExchangeApi{
			Name: consts.Bybit,
			HostType: consts.Future,
			BaseURL: LINEAR_MAINNET,
		}
		default:
		return nil, fmt.Errorf("category error")

	}
	client := &WsBybitClient{
		WsClient: wsClient.NewWsClient(exchInfo, nil), 
		maxAliveTime: "",
	}
	opts = append(opts, WithWsHandle(quoteHandle), WithPublicPingfunction())
	for _, opt := range opts {
		opt(client)
	}
	return client, nil
}

func NewWsSpotQuoteClient(quoteHandle func([]byte), opts ...func(*WsBybitClient)) (*WsBybitClient, error) {
	return NewWsQuoteClient(consts.Spot, quoteHandle, opts...)
}

func NewWsPerpQuoteClient(quoteHandle func([]byte), opts ...func(*WsBybitClient)) (*WsBybitClient, error) {
	return NewWsQuoteClient(consts.Future, quoteHandle, opts...)
}

// 針對私有頻道和交易, 您可以自定義連接存活時長, 通過增加參數max_active_time, 最小支持30s (30秒), 最大支持600s (10分鐘)
func WithMaxAliveTime(maxAliveTime int) func(*WsBybitClient) {
	return func(wsc *WsBybitClient) {
		wsc.maxAliveTime = strconv.Itoa(maxAliveTime)
	}
}

func WithPublicPingfunction() func(client *WsBybitClient) {
	return func(client *WsBybitClient) {
		client.PingFunc = func () (err error) {
			client.ReqMap["ping"] = make(chan []byte, 2)
			client.Send([]byte(`{"op":"ping"}`))
			select {
			// case <- time.After(5 * time.Second):
			// 	err = fmt.Errorf("%s: Ping server timeout", wsc.ExchangeInfo.HostType)
			// 	delete(client.ReqMap, "pong")
			// 	wsc.PingTimeout.Stop()
			// 	wsc.Conn.NetConn().Close()
			// 	// wsc.OnClose(wsc.Conn, err)
			case respData := <-client.ReqMap["ping"]:
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