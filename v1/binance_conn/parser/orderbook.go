package parser

import (
	// "encoding/json"
	"github.com/lianyun0502/exchange_conn/v1"
	"github.com/valyala/fastjson"
	"strconv"
	"sync"
)

var p fastjson.Parser

type OrderBook struct {
	exchange_conn.OrderBookStream

	orderCache []*fastjson.Value `json:"-"`
	L          *sync.Mutex 	 `json:"-"`
	isInit     bool 		  `json:"-"`
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		OrderBookStream: exchange_conn.OrderBookStream{
			Bids:       make(map[string]string),
			Asks:       make(map[string]string),
		},
		orderCache: make([]*fastjson.Value, 0),
		isInit:     false,
		L:          &sync.Mutex{},
	}
}

func (ob *OrderBook) Init(rawData []byte) (err error) {

	snapshot, err := p.ParseBytes(rawData)
	if err != nil {
		return err
	}
	lastID := snapshot.GetInt64("lastUpdateId")
	bids := snapshot.GetArray("bids")
	for _, v := range bids {
		ob.Bids[string(v.GetStringBytes("0"))] = string(v.GetStringBytes("1"))
	}
	asks := snapshot.GetArray("asks")
	for _, v := range asks {
		ob.Asks[string(v.GetStringBytes("0"))] = string(v.GetStringBytes("1"))
	}

	ob.L.Lock()
	for i := 0; i < len(ob.orderCache); i++ {
		if ob.orderCache[i].GetInt64("u") <= lastID {
			continue
		}
		// if ob.cache[i].GetInt64("U") < lastID + 1 && ob.cache[i].GetInt64("u") >= lastID + 1 {
		if ob.orderCache[i].GetInt64("U") <= lastID+1 {
			UpdateCurrentOrder(ob.orderCache[i].GetArray("b"), ob.Bids)
			UpdateCurrentOrder(ob.orderCache[i].GetArray("a"), ob.Asks)
		}
	}
	ob.isInit = true
	ob.L.Unlock()
	return

}

func (ob *OrderBook) cache(rawData []byte) error {
	v, err := p.ParseBytes(rawData)
	if err != nil {
		return err
	}
	ob.orderCache = append(ob.orderCache, v)
	return nil
}

func (ob *OrderBook) Update(rawData []byte) (Rawdata *OrderBook, err error) {
	ob.L.Lock()
	if ob.isInit {
		data, err := p.ParseBytes(rawData)
		if err != nil {
			return nil, err
		}
		UpdateCurrentOrder(data.GetArray("b"), ob.Bids)
		UpdateCurrentOrder(data.GetArray("a"), ob.Asks)

		ob.Topic = string(data.GetStringBytes("e"))
		ob.Time = data.GetInt64("E")
		ob.Symbol = string(data.GetStringBytes("s"))
		return ob, nil
	}
	// if orderbook is not initialized, cache the data
	err = ob.cache(rawData)
	if err != nil {
		return nil, err
	}
	ob.L.Unlock()
	return nil, nil
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
