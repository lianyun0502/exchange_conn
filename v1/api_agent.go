package exchange_conn

import (
	"net/http"
)

type IRequsetOption[R IRequest] interface {
	func(*R)
}

type IExchange[R IRequest] interface {
	Request(string, string, bool, bool, ...func(R)) R
	SetRequest(R) (*http.Request, error)
	Call(*http.Request) ([]byte, error)
}

type IRequest interface {
	SetQuery(key string, value interface{})
	SetParam(key string, value interface{})
	SetQueries(map[string]interface{})
	SetParams(map[string]interface{})
}

type apiAgent[E IExchange[R], R IRequest] struct {
	Client  E
	request R
}

func NewAgent[E IExchange[R], R IRequest](client E, opts ...func(E)) *apiAgent[E, R] {
	for _, opt := range opts {
		opt(client)
	} 
	return &apiAgent[E, R]{
		Client: client,
	}
}

func (a *apiAgent[E, R]) Request(method string, endpoint string, key bool, signed bool, reqOpts ...func(R)) *apiAgent[E, R] {
	a.request = a.Client.Request(method, endpoint, key, signed, reqOpts...)
	return a
}

func (a *apiAgent[E, R]) SetQuery(key string, value any) *apiAgent[E, R] {
	a.request.SetQuery(key, value)
	return a
}

func (a *apiAgent[E, R]) SetQueries(params map[string]any) *apiAgent[E, R] {
	a.request.SetQueries(params)
	return a
}

func (a *apiAgent[E, R]) SetParam(key string, value any) *apiAgent[E, R] {
	a.request.SetParam(key, value)
	return a
}

func (a *apiAgent[E, R]) SetParams(params map[string]any) *apiAgent[E, R] {
	a.request.SetParams(params)
	return a
}

func (a *apiAgent[E, R]) Send() (data []byte, err error) {
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
