package binance_conn

import (
	"net/url"
	"io"
	"fmt"

	"github.com/lianyun0502/exchange_conn/v1"
)

type Request struct {
	Method   string // http method
	Endpoint string // every api specific url
	SercType SecurityT // security type

	Body     io.Reader 
	Query    url.Values // query string
	Form     url.Values // extually is form data, covert to body in the end
	// header   http.Header
}
func NewBinanceRequest(method, endpoint string, sercType SecurityT) *Request {
	return &Request{
		Method: method,
		Endpoint: endpoint,
		SercType: sercType,

		Query: url.Values{},
		Form: url.Values{},
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
func (r *Request) SetParam(key string, value interface{}) exchange_conn.IRequest {
	if r.Form.Get(key) == "" {
		r.Form.Add(key, fmt.Sprintf("%v", value))
		return r
	}
	r.Form.Set(key, fmt.Sprintf("%v", value))
	return r
}

type RequsetOption func(*url.Values)


