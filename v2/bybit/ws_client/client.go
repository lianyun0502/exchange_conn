package bybit

import (
	"net/http"
	"strconv"
	"time"

	"github.com/lianyun0502/exchange_conn/v2/http_client/consts"
	"github.com/lianyun0502/exchange_conn/v2/ws_client"
	"github.com/lxzan/gws"
	"github.com/sirupsen/logrus"
	"github.com/lianyun0502/exchange_conn/v2/common"
)

type WsBybitClient struct {
	*exchange_conn.WsClient
	maxAliveTime string
}

func (wsc *WsBybitClient) Connect() (resp *http.Response, err error) {
	url := wsc.ExchangeInfo.BaseURL
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
	if err != nil {
		wsc.Logger.WithFields(logrus.Fields{"respone": resp}).Error(err)
	}
	return resp, err
}

func (wsc *WsBybitClient) Subscribe(topics string) {
	uuid := common.GetUUID()
	wsc.Send([]byte(`{"req_id": `+ uuid +`,"op":"subscribe","args":` + topics + `}`))
}

func NewWsSpotClient(opts ...func(*WsBybitClient)) (client *WsBybitClient) {
	exchInfo := &exchange_conn.ExchangeApi{
		Name: consts.Bybit,
		HostType: consts.Spot,
		BaseURL: SPOT_MAINNET,
	}
	client = &WsBybitClient{
		WsClient: exchange_conn.NewWsClient(exchInfo, nil), 
		maxAliveTime: "",
	}
	for _, opt := range opts {
		opt(client)
	}
	return client
}

func NewWsFutureClient(opts ...func(*WsBybitClient)) (client *WsBybitClient) {
	exchInfo := &exchange_conn.ExchangeApi{
		Name: consts.Bybit,
		HostType: consts.Future,
		BaseURL: LINEAR_MAINNET,
	}
	client = &WsBybitClient{
		WsClient: exchange_conn.NewWsClient(exchInfo, nil), 
		maxAliveTime: "",
	}
	for _, opt := range opts {
		opt(client)
	}
	return client
}

// 針對私有頻道和交易, 您可以自定義連接存活時長, 通過增加參數max_active_time, 最小支持30s (30秒), 最大支持600s (10分鐘)
func WithMaxAliveTime(maxAliveTime int) func(*WsBybitClient) {
	return func(wsc *WsBybitClient) {
		wsc.maxAliveTime = strconv.Itoa(maxAliveTime)
	}
}

func WithTestNet() func(*WsBybitClient) {
	return func(wsc *WsBybitClient) {
		switch wsc.ExchangeInfo.HostType {
		case consts.Future :
			wsc.ExchangeInfo.BaseURL = LINEAR_TESTNET
		case consts.Spot :
			wsc.ExchangeInfo.BaseURL = SPOT_TESTNET
		case consts.Trade :
			wsc.ExchangeInfo.BaseURL = WEBSOCKET_TRADE_TESTNET
		}
	}
}


func TradeHandler(message []byte){
	
}