package httpClient

import (
	"net/http"
)


func (c *HttpClient[R]) MakeOrder(params map[string]interface{}) (*OrderResponse, *http.Response, error) {
	c.Log.Warnf("NotSupport:HttpRequest:%s:%s", c.Exchange.Name, c.Exchange.HostType)
	return nil, nil, nil
}
func (c *HttpClient[R]) GetOrderDetail(params map[string]interface{}) (*OrderMessage, *http.Response, error) {
	c.Log.Warnf("NotSupport:HttpRequest:%s:%s", c.Exchange.Name, c.Exchange.HostType)
	return nil, nil, nil
}
func (c *HttpClient[R]) CancelOrder(params map[string]interface{}) (bool, *http.Response, error) {
	c.Log.Warnf("NotSupport:HttpRequest:%s:%s", c.Exchange.Name, c.Exchange.HostType)
	return false, nil, nil
}
func (c *HttpClient[R]) GetOpeningOrders(symbol string) ([]*OrderMessage, *http.Response, error) {
	c.Log.Warnf("NotSupport:HttpRequest:%s:%s", c.Exchange.Name, c.Exchange.HostType)
	return nil, nil, nil
}
func (c *HttpClient[R]) GetAccountBalance() (interface{}, *http.Response, error) {
	c.Log.Warnf("NotSupport:HttpRequest:%s:%s", c.Exchange.Name, c.Exchange.HostType)
	return nil, nil, nil
}
func (c *HttpClient[R]) GetPositions() (interface{}, *http.Response, error) {
	c.Log.Warnf("NotSupport:HttpRequest:%s:%s", c.Exchange.Name, c.Exchange.HostType)
	return nil, nil, nil
}
func (c *HttpClient[R]) Volume24h(symbol string) (*Volume24hData, *http.Response, error) {
	c.Log.Warnf("NotSupport:HttpRequest:%s:%s", c.Exchange.Name, c.Exchange.HostType)
	return nil, nil, nil
}
func (c *HttpClient[R]) Ticker24h(symbol string) (*Ticker24hData, *http.Response, error) {
	c.Log.Warnf("NotSupport:HttpRequest:%s:%s", c.Exchange.Name, c.Exchange.HostType)
	return nil, nil, nil
}
func (c *HttpClient[R]) BatchOperation(opType string, params map[string]interface{}) (interface{}, error) {
	c.Log.Warnf("NotSupport:HttpRequest:%s:%s", c.Exchange.Name, c.Exchange.HostType)
	return nil, nil
}
func (c *HttpClient[R]) BatchMakeOrder(orders []map[string]interface{}) ([]*ReqResponseItem[*OrderResponse], error) {
	c.Log.Warnf("NotSupport:HttpRequest:%s:%s", c.Exchange.Name, c.Exchange.HostType)
	return nil, nil
}
func (c *HttpClient[R]) BatchCancelOrder(symbol string, orderIDs []string, clientIDs []string) ([]*ReqResponseItem[bool], error) {
	c.Log.Warnf("NotSupport:HttpRequest:%s:%s", c.Exchange.Name, c.Exchange.HostType)
	return nil, nil
}
func (c *HttpClient[R]) GetClockDiff(int) (int64, error) {
	c.Log.Warnf("NotSupport:HttpRequest:%s:%s", c.Exchange.Name, c.Exchange.HostType)
	return 0, nil
}
