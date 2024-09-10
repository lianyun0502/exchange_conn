package bybit

import (
	"net/http"

	"github.com/lianyun0502/exchange_conn/v2/http_client"
	"github.com/lianyun0502/exchange_conn/v2/http_client/consts"
	"github.com/sirupsen/logrus"
)


type ByBitClient struct{
	*exchange_conn.HttpClient[*Request]
}

func (c *ByBitClient) Request(method, endpoint string, reqOpts ...func(*Request)) *Request {
	req := c.HttpClient.Request(method, endpoint, reqOpts...)
	req.SetGenHttpRequest(WithBybitRequest(req, c.HttpClient.Exchange, c.Log))
	return req
}


func NewClient(apiKey, secretKey string, apiOpts...func(*ByBitClient)) *ByBitClient {
	client := &ByBitClient{
		HttpClient: &exchange_conn.HttpClient[*Request]{
			Client: http.DefaultClient,
			Exchange: &exchange_conn.ExchangeApi{
				Name: consts.Bybit,
				HostType: consts.Spot,
				APIKey: apiKey,
				SecretKey: secretKey,
				BaseURL: "https://api.bybit.com",
			},
			NewRequest: NewRequest,
			Log: logrus.New(),
		},
	}

	for _, opt := range apiOpts {
		opt(client)
	}
	return client
}

func WithTestNet(isTest bool) func(*ByBitClient) {
	return func(c *ByBitClient) {
		if isTest {
			c.HttpClient.Exchange.BaseURL = "https://api-testnet.bybit.com"
		}else{
			c.HttpClient.Exchange.BaseURL = "https://api.bybit.com"
		}
	}
}