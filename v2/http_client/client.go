package exchange_conn

import (
	"net/http"
	"net/url"

	// "github.com/lianyun0502/exchange_conn/v1/http_client/consts"
	"github.com/sirupsen/logrus"
)

type ExchangeApi struct {
	Name      string
	HostType  string
	APIKey    string // API key
	SecretKey string // Secret key
	BaseURL   string // Base URL for API requests
}

type IRequest interface {
	B() []byte
	Q() url.Values
	SetB([]byte)
	SetQ(url.Values)
	SetSend(func(*http.Request) ([]byte, error))
	SetGenHttpRequest(func() (*http.Request, error))
}

type HttpClient[R IRequest] struct {
	Client *http.Client

	Exchange *ExchangeApi

	NewRequest     func(method string, endPoint string, reqOpts ...func(R)) R
	CurrentRequest R

	Log *logrus.Logger
}

func (c *HttpClient[R]) Request(method string, endpoint string, reqOpts ...func(R)) R {
	r := c.NewRequest(method, endpoint, reqOpts...)
	r.SetSend(WithSendFunction(c.Client, c.Log))
	return r
}
