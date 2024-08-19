package data_stream

import (
	"errors"
	"strconv"

	"github.com/lianyun0502/exchange_conn/v1"
	"github.com/valyala/fastjson"
)

type OrderBook exchange_conn.OrderBookStream

func NewOrderBook() *OrderBook {
	return &OrderBook{
		Bids: make(map[string]string),
		Asks: make(map[string]string),
	}
}

func (ob *OrderBook) Update(rawdata []byte) (*OrderBook, error) {
	raw := fastjson.MustParseBytes(rawdata)
	ob.Time = raw.GetInt64("ts")
	types := string(raw.GetStringBytes("type"))
	switch types {
	case "snapshot":
		ob.Topic = string(raw.GetStringBytes("topic"))
		ob.Symbol = string(raw.GetStringBytes("data", "s"))
		clear(ob.Bids)
		clear(ob.Asks)
		UpdateCurrentOrder(raw.GetArray("data", "b"), ob.Bids)
		UpdateCurrentOrder(raw.GetArray("data", "a"), ob.Asks)
	case "delta":
		UpdateCurrentOrder(raw.GetArray("data", "b"), ob.Bids)
		UpdateCurrentOrder(raw.GetArray("data", "a"), ob.Asks)
	default:
		return nil, errors.New("unknown type")
	}
	return ob, nil
}

func UpdateCurrentOrder(srcOrders []*fastjson.Value, curOrders map[string]string) {
	for i := 0; i < len(srcOrders); i++ {
		quantity := string(srcOrders[i].GetStringBytes("1"))
		price := string(srcOrders[i].GetStringBytes("0"))
		if f, _ := strconv.ParseFloat(quantity, 64); f == 0 {
			delete(curOrders, price)
			continue
		}
		curOrders[price] = quantity
	}

}
