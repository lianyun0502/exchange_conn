package parser

import (
	"encoding/json"

	"github.com/lianyun0502/exchange_conn/v1"
)


type BinanceTradeStreams struct {
	EventType string `json:"e"`
	EventTime int64 `json:"E"`
	Symbol string `json:"s"`
	TradeId int64 `json:"t"`
	TradeTime int64 `json:"T"`
	Price string `json:"p"`
	Quantity string `json:"q"`
	IsMaker bool `json:"m"`
}
type BinanceAggregateTradeStreams struct {
	EventType string `json:"e"`
	EventTime int64 `json:"E"`
	Symbol string `json:"s"`
	TradeId int64 `json:"a"`
	TradeTime int64 `json:"T"`
	Price string `json:"p"`
	Quantity string `json:"q"`
	FirstTradeId int64 `json:"f"`
	LastTradeId int64 `json:"l"`
	IsMaker bool `json:"m"`
}

func ToNormalTradeData(rawData []byte) (data *exchange_conn.TradeSrteam, err error) {
	
	raw := new(BinanceTradeStreams)
	err = json.Unmarshal(rawData, raw)
	if err != nil {
		return 
	}
	var side string
	if raw.IsMaker {
		side = "sell"
	} else {
		side = "buy"
	}	
	data = &exchange_conn.TradeSrteam{
		Topic: raw.EventType,
		Time: raw.EventTime,
		Symbol: raw.Symbol,
		TradeId: raw.TradeId,
		TradeTime: raw.TradeTime,
		Price: raw.Price,
		Quantity: raw.Quantity,
		Side: side,
	}
	

	return data, nil
}

func ToNormalAggregateTradeData(rawData []byte) (data *exchange_conn.TradeSrteam, err error) {
	
	raw := new(BinanceAggregateTradeStreams)
	err = json.Unmarshal(rawData, raw)
	if err != nil {
		return 
	}
	var side string
	if raw.IsMaker {
		side = "sell"
	} else {
		side = "buy"
	}	
	data = &exchange_conn.TradeSrteam{
		Topic: raw.EventType,
		Time: raw.EventTime,
		Symbol: raw.Symbol,
		TradeId: raw.TradeId,
		TradeTime: raw.TradeTime,
		Price: raw.Price,
		Quantity: raw.Quantity,
		Side: side,
	}
	

	return data, nil
}