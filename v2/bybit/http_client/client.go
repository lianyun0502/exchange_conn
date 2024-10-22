package bybit

import (
	"net/http"

	"github.com/lianyun0502/exchange_conn/v2/http_client"
	"github.com/lianyun0502/exchange_conn/v2/consts"
	"github.com/sirupsen/logrus"
)


type ByBitClient struct{
	*httpClient.HttpClient[*Request]
}

func (c *ByBitClient) Request(method, endpoint string, reqOpts ...func(*Request)) *Request {
	c.Log.Debug("Gen Request")
	req := c.HttpClient.Request(method, endpoint, reqOpts...)
	req.SetGenHttpRequest(WithBybitRequest(req, c.HttpClient.Exchange, c.Log))
	return req
}


func NewSpotClient(apiKey, secretKey string, apiOpts...func(*ByBitClient)) *ByBitClient {
	client := &ByBitClient{
		HttpClient: &httpClient.HttpClient[*Request]{
			Client: http.DefaultClient,
			Exchange: &httpClient.ExchangeApi{
				Name: consts.Bybit,
				HostType: consts.Spot,
				APIKey: apiKey,
				SecretKey: secretKey,
				BaseURL: MAINNET,
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

func IsTestNet() func(*ByBitClient) {
	return func(c *ByBitClient) {
		switch c.Exchange.HostType{
		case consts.Spot:
			c.HttpClient.Exchange.BaseURL = TESTNET
		}
	}
}