package bybit_test

import (
	"testing"

	"github.com/lianyun0502/exchange_conn/v2/bybit/http_client"
)


func TestMarket_PremiumIndexPrice(t *testing.T) {
	api := bybit.NewSpotClient(apiKey, secretKey, bybit.IsTestNet())
	resp, err := api.Market_PremiumIndexPrice("WIFUSDT", "15")
	if err != nil {
		t.Error(resp)
		return
	}
}