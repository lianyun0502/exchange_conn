package binance

import (
	"fmt"
	"github.com/lianyun0502/exchange_conn/v2/consts"
	"github.com/lianyun0502/exchange_conn/v2/ws_client"
)

func NewWsQuoteClient(category string, quoteHandle func([]byte), opts ...func(*WsBinanceClient)) (*WsBinanceClient, error) {
	var exchInfo *wsClient.ExchangeApi
	switch category {
		case consts.Spot:
		exchInfo = &wsClient.ExchangeApi{
			Name: string(consts.Binance),
			HostType: consts.Spot,
			BaseURL: SPOT_QUOTE_MAINNET,
		}
		case consts.Future:
		exchInfo = &wsClient.ExchangeApi{
			Name: string(consts.Binance),
			HostType: consts.Future,
			BaseURL: USD_QUOTE_MAINNET,
		}
		default:
		return nil, fmt.Errorf("category error")

	}
	client := &WsBinanceClient{
		WsClient: wsClient.NewWsClient(exchInfo, nil), 
		ReceiveWindow: "5000",
	}
	opts = append(opts, WithWsHandle(quoteHandle))
	for _, opt := range opts {
		opt(client)
	}
	return client, nil
}