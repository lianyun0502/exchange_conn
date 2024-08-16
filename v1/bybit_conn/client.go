package bybit_conn

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/lianyun0502/exchange_conn/v1/common"
	log "github.com/sirupsen/logrus"
)

// Client define API client
type Client struct {
	APIKey     string
	SecretKey  string
	BaseURL    string
	HTTPClient *http.Client

	Debug  bool
	Logger *log.Logger
	do     func(req *http.Request) (*http.Response, error)
}

// Client factory function
func NewClient(apiKey, secretKey, baseURL string) *Client {
	url := baseURL
	if baseURL == "" {
		url = "https://api.bybit.com"
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
	req := NewByBitRequest(method, endpoint, sercType)
	return req
}

func (c *Client) SetRequest(r *request) (req *http.Request, err error) {
	fullURL := fmt.Sprintf("%s%s", c.BaseURL, r.Endpoint)

	bodyString := r.Form.Encode()
	queryString := r.Query.Encode()

	if bodyString != "" {
		r.Body = bytes.NewBufferString(bodyString)
	}
	if queryString != "" {
		fullURL = fmt.Sprintf("%s?%s", fullURL, queryString)
	}

	req, err = http.NewRequest(r.Method, fullURL, r.Body)
	if err != nil {
		return
	}
	c.Logger.WithFields(log.Fields{"url": fullURL}).Debug("Compose URL")
	c.Logger.WithFields(log.Fields{"body": r.Form}).Debug("Body Ready")

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
		c.Logger.WithFields(log.Fields{"signatureBase": signatureBase}).Debug("Signature Base")
		signature := common.GetSignature(c.SecretKey, signatureBase)
		req.Header.Set(signatureKey, signature)
	}
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
