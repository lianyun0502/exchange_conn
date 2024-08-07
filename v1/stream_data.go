package exchange_conn

import (
)

type TradeSrteam struct {
	Topic *string `json:"topic"`
	Time *int64 `json:"time"`
	TradeTime *int64 `json:"tradeTime"`
	Symbol *string `json:"symbol"`
	TradeId *int64 `json:"tradeId"`
	Price *string `json:"price"`
	Quantity *string `json:"quantity"`
	Side *string `json:"side"`
}

type Order struct {
}

type OrderBookStream struct {
	Topic *string `json:"topic"`
	Time *int64 `json:"time"`
	Symbol *string `json:"symbol"`
	Bids *[][]string `json:"bids"`
	Asks *[][]string `json:"asks"`
}


type KLineStream struct {
	Topic *string `json:"topic"`
	Time *int64 `json:"time"`
	Symbol *string `json:"symbol"`
	Interval *string `json:"interval"`
	StartTime *int64 `json:"startTime"`
	EndTime *int64 `json:"endTime"`
	Open *string `json:"open"`
	Close *string `json:"close"`
	High *string `json:"high"`
	Low *string `json:"low"`
	Volume *string `json:"volume"`
}