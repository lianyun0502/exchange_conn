package data_stream

import (
	// "encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"sync"

	"github.com/duke-git/lancet/v2/maputil"
	"github.com/lianyun0502/exchange_conn/v2"
	"github.com/valyala/fastjson"
)

var p fastjson.Parser

type RingQueue[T any] struct {
	queue  []T
	isFull bool
	Start  int
	End    int
}

func NewRingQueue[T any](size int) *RingQueue[T] {
	return &RingQueue[T]{
		queue:  make([]T, size),
		isFull: false,
		Start:  0,
		End:    0,
	}
}

func (q *RingQueue[T]) Size() int {
	return len(q.queue)
}

func (q *RingQueue[T]) Push(data T) {
	q.queue[q.End] = data
	q.End = (q.End + 1) % q.Size()
	q.isFull = (q.End == q.Start)
	if q.isFull {
		q.Start = (q.Start + 1) % q.Size()
	}
}

func (q *RingQueue[T]) Get(index int) T {
	return q.queue[index]
}

func (q *RingQueue[T]) Clear() {
	q.Start = 0
	q.End = 0
	q.isFull = false
	clear(q.queue)
}

type OrderBookParser struct {
	orderCache *RingQueue[*fastjson.Value] `json:"-"`
	L          *sync.Mutex                 `json:"-"`
	isInit     bool                        `json:"-"`
	snapshot   *fastjson.Value             `json:"-"`
	lastID     int                         `json:"-"`
	BestDepth  int                         `json:"-"`
	Bids       map[string]string
	Asks       map[string]string
}

func NewOrderBookParser(depth int) *OrderBookParser {
	return &OrderBookParser{
		Bids: make(map[string]string),
		Asks: make(map[string]string),
		orderCache: NewRingQueue[*fastjson.Value](30),
		isInit:     false,
		L:          &sync.Mutex{},
		BestDepth: depth,
	}
}
func (ob *OrderBookParser) IsInit() bool {
	return ob.isInit
}

func (ob *OrderBookParser) SetSnapshot(snapshot []byte) (err error) {
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

func (ob *OrderBookParser) Init(rawData []byte) (err error) {
	cache := ob.orderCache
	i := cache.Start
	for ; i != cache.End; i = (i + 1) % cache.Size() {
		order := cache.Get(i)
		if order.GetInt("U") <= ob.lastID+1 && order.GetInt("u") >= ob.lastID+1 {
			fmt.Printf("%v (last ID) +1 >= %v (first ID) and <= %v (final ID) \n", ob.lastID, order.GetInt("U"), order.GetInt("u"))
			break
		}
	}
	if i == cache.End {
		return errors.New("no snapshot")
	}
	for ; i != cache.End; i = (i + 1) % cache.Size() {
		order := cache.Get(i)
		UpdateCurrentOrder(order.GetArray("b"), ob.Bids)
		UpdateCurrentOrder(order.GetArray("a"), ob.Asks)
	}
	return

}

func (ob *OrderBookParser) cache(rawData []byte) error {
	v, err := p.ParseBytes(rawData)
	if err != nil {
		return err
	}
	// fmt.Println(v)
	ob.orderCache.Push(v)
	return nil
}

func (ob *OrderBookParser) Update(rawData []byte) (data *exchange_conn.OrderBookStream, err error) {
	v, err := fastjson.ParseBytes(rawData)
	if err != nil {
		return nil, err
	}
	// 是否為深度更新
	if string(v.GetStringBytes("e")) == "depthUpdate" {
		err = ob.cache(rawData)
		if err != nil {
			return nil, err
		}
		if ob.isInit {
			ret := new(exchange_conn.OrderBookStream)
			data, err := p.ParseBytes(rawData)
			if err != nil {
				return nil, err
			}
			UpdateCurrentOrder(data.GetArray("b"), ob.Bids)
			UpdateCurrentOrder(data.GetArray("a"), ob.Asks)

			ret.Topic = string(data.GetStringBytes("e"))
			ret.Time = data.GetInt64("E")
			ret.Symbol = string(data.GetStringBytes("s"))
			ret.Bids = BestMap(ob.Bids, -ob.BestDepth)
			ret.Asks = BestMap(ob.Asks, ob.BestDepth)
			return ret, nil
		} else {
			if ob.snapshot != nil {
				err = ob.Init(nil)
				if err != nil {
					return nil, err
				}
				ob.isInit = true
			}
			return nil, nil
		}
	} else if v.GetInt("result", "lastUpdateId") != 0 && !ob.isInit {
		err = ob.SetSnapshot(v.GetStringBytes("result"))
		return nil, err
	} else {
		return nil, errors.New("unknown data")
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

type PartialOrderBook struct {
	BestDepth int
}

func NewPartialOrderBook(bestDepth int) *PartialOrderBook {
	return &PartialOrderBook{
		BestDepth: bestDepth,
	}
}

func (ob* PartialOrderBook) Update(rawData []byte) (data *exchange_conn.OrderBookStream, err error) {
	v, err := p.ParseBytes(rawData)
	if err != nil {
		return nil, err
	}
	data = &exchange_conn.OrderBookStream{
		Bids: make(map[string]string),
		Asks: make(map[string]string),
	}

	UpdateCurrentOrder(v.GetArray("b"), data.Bids)
	UpdateCurrentOrder(v.GetArray("a"), data.Asks)
	UpdateCurrentOrder(v.GetArray("bids"), data.Bids)
	UpdateCurrentOrder(v.GetArray("asks"), data.Asks)


	data.Topic = string(v.GetStringBytes("e"))
	data.Time = v.GetInt64("E")
	data.Symbol = string(v.GetStringBytes("s"))
	data.Bids = BestMap(data.Bids, -ob.BestDepth)
	data.Asks = BestMap(data.Asks, ob.BestDepth)
	return data, nil
}

func BestMap(src map[string]string, num int) map[string]string {
	var bestMap map[string]string
	keys := sort.StringSlice(maputil.Keys(src))
	keys.Sort()
	l := len(keys)
	if num > 0 {
		if 2*num <= l {
			keys = keys[:2*num]
			src = maputil.FilterByKeys(src, keys)
		}
		bestMap = maputil.FilterByKeys(src, keys[:num])

	} else if num < 0 {
		if 2*-num <= l {
			keys = keys[l+2*num:]
			src = maputil.FilterByKeys(src, keys)
		}
		bestMap = maputil.FilterByKeys(src, keys[l+num:])
	}
	return bestMap
}
