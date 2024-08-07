package exchange_conn

import (
	"net/http"
)



type IExchange interface {
	Request(string, string, bool, bool, ...any) (IRequest)
	SetRequest(IRequest) (*http.Request, error)
	Call(*http.Request) ([]byte, error)
}

type IRequest interface {
	SetQuery(key string, value interface{}) (IRequest)
	SetParam(key string, value interface{}) (IRequest)
	SetQueries(map[string]interface{}) (IRequest)
	SetParams(map[string]interface{}) (IRequest)
}


type APIAgent struct {
	Exchange IExchange
	request  IRequest
}

func NewAgent(ex IExchange) *APIAgent {
	return &APIAgent{
		Exchange: ex,
	}
}

func (a *APIAgent) Request(method string, endpoint string, key bool, signed bool, args ...any) *APIAgent {
	a.request = a.Exchange.Request(method, endpoint, key, signed, args)
	return a
}

func (a *APIAgent) SetQuery(key string, value any) *APIAgent {
	a.request.SetQuery(key, value)
	return a
}

func (a *APIAgent) SetQueries(params map[string]any) *APIAgent {
	a.request.SetQueries(params)
	return a
}

func (a *APIAgent) SetParam(key string, value any) *APIAgent {
	a.request.SetParam(key, value)
	return a
}

func (a *APIAgent) SetParams(params map[string]any) *APIAgent {
	a.request.SetParams(params)
	return a
}

func (a *APIAgent) Send() (data []byte, err error) {
	req, err := a.Exchange.SetRequest(a.request)
	if err != nil {
		return
	}
	data, err = a.Exchange.Call(req)
	if err != nil {
		return
	}
	return
}
