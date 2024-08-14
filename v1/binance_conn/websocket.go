package binance_conn

import (
	// "fmt"
	// "log"
	"net/http"
	"time"

	"github.com/lianyun0502/exchange_conn/v1"
	"github.com/lianyun0502/exchange_conn/v1/common"
	"github.com/lxzan/gws"
	log "github.com/sirupsen/logrus"
)

type ErrHandler func(err error)
type WsHandler func(message []byte)

type WebSocketEvent struct {
	Err_Handler   func(err error)
	Ws_Handler    func(message []byte)
	Close_Handler func()

	pingTimer common.Timer

	isClosed bool

	Logger *log.Logger
}

func (conn *WebSocketEvent) OnOpen(socket *gws.Conn) {
	conn.Logger.Infoln("OnOpen")
	conn.isClosed = false
	conn.pingTimer = common.Timer{
		Interval: 10 * time.Second,
		Handler: func() {
			conn.Logger.Warningln("Ping server timeout")
			socket.NetConn().Close()
		},
	}
	conn.pingTimer.Start(nil)
	socket.WritePing([]byte("ping"))
}
func (conn *WebSocketEvent) OnPing(socket *gws.Conn, message []byte) {
	conn.Logger.Infoln("OnPing")
	socket.WritePong(message)
}
func (conn *WebSocketEvent) OnPong(socket *gws.Conn, message []byte) {
	conn.Logger.Infoln("OnPong")
	go func() {
		time.Sleep(5 * time.Second)
		socket.WritePing([]byte("ping"))
		conn.pingTimer.Reset()
	}()
}
func (conn *WebSocketEvent) OnMessage(socket *gws.Conn, message *gws.Message) {
	defer message.Close()
	conn.Logger.Debugln("OnMessage")
	if conn.Ws_Handler == nil {
		return
	}
	conn.Ws_Handler(message.Data.Bytes())
}
func (conn *WebSocketEvent) OnClose(socket *gws.Conn, err error) {
	conn.Logger.Infoln("OnClose")
	conn.isClosed = true
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
	Conn         *gws.Conn

	ApiKey      string
	SecretKey   string
	reconnTimes int
	eventLoop   *exchange_conn.EventEngine

	DoneSignal  chan struct{}
	StartSignal chan struct{}
}

// override the OnClose method
func (wsc *WsClient) OnClose(socket *gws.Conn, err error) {
	wsc.WebSocketEvent.OnClose(socket, err)
	wsc.AddEvent(&exchange_conn.Event{
		Name:    "reconnect",
		Handler: wsc.Reconnect,
		IsBlock: true,
	})
}

func (wsc *WsClient) AddEvent(e *exchange_conn.Event) {
	wsc.eventLoop.AddEvent(e)
}

func (wsc *WsClient) StartLoop() {
	wsc.StartSignal <- struct{}{}
	wsc.Conn.ReadLoop()
}

func (wsc *WsClient) Stop() (err error) {
	wsc.eventLoop.Stop()
	if !wsc.isClosed {
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
	wsc.Logger.Infoln("reconnect")
	if wsc.reconnTimes < 0 {
		for {
			conn, _, err := gws.NewClient(wsc, wsc.ClientOption)
			wsc.Conn = conn
			if err == nil {
				wsc.AddEvent(&exchange_conn.Event{
					Name:    "restart",
					IsBlock: false,
					Handler: wsc.StartLoop,
				})

			}
		}
	} else {
		for i := 0; i < wsc.reconnTimes; i++ {
			wsc.Logger.Infof("reconnect times {%d}", i+1)
			conn, _, err := gws.NewClient(wsc, wsc.ClientOption)
			wsc.Conn = conn
			if err == nil {
				wsc.AddEvent(&exchange_conn.Event{
					Name:    "restart",
					IsBlock: false,
					Handler: wsc.StartLoop,
				})
				return
			}
		}
	}
	go wsc.Stop()
}

func (wsc *WsClient) Connect(url string) (resp *http.Response, err error) {
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

	wsc.Conn, resp, err = gws.NewClient(wsc, wsc.ClientOption)
	return resp, err
}

func NewWsClient(messageHandle WsHandler, errHandle ErrHandler, reconnectTimes int) (client *WsClient) {
	engine := exchange_conn.NewEventEngine()
	logger := engine.Logger
	engine.Luanch()
	return &WsClient{
		reconnTimes: reconnectTimes,
		eventLoop:   engine,
		WebSocketEvent: WebSocketEvent{
			Err_Handler: errHandle,
			Ws_Handler:  messageHandle,
			Logger:      logger,
		},
		DoneSignal:  make(chan struct{}),
		StartSignal: make(chan struct{}),
	}
}
