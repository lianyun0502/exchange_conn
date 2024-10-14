package bybit

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/lianyun0502/exchange_conn/v2/common"
	"github.com/lianyun0502/exchange_conn/v2/consts"
	"github.com/lianyun0502/exchange_conn/v2/ws_client"
	"github.com/valyala/fastjson"
)

type WsBybitClient struct {
	*wsClient.WsClient
	maxAliveTime string
	IsSubscribed bool
}

func (wsc *WsBybitClient) Connect() (resp *http.Response, err error) {
	if wsc.maxAliveTime != "" {
		wsc.ExchangeInfo.BaseURL += ("?max_alive_time=" + wsc.maxAliveTime)
	}
	return wsc.WsClient.Connect()
}

func (wsc *WsBybitClient) Subscribe(topics []string) (respData []byte, err error) {
	id := common.GetUUID()
	jTopics, _ := json.Marshal(topics)
	msg := fmt.Sprintf(`{"req_id":"%s","op":"subscribe","args":%s}`, id, string(jTopics))
	respCh := make(chan []byte, 2)
	wsc.ReqMap[id] = respCh
	err = wsc.Send([]byte(msg))
	var resp []byte
	select {
	case <-time.After(5 * time.Second):
		wsc.Logger.Warning("Subscribe Request timeout")
		err = errors.New("Request timeout")
		close(respCh)
	case resp = <-respCh:
		wsc.Logger.Debugf(`Response: %s`, string(resp))
		wsc.IsSubscribed = true
	}
	delete(wsc.ReqMap, id)
	return resp, err
}

func (wsc *WsBybitClient) Auth() (respData []byte, err error) {
	expires := time.Now().Unix()*1000 + 10000
	param := []string{
		wsc.ExchangeInfo.APIKey,
		strconv.FormatInt(expires, 10),
		common.GetSignature(wsc.ExchangeInfo.SecretKey, fmt.Sprintf("GET/realtime%d", expires)),
	}
	resp, err := wsc.Request("auth", nil, param, 10)
	if err != nil {
		wsc.Logger.Warningf("%s: Auth Request timeout", wsc.ExchangeInfo.HostType)
		return nil, err
	}
	return resp, nil
}

func (wsc *WsBybitClient) Request(op string, header any, args any, timeOut time.Duration) (respData []byte, err error) {
	id := common.GetUUID()
	req := &Request{
		ReqID:  id,
		Header: header,
		Op:     op,
		Args:   args,
	}
	reqByte, err := json.Marshal(req)
	if err != nil {
		return
	}
	respCh := make(chan []byte, 2)
	wsc.ReqMap[id] = respCh
	wsc.Logger.Debug("Send request")
	wsc.Send(reqByte)
	var resp []byte
	select {
	case <-time.After(timeOut * time.Second):
		wsc.Logger.Warning("Request timeout")
		err = errors.New("Request timeout")
		close(respCh)
	case resp = <-respCh:
		wsc.Logger.Debugf(`Response: %s`, string(resp))
	}
	delete(wsc.ReqMap, id)
	return resp, err
}

func (wsc *WsBybitClient) WithPingServer() func(msg []byte) error {
	wsc.ReqMap["pong"] = make(chan []byte, 2)
	return wsc.PingServer
}

func (wsc *WsBybitClient) PingServer(msg []byte) (err error) {
	wsc.ReqMap["pong"] = make(chan []byte, 2)
	wsc.Send([]byte(`{"op":"ping"}`))
	select {
	case <- time.After(5 * time.Second):
		err = fmt.Errorf("%s: Ping server timeout", wsc.ExchangeInfo.HostType)
		delete(wsc.ReqMap, "pong")
		wsc.OnClose(wsc.Conn, err)
	case respData := <-wsc.ReqMap["pong"]:
		resp := fastjson.MustParseBytes(respData)
		delete(wsc.ReqMap, "pong")
		if retCode := resp.GetInt("retCode"); retCode != 0 {
			err = fmt.Errorf("ping server failed, retCode=%d, retMsg:%s", retCode, string(resp.GetStringBytes("retMsg")))
			wsc.OnClose(wsc.Conn, err)
		}else{
			wsc.OnPong(wsc.Conn, respData)
		}
	}
	return err
}

func IsTestNet() func(*WsBybitClient) {
	return func(wsc *WsBybitClient) {
		switch wsc.ExchangeInfo.HostType {
		case consts.Future:
			wsc.ExchangeInfo.BaseURL = LINEAR_TESTNET
		case consts.Spot:
			wsc.ExchangeInfo.BaseURL = SPOT_TESTNET
		case consts.Trade:
			wsc.ExchangeInfo.BaseURL = WEBSOCKET_TRADE_TESTNET
		case consts.Private:
			wsc.ExchangeInfo.BaseURL = WEBSOCKET_PRIVATE_TESTNET
		}
	}
}

func WithWsHandle(qouteHandler func(message []byte)) func(*WsBybitClient) {
	return func(client *WsBybitClient) {
		client.Ws_Handler = func(rawData []byte) {
			v, _ := fastjson.ParseBytes(rawData)
			if reqID := string(v.GetStringBytes("reqId")); reqID != "" {
				if respCh, ok := client.ReqMap[reqID]; ok {
					respCh <- rawData
					return
				}
			}
			if reqID := string(v.GetStringBytes("req_id")); reqID != "" {
				if respCh, ok := client.ReqMap[reqID]; ok {
					respCh <- rawData
					return
				}
			}
			if op := string(v.GetStringBytes("op")); op == "pong" {
				client.ReqMap["pong"] <- rawData
				return
			}
			if qouteHandler != nil {
				if client.IsSubscribed {
					qouteHandler(rawData)
				}
			}
		}
	}
}


type Request struct {
	ReqID  string `json:"reqId"`
	Header any    `json:"header,omitempty"`
	Op     string `json:"op"`
	Args   any    `json:"args,omitempty"`
}
