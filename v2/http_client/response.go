package httpClient

import "net/http"


type OrderResponse struct {
	Contract string 
	Exch string
	Market string
	OrderID string
	ClientID string
	Raw any
}

type OdrMsgTag string
const (
	WS     OdrMsgTag = "WS"
	Cancel OdrMsgTag = "Cancel"
	Query  OdrMsgTag = "Query"
)

type OrderMessage struct {
	Contract   string
	Exch         string
	Market        string
	Price      float64
	Qty        float64
	Side       int
	FilledQty  float64
	AvgFillPrc float64
	OrderID    string
	ClientID   string
	LocalTime  int64
	ServerTime int64
	Status     int
	OrderType  string
	Tag        OdrMsgTag
	Meta       map[string]any
	Raw        any
}


type Volume24hData struct {
	Volume     float64
	BaseVolume float64
}


type Ticker24hData struct {
	Symbol   string
	Open     float64
	High     float64
	Low      float64
	Close    float64
	Latest   float64
	Volume   float64
	Notional float64
	AskPrice float64
	BidPrice float64
	Raw      any
}

type ReqResponseItem[T any] struct {
	Data         []T
	HttpResponse *http.Response
	Error        error
}