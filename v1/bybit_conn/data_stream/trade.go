package data_stream

import (
	"github.com/lianyun0502/exchange_conn/v1"
	"github.com/valyala/fastjson"
	"strconv"
)

type Trade struct {}

func NewTrade() *Trade {
	return &Trade{}
}

func (t *Trade) Update(rawdata []byte) ([]*exchange_conn.TradeStream, error) {
	raw := fastjson.MustParseBytes(rawdata)
	trades := make([]*exchange_conn.TradeStream, 0)
	for i := 0; i < len(raw.GetArray("data")); i++ {
		trades = append(trades, &exchange_conn.TradeStream{
			Topic:     string(raw.GetStringBytes("topic")),
			Time:      raw.GetInt64("ts"),
			Symbol:    string(raw.GetStringBytes("data", strconv.Itoa(i), "s")),
			TradeId:   string(raw.GetStringBytes("data", strconv.Itoa(i), "i")),
			TradeTime: raw.GetInt64("data", strconv.Itoa(i), "T"),
			Price:     string(raw.GetStringBytes("data", strconv.Itoa(i), "p")),
			Quantity:  string(raw.GetStringBytes("data", strconv.Itoa(i), "v")),
			Side:      string(raw.GetStringBytes("data", strconv.Itoa(i), "S")),
		})
	}
	return trades, nil
}
