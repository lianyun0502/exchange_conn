package binance

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"strconv"

	"github.com/lianyun0502/exchange_conn/v2/common"
	"github.com/lianyun0502/exchange_conn/v2/http_client"
	"github.com/sirupsen/logrus"
)

const (
	Signed = 0b01
	ApiKey = 0b10
)

type Request struct {
	*exchange_conn.Request

	SercType   int
	recvWindow string
}

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

func (r *Request) SetParam(param exchange_conn.ParamMap) *Request {
	// override the SetParam function
	return exchange_conn.WithBody(r)(param)
}

func (r *Request) SetQuery(query exchange_conn.QueryMap) *Request {
	// override the SetQuery function
	return exchange_conn.WithQuery(r)(query)
}

func NewRequest(method, endpoint string, reqOpts ...func(*Request)) *Request {
	req := &Request{
		Request:    exchange_conn.NewRequest(method, endpoint),
		recvWindow: "5000",
	}
	for _, opt := range reqOpts {
		opt(req)
	}
	return req
}

/*
function that set the security type for the request

  - signed : enable to carry signature
  - haskey : enable to carry apikey
*/
func SetSercurityType(signed, haskey bool) func(*Request) {
	return func(r *Request) {
		switch {
		case haskey && signed:
			r.SercType = 3
		case haskey && !signed:
			r.SercType = 2
		}
	}
}

/*
Default recvWindow is 5000

function that set the recive window time(ms) for the request
*/
func SetRecvWindow(time int) func(*Request) {
	return func(r *Request) {
		r.recvWindow = strconv.Itoa(time)
	}
}

func WithBinanceRequest(req *Request, exchInfo *exchange_conn.ExchangeApi, log *logrus.Logger) func() (*http.Request, error) {
	Exch := exchInfo
	Log := log
	r := req
	return func() (*http.Request, error) {
		if r.SercType&Signed != 0 {
			r.Q().Set("timestamp", fmt.Sprintf("%v", time.Now().UnixNano()/int64(time.Millisecond)))
		}

		queryString := r.Q().Encode()
		bodyString := string(r.B())

		if r.SercType == Trade || r.SercType == UserData {
			signatureBase := queryString + bodyString
			Log.WithFields(logrus.Fields{"signatureBase": signatureBase}).Debug("Signature Base")
			r.Q().Set("signature", common.GetSignature(Exch.SecretKey, signatureBase))
			queryString = r.Q().Encode()
		}
		fullURL := fmt.Sprintf("%s%s", Exch.BaseURL, r.Endpoint)
		if queryString != "" {
			fullURL = fmt.Sprintf("%s?%s", fullURL, queryString)
		}
		Log.WithFields(logrus.Fields{"url": fullURL}).Debug("Compose URL")
		if body := r.B(); body != nil {
			Log.WithFields(logrus.Fields{"body": string(body)}).Debug("Body Ready")
		}
		req, err := http.NewRequest(r.Method, fullURL, bytes.NewBuffer(r.B()))
		if err != nil {
			return nil, err
		}
		req.Header.Set("User-Agent", fmt.Sprintf("%s/%s", "binance_connect", "v1"))
		if bodyString != "" {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
		if r.SercType != 0 {
			req.Header.Set("X-MBX-APIKEY", Exch.APIKey)
		}
		Log.WithFields(logrus.Fields{"header": req.Header}).Debug("Header Ready")
		return req, err
	}
}
