package exchange_conn

import (
// "log"

// "github.com/lxzan/gws"
)

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

func (a *WebSocketAgent[T]) Send(msg []byte) {
	a.Client.Send(msg)
}

func (a *WebSocketAgent[T]) Connect(url string) error {
	return a.Client.Connect(url)
}

func (a *WebSocketAgent[T]) StartLoop() {
	a.Client.StartLoop()
}

func (a *WebSocketAgent[T]) Stop() error {
	return a.Client.Stop()
}

func (a *WebSocketAgent[T]) Reconnect() {
	a.Client.Reconnect()
}

func (a *WebSocketAgent[T]) SendString(msg string) {
	a.Client.Send([]byte(msg))
}
