package data_stream

import (
	"strconv"
	"strings"

	"github.com/lianyun0502/exchange_conn/v2/data_format"
	"github.com/valyala/fastjson"
)

type Trade struct {}

func NewTrade() *Trade {
	return &Trade{}
}

func (t *Trade) Update(rawdata []byte) (*format.MultiTradeStream, error) {
	raw := fastjson.MustParseBytes(rawdata)
	data := &format.MultiTradeStream{
		Trades: make([]*format.TradeStream, 0),
	}
	for i := 0; i < len(raw.GetArray("data")); i++ {
		data.Trades = append(data.Trades, &format.TradeStream{
			Topic:     string(raw.GetStringBytes("topic")),
			Time:      raw.GetInt64("ts"),
			Symbol:    string(raw.GetStringBytes("data", strconv.Itoa(i), "s")),
			TradeId:   string(raw.GetStringBytes("data", strconv.Itoa(i), "i")),
			TradeTime: raw.GetInt64("data", strconv.Itoa(i), "T"),
			Price:     string(raw.GetStringBytes("data", strconv.Itoa(i), "p")),
			Quantity:  string(raw.GetStringBytes("data", strconv.Itoa(i), "v")),
			Side:      strings.ToUpper(string(raw.GetStringBytes("data", strconv.Itoa(i), "S"))),
		})
	}
	return data, nil
}
