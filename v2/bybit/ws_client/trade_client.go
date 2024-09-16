package bybit

import (
	"encoding/json"
	"time"
	"strconv"

	"github.com/lianyun0502/exchange_conn/v2/common"
	"github.com/lianyun0502/exchange_conn/v2/http_client/consts"
	"github.com/lianyun0502/exchange_conn/v2/ws_client"
)
func NewWsTradeClient(apiKey, secretKey string, opts ...func(*WsBybitClient)) (client *WsBybitClient) {
	exchInfo := &exchange_conn.ExchangeApi{
		Name: consts.Bybit,
		HostType: consts.Trade,
		APIKey: apiKey,
		SecretKey: secretKey,
		BaseURL: WEBSOCKET_TRADE_MAINNET,
	}
	client = &WsBybitClient{
		WsClient: exchange_conn.NewWsClient(exchInfo, nil), 
		maxAliveTime: "",
	}
	for _, opt := range opts {
		opt(client)
	}
	return client
}


func (wsc *WsBybitClient) Order(op string) {
	uuid := common.GetUUID()

	request := &Request{
		ReqID: uuid,
		Header: &RequestHeader{
			Timestamp: time.Now().UnixMicro(),
			RecvWindow: 5000,
		},
		Op: op,
	}
	req, err := json.Marshal(request)
	if err != nil {
		return
	}
	wsc.Send(req)
}


type RequestHeader struct {	
	Timestamp int64 `json:"X-BAPI-TIMESTAMP"`
	RecvWindow int64 `json:"X-BAPI-RECV-WINDOW"`
	Referer string `json:"Referer"`
}

type Request struct {
	ReqID string `json:"reqId"`
	Header *RequestHeader `json:"header"`
	Op string `json:"op"`
	Args string `json:"args"`
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


func Order(category, symbol, side, orderType, qty string, orderOpts...func(map[string]string)) map[string]string {
	args := make(map[string]string)
	args["category"] = category
	args["symbol"] = symbol
	args["side"] = side
	args["orderType"] = orderType
	args["qty"] = qty
	for _, opt := range orderOpts {
		opt(args)
	}
	return args
}

/*
是否借貸. 僅統一帳戶的現貨交易有效. 

	0(default): 否，則是幣幣訂單
	1: 是，則是槓桿訂單
*/
func isLeverage(isLeverage int) func(map[string]string) {
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
/*
direction:

	1: 當市場價上漲到了triggerPrice時觸發條件單
	2: 當市場價下跌到了triggerPrice時觸發條件單

*/
func Trigger(price, direction string) func(map[string]string) {
	return func(args map[string]string) {
		if (args["category"] != "linear") || (args["category"] != "inverse") {
			return
		}
		args["triggerDirection"] = direction
	}
}
