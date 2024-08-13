package exchange_conn

import (
	"net/http"
	"log"
	"io"
)


type Client struct {
	APIKey     string // API key
	SecretKey  string // Secret key
	BaseURL    string // Base URL for API requests
	HTTPClient *http.Client
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

	data, err = io.ReadAll(resp.Body)

	return data, err
}
