package exchange_conn

import (
	"log"

	"github.com/lxzan/gws"
)

type WsClient struct {
	ClientOption *gws.ClientOption
	WsEvent gws.Event
	Conn *gws.Conn
}

type WebSocketAgent struct {
	reconnTimes int

	engine *EventEngine
	client *WsClient

	DoneSignal chan struct{}
}

func NewWebSocketAgent(client *WsClient, reconnTimes int) *WebSocketAgent {
	agent := &WebSocketAgent{
		reconnTimes: reconnTimes,
		engine: new(EventEngine),
		client: client,
	}
	return agent
}

func (a *WebSocketAgent) AddEvent(e *Event){
	a.engine.AddEvent(e)
}

func (a *WebSocketAgent) Start () {
	a.engine.Luanch()
	go a.client.Conn.ReadLoop()
}

func (a *WebSocketAgent) Reconnect () {
	log.Printf("reconnect")
	if a.reconnTimes < 0{
		for {
			conn, _, err := gws.NewClient(a.client.WsEvent, a.client.ClientOption)
			a.client.Conn = conn
			if err == nil{
				a.engine.AddEvent(&Event{
					Name: "restart",
					IsBlock: false,
					Handler: a.client.Conn.ReadLoop,
				})

			}	
		}
	} else {
		for i := 0; i < a.reconnTimes; i++ {
			log.Printf("reconnect times {%d}", i+1)
			conn, _, err := gws.NewClient(a.client.WsEvent, a.client.ClientOption)
			a.client.Conn = conn
			if err == nil{
				a.engine.AddEvent(&Event{
					Name: "restart",
					IsBlock: false,
					Handler: a.client.Conn.ReadLoop,
				})
			return
			}
		}
	}
}


func (a *WebSocketAgent) Stop () {
	a.engine.Stop()
	a.client.Conn.NetConn().Close()
	a.DoneSignal <- struct{}{}
}