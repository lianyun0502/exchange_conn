package bybit_test

import (
	"testing"

	"github.com/lianyun0502/exchange_conn/v1/common"
	"github.com/lianyun0502/exchange_conn/v2/bybit/http_client"
)

var api = bybit.NewSpotClient(apiKey, secretKey, bybit.IsTestNet())

func TestMarket_PremiumIndexPrice(t *testing.T) {
	resp, err := api.Market_PremiumIndexPrice("WIFUSDT", "15")
	if err != nil {
		t.Error(resp)
		return
	}
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

