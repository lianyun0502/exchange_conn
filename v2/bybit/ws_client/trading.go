package bybit

import (
	"strconv"
	"time"
	"fmt"
	"encoding/json"

	"github.com/valyala/fastjson"
	"github.com/lianyun0502/exchange_conn/v2/consts"
	"github.com/lianyun0502/exchange_conn/v2/ws_client"
)
func NewWsAPIClient(hostType string, apiKey, secretKey string, opts ...func(*WsBybitClient)) (*WsBybitClient, error) {
	var exchInfo *wsClient.ExchangeApi
	switch hostType {
		case consts.Trade:
		exchInfo = &wsClient.ExchangeApi{
			Name: consts.Bybit,
			HostType: consts.Trade,
			APIKey: apiKey,
			SecretKey: secretKey,
			BaseURL: WEBSOCKET_TRADE_MAINNET,
		}
		default:
			return nil, fmt.Errorf("hostType error")
	}
	client := &WsBybitClient{
		WsClient: wsClient.NewWsClient(exchInfo, nil), 
		maxAliveTime: "",
	}
	opts = append(opts, WithWsHandle(nil))
	for _, opt := range opts {
		opt(client)
	}
	return client, nil
}



func NewWsTradeClient(apiKey, secretKey string, opts ...func(*WsBybitClient)) (*WsBybitClient, error) {
	client, err := NewWsAPIClient(consts.Trade, apiKey, secretKey, opts...)
	if err != nil {
		return nil, err
	}
	client.PingMessage = `{"op":"ping"}`
	client.Ping = client.PingServer
	return client, nil
}

func (wsc *WsBybitClient) order(op string, args any) (respData []byte, err error) {
	header := &RequestHeader{
		Timestamp: time.Now().UnixMilli(),
		RecvWindow: 8000,
		Referer: "bot-001",
	}
	resp, err := wsc.Request(op, header, args, 5)
	if err != nil {
		wsc.Logger.WithField("error", err).Error("Request failed")
		return nil, err
	}
	v, _ := fastjson.ParseBytes(resp)
	if fastjson.GetInt(resp, "retCode") != 0 {
		retMsg := string(v.GetStringBytes("retMsg"))
		wsc.Logger.WithField("retMsg", retMsg).Warn("Request failed")
	}
	return resp, nil
}

func (wsc *WsBybitClient) Order(op string, args any) (*CreateOrderResponse, error) {
	respData, err := wsc.order(op, args)
	if err != nil {
		return nil, err
	}
	resp := new(WsResponse[*CreateOrderResponse])
	if err = json.Unmarshal(respData, resp); err != nil {
		return nil, err
	}
	if resp.RetCode != 0 {
		return nil, fmt.Errorf("retCode=%d, retMsg=%s", resp.RetCode, resp.RetMsg)
	}
	return resp.Data, nil
}


type RequestHeader struct {	
	Timestamp int64 `json:"X-BAPI-TIMESTAMP"`
	RecvWindow int64 `json:"X-BAPI-RECV-WINDOW,string"`
	Referer string `json:"Referer"`
}



type ResponseHeader struct {
	TraceId string `json:"TraceId"`
	TimeNow int64 `json:"Timenow"`
	Limit int `json:"X-Bapi-Limit"`
	LimitStatus int `json:"X-Bapi-Limit-Status"`
	ResetTimestamp int64 `json:"X-Bapi-Limit-Reset"`
}


type Response struct {
	ReqID string `json:"reqId"`
	RetCode int `json:"retCode"`
	RetMsg string `json:"retMsg"`
	Op string `json:"op"`
	Header ResponseHeader `json:"header"`
	Data any `json:"data"`
	ConnId string `json:"connId"`
}


// func Order(category, symbol, side, orderType, qty string, orderOpts...func(map[string]string)) map[string]string {
// 	args := make(map[string]string)
// 	args["category"] = category
// 	args["symbol"] = symbol
// 	args["side"] = side
// 	args["orderType"] = orderType
// 	args["qty"] = qty
// 	for _, opt := range orderOpts {
// 		opt(args)
// 	}
// 	return args
// }

/*
是否借貸. 僅統一帳戶的現貨交易有效. 

	0(default): 否，則是幣幣訂單
	1: 是，則是槓桿訂單
*/
func IsLeverage(isLeverage int) func(map[string]string) {
	return func(args map[string]string) {
		if args["category"] != "spot" {
			return
		}
		args["leverage"] = strconv.Itoa(isLeverage)
	}
}

func MarketUnit(unit string) func(map[string]string) {
	return func(args map[string]string) {
		if (args["orderType"] != "Market") || (args["category"] != "spot") {
			return
		}
		args["unit"] = unit
	}
}

func Price(price string) func(map[string]string) {
	return func(args map[string]string) {
		if (args["orderType"] != "Limit"){
			return
		}
		args["price"] = price
	}
}
// /*
// direction:

// 	1: 當市場價上漲到了triggerPrice時觸發條件單
// 	2: 當市場價下跌到了triggerPrice時觸發條件單

// */
// func Trigger(price, direction string) func(map[string]string) {
// 	return func(args map[string]string) {
// 		if (args["category"] != "linear") || (args["category"] != "inverse") {
// 			return
// 		}
// 		args["triggerDirection"] = direction
// 	}
// }

type CreateOrderResponse struct {
	OrderID     string `json:"orderId,omitempty"`
	OrderLinkID string `json:"orderLinkId,omitempty"`
}

type WsResponse[Res any] struct {
	RetCode int    `json:"retCode"`
	RetMsg  string `json:"retMsg"`
	Results Res    `json:"result"`
	Data    Res    `json:"data"`
}