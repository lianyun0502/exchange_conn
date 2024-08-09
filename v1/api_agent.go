package exchange_conn

import (
	"net/http"
)

type IExchange interface {
	Request(string, string, bool, bool, ...any) IRequest
	SetRequest(IRequest) (*http.Request, error)
	Call(*http.Request) ([]byte, error)
}

type IRequest interface {
	SetQuery(key string, value interface{}) IRequest
	SetParam(key string, value interface{}) IRequest
	SetQueries(map[string]interface{}) IRequest
	SetParams(map[string]interface{}) IRequest
}

type APIAgent[E IExchange] struct {
	Client  E
	request IRequest
}

func NewAgent[E IExchange](client E) *APIAgent[E] {
	return &APIAgent[E]{
		Client: client,
	}
}

func (a *APIAgent[E]) Request(method string, endpoint string, key bool, signed bool, args ...any) *APIAgent[E] {
	a.request = a.Client.Request(method, endpoint, key, signed, args...)
	return a
}

func (a *APIAgent[E]) SetQuery(key string, value any) *APIAgent[E] {
	a.request.SetQuery(key, value)
	return a
}

func (a *APIAgent[E]) SetQueries(params map[string]any) *APIAgent[E] {
	a.request.SetQueries(params)
	return a
}

func (a *APIAgent[E]) SetParam(key string, value any) *APIAgent[E] {
	a.request.SetParam(key, value)
	return a
}

func (a *APIAgent[E]) SetParams(params map[string]any) *APIAgent[E] {
	a.request.SetParams(params)
	return a
}

func (a *APIAgent[E]) Send() (data []byte, err error) {
	req, err := a.Client.SetRequest(a.request)
	if err != nil {
		return
	}
	data, err = a.Client.Call(req)
	if err != nil {
		return
	}
	return
}
