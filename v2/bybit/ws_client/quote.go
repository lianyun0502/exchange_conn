package bybit

import (
	"fmt"
	"strconv"

	"github.com/lianyun0502/exchange_conn/v2/consts"
	"github.com/lianyun0502/exchange_conn/v2/ws_client"
)

func NewWsQuoteClient(category string, quoteHandle func([]byte), opts ...func(*WsBybitClient)) (*WsBybitClient, error) {
	var exchInfo *exchange_conn.ExchangeApi
	switch category {
		case consts.Spot:
		exchInfo = &exchange_conn.ExchangeApi{
			Name: consts.Bybit,
			HostType: consts.Spot,
			BaseURL: SPOT_MAINNET,
		}
		case consts.Future:
		exchInfo = &exchange_conn.ExchangeApi{
			Name: consts.Bybit,
			HostType: consts.Future,
			BaseURL: LINEAR_MAINNET,
		}
		default:
		return nil, fmt.Errorf("category error")

	}
	client := &WsBybitClient{
		WsClient: exchange_conn.NewWsClient(exchInfo, nil), 
		maxAliveTime: "",
	}
	opts = append(opts, WithWsHandle(quoteHandle))
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