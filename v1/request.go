package exchange_conn

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

type Params map[string]interface{}

func (p *Params) Set(params map[string]interface{}) {
	for k, v := range params {
		(*p)[k] = v
	}
}

func (p *Params) SetString(params string) {
	j := make(map[string]interface{})
	json.Unmarshal([]byte(params), &j)
	p.Set(j)
}

func (p *Params) Add(key string, value interface{}) {
	(*p)[key] = value
}

func (p *Params) Get(key string) any {
	return (*p)[key]
}

func (p *Params) Encode() string {
	data, _ := json.Marshal(p)
	ret := string(data)
	if ret == "{}" {
		return ""
	}
	return ret
}

type Request struct {
	Method   string // http method
	Endpoint string // every api specific url

	Body  io.Reader
	Query url.Values // query string
	Param url.Values // body string
	Form  Params     // extually is form data, covert to body in the end
}

func NewRequest(method, endpoint string) *Request {
	return &Request{
		Method:   method,
		Endpoint: endpoint,
		// SercType: sercType,

		Query: make(url.Values),
		Form:  make(Params),
	}
}

func (r *Request) SetQuery(key string, value any) {
	if r.Query.Get(key) == "" {
		r.Query.Add(key, fmt.Sprintf("%v", value))
	}
	r.Query.Set(key, fmt.Sprintf("%v", value))

}
func (r *Request) SetQueries(params map[string]any) {
	for k, v := range params {
		r.SetQuery(k, v)
	}

}
func (r *Request) SetParam(key string, value any) {
	if r.Method != http.MethodPost {
		return
	}
	r.Form.Add(key, fmt.Sprintf("%v", value))

}
func (r *Request) SetParamsString(params string) {
	if r.Method != http.MethodPost {
		return
	}
	r.Form.SetString(params)

}

func (r *Request) SetParams(params map[string]interface{}) {
	if r.Method != http.MethodPost {
		return
	}
	r.Form.Set(params)
}
