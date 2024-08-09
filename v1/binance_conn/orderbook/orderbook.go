package orderbook

import (
	// "encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/lianyun0502/exchange_conn/v1"
	"github.com/valyala/fastjson"
)

var p fastjson.Parser

type OrderBook struct {
	exchange_conn.OrderBookStream

	orderCache *fastjson.Value `json:"-"`
	L          *sync.Mutex     `json:"-"`
	isInit     bool            `json:"-"`
	snapshot   *fastjson.Value `json:"-"`
	lastID     int          `json:"-"`
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		OrderBookStream: exchange_conn.OrderBookStream{
			Bids: make(map[string]string),
			Asks: make(map[string]string),
		},
		orderCache: fastjson.MustParse("[]"),
		isInit:     false,
		L:          &sync.Mutex{},
	}
}
func (ob *OrderBook) IsInit() bool {
	return ob.isInit
}

func (ob *OrderBook) SetSnapshot(snapshot []byte) (err error) {
	ob.L.Lock()
	ob.snapshot, err = p.ParseBytes(snapshot)
	if err != nil {
		return err
	}
	bids := ob.snapshot.GetArray("bids")
	UpdateCurrentOrder(bids, ob.Bids)
	asks := ob.snapshot.GetArray("asks")
	UpdateCurrentOrder(asks, ob.Asks)
	ob.lastID = ob.snapshot.GetInt("lastUpdateId")
	ob.L.Unlock()
	return
}

func (ob *OrderBook) Init(rawData []byte) (err error) {
	cache, _ := ob.orderCache.Array()
	var i = 0
	for i=0; i < len(cache); i++ {
		order := cache[i]
		if order.GetInt("U") <= ob.lastID +1 && order.GetInt("u") >= ob.lastID + 1 {
			fmt.Printf("%v (last ID) +1 >= %v (first ID) and <= %v (final ID) \n", ob.lastID, order.GetInt("U"), order.GetInt("u"))
			break
		}
	}
	cache = cache[i:]
	// println(len(cache), "i=", i)
	if len(cache) == 0 {
		return errors.New("no snapshot")	
	}
	for i := 0; i < len(cache); i++ {
		order := cache[i]
		UpdateCurrentOrder(order.GetArray("b"), ob.Bids)
		UpdateCurrentOrder(order.GetArray("a"), ob.Asks)
		// ob.isInit = true
	}

	return

}

func (ob *OrderBook) cache(rawData []byte) error {
	v, err := p.ParseBytes(rawData)
	if err != nil {
		return err
	}
	// fmt.Println(v)
	cache, _ := ob.orderCache.Array()
	ob.orderCache.SetArrayItem(len(cache), v)

	return nil
}

func (ob *OrderBook) Update(rawData []byte) (data *OrderBook, err error) {
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
	} else {
		// if orderbook is not initialized, cache the data
		err = ob.cache(rawData)
		if err != nil {
			return nil, err
		}
		if ob.snapshot != nil {
			// fmt.Println("init orderbook")
			err = ob.Init(nil)
			if err != nil {
				return nil, err
			}
			ob.isInit = true
		}

		return nil, nil
	}
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
