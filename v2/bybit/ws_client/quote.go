package bybit


import (
	"strconv"

	"github.com/lianyun0502/exchange_conn/v2/ws_client"
	"github.com/lianyun0502/exchange_conn/v2/consts"

)

func NewWsSpotClient(opts ...func(*WsBybitClient)) (client *WsBybitClient) {
	exchInfo := &exchange_conn.ExchangeApi{
		Name: consts.Bybit,
		HostType: consts.Spot,
		BaseURL: SPOT_MAINNET,
	}
	client = &WsBybitClient{
		WsClient: exchange_conn.NewWsClient(exchInfo, nil), 
		maxAliveTime: "",
	}
	client.Ws_Handler = WithBybitHandler(client.ReqMap, nil)
	for _, opt := range opts {
		opt(client)
	}
	return client
}

func NewWsFutureClient(opts ...func(*WsBybitClient)) (client *WsBybitClient) {
	exchInfo := &exchange_conn.ExchangeApi{
		Name: consts.Bybit,
		HostType: consts.Future,
		BaseURL: LINEAR_MAINNET,
	}
	client = &WsBybitClient{
		WsClient: exchange_conn.NewWsClient(exchInfo, nil), 
		maxAliveTime: "",
	}
	client.Ws_Handler = WithBybitHandler(client.ReqMap, nil)
	for _, opt := range opts {
		opt(client)
	}
	return client
}

// 針對私有頻道和交易, 您可以自定義連接存活時長, 通過增加參數max_active_time, 最小支持30s (30秒), 最大支持600s (10分鐘)
func WithMaxAliveTime(maxAliveTime int) func(*WsBybitClient) {
	return func(wsc *WsBybitClient) {
		wsc.maxAliveTime = strconv.Itoa(maxAliveTime)
	}
}