package bybit

import (
	"net/http"
	"time"

	"github.com/lianyun0502/exchange_conn/v2/consts"
	"github.com/lianyun0502/exchange_conn/v2/ws_client"
	"github.com/lxzan/gws"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fastjson"
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


func IsTestNet() func(*WsBybitClient) {
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


func WithBybitHandler(reqMap map[string]chan []byte, qouteHandler func(message []byte)) func([]byte) {
	return func(rawData []byte){ {
		v, _ := fastjson.ParseBytes(rawData)
		if reqID := string(v.GetStringBytes("reqId")); reqID != "" {
			if respCh, ok := reqMap[reqID]; ok {
				respCh <- rawData
				return
			}
		}
		if reqID := string(v.GetStringBytes("req_id")); reqID != "" {
			if respCh, ok := reqMap[reqID]; ok {
				respCh <- rawData
				return
			}
		}
		if qouteHandler != nil {
			qouteHandler(rawData)
		}
	}
}
	
}




