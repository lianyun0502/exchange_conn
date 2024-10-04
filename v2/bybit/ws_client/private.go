package bybit

import (
	"fmt"
	"github.com/lianyun0502/exchange_conn/v2/consts"
	"github.com/lianyun0502/exchange_conn/v2/ws_client"
	
)

func NewWsStreamClient(hostType, apiKey, secretKey string, opts ...func(*WsBybitClient)) (*WsBybitClient, error) {
	var client *WsBybitClient
	switch hostType {
		case consts.Private:
		client = &WsBybitClient{
			WsClient: wsClient.NewWsClient(&wsClient.ExchangeApi{
				Name: consts.Bybit,
				HostType: consts.Private,
				APIKey: apiKey,
				SecretKey: secretKey,
				BaseURL: WEBSOCKET_PRIVATE_MAINNET,
			}, nil),
			maxAliveTime: "",
		}
		default:
			return nil, fmt.Errorf("hostType error")
	}
	for _, opt := range opts {
		opt(client)
	}
	return client, nil
}

func NewWsPrivateClient(apiKey, secretKey string, opts ...func(*WsBybitClient)) (*WsBybitClient, error) {
	client, err := NewWsStreamClient(consts.Private, apiKey, secretKey, opts...)
	if err != nil {
		return nil, err
	}
	client.PingMessage = `{"op":"ping"}`
	client.Ping = client.PingServer
	return client, nil
}