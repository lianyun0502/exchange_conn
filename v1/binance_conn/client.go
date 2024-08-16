package binance_conn

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/lianyun0502/exchange_conn/v1/common"
	log "github.com/sirupsen/logrus"
)

var BaseURL = [6]string{
	"https://api.binance.com",
	"https://api-gcp.binance.com",
	"https://api1.binance.com",
	"https://api2.binance.com",
	"https://api3.binance.com",
	"https://api4.binance.com",
}

type Client struct {
	APIKey     string // API key
	SecretKey  string // Secret key
	BaseURL    string // Base URL for API requests
	HTTPClient *http.Client

	Logger *log.Logger
}

// Client factory function
func NewClient(apiKey, secretKey, baseURL string) *Client {
	url := baseURL
	if baseURL == "" {
		url = "https://api.binance.com"
	}
	logger := log.New()
	logger.SetFormatter(&log.TextFormatter{TimestampFormat: "2006-01-02 15:04:05.000000", FullTimestamp: true})
	return &Client{
		APIKey:     apiKey,
		SecretKey:  secretKey,
		BaseURL:    url,
		HTTPClient: http.DefaultClient,
		Logger:     logger,
	}
}

func (c *Client) Request(method, endpoint string, key, signed bool, opts ...func(*request)) *request {
	sercType := None
	switch {
	case key && signed:
		sercType = Trade
	case key && !signed:
		sercType = UserStream
	}
	req := NewRequest(method, endpoint, sercType)
	return req
}

func (c *Client) SetRequest(r *request) (req *http.Request, err error) {
	if r.SercType == Trade || r.SercType == UserData {
		r.Query.Set("timestamp", fmt.Sprintf("%v", time.Now().UnixNano()/int64(time.Millisecond)))
	}

	bodyString := r.Param.Encode()
	queryString := r.Query.Encode()

	if bodyString != "" {
		r.Body = bytes.NewBufferString(bodyString)
	}

	if r.SercType == Trade || r.SercType == UserData {
		signatureBase := queryString + bodyString
		c.Logger.WithFields(log.Fields{"signatureBase": signatureBase}).Debug("Signature Base")
		r.Query.Set("signature", common.GetSignature(c.SecretKey, signatureBase))
		queryString = r.Query.Encode()
	}
	fullURL := fmt.Sprintf("%s%s", c.BaseURL, r.Endpoint)
	if queryString != "" {
		fullURL = fmt.Sprintf("%s?%s", fullURL, queryString)
	}
	c.Logger.WithFields(log.Fields{"url": fullURL}).Debug("Compose URL")
	c.Logger.WithFields(log.Fields{"body": r.Form}).Debug("Body Ready")
	req, err = http.NewRequest(r.Method, fullURL, r.Body)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", fmt.Sprintf("%s/%s", "binance_connect", "v1"))
	if bodyString != "" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	if r.SercType != None {
		req.Header.Set("X-MBX-APIKEY", c.APIKey)
	}
	c.Logger.WithFields(log.Fields{"header": req.Header}).Debug("Header Ready")
	return
}

func (c *Client) Call(r *http.Request) (data []byte, err error) {
	c.Logger.WithFields(log.Fields{"method": r.Method, "url": r.URL}).Info("Send Request")
	c.Logger.WithFields(log.Fields{"header": r.Header, "body": r.Body}).Debug("Request Content")
	resp, err := c.HTTPClient.Do(r)
	if err != nil {
		c.Logger.Errorf("Error: %s", err)
		return
	}
	defer resp.Body.Close()
	c.Logger.WithFields(log.Fields{"status": resp.Status}).Info("Response Status")
	return io.ReadAll(resp.Body)
}

type ClientOption func(*Client)

func WithLogger(logger *log.Logger) ClientOption {
	return func(c *Client) {
		c.Logger = logger
	}
}
