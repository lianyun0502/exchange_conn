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
	"github.com/lxzan/gws"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fastjson"
)

type WsBybitClient struct {
	*exchange_conn.WsClient
	maxAliveTime string
}

func (wsc *WsBybitClient) Connect() (resp *http.Response, err error) {
	url := wsc.ExchangeInfo.BaseURL
	if wsc.maxAliveTime != "" {
		url += "?max_alive_time=" + wsc.maxAliveTime
	}

	wsc.ClientOption = &gws.ClientOption{
		ReadBufferSize:   655350,
		Addr:             url,
		HandshakeTimeout: 45 * time.Second,
		PermessageDeflate: gws.PermessageDeflate{
			Enabled:               true,
			ServerContextTakeover: true,
			ClientContextTakeover: true,
		},
	}

	wsc.Conn, resp, err = gws.NewClient(wsc, wsc.ClientOption)
	if err != nil {
		wsc.Logger.WithFields(logrus.Fields{"respone": resp}).Error(err)
	}
	return resp, err
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
		wsc.Logger.Warning("Request timeout")
		err = errors.New("Request timeout")
		close(respCh)
	case resp = <-respCh:
		wsc.Logger.Debugf(`Response: %s`, string(resp))
	}
	delete(wsc.ReqMap, id)
	return resp, err
}

func (wsc *WsBybitClient) GetSignature() (respData []byte, err error) {
	expires := time.Now().Unix()*1000 + 10000
	param := []string{
		wsc.ExchangeInfo.APIKey,
		strconv.FormatInt(expires, 10),
		common.GetSignature(wsc.ExchangeInfo.SecretKey, fmt.Sprintf("GET/realtime%d", expires)),
	}
	resp, err := wsc.Request("auth", nil, param)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (wsc *WsBybitClient) Request(op string, header any, args any) (respData []byte, err error) {
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
	case <-time.After(5 * time.Second):
		wsc.Logger.Warning("Request timeout")
		err = errors.New("Request timeout")
		close(respCh)
	case resp = <-respCh:
		wsc.Logger.Debugf(`Response: %s`, string(resp))
	}
	delete(wsc.ReqMap, id)
	return resp, err
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
		}
	}
}

func WithBybitHandler(reqMap map[string]chan []byte, qouteHandler func(message []byte)) func([]byte) {
	return func(rawData []byte) {
		v, _ := fastjson.ParseBytes(rawData)
		if reqID := string(v.GetStringBytes("reqId")); reqID != "" {
			if respCh, ok := reqMap[reqID]; ok {
				respCh <- rawData
				return
			}
		}
		if reqID := string(v.GetStringBytes("req_id")); reqID != "" {
			if respCh, ok := reqMap[reqID]; ok {
				respCh <- rawData
				return
			}
		}
		if qouteHandler != nil {
			qouteHandler(rawData)
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
			if qouteHandler != nil {
				qouteHandler(rawData)
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
