package exchange_conn

import (
	"net/http"
)

type IRequsetOption[R IRequest] interface {
	func(*R)
}

type IClient[R IRequest] interface {
	Request(string, string, bool, bool, ...func(R)) R
	SetRequest(R) (*http.Request, error)
	Call(*http.Request) ([]byte, error)
	// MakeOrder(params map[string]interface{}) (*utils.OrderResponse, *http.Response, error)
	// GetOrderDetail(params map[string]interface{}) (*utils.OrderMessage, *http.Response, error)
	// CancelOrder(params map[string]interface{}) (bool, *http.Response, error)
	// GetOpeningOrders(symbol string) ([]*utils.OrderMessage, *http.Response, error)
	// GetAccountBalance() (interface{}, *http.Response, error)
	// GetPositions() (interface{}, *http.Response, error)
	// Volume24h(symbol string) (*utils.Volume24hData, *http.Response, error)
	// Ticker24h(symbol string) (*utils.Ticker24hData, *http.Response, error)
	// BatchOperation(opType string, params map[string]interface{}) (interface{}, error)
	// BatchMakeOrder(orders []map[string]interface{}) ([]*ReqResponseItem[*utils.OrderResponse], error)
	// BatchCancelOrder(symbol string, orderIDs []string, clientIDs []string) ([]*ReqResponseItem[bool], error)
	// GetClockDiff(int) (int64, error)
	// GetTransferRecords(map[string]interface{}) (*utils.TransferRecordsResult, *http.Response, error)
}

type IRequest interface {
	SetQuery(key string, value interface{})
	SetParam(key string, value interface{})
	SetQueries(map[string]interface{})
	SetParams(map[string]interface{})
}

type apiAgent[E IClient[R], R IRequest] struct {
	Client  E
	request R
}

func NewAgent[E IClient[R], R IRequest](client E, opts ...func(E)) *apiAgent[E, R] {
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
