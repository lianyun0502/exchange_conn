package data_stream

import (
	"errors"
	"strconv"
	"sort"

	"github.com/lianyun0502/exchange_conn/v2"
	"github.com/valyala/fastjson"
	"github.com/duke-git/lancet/v2/maputil"
)

type OrderBook struct {
	BestDepth int
	Bids      map[string]string
	Asks      map[string]string
}

func NewOrderBook(bestDepth int) *OrderBook {
	return &OrderBook{
		BestDepth: bestDepth,
		Bids: make(map[string]string),
		Asks: make(map[string]string),
	}
}

func (ob *OrderBook) Update(rawdata []byte) (*exchange_conn.OrderBookStream, error) {
	raw := fastjson.MustParseBytes(rawdata)
	ret := &exchange_conn.OrderBookStream{
		Time: raw.GetInt64("ts"),
	}
	types := string(raw.GetStringBytes("type"))
	switch types {
	case "snapshot":
		ret.Topic = string(raw.GetStringBytes("topic"))
		ret.Symbol = string(raw.GetStringBytes("data", "s"))
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
	ret.Bids = BestMap(ob.Bids, -ob.BestDepth)
	ret.Asks = BestMap(ob.Asks, ob.BestDepth)
	return ret, nil
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

func BestMap(src map[string]string, num int) map[string]string {
	bestMap := make(map[string]string)
	keys := maputil.Keys(src)
	sort.Strings(keys)
	if num >0 && num < len(keys) {
		keys = keys[:num]
	} else if num < 0 && num > -len(keys) {
		keys = keys[len(keys)+num:]
	}
	for _, key := range keys {
		bestMap[key] = src[key]
	}
	return bestMap
}
