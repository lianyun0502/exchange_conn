package binance

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sort"

	// "net/http"
	"strconv"
	"time"

	"github.com/lianyun0502/exchange_conn/v2/common"
	// "github.com/lianyun0502/exchange_conn/v2/consts"
	"github.com/lianyun0502/exchange_conn/v2/ws_client"
	// "github.com/lxzan/gws"
	// "github.com/sirupsen/logrus"
	"github.com/valyala/fastjson"
)

type WsBinanceClient struct {
	*wsClient.WsClient
	ReceiveWindow string
}

func (wsc *WsBinanceClient) Subscribe(topics []string) (respData []byte, err error) {
	id := common.GetUUID()
	jTopics, _ := json.Marshal(topics)
	msg := fmt.Sprintf(`{"id":"%s","method": "SUBSCRIBE","params":%s}`, id, string(jTopics))
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

func (wsc *WsBinanceClient) Auth() (respData []byte, err error) {
	timeStamp := time.Now().UnixMilli()
	body := fmt.Sprintf(`apiKey=%s&recvWindow=%s&timestamp=%d`,wsc.ExchangeInfo.APIKey, wsc.ReceiveWindow, timeStamp)
	param := map[string]string{
		"apiKey":     wsc.ExchangeInfo.APIKey,
		"timestamp":  strconv.FormatInt(timeStamp, 10),
		"recvWindow": wsc.ReceiveWindow,
		"signature":  common.GetSignature(wsc.ExchangeInfo.SecretKey, body),
	}
	resp, err := wsc.Request("session.logon", param)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (wsc *WsBinanceClient) Request(method string, args any) (respData []byte, err error) {
	id := common.GetUUID()
	req := &Request{
		ReqID:  id,
		Method: method,
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

type Request struct {
	ReqID  string `json:"id"`
	Method string `json:"method"`
	Args   any    `json:"params,omitempty"`
}

func WithWsHandle(qouteHandler func(message []byte)) func(*WsBinanceClient) {
	return func(client *WsBinanceClient) {
		client.Ws_Handler = func(rawData []byte) {
			v, _ := fastjson.ParseBytes(rawData)
			if reqID := string(v.GetStringBytes("id")); reqID != "" {
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

func IsTestNet() func(*WsBinanceClient) {
	return func(c *WsBinanceClient) {
		switch c.ExchangeInfo.BaseURL {
		case SPOT_QUOTE_MAINNET:
			c.WsClient.ExchangeInfo.BaseURL = SPOT_QUOTE_TESTNET
		case SPOT_MAINNET:
			c.WsClient.ExchangeInfo.BaseURL = SPOT_TESTNET
		case USD_QUOTE_MAINNET:
			c.WsClient.ExchangeInfo.BaseURL = USD_QUOTE_TESTNET
		case USD_MAINNET:
			c.WsClient.ExchangeInfo.BaseURL = USD_TESTNET
		}
	}
}
type ParamMap map[string]string	

func (pm ParamMap) Sign(apiKey, secretKey string) {
	pm["apiKey"] = apiKey
	pm["recvWindow"] = "5000"
	Keys := make([]string, 0)
	for k := range pm {
		Keys = append(Keys, k)
	}
	sort.Strings(Keys)
	queryStrs := make([]string, 0)
	for _, k := range Keys {
		queryStrs = append(queryStrs, fmt.Sprintf(`%s=%s`, k, pm[k]))
	}
	query := strings.Join(queryStrs, "&")
	pm["signature"] = common.GetSignature(secretKey, query)
}