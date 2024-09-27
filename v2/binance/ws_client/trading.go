package binance

import (
	// "strconv"
	// "time"
	"fmt"

	// "github.com/valyala/fastjson"
	"github.com/lianyun0502/exchange_conn/v2/consts"
	"github.com/lianyun0502/exchange_conn/v2/ws_client"
)

func NewWsAPIClient(hostType string, apiKey, secretKey string, opts ...func(*WsBinanceClient)) (*WsBinanceClient, error) {
	var exchInfo *wsClient.ExchangeApi
	switch hostType {
		case consts.Spot:
			exchInfo = &wsClient.ExchangeApi{
				Name: string(consts.Binance),
				HostType: consts.Spot,
				APIKey: apiKey,
				SecretKey: secretKey,
				BaseURL: SPOT_MAINNET,
			}
		case consts.Future:
			exchInfo = &wsClient.ExchangeApi{
				Name: string(consts.Binance),
				HostType: consts.Future,
				APIKey: apiKey,
				SecretKey: secretKey,
				BaseURL: USD_MAINNET,
			}
		default:
			return nil, fmt.Errorf("hostType error")
	}
	client := &WsBinanceClient{
		WsClient: wsClient.NewWsClient(exchInfo, nil), 
		ReceiveWindow: "5000",
	}
	opts = append(opts, WithWsHandle(nil))
	for _, opt := range opts {
		opt(client)
	}
	return client, nil
}