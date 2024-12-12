package bybit_test

import (
	"testing"

	"github.com/lianyun0502/exchange_conn/v1/common"
	"github.com/lianyun0502/exchange_conn/v2/bybit/http_client"
)

var api = bybit.NewSpotClient(apiKey, secretKey, bybit.IsTestNet())

func TestMarket_PremiumIndexPrice(t *testing.T) {
	resp, err := api.Market_PremiumIndexPrice("WIFUSDT", "1")
	if err != nil {
		t.Error(resp)
		return
	}
	t.Log(common.PrettyPrint(resp))
}

func TestMarket_InstrumentInfo(t *testing.T) {
	resp, err := api.Market_InstrumentsInfo("spot", bybit.WithQuery(map[string]string{"symbol": "BOMEUSDT"}))
	if err != nil {
		t.Error(resp)
		return
	}
	resp2, err := api.Market_InstrumentsInfo("linear", bybit.WithQuery(map[string]string{"symbol": "BOMEUSDT"}))
	if err != nil {
		t.Error(resp2)
		return
	}
}

func TestTrade_ExecutionList(t *testing.T) {
	resp, err := api.Trade_ExecutionList("spot")
	if err != nil {
		t.Error(resp)
		return
	}
	t.Log(common.PrettyPrint(resp))
}


func TestAccount_WalletBalance(t *testing.T) {
	resp, err := api.Account_WalletBalance("UNIFIED")
	if err != nil {
		t.Error(resp)
		return
	}
	t.Log(common.PrettyPrint(resp))
}

func TestPositionList(t *testing.T) {
	resp, err := api.Position_List("linear", bybit.WithQuery(map[string]string{"settleCoin": "USDT"}))
	if err != nil {
		t.Error(resp)
		return
	}
	t.Log(common.PrettyPrint(resp))
}


func TestMarket_RiskLimit(t *testing.T) {
	resp, err := api.Market_RiskLimit("linear", bybit.WithQuery(map[string]string{"symbol": "BTCUSDT"}))
	if err != nil {
		t.Error(resp)
		return
	}
	t.Log(common.PrettyPrint(resp))
}

func TestMarket_FundingHistory(t *testing.T) {
	resp, err := api.Market_FundingHistory("linear", "BTCUSDT")
	if err != nil {
		t.Error(resp)
		return
	}
	t.Log(common.PrettyPrint(resp))
}