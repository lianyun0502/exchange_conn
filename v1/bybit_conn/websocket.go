package bybit_conn

import (
	// "fmt"
	"log"
	"time"

	"github.com/lianyun0502/exchange_conn/v1"
	"github.com/lianyun0502/exchange_conn/v1/common"
	"github.com/lxzan/gws"
)

type ErrHandler func(err error)
type WsHandler func(message []byte)

type WebSocketEvent struct {
	Err_Handler   func(err error)
	Ws_Handler    func(message []byte)
	Close_Handler func()

	pingTimer common.Timer

	isClosed bool
}

func (conn *WebSocketEvent) OnOpen(socket *gws.Conn) {
	log.Println("OnOpen")
	conn.isClosed =false
	conn.pingTimer = common.Timer{
		Interval: 10 * time.Second,
		Handler: func() {
			log.Println("Ping server timeout")
			socket.NetConn().Close()
		},
	}
	conn.pingTimer.Start(nil)
	socket.WritePing([]byte("ping"))
}
func (conn *WebSocketEvent) OnPing(socket *gws.Conn, message []byte) {
	log.Println("OnPing")
	socket.WritePong(message)
}
func (conn *WebSocketEvent) OnPong(socket *gws.Conn, message []byte) {
	log.Println("OnPong")
	go func() {
		time.Sleep(5 * time.Second)
		socket.WritePing([]byte("ping"))
		conn.pingTimer.Reset()	
	}()
}
func (conn *WebSocketEvent) OnMessage(socket *gws.Conn, message *gws.Message) {
	defer message.Close()
	log.Println("OnMessage")
	if conn.Ws_Handler == nil {
		return
	}
	conn.Ws_Handler(message.Data.Bytes())
}
func (conn *WebSocketEvent) OnClose(socket *gws.Conn, err error) {
	log.Println("OnClose")
	conn.isClosed =true
	conn.pingTimer.Stop()
	if conn.Err_Handler == nil {
		return
	}
	if err != nil {
		conn.Err_Handler(err)
	}
}

type WsClient struct {
	WebSocketEvent
	ClientOption *gws.ClientOption
	Conn *gws.Conn

	ApiKey   string
	SecretKey   string
	reconnTimes int
	maxAliveTime string
	eventLoop *exchange_conn.EventEngine

	DoneSignal chan struct{}
}
// override the OnClose method
func (wsc *WsClient) OnClose(socket *gws.Conn, err error){
	wsc.WebSocketEvent.OnClose(socket, err)
	wsc.AddEvent(&exchange_conn.Event{
		Name: "reconnect",
		Handler: wsc.Reconnect,
		IsBlock: true,
	})
}

func (wsc *WsClient) AddEvent(e *exchange_conn.Event) {
	wsc.eventLoop.AddEvent(e)
}

func (wsc *WsClient) StartLoop() {
	wsc.Conn.ReadLoop()
}

func (wsc *WsClient) Stop() (err error) {
	wsc.eventLoop.Stop()
	if !wsc.isClosed{
		err = wsc.Conn.NetConn().Close()
		if err != nil {
			return
		}
	}
	wsc.DoneSignal <- struct{}{}
	return
}

func (wsc *WsClient) Send(msg []byte) {
	wsc.Conn.WriteMessage(gws.OpcodeText, msg)
}

func (wsc *WsClient) Reconnect() {
	log.Printf("reconnect")
	if wsc.reconnTimes < 0{
		for {
			conn, _, err := gws.NewClient(wsc, wsc.ClientOption)
			wsc.Conn = conn
			if err == nil{
				wsc.AddEvent(&exchange_conn.Event{
					Name: "restart",
					IsBlock: false,
					Handler: wsc.StartLoop,
				})

			}	
		}
	} else {
		for i := 0; i < wsc.reconnTimes; i++ {
			log.Printf("reconnect times {%d}", i+1)
			conn, _, err := gws.NewClient(wsc, wsc.ClientOption)
			wsc.Conn = conn
			if err == nil{
				wsc.AddEvent(&exchange_conn.Event{
					Name: "restart",
					IsBlock: false,
					Handler: wsc.StartLoop,
				})
			return
			}
		}
	} 
	// wsc.AddEvent(&exchange_conn.Event{
	// 	Name: "reconnect fail",
	// 	IsBlock: true,
	// 	Handler: func () {wsc.Stop()},
	// })
	go wsc.Stop()
}

func (wsc *WsClient) Connect(url string) (err error) {

	if wsc.maxAliveTime != "" {
		url += "?max_alive_time=" + wsc.maxAliveTime
	}

	wsc.ClientOption = &gws.ClientOption{
		ReadBufferSize:   655350,
		Addr:             url,
		HandshakeTimeout: 45 * time.Second,
		PermessageDeflate: gws.PermessageDeflate{
			Enabled:               true,
			ServerContextTakeover: true,
			ClientContextTakeover: true,
		},
	}

	conn, _, err := gws.NewClient(
		wsc,
		wsc.ClientOption,
	)
	if err != nil {
		return err
	}
	wsc.Conn = conn
	return
}

func NewWsClient(messageHandle WsHandler, errHandle ErrHandler, reconnectTimes int) (client *WsClient) {
	engine := exchange_conn.NewEventEngine()
	engine.Luanch()
	return &WsClient{
		reconnTimes: reconnectTimes,
		eventLoop: engine,
		WebSocketEvent: WebSocketEvent{
			Err_Handler: errHandle,
			Ws_Handler: messageHandle,
		},
		DoneSignal: make(chan struct{}),
	}
}

type Object exchange_conn.TestObject

func (t Object) Test() {
	log.Println(" binance test ")
}