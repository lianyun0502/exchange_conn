package binance_conn

import (
	// "fmt"
	"encoding/json"
	"log"
	"time"
	"errors"

	"github.com/lxzan/gws"
	"github.com/lianyun0502/exchange_conn/v1/common"
)

type WsAPIErrorResponse struct {
	ID string `json:"id"`
	Code int `json:"code"`
	Msg string `json:"msg"`
}

type WebSocketAPI struct {
	APIKey     string // API key
	SecretKey  string // Secret key

	Connector *gws.Conn
	DoneCh, StopCh chan struct{}
	WriteCh chan []byte
	resp_handlers map[string]chan []byte
}

func (ws *WebSocketAPI) StartLoop(){
	go ws.Connector.ReadLoop()

	go func() { // write loop
		for {
			err := ws.Connector.WriteMessage(gws.OpcodeText, <-ws.WriteCh)
			if err != nil {
				log.Println(err)
				break
			}
		}
	}()
}

func (ws *WebSocketAPI) StopLoop(){
	ws.Connector.NetConn().Close()
}
func (ws *WebSocketAPI) SendMessage(id string, resp_chan chan []byte, data []byte){
	ws.resp_handlers[id] = resp_chan
	ws.WriteCh <- data
}

type WsApiPingResponse struct {
	ID string `json:"id"`
	Method string `json:"method"`
}

func (ws *WebSocketAPI) PingServer() (resp interface{}, err error){
	res := &WsApiPingResponse{ID: common.GetUUID(), Method: "ping"}
	respCh := make(chan []byte)

	data, _ := json.Marshal(res)

	ws.SendMessage(res.ID, respCh, data)

	data = <- respCh

	err = json.Unmarshal(data, &resp)
	if err != nil {
		log.Println(err)
		return
	}
	return
}


type WebSocketAPIEvent struct {
	Err_Handler   func(err error)
	Resp_Handlers map[string]chan []byte
}

func (event *WebSocketAPIEvent) OnOpen(socket *gws.Conn) {
	log.Println("OnOpen")
}
func (event *WebSocketAPIEvent) OnPing(socket *gws.Conn, message []byte) {
	log.Println("OnPing")
	socket.WritePong(message)
}
func (event *WebSocketAPIEvent) OnPong(socket *gws.Conn, message []byte) {
	log.Println("OnPong")
}
func (event *WebSocketAPIEvent) OnMessage(socket *gws.Conn, message *gws.Message) {
	defer message.Close()
	log.Println("OnMessage")
	resp := new(WsAPIErrorResponse)
	data := message.Data.Bytes()
	json.Unmarshal(data, &resp)
	if resp.Code != 0 {
		event.Err_Handler(errors.New(resp.Msg))
		if channel, ok := event.Resp_Handlers[resp.ID]; ok {
			close(channel)
			delete(event.Resp_Handlers, resp.ID)
		}
	}

	if channel, ok := event.Resp_Handlers[resp.ID]; ok {
		channel <- data
		delete(event.Resp_Handlers, resp.ID)
	}
}
func (event *WebSocketAPIEvent) OnClose(socket *gws.Conn, err error) {
	log.Println("OnClose")
	if err != nil {
		event.Err_Handler(err)
	}
}

func NewWebSocketAPI(apiKey, secretKey, url string, errHandler ErrHandler) (ws *WebSocketAPI, err error) {
	resp_handlers := make(map[string]chan []byte)
	conn, _, err := gws.NewClient(
		&WebSocketAPIEvent{
			Err_Handler:   errHandler,
			Resp_Handlers: resp_handlers,
		},
		&gws.ClientOption{
			ReadBufferSize:   655350,
			Addr:             url,
			HandshakeTimeout: 45 * time.Second,
			PermessageDeflate: gws.PermessageDeflate{
				Enabled:               true,
				ServerContextTakeover: true,
				ClientContextTakeover: true,
			},
		},
	)

	if err != nil {
		return nil, err
	}
	ws = &WebSocketAPI{
		APIKey: apiKey,
		SecretKey: secretKey,
		Connector: conn,
		DoneCh: make(chan struct{}),
		StopCh: make(chan struct{}),
		WriteCh: make(chan []byte),
		resp_handlers: resp_handlers,
	}
	return ws, nil
}
