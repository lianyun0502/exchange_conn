package bybit

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/lianyun0502/exchange_conn/v2/common"
	"github.com/lianyun0502/exchange_conn/v2/http_client"
	"github.com/sirupsen/logrus"
	"strconv"
)

const (
	Signed = 0b01
	ApiKey = 0b10
)

type Request struct {
	*httpClient.Request

	SecType    int
	recvWindow string
}

func (r *Request) SetParam(param httpClient.ParamMap) *Request {
	// override the SetParam function
	return httpClient.WithParam(r)(param)
}

func (r *Request) SetQuery(query httpClient.QueryMap) *Request {
	// override the SetQuery function
	return httpClient.WithQuery(r)(query)
}

func NewRequest(method, endpoint string, reqOpts ...func(*Request)) *Request {
	req := &Request{
		Request:    httpClient.NewRequest(method, endpoint),
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
			r.SecType = 3
		case haskey && !signed:
			r.SecType = 2
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

func WithBybitRequest(req *Request, exchInfo *httpClient.ExchangeApi, log *logrus.Logger) func() (*http.Request, error) {
	Exch := exchInfo
	Log := log
	r := req
	return func() (*http.Request, error) {
		fullURL := fmt.Sprintf("%s%s", Exch.BaseURL, r.Endpoint)

		queryString := r.Q().Encode()
		if queryString != "" {
			fullURL = fmt.Sprintf("%s?%s", fullURL, queryString)
		}

		req, err := http.NewRequest(r.Method, fullURL, bytes.NewBuffer(r.B()))
		if err != nil {
			return nil, err
		}
		Log.WithFields(logrus.Fields{"url": fullURL}).Debug("Compose URL")
		if body := r.B(); body != nil {
			Log.Debugf(`Request Body : %s`, string(body))
		}
		req.Header.Set("User-Agent", fmt.Sprintf("%s/%s", "bybit_connect", "v1"))
		if r.B() != nil {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}

		if (r.SecType & Signed) != 0x00 {
			timeStamp := common.GetCurrentTime()
			req.Header.Set(signTypeKey, "2")
			req.Header.Set(apiRequestKey, Exch.APIKey)
			req.Header.Set(timestampKey, strconv.FormatInt(timeStamp, 10))
			req.Header.Set(recvWindowKey, r.recvWindow)

			var signatureBase string
			if r.Method == "POST" {
				req.Header.Set("Content-Type", "application/json")
				signatureBase = strconv.FormatInt(timeStamp, 10) + Exch.APIKey + r.recvWindow + string(r.B())
			} else {
				signatureBase = strconv.FormatInt(timeStamp, 10) + Exch.APIKey + r.recvWindow + queryString
			}
			Log.Debugf(`Signature : %s`, signatureBase)
			signature := common.GetSignature(Exch.SecretKey, signatureBase)
			req.Header.Set(signatureKey, signature)
		}
		return req, err
	}
}
