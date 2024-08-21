package data_stream

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

type OrderBook struct {
	exchange_conn.OrderBookStream
	orderCache *RingQueue[*fastjson.Value] `json:"-"`
	L          *sync.Mutex                 `json:"-"`
	isInit     bool                        `json:"-"`
	snapshot   *fastjson.Value             `json:"-"`
	lastID     int                         `json:"-"`
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		OrderBookStream: exchange_conn.OrderBookStream{
			Bids: make(map[string]string),
			Asks: make(map[string]string),
		},
		orderCache: NewRingQueue[*fastjson.Value](30),
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
	cache:= ob.orderCache
	i := cache.Start
	for ; i != cache.End; i=(i+1)%cache.Size() {
		order := cache.Get(i)
		if order.GetInt("U") <= ob.lastID+1 && order.GetInt("u") >= ob.lastID+1 {
			fmt.Printf("%v (last ID) +1 >= %v (first ID) and <= %v (final ID) \n", ob.lastID, order.GetInt("U"), order.GetInt("u"))
			break
		}
	}
	if i == cache.End {
		return errors.New("no snapshot")
	}
	for ; i != cache.End; i=(i+1)%cache.Size() {
		order := cache.Get(i)
		UpdateCurrentOrder(order.GetArray("b"), ob.Bids)
		UpdateCurrentOrder(order.GetArray("a"), ob.Asks)
	}
	return

}

func (ob *OrderBook) cache(rawData []byte) error {
	v, err := p.ParseBytes(rawData)
	if err != nil {
		return err
	}
	// fmt.Println(v)
	ob.orderCache.Push(v)
	return nil
}

func (ob *OrderBook) Update(rawData []byte) (data *OrderBook, err error) {
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
	}else{
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

type BinancePartialOrderBook struct {
	LastUpdateID int64      `json:"lastUpdateId"`
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`
}

func ToConsistentOrderBook(rawData []byte) (data *exchange_conn.OrderBookStream, err error) {
	v, err := p.ParseBytes(rawData)
	if err != nil {
		return nil, err
	}
	bids := make(map[string]string)
	asks := make(map[string]string)

	bidsArray := v.GetArray("bids")
	for i := 0; i < len(bidsArray); i++ {
		price := string(bidsArray[i].GetStringBytes("0"))
		quantity := string(bidsArray[i].GetStringBytes("1"))
		bids[price] = quantity
	}
	asksArray := v.GetArray("asks")
	for i := 0; i < len(asksArray); i++ {
		price := string(asksArray[i].GetStringBytes("0"))
		quantity := string(asksArray[i].GetStringBytes("1"))
		asks[price] = quantity
	}

	for i := 0; i < len(v.GetArray("bids")); i++ {

		data = &exchange_conn.OrderBookStream{
			Bids: bids,
			Asks: asks,
		}
	}
	return data, nil
}
