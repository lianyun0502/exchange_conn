package bybit_conn

import (
	"net/http"
	"strconv"
	"time"

	"github.com/lianyun0502/exchange_conn/v1"
	"github.com/lxzan/gws"
)

type ErrHandler func(err error)
type WsHandler func(message []byte)
	

type WsClient struct {
	*exchange_conn.WsClient
	maxAliveTime string
}

func (wsc *WsClient) Connect(url string) (resp *http.Response, err error) {

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

	wsc.Conn, resp, err = gws.NewClient(wsc, wsc.ClientOption)
	return resp, err
}

func (wsc *WsClient) Subscribe(topics string) {
	wsc.Send([]byte(`{"req_id": "test","op":"subscribe","args":`+ topics +`}`))
}

func NewWsClient(messageHandle func(message []byte), errHandle func(err error), reconnectTimes int, opts ...func(*WsClient)) (client *WsClient) {
	client = &WsClient{WsClient: exchange_conn.NewWsClient(messageHandle, errHandle, reconnectTimes), maxAliveTime: ""}
	for _, opt := range opts {
		opt(client)
	}
	return client
}
//針對私有頻道和交易, 您可以自定義連接存活時長, 通過增加參數max_active_time, 最小支持30s (30秒), 最大支持600s (10分鐘)
func WithMaxAliveTime(maxAliveTime int) func(*WsClient) {
	return func(wsc *WsClient) {
		wsc.maxAliveTime = strconv.Itoa(maxAliveTime)
	}
}


