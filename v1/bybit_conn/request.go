package bybit_conn

import (
	"fmt"
	"io"
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
	None        SecurityT = 0               // all public access
	Trade                 = Signed | ApiKey // API-key and Singnature required
	UserData              = Signed | ApiKey // API-key and Singnature required
	UserStream            = ApiKey          // API-key required
	MARKET_DATA           = ApiKey          // API-key required
)

type Request struct {
	Method   string    // http method
	Endpoint string    // every api specific url
	SercType SecurityT // security type

	Body  io.Reader
	Query url.Values // query string
	Form  url.Values // extually is form data, covert to body in the end
	recvWindow string
}

func NewByBitRequest(method, endpoint string, sercType SecurityT) *Request {
	return &Request{
		Method:   method,
		Endpoint: endpoint,
		SercType: sercType,

		Query: url.Values{},
		Form:  url.Values{},
	}
}

func (r *Request) SetQuery(key string, value interface{}) exchange_conn.IRequest {
	if r.Query.Get(key) == "" {
		r.Query.Add(key, fmt.Sprintf("%v", value))
		return r
	}
	r.Query.Set(key, fmt.Sprintf("%v", value))
	return r
}
func (r *Request) SetQueries(params map[string]interface{}) exchange_conn.IRequest {
	for k, v := range params {
		r.SetQuery(k, v)
	}
	return r
}
func (r *Request) SetParam(key string, value interface{}) exchange_conn.IRequest {
	if r.Form.Get(key) == "" {
		r.Form.Add(key, fmt.Sprintf("%v", value))
		return r
	}
	r.Form.Set(key, fmt.Sprintf("%v", value))
	return r
}

func (r *Request) SetParams(params map[string]interface{}) exchange_conn.IRequest {
	for k, v := range params {
		r.SetQuery(k, v)
	}
	return r
}

type RequestOption func(*Request)
