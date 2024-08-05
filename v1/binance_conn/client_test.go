package binance_conn_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/lianyun0502/exchange_conn/v1/common"
	"github.com/lianyun0502/exchange_conn/v1/binance_conn"
	
)

var (
	apiKey = "xTz5GK0rSyPANKeQTR5o1zohOdF7DmPRGR2ReAFKNLj0GjoIacB2Ld5Sjzd2p8Wk"
	secretKey = "Hvsqtth66iAyXw7lnbzQGdw0ZCLPru5MWZPllLbcAuHpGMPNiuWoxXAE6LjpKqNg"
	// testURL = "https://testnet.binance.vision"
	testURL = "https://api1.binance.com"
)


func TestPingAPIServer(t *testing.T) {
	client := binance_conn.NewClient(
		apiKey, 
		secretKey, 
		testURL,
	)

	req := binance_conn.NewBinanceRequest(
		http.MethodGet,
		"/api/v3/ping",
		binance_conn.None,
	)

	b_req, err := client.SetRequest(req)
	if err != nil {
		t.Error(err)
		return
	}

	data, err := client.Call(b_req)
	if err != nil {
		t.Error(err)
		return
	}
	var j interface{}

	if string(data) != "{}" {
		t.Errorf(common.PrettyPrint(j))
	}
}

type CheckServerTimeResponce struct {
	ServerTime int64 `json:"serverTime"`
}

func TestCheckServerTime(t *testing.T) {
	client := binance_conn.NewClient(
		apiKey, 
		secretKey, 
		testURL,
	)

	req := binance_conn.NewBinanceRequest(
		http.MethodGet,
		"/api/v3/time",
		binance_conn.None,
	)

	b_req, err := client.SetRequest(req)
	if err != nil {
		t.Error(err)
		return
	}

	data, err := client.Call(b_req)
	if err != nil {
		t.Error(err)
		return
	}

	j := new(CheckServerTimeResponce)
	err = json.Unmarshal(data, j)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf(common.PrettyPrint(j))
}

func TestGetExchangeInfo(t *testing.T) {
	client := binance_conn.NewClient(
		apiKey, 
		secretKey, 
		testURL,
	)

	req := binance_conn.NewBinanceRequest(
		http.MethodGet,
		"/api/v3/exchangeInfo",
		binance_conn.None,
	)

	req.SetQuery("symbols", `["BTCUSDT", "ETHUSDT"]`)

	b_req, err := client.SetRequest(req)
	if err != nil {
		t.Error(err)
		return
	}

	data, err := client.Call(b_req)
	if err != nil {
		t.Error(err)
		return
	}

	j := new(interface{})
	err = json.Unmarshal(data, j)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf(common.PrettyPrint(j))
}