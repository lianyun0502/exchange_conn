package exchange_conn

import "net/http"

type IWsClient interface {
	AddEvent(e *Event)
	Reconnect()
	StartLoop()
	Connect(string) (*http.Response, error)
	Stop() error
	Send([]byte) error
	Subscribe(string)
}

type webSocketAgent[T IWsClient] struct {
	Client T
}

func NewWebSocketAgent[T IWsClient](client T) *webSocketAgent[T] {
	return &webSocketAgent[T]{
		Client: client,
	}
}

func (a *webSocketAgent[T]) Subscribe(params string) {
	a.Client.Subscribe(params)
}

func (a *webSocketAgent[T]) Send(request []byte) error {
	return a.Client.Send(request)
}

func (a *webSocketAgent[T]) Connect(url string) (*http.Response, error) {
	return a.Client.Connect(url)
}

func (a *webSocketAgent[T]) Stop() error {
	return a.Client.Stop()
}

func (a *webSocketAgent[T]) StartLoop() {
	a.Client.StartLoop()
}

func (a *webSocketAgent[T]) SendString(msg string) {
	a.Send([]byte(msg))
}

func WithSubscribe[T IWsClient](c T) func(string) {
	return c.Subscribe
}
