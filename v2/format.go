package exchange_conn

import ()

type MultiTradeStream struct {
	Trades []*TradeStream `json:"Trds"`
}

type TradeStream struct {
	Topic     string `json:"Top"`
	Time      int64  `json:"T"`
	TradeTime int64  `json:"TrdT"`
	Symbol    string `json:"S"`
	TradeId   string `json:"TrdId"`
	Price     string `json:"P"`
	Quantity  string `json:"Q"`
	Side      string `json:"Sd"`
}

type Order struct {
}

type OrderBookStream struct {
	Topic  string            `json:"Top"`
	Time   int64             `json:"T"`
	Symbol string            `json:"S"`
	Bids   map[string]string `json:"Bids"`
	Asks   map[string]string `json:"Asks"`
}

type KLineStream struct {
	Topic     *string `json:"Top"`
	Time      *int64  `json:"T"`
	Symbol    *string `json:"S"`
	Interval  *string `json:"Intv"`
	StartTime *int64  `json:"ST"`
	EndTime   *int64  `json:"ET"`
	Open      *string `json:"O"`
	Close     *string `json:"C"`
	High      *string `json:"H"`
	Low       *string `json:"L"`
	Volume    *string `json:"V"`
}

type MarKetPriceStream struct {
	Topic           string `json:"Top"`
	Time            int64  `json:"T"`
	Symbol          string `json:"S"`
	Category        string `json:"C,omitempty"`
	LastPrice	   	string `json:"P,omitempty"`
	MarketPrice     string `json:"MrkP,omitempty"`
	IndexPrice      string `json:"IdxP,omitempty"`
	FundingRate     string `json:"FndR,omitempty"`
	NextFundingTime int64  `json:"NFT,omitempty"`
}
