package exchange_conn

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
)

type Request struct {
	Method   string
	Endpoint string
	body     []byte
	query    url.Values

	send           func(*http.Request) ([]byte, error)
	genHttpRequest func() (*http.Request, error)
}

func NewRequest(method, endpoint string) *Request {
	return &Request{
		Method:   method,
		Endpoint: endpoint,
		body:     nil,
		query:    make(url.Values),
	}
}

func (r *Request) SetSend(send func(*http.Request) ([]byte, error)) {
	r.send = send
}

func (r *Request) SetGenHttpRequest(gen func() (*http.Request, error)) {
	r.genHttpRequest = gen
}

func (r *Request) Send() ([]byte, error) {
	req, err := r.genHttpRequest()
	if err != nil {
		return nil, err
	}
	return r.send(req)
}

func (r *Request) SetParam(param ParamMap) *Request {
	return WithParam(r)(param)
}

func (r *Request) SetQuery(query QueryMap) *Request {
	return WithQuery(r)(query)
}

func (r *Request) SetB(body []byte) {
	r.body = body
}

func (r *Request) SetQ(query url.Values) {
	r.query = query
}

func (r *Request) B() []byte {
	return r.body
}

func (r *Request) Q() url.Values {
	return r.query
}

func WithSendFunction(client *http.Client, log *logrus.Logger) func(*http.Request) ([]byte, error) {
	Client := client
	Log := log
	return func(r *http.Request) ([]byte, error) {
		resp, err := Client.Do(r)
		if err != nil {
			Log.Printf("Error: %s", err)
			return nil, err
		}
		defer func() {
			err = resp.Body.Close()
		}()
		return io.ReadAll(resp.Body)

	}
}

type ParamMap map[string]any
type QueryMap map[string]string

func WithQuery[R IRequest](r R) func(QueryMap) R {
	return func(query QueryMap) R {
		r.SetQ(make(url.Values))
		for k, v := range query {
			r.Q().Add(k, v)
		}
		return r
	}
}

// Json like: {"key": "value"}
func WithParam[R IRequest](r R) func(ParamMap) R {
	return func(param ParamMap) R {
		if len(param) == 0 {
			r.SetB(nil)
		} else {
			data, _ := json.Marshal(param)
			r.SetB(data)
		}
		return r
	}
}

// Query like: key1=value1&key2=value2
func WithBody[R IRequest](r R) func(ParamMap) R {
	return func(body ParamMap) R {
		p := make(url.Values)
		for k, v := range body {
			p.Add(k, fmt.Sprintf("%v", v))
		}
		r.SetB([]byte(p.Encode()))
		return r
	}
}
