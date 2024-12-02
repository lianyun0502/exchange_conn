package bybit

import (
	"fmt"
	"github.com/lianyun0502/exchange_conn/v2/consts"
	"github.com/lianyun0502/exchange_conn/v2/ws_client"
	
)

func NewWsStreamClient(hostType, apiKey, secretKey string, opts ...func(*WsBybitClient)) (*WsBybitClient, error) {
	var client *WsBybitClient
	switch hostType {
		case consts.Private:
		client = &WsBybitClient{
			WsClient: wsClient.NewWsClient(&wsClient.ExchangeApi{
				Name: consts.Bybit,
				HostType: consts.Private,
				APIKey: apiKey,
				SecretKey: secretKey,
				BaseURL: WEBSOCKET_PRIVATE_MAINNET,
			}, nil),
			maxAliveTime: "",
		}
		default:
			return nil, fmt.Errorf("hostType error")
	}
	for _, opt := range opts {
		opt(client)
	}
	return client, nil
}

func NewWsPrivateClient(apiKey, secretKey string, opts ...func(*WsBybitClient)) (*WsBybitClient, error) {
	client, err := NewWsStreamClient(consts.Private, apiKey, secretKey, opts...)
	if err != nil {
		return nil, err
	}
	client.PingMessage = `{"op":"ping"}`
	client.Ping = client.PingServer
	return client, nil
}


type PrivateQuote[D any] struct {
	ID    string `json:"id,omitempty"`
	Topic string `json:"topic,omitempty"`
	Time  int64  `json:"time,omitempty"`
	Data  []D    `json:"data,omitempty"`
}
type Execution struct {
	Category    string  `json:"category,omitempty"`
	Symbol      string  `json:"symbol,omitempty"`
	IsLeverage  string  `json:"isLeverage,omitempty"`
	Side        string  `json:"side,omitempty"`
	OrderID     string  `json:"orderId,omitempty"`
	OrderLinkID string  `json:"orderLinkId,omitempty"`
	OrderQty    float64 `json:"orderQty,omitempty,string"`
	LeavesQty   float64 `json:"leavesQty,omitempty,string"`
	OrderType   string  `json:"orderType,omitempty"`
	ExecType    string  `json:"execType,omitempty"`
	ExecPnL     string  `json:"execPnL,omitempty"`
	ExecFee     string  `json:"execFee,omitempty"`
	ExecID      string  `json:"execId,omitempty"`
	ExecPrice   float64 `json:"execPrice,omitempty,string"`
	ExecQty     float64 `json:"execQty,omitempty,string"`
	ExecTime    string  `json:"execTime,omitempty"`
	FeeRate     string  `json:"feeRate,omitempty"`
	IsMaker     bool    `json:"isMaker,omitempty"`
	SeqNum      int     `json:"seq,omitempty"`
}

type Position struct {
	Symbol        string `json:"symbol,omitempty"`
	Side          string `json:"side,omitempty"`
	Size          string `json:"size,omitempty"`
	TradeMode     int    `json:"tradeMode,omitempty"`
	PositionValue string `json:"positionValue,omitempty"`
	Leverage      string `json:"leverage,omitempty"`
	PositionIM    string `json:"positionIM,omitempty"`
	PositionMM    string `json:"positionMM,omitempty"`
	CreateTime    string `json:"createdTime,omitempty"`
	UpdateTime    string `json:"updatedTime,omitempty"`
	SeqNum        int    `json:"seq,omitempty"`
}

type Wallet struct {
	AccountType            string  `json:"accountType"`
	AccountIMRate          string  `json:"accountIMRate"`
	AccountMMRate          string  `json:"accountMMRate"`
	TotalEquity            string  `json:"totalEquity"`
	TotalBalance           string  `json:"totalWalletBalance"`
	TotalMarginBalance     string  `json:"totalMarginBalance"`
	TotalAvailableBalance  string  `json:"totalAvailableBalance"`
	TotalPerpUPL           string  `json:"totalPerpUPL"`
	TotalInitialMargin     string  `json:"totalInitialMargin"`
	TotalMaintenanceMargin string  `json:"totalMaintenanceMargin"`
	Coins                  []coin `json:"coin"`
}

type coin struct {
	Coin            string `json:"coin"`
	Equity          string `json:"equity"`
	USDValue        string `json:"usdValue"`
	Walletbalance   string `json:"walletbalance"`
	BorrowAmount    string `json:"borrowAmount"`
	TotalPositionIM string `json:"totalPositionIM"`
	TotalPositionMM string `json:"totalPositionMM"`
	AccruedInterest string `json:"accruedInterest"`
	Free            string `json:"free"`
}

type Order struct {
	Symbol       string `json:"symbol"`
	OrderID      string `json:"orderId"`
	OrderType    string `json:"orderType"`
	CancelType   string `json:"cancelType"`
	Price        string `json:"price"`
	Qty          string `json:"qty"`
	TimeInForce  string `json:"timeInForce"`
	OrderStatus  string `json:"orderStatus"`
	ReduceOnly   bool   `json:"reduceOnly"`
	Side         string `json:"side"`
	RejectReason string `json:"rejectReason"`
	CreateTime   string `json:"createdTime"`
	UpdateTime   string `json:"updatedTime"`
}