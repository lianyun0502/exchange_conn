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
	Subscribe func(string)
	Send func([]byte) error
	Connect func(string) (*http.Response, error)
	StartLoop func()
	Stop func() error
}

func NewWebSocketAgent[T IWsClient](client T) *webSocketAgent[T] {
	agent := &webSocketAgent[T]{
		Client: client,
		Subscribe: WithSubscribe(client),
		Send: client.Send,
		Connect: client.Connect,
		StartLoop: client.StartLoop,
		Stop: client.Stop,
	}
	return agent
}

func (a *webSocketAgent[T]) SendString(msg string) {
	a.Send([]byte(msg))
}

func WithSubscribe[T IWsClient](c T) func(string) {
	return c.Subscribe
}
