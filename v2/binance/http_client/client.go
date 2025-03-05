package binance

import (
	"fmt"
	"net/http"

	"github.com/lianyun0502/exchange_conn/v2/consts"
	"github.com/lianyun0502/exchange_conn/v2/http_client"
	"github.com/sirupsen/logrus"
)

type BinanceClient struct {
	*httpClient.HttpClient[*Request]
}

func (c *BinanceClient) Request(method, endpoint string, reqOpts ...func(*Request)) *Request {
	req := c.HttpClient.Request(method, endpoint, reqOpts...)
	req.SetGenHttpRequest(WithBinanceRequest(req, c.HttpClient.Exchange, c.Log))
	return req
}

func NewAPIClient(hostType string, apiKey, secretKey string, opts ...func(*BinanceClient)) (*BinanceClient, error) {
	var exchInfo *httpClient.ExchangeApi
	switch hostType {
	case consts.Spot:
		exchInfo = &httpClient.ExchangeApi{
			Name:      string(consts.Binance),
			HostType:  consts.Spot,
			APIKey:    apiKey,
			SecretKey: secretKey,
			BaseURL:   SPOT_MAINNET,
		}
	case consts.Future:
		exchInfo = &httpClient.ExchangeApi{
			Name:      string(consts.Binance),
			HostType:  consts.Future,
			APIKey:    apiKey,
			SecretKey: secretKey,
			BaseURL:   UFUTURE_MAINNET,
		}
	default:
		return nil, fmt.Errorf("hostType %s not supported", hostType)
	}
	client := &BinanceClient{
		HttpClient: &httpClient.HttpClient[*Request]{
			Client:     http.DefaultClient,
			Exchange:   exchInfo,
			NewRequest: NewRequest,
			Log:        logrus.New(),
		},
	}
	for _, opt := range opts {
		opt(client)
	}
	return client, nil
}

func NewAPISpotClient(apiKey, secretKey string, opts ...func(*BinanceClient)) (*BinanceClient, error) {
	return NewAPIClient(consts.Spot, apiKey, secretKey, opts...)
}

func NewAPIFutureClient(apiKey, secretKey string, opts ...func(*BinanceClient)) (*BinanceClient, error) {
	return NewAPIClient(consts.Future, apiKey, secretKey, opts...)
}

func IsTestNet() func(*BinanceClient) {
	return func(c *BinanceClient) {
		c.HttpClient.Exchange.BaseURL = SPOT_TESTNET
	}
}
