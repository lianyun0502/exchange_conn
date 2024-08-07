package bybit_conn

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/lianyun0502/exchange_conn/v1/common"
	"github.com/lianyun0502/exchange_conn/v1"
)

// Client define API client
type Client struct {
	APIKey     string
	SecretKey  string
	BaseURL    string
	HTTPClient *http.Client
	Debug      bool
	Logger     *log.Logger
	do         func(req *http.Request) (*http.Response, error)
}

// Client factory function
func NewClient(apiKey, secretKey, baseURL string) *Client {
	url := baseURL
	if baseURL == "" {
		url = "bybit.com"
	}
	return &Client{
		APIKey:     apiKey,
		SecretKey:  secretKey,
		BaseURL:    url,
		HTTPClient: http.DefaultClient,
	}
}

func (c *Client) Request(method, endpoint string, key, signed bool, opts ...any) exchange_conn.IRequest {
	sercType := None
	if key {sercType &= 0b10}
	if signed {sercType &= 0b01}
	req := NewByBitRequest(method, endpoint, sercType)
	return req
}

func (c *Client) parseRequest(r *Request, opts ...RequestOption) (req *http.Request, err error) {
	// set request options from user
	for _, opt := range opts {
		opt(r)
	}
	// err = r.validate()
	// if err != nil {
	// 	return err
	// }

	fullURL := fmt.Sprintf("%s%s", c.BaseURL, r.Endpoint)

	bodyString := r.Form.Encode()
	queryString := r.Query.Encode()

	if bodyString != "" {
		r.Body = bytes.NewBufferString(bodyString)
	}
	// header := http.Header{}

	req, err = http.NewRequest(r.Method, fullURL, r.Body)
	if err != nil {
		return
	}
	log.Printf("full url: %s\nrequese body: %s", req.URL.String(), common.PrettyPrint(r.Form))

	req.Header.Set("User-Agent", fmt.Sprintf("%s/%s", "bybit_connect", "v1"))
	if bodyString != "" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	if (r.SercType & Signed) != 0x00 {
		timeStamp := common.GetCurrentTime()
		req.Header.Set(signTypeKey, "2")
		req.Header.Set(apiRequestKey, c.APIKey)
		req.Header.Set(timestampKey, strconv.FormatInt(timeStamp, 10))
		if r.recvWindow == "" {
			r.recvWindow = "5000"
		}
		req.Header.Set(recvWindowKey, r.recvWindow)

		var signatureBase string
		if r.Method == "POST" {
			req.Header.Set("Content-Type", "application/json")
			signatureBase = strconv.FormatInt(timeStamp, 10) + c.APIKey + r.recvWindow + bodyString
		} else {
			signatureBase = strconv.FormatInt(timeStamp, 10) + c.APIKey + r.recvWindow + queryString
		}
		signature := common.GetSignature(c.SecretKey, signatureBase)
		req.Header.Set(signatureKey, signature)
	}
	if queryString != "" {
		fullURL = fmt.Sprintf("%s?%s", fullURL, queryString)
	}
	// c.debug("full url: %s, body: %s", fullURL, body)
	// r.fullURL = fullURL
	return 
}
func (c *Client) Call(req *http.Request, ctx context.Context) (data []byte, err error) {

	req = req.WithContext(ctx)

	f := c.do
	if f == nil {
		f = c.HTTPClient.Do
	}

	resp, err := f(req)
	if err != nil {
		log.Printf("Error: %s", err)
		return
	}
	defer func() {
		err = resp.Body.Close()
	}()


	if resp.StatusCode != 200 {
		log.Printf("Error: %s", resp.Status)
		return
	}
	// c.debug("response: %#v", res)
	// c.debug("response body: %s", string(data))
	// c.debug("response status code: %d", res.StatusCode)

	// log.Printf("response header: %s", PrettyPrint(resp.Header))

	return io.ReadAll(resp.Body)
}