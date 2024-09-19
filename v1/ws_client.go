package exchange_conn

import (
	"net/http"
	"time"

	"github.com/lxzan/gws"
	log "github.com/sirupsen/logrus"
)

type ErrHandler func(err error)
type WsHandler func(message []byte)

type WebSocketEvent struct {
	Err_Handler   func(err error)
	Ws_Handler    func(message []byte)
	Close_Handler func()

	pingTimeout *time.Timer
	msgTimout   *time.Timer

	isClosed   bool
	doneSignal chan struct{}

	Logger *log.Logger
}

func (conn *WebSocketEvent) OnOpen(socket *gws.Conn) {
	conn.Logger.Info("OnOpen")
	conn.isClosed = false
	conn.doneSignal = make(chan struct{})
	conn.pingTimeout = time.NewTimer(3 * time.Second)
	conn.msgTimout = time.NewTimer(5 * time.Minute)
	go func() {
		for {
			select {
			case <-conn.pingTimeout.C:
				conn.Logger.Warning("Ping server timeout")
				conn.pingTimeout.Stop()
				socket.NetConn().Close()
			case <-conn.msgTimout.C:
				conn.Logger.Warning("OnMessage timeout")
				conn.msgTimout.Stop()
				socket.NetConn().Close()
			case <-conn.doneSignal:
				return
			}
		}
	}()
	socket.WritePing([]byte("ping"))
}
func (conn *WebSocketEvent) OnPing(socket *gws.Conn, message []byte) {
	conn.Logger.Info("OnPing")
	socket.WritePong(message)
}
func (conn *WebSocketEvent) OnPong(socket *gws.Conn, message []byte) {
	conn.Logger.Info("OnPong")
	conn.pingTimeout.Reset(6 * time.Second)
	go func() {
		time.Sleep(5 * time.Second)
		socket.WritePing([]byte("ping"))
	}()
}
func (conn *WebSocketEvent) OnMessage(socket *gws.Conn, message *gws.Message) {
	defer message.Close()
	conn.Logger.Debug("OnMessage")
	conn.msgTimout.Reset(5 * time.Minute)
	if conn.Ws_Handler != nil {
		conn.Ws_Handler(message.Data.Bytes())
	}
}
func (conn *WebSocketEvent) OnClose(socket *gws.Conn, err error) {
	conn.Logger.Info("OnClose")
	conn.isClosed = true
	conn.doneSignal <- struct{}{}
	if conn.Err_Handler != nil {
		if err == nil {
			return
		}
		conn.Logger.Error(err)
		conn.Err_Handler(err)
	}
}

type WsClient struct {
	*WebSocketEvent
	ClientOption *gws.ClientOption
	Conn         *gws.Conn

	ApiKey      string
	SecretKey   string
	reconnTimes int
	eventLoop   *EventEngine

	DoneSignal  chan struct{}
	StartSignal chan struct{}
}

// override the OnClose method
func (wsc *WsClient) OnClose(socket *gws.Conn, err error) {
	wsc.WebSocketEvent.OnClose(socket, err)
	wsc.AddEvent(&Event{
		Name:    "reconnect",
		Handler: wsc.Reconnect,
		IsBlock: true,
	})
}

func (wsc *WsClient) AddEvent(e *Event) {
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
			wsc.Logger.Error(err)
			return
		}
	}
	wsc.DoneSignal <- struct{}{}
	return
}

func (wsc *WsClient) Send(msg []byte) (err error) {
	err = wsc.Conn.WriteMessage(gws.OpcodeText, msg)
	if err != nil {
		wsc.Logger.Error(err)
	}
	return
}

func (wsc *WsClient) Reconnect() {
	if wsc.reconnTimes < 0 {
		for i := 0; ; i++ {
			wsc.Logger.WithField("times", i+1).Info("Reconnect...")
			conn, _, err := gws.NewClient(wsc, wsc.ClientOption)
			wsc.Conn = conn
			if err == nil {
				wsc.AddEvent(&Event{
					Name:    "restart",
					IsBlock: false,
					Handler: wsc.StartLoop,
				})

			}
		}
	} else {
		for i := 0; i < wsc.reconnTimes; i++ {
			wsc.Logger.WithField("times", i+1).Info("Reconnect...")
			conn, _, err := gws.NewClient(wsc, wsc.ClientOption)
			wsc.Conn = conn
			if err == nil {
				wsc.AddEvent(&Event{
					Name:    "restart",
					IsBlock: false,
					Handler: wsc.StartLoop,
				})
				return
			}
		}
	}
	wsc.Logger.Warning("Reconnection failed")
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
	if err != nil {
		wsc.Logger.WithFields(log.Fields{"respone": resp}).Error(err)
	}
	return resp, err
}

func NewWsClient(messageHandle WsHandler, errHandle ErrHandler, reconnectTimes int) (client *WsClient) {
	engine := NewEventEngine()
	engine.Luanch()
	return &WsClient{
			reconnTimes: reconnectTimes,
			eventLoop:   engine,
			WebSocketEvent: &WebSocketEvent{
				Err_Handler: errHandle,
				Ws_Handler:  messageHandle,
				Logger:      engine.Logger,
			},
		DoneSignal:  make(chan struct{}),
		StartSignal: make(chan struct{}, 5),
	}
}
