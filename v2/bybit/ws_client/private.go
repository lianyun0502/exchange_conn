package bybit

import (
	"github.com/lianyun0502/exchange_conn/v2/consts"
	
)

func NewWsPrivateClient(apiKey, secretKey string, opts ...func(*WsBybitClient)) (*WsBybitClient, error) {
	return NewWsAPIClient(consts.Private, apiKey, secretKey, opts...)
}