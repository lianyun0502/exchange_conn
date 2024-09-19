package binance

import (
	"net/http"

	"github.com/lianyun0502/exchange_conn/v2/http_client"
	"github.com/lianyun0502/exchange_conn/v2/consts"
	"github.com/sirupsen/logrus"
)

type BinanceSpotClient struct {
	*exchange_conn.HttpClient[*Request]
}

func (c *BinanceSpotClient) Request(method, endpoint string, reqOpts ...func(*Request)) *Request {
	req := c.HttpClient.Request(method, endpoint, reqOpts...)
	req.SetGenHttpRequest(WithBinanceRequest(req, c.HttpClient.Exchange, c.Log))
	return req
}

func NewSpotClient(apiKey, secretKey string, apiOpts ...func(*BinanceSpotClient)) *BinanceSpotClient {
	client := &BinanceSpotClient{
		HttpClient: &exchange_conn.HttpClient[*Request]{
			Client: http.DefaultClient,
			Exchange: &exchange_conn.ExchangeApi{
				Name:      string(consts.Binance),
				HostType:  consts.Spot,
				APIKey:    apiKey,
				SecretKey: secretKey,
				BaseURL:   SPOT_MAINNET,
			},
			NewRequest: NewRequest,
			Log:        logrus.New(),
		},
	}

	for _, opt := range apiOpts {
		opt(client)
	}
	return client
}

func WithTestNet() func(*BinanceSpotClient) {
	return func(c *BinanceSpotClient) {
		c.HttpClient.Exchange.BaseURL = SPOT_TESTNET
	}
}
