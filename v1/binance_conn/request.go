package binance_conn

import (
	// "fmt"
	// "io"
	"fmt"
	"net/url"

	"github.com/lianyun0502/exchange_conn/v1"
)

const (
	Signed = 0b01
	ApiKey = 0b10
)

// Endpoint security type
//   - If no security type is stated, assume the security type is NONE.
//   - API-keys are passed into the Rest API via the X-MBX-APIKEY header.
//   - API-keys and secret-keys are case sensitive.
//   - API-keys can be configured to only access certain types of secure endpoints.
//     For example, one API-key could be used for TRADE only,
//     while another API-key can access everything except for TRADE routes.
//   - By default, API-keys can access all secure routes.
type SecurityT int

const (
	None       SecurityT = 0               // all public access
	Trade                = Signed | ApiKey // API-key and Singnature required
	UserData             = Signed | ApiKey // API-key and Singnature required
	UserStream           = ApiKey          // API-key required
	MarketData           = ApiKey          // API-key required
)

type request struct {
	exchange_conn.Request
	SercType SecurityT // security type
}

func NewRequest(method, endpoint string, sercType SecurityT) *request {
	return &request{Request: exchange_conn.Request{
		Method:   method,
		Endpoint: endpoint,

		Query: make(url.Values),
		Param: make(url.Values),
	},
		SercType: sercType,
	}
}

func (r *request) SetParam(key string, value interface{}) {
	r.Param.Set(key, fmt.Sprintf("%v", value))
}

func (r *request) SetParams(params map[string]interface{}) {
	for k, v := range params {
		r.SetParam(k, v)
	}
}

func WithTest(test string) func(*request) {
	return func(r *request) {
		r.Query.Set("test", test)
	}
}
