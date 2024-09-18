package exchange_conn

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/lianyun0502/exchange_conn/v2/common"
	"github.com/lxzan/gws"
	"github.com/sirupsen/logrus"
)

type ExchangeApi struct {
	Name      string
	HostType  string
	APIKey    string // API key
	SecretKey string // Secret key
	BaseURL   string // Base URL for API requests
}

type WsClient struct {
	// *WebSocketEvent
	ClientOption *gws.ClientOption
	Conn         *gws.Conn
	ExchangeInfo *ExchangeApi

	reconnTimes int
	// eventLoop   *EventEngine

	Logger      *logrus.Logger
	pingTimeout *time.Timer
	msgTimout   *time.Timer

	Ws_Handler func(message []byte)

	ReqMap	  map[string]chan []byte

	StopSignal  chan struct{}
	StartSignal chan struct{}
}

func (wsc *WsClient) OnOpen(socket *gws.Conn) {
	wsc.Logger.Info("OnOpen")
	wsc.StopSignal = make(chan struct{})
	wsc.pingTimeout = time.NewTimer(3 * time.Second)
	wsc.msgTimout = time.NewTimer(5 * time.Minute)
	go func() {
		for {
			select {
			case <-wsc.pingTimeout.C:
				wsc.Logger.Warning("Ping server timeout")
				wsc.pingTimeout.Stop()
				socket.NetConn().Close()
			case <-wsc.msgTimout.C:
				wsc.Logger.Warning("OnMessage timeout")
				wsc.msgTimout.Stop()
				socket.NetConn().Close()
			case <-wsc.StopSignal:
				return
			}
		}
	}()
	socket.WritePing([]byte("ping"))
}
func (wsc *WsClient) OnPing(socket *gws.Conn, message []byte) {
	wsc.Logger.Info("OnPing")
	socket.WritePong(message)
}
func (wsc *WsClient) OnPong(socket *gws.Conn, message []byte) {
	wsc.Logger.Info("OnPong")
	wsc.pingTimeout.Reset(6 * time.Second)
	go func() {
		time.Sleep(5 * time.Second)
		socket.WritePing([]byte("ping"))
	}()
}
func (wsc *WsClient) OnMessage(socket *gws.Conn, message *gws.Message) {
	defer message.Close()
	wsc.Logger.Debug("OnMessage")
	wsc.msgTimout.Reset(5 * time.Minute)
	rawData := make([]byte, message.Data.Len())
	copy(rawData, message.Data.Bytes())
	if wsc.Ws_Handler != nil {
		wsc.Ws_Handler(rawData)
	}
}
func (wsc *WsClient) OnClose(socket *gws.Conn, err error) {
	wsc.Logger.Info("OnClose")
	if err != nil {
		wsc.Logger.Error(err)
	}
	if _, ok := <- wsc.StopSignal; ok {
		close(wsc.StopSignal)
		wsc.Reconnect()
	}
}

func (wsc *WsClient) StartLoop() {
	wsc.StartSignal <- struct{}{}
	wsc.Conn.ReadLoop()
}

func (wsc *WsClient) Stop() (err error) {
	wsc.Logger.Info("Stop client")
	close(wsc.StopSignal)
	err = wsc.Conn.NetConn().Close()
	if err != nil {
		wsc.Logger.Error(err)
		return err
	}
	return nil
}

func (wsc *WsClient) Subscribe(topics []string) (respData []byte, err error) {
	id := common.GetUUID()
	jTopics, _ := json.Marshal(topics)
	msg := fmt.Sprintf(`{"req_id":"%s","op":"subscribe","args":%s}`, id, string(jTopics))
	respCh := make(chan []byte, 2)
	wsc.ReqMap[id] = respCh
	err = wsc.Send([]byte(msg))
	var resp []byte
	select {
	case <-time.After(5 * time.Second):
		wsc.Logger.Warning("Request timeout")
		err = errors.New("Request timeout")
		close(respCh)
	case resp = <-respCh:
		wsc.Logger.Debugf(`Response: %s`, string(resp))
	}
	delete(wsc.ReqMap, id)
	return resp, err
}

func (wsc *WsClient) GetSignature() (respData []byte, err error) {
	expires := time.Now().Unix()*1000 + 10000
	param := []string{
		wsc.ExchangeInfo.APIKey,
		strconv.FormatInt(expires, 10),
		common.GetSignature(wsc.ExchangeInfo.SecretKey, fmt.Sprintf("GET/realtime%d", expires)),
	}
	resp, err := wsc.Request("auth", nil, param)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (wsc *WsClient) Request(op string, header any, args any) (respData []byte, err error) {
	id := common.GetUUID()
	req := &Request{
		ReqID: id,
		Header: header,
		Op: op,
		Args: args,
	}
	reqByte, err := json.Marshal(req)
	if err != nil {
		return
	}
	respCh := make(chan []byte, 2)
	wsc.ReqMap[id] = respCh
	wsc.Logger.Debug("Send request")
	wsc.Send(reqByte)
	var resp []byte
	select {
	case <-time.After(5 * time.Second):
		wsc.Logger.Warning("Request timeout")
		err = errors.New("Request timeout")
		close(respCh)
	case resp = <-respCh:
		wsc.Logger.Debugf(`Response: %s`, string(resp))
	}
	delete(wsc.ReqMap, id)
	return resp, err
}

func (wsc *WsClient) Send(msg []byte) (err error) {
	wsc.Logger.Info("Send message")
	wsc.Logger.Debugf(`Send message: %s`, string(msg))
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
			if _, err := wsc.Connect(); err == nil {
				go wsc.StartLoop()
				wsc.Logger.Info("Reconnection success")
				return
			}
		}
	} else {
		for i := 0; i < wsc.reconnTimes; i++ {
			wsc.Logger.WithField("times", i+1).Info("Reconnect...")

			if _, err := wsc.Connect(); err == nil {
				go wsc.StartLoop()
				wsc.Logger.Info("Reconnection success")
				return
			}
		}
	}
	wsc.Logger.Warning("Reconnection failed")
	go wsc.Stop()
}

func (wsc *WsClient) Connect() (resp *http.Response, err error) {
	wsc.ClientOption = &gws.ClientOption{
		ReadBufferSize:   655350,
		Addr:             wsc.ExchangeInfo.BaseURL,
		HandshakeTimeout: 45 * time.Second,
		PermessageDeflate: gws.PermessageDeflate{
			Enabled:               true,
			ServerContextTakeover: true,
			ClientContextTakeover: true,
		},
		Logger: wsc.Logger,
	}

	wsc.Conn, resp, err = gws.NewClient(wsc, wsc.ClientOption)
	if err != nil {
		wsc.Logger.WithFields(logrus.Fields{"respone": resp}).Error(err)
	}
	return resp, err
}

func NewWsClient(exchangeInfo *ExchangeApi, wsHandler func(message []byte), clientOpts ...func(*WsClient)) *WsClient {
	client := &WsClient{
		ExchangeInfo: exchangeInfo,
		reconnTimes:  -1,
		Logger:       logrus.New(),
		Ws_Handler:   wsHandler,
		StartSignal:  make(chan struct{}, 5),
		ReqMap: make(map[string]chan []byte),
	}
	for _, opt := range clientOpts {
		opt(client)
	}
	return client
}

/*
Reconnection times when connection is lost

	-1: infinite reconnection
	 0: no reconnection
	 n: reconnection n times
*/
func WithReconnectionTimes(times int) func(*WsClient) {
	return func(client *WsClient) {
		client.reconnTimes = times
	}
}


type Request struct {
	ReqID string `json:"reqId"`
	Header any `json:"header,omitempty"`
	Op string `json:"op"`
	Args any `json:"args,omitempty"`
}