package binance_conn

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/lianyun0502/exchange_conn/v1"
	"github.com/lianyun0502/exchange_conn/v1/common"
)

var BaseURL = [6]string{
	"https://api.binance.com",
	"https://api-gcp.binance.com",
	"https://api1.binance.com",
	"https://api2.binance.com",
	"https://api3.binance.com",
	"https://api4.binance.com",
}

// type SecurityT int
// const (
// 	None SecurityT = iota // all public access
// 	Trade // API-key and Singnature required
// 	UserData // API-key and Singnature required
// 	UserStream // API-key required
// 	MARKET_DATA // API-key required
// )

type Client struct {
	APIKey     string // API key
	SecretKey  string // Secret key
	BaseURL    string // Base URL for API requests
	HTTPClient *http.Client
}

// Client factory function
func NewClient(apiKey, secretKey, baseURL string) *Client {
	url := baseURL
	if baseURL == "" {
		url = "https://api.binance.com"
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
	switch {
	case key && signed:
		sercType = Trade
	case key && !signed:
		sercType = UserStream
	}
	req := NewBinanceRequest(method, endpoint, sercType)
	return req
}

func (c *Client) SetRequest(r exchange_conn.IRequest) (req *http.Request, err error) {
	r_ptr := r.(*Request)
	return c.setBinanceRequest(r_ptr)
}

func (c *Client) setBinanceRequest(r *Request, opts ...RequsetOption) (req *http.Request, err error) {
	if r.SercType == Trade || r.SercType == UserData {
		r.Query.Set("timestamp", fmt.Sprintf("%v", time.Now().UnixNano()/int64(time.Millisecond)))
	}

	bodyString := r.Form.Encode()
	queryString := r.Query.Encode()

	if bodyString != "" {
		r.Body = bytes.NewBufferString(bodyString)
	}
	if r.SercType == Trade || r.SercType == UserData {
		r.Query.Set("signature", common.GetSignature(c.SecretKey, fmt.Sprintf("%s%s", queryString, bodyString)))
		queryString = r.Query.Encode()
	}

	fullURL := fmt.Sprintf("%s%s?%s", c.BaseURL, r.Endpoint, queryString)

	req, err = http.NewRequest(r.Method, fullURL, r.Body)
	if err != nil {
		return
	}
	log.Printf("full url: %s\nrequese body: %s", req.URL.String(), common.PrettyPrint(r.Form))

	req.Header.Set("User-Agent", fmt.Sprintf("%s/%s", "binance_connect", "v1"))
	if bodyString != "" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	if r.SercType != None {
		req.Header.Set("X-MBX-APIKEY", c.APIKey)
	}
	return
}

func (c *Client) Call(r *http.Request) (data []byte, err error) {
	resp, err := c.HTTPClient.Do(r)
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

	// log.Printf("response header: %s", PrettyPrint(resp.Header))

	data, err = io.ReadAll(resp.Body)

	return data, err
}
