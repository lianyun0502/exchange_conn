package exchange_conn

import (
	// "log"

	// "github.com/lxzan/gws"
)

type TestObject struct {
	A int
	B string
}

type IWsClient interface {
	AddEvent(e *Event)
	Reconnect()
	StartLoop()
	Connect(string) error
	Stop() error
	Send([]byte)
}

type WebSocketAgent[T IWsClient] struct {
	Client T
	
}

func NewWebSocketAgent[T IWsClient](client T) *WebSocketAgent[T] {
	agent := &WebSocketAgent[T]{
		Client: client,
	}
	return agent
}