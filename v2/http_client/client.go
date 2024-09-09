package exchange_conn

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	// "github.com/lianyun0502/exchange_conn/v1/http_client/consts"
	"github.com/sirupsen/logrus"
)

type ExchangeApi struct {
	Name     string
	HostType string
	APIKey    string // API key
	SecretKey string // Secret key
}

type IRequest interface {
	Body() []byte
	Query() url.Values
	SetBody([]byte)
	SetQuery(url.Values)
}

type HttpClient[R IRequest] struct {
	Client *http.Client

	Exchange *ExchangeApi
	BaseURL string // Base URL for API requests

	NewRequest func(method string,  endPoint string, reqOpts ...func(R)) R
	SetRequest func(r R, c *HttpClient[R]) (*http.Request, error)

	Log *logrus.Logger
}


func (c *HttpClient[R]) Request(method string, endpoint string, reqOpts ...func(R)) (data []byte, err error) {
	req, err := c.SetRequest(c.NewRequest(method, endpoint, reqOpts...), c)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// func (c *HttpClient[R]) setRequest(r R) (req *http.Request, err error) {
// 	return nil, fmt.Errorf("setReq method not implemented")
// }

func (c *HttpClient[R]) Do(r *http.Request) (data []byte, err error) {
	resp, err := c.Client.Do(r)
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

func SetQuery[R IRequest](query map[string]string) func(r R) {
	return func(r R) {
		r.SetQuery(make(url.Values))
		for k, v := range query {
			r.Query().Add(k, v)
		}
	}
}

// Json like: {"key": "value"}
func SetParam[R IRequest](param map[string]any) func(r R) {
	return func(r R) {
		if len(param) == 0 {
			r.SetBody(nil)
		} else {
			data, _ := json.Marshal(param)
			r.SetBody(data)
		}
	}
}

// Query like: key1=value1&key2=value2
func SetBody[R IRequest](body map[string]any) func(r R) {
	return func(r R) {
		p := make(url.Values)
		for k, v := range body {
			p.Add(k, fmt.Sprintf("%v", v))
		}
		r.SetBody([]byte(p.Encode()))
	}
}
