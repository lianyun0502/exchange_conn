package data_stream

import (
	"strconv"
	"strings"

	"github.com/valyala/fastjson"

	"github.com/lianyun0502/exchange_conn/v1"
)

type BinanceTradeStreams struct {
	EventType string `json:"e"`
	EventTime int64  `json:"E"`
	Symbol    string `json:"s"`
	TradeId   int64  `json:"t"`
	TradeTime int64  `json:"T"`
	Price     string `json:"p"`
	Quantity  string `json:"q"`
	IsMaker   bool   `json:"m"`
}
type BinanceAggregateTradeStreams struct {
	EventType    string `json:"e"`
	EventTime    int64  `json:"E"`
	Symbol       string `json:"s"`
	TradeId      int64  `json:"a"`
	TradeTime    int64  `json:"T"`
	Price        string `json:"p"`
	Quantity     string `json:"q"`
	FirstTradeId int64  `json:"f"`
	LastTradeId  int64  `json:"l"`
	IsMaker      bool   `json:"m"`
}

func ToNormalTradeData(rawData []byte) (data *exchange_conn.MultiTradeStream, err error) {
	v := fastjson.MustParseBytes(rawData)
	var side string
	if v.GetBool("m") {
		side = "sell"
	} else {
		side = "buy"
	}
	data = &exchange_conn.MultiTradeStream{Trades: make([]*exchange_conn.TradeStream, 0)}
	data.Trades = append(data.Trades, &exchange_conn.TradeStream{
		Topic:     string(v.GetStringBytes("e")),
		Time:      v.GetInt64("E"),
		Symbol:    string(v.GetStringBytes("s")),
		TradeId:   strconv.FormatInt(v.GetInt64("t"), 10),
		TradeTime: v.GetInt64("T"),
		Price:     string(v.GetStringBytes("p")),
		Quantity:  string(v.GetStringBytes("q")),
		Side:      strings.ToUpper(side),
	})

	return data, nil
}

func ToNormalAggregateTradeData(rawData []byte) (data *exchange_conn.MultiTradeStream, err error) {

	v := fastjson.MustParseBytes(rawData)
	var side string
	if v.GetBool("m") {
		side = "sell"
	} else {
		side = "buy"
	}
	data = &exchange_conn.MultiTradeStream{Trades: make([]*exchange_conn.TradeStream, 0)}
	data.Trades = append(data.Trades, &exchange_conn.TradeStream{
		Topic:     string(v.GetStringBytes("e")),
		Time:      v.GetInt64("E"),
		Symbol:    string(v.GetStringBytes("s")),
		TradeId:   strconv.FormatInt(v.GetInt64("a"), 10),
		TradeTime: v.GetInt64("T"),
		Price:     string(v.GetStringBytes("p")),
		Quantity:  string(v.GetStringBytes("q")),
		Side:      strings.ToUpper(side),
	})

	return data, nil
}
