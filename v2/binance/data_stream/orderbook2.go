package data_stream

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	http_client "github.com/lianyun0502/exchange_conn/v2/binance/http_client"
	"github.com/lianyun0502/exchange_conn/v2/data_format"
	"github.com/puzpuzpuz/xsync/v3"
)

type DepthUpdate struct {
	EventType string `json:"e"`
	EventTime int64  `json:"E"`
	Symbol    string `json:"s"`
	FirstId   int64  `json:"U"`
	LastId    int64  `json:"u"`
	Bids      [][]string `json:"b"`
	Asks      [][]string `json:"a"`
}

type Depth struct {
	LastUpdateId int64 `json:"lastUpdateId"`
	Symbol       string
	Bids         [][]string `json:"bids"`
	Asks 		[][]string `json:"asks"`
}

type OrderBook struct{
	Depth *Depth
	// PriceMap map[string][]string
	PriceMap *xsync.MapOf[string, []string]
}

func NewOrderBook(d *Depth, symbol string) *OrderBook {
	d.Symbol = symbol
	ob := &OrderBook{
		Depth: d,
		PriceMap: xsync.NewMapOf[string, []string](),
	}
	for _, bid := range d.Bids {
		// ob.PriceMap[bid[0]] = bid
		ob.PriceMap.Store(bid[0], bid)
	}
	for _, ask := range d.Asks {
		// ob.PriceMap[ask[0]] = ask
		ob.PriceMap.Store(ask[0], ask)
	}
	return ob
}


type OrderBooks struct {
	API *http_client.BinanceClient
	*xsync.MapOf[string, *OBObject]
}



func NewOrderBookMap(host_type string) (*OrderBooks, error){
	api, err := http_client.NewAPIClient(host_type, "", "")
	if err != nil {
		return nil , err
	}
	ob := &OrderBooks{
		API: api,
		MapOf: xsync.NewMapOf[string, *OBObject](),
	}
	return ob, nil
}

type OBObject struct{
	ob *OrderBook
	EvQueue *xsync.SPSCQueueOf[*DepthUpdate]
	LastId int64
}

func NewOBObject() *OBObject {
	return &OBObject{
		ob: nil,
		EvQueue: xsync.NewSPSCQueueOf[*DepthUpdate](100),
	}
}

func (obs *OrderBooks) Update(rawData []byte, opts... func(*OrderBooks)) (data *format.OrderBookStream, err error){
	for _, opt := range opts {
		opt(obs)
	}
	// fmt.Println(string(rawData))
	var depthUpdate = new(DepthUpdate)
	json.Unmarshal(rawData, depthUpdate)
	ob, ok := obs.Load(depthUpdate.Symbol)
	if ok { 
		if ob.EvQueue.TryEnqueue(depthUpdate) {
			// fmt.Println("load success")
			// fmt.Println(depthUpdate)
		}
	} else {
		obj := NewOBObject()
		obs.Store(depthUpdate.Symbol, obj)
		if obj.EvQueue.TryEnqueue(depthUpdate) {
			// fmt.Println("store success")
		}
		go func () {
			for {
				if obs.Init(depthUpdate.Symbol) {
					// fmt.Println("init success")
					break
				}
			}
		}()
		return nil, nil
	}
	if ob.ob == nil {
		return nil , fmt.Errorf("orderbook is nil")
	}
	if ob.ob.Depth == nil {
		return nil , fmt.Errorf("depth is nil")
	}
	for {
		if depth, ok := ob.EvQueue.TryDequeue(); !ok {
			break
		}else{
			depthUpdate = depth
			if depthUpdate.FirstId - ob.ob.Depth.LastUpdateId > 10 {
				fmt.Println("depthUpdate.FirstId - depth.LastUpdateId > 10")
				obs.Delete(depthUpdate.Symbol)
				return nil, fmt.Errorf("depthUpdate.FirstId - depth.LastUpdateId > 10")
			}
			ob.ob.Depth.LastUpdateId = depthUpdate.LastId
		}
		
		for _, bid := range depthUpdate.Bids {
			if b, ok := ob.ob.PriceMap.Load(bid[0]); ok {
				// fmt.Printf("b[1]: %s, bid[1]: %s\n", b[1], bid[1])
				b[1] = bid[1]
			}else{
				ob.ob.Depth.Bids = append(ob.ob.Depth.Bids, bid)
				ob.ob.PriceMap.Store(bid[0], bid)
			}

		}
		// fmt.Printf("ob.ob.Depth.Bids len : %d\n", len(ob.ob.Depth.Bids))
		sort.Slice(ob.ob.Depth.Bids, func(i, j int) bool {
			if qtyi, _ := strconv.ParseFloat(ob.ob.Depth.Bids[i][1], 64); qtyi == 0 {
				return false
			}
			if qtyj, _ := strconv.ParseFloat(ob.ob.Depth.Bids[j][1], 64); qtyj == 0 {
				return true
			}
			pricei, _ := strconv.ParseFloat(ob.ob.Depth.Bids[i][0], 64)
			pricej, _ := strconv.ParseFloat(ob.ob.Depth.Bids[j][0], 64)
			// fmt.Printf("ob.ob.Depth.Bids[i][0]: %d, ob.ob.Depth.Bids[j][0]: %d\n", i, j)
			return pricei > pricej
		})
		// for i := 0; i < 10; i++ {
		// 	fmt.Println(ob.ob.Depth.Bids[i])
		// }
		for _, bid := range ob.ob.Depth.Bids[10:] {
			// delete(ob.ob.PriceMap, bid[0])
			ob.ob.PriceMap.Delete(bid[0])
		}
		ob.ob.Depth.Bids = ob.ob.Depth.Bids[:10]

		for _, ask := range depthUpdate.Asks {
			if a, ok := ob.ob.PriceMap.Load(ask[0]); ok {
				a[1] = ask[1]
			}else{
				ob.ob.Depth.Asks = append(ob.ob.Depth.Asks, ask)
				ob.ob.PriceMap.Store(ask[0], ask)
			}
		}
		// fmt.Printf("ob.ob.Depth.Asks len: %d\n", len(ob.ob.Depth.Asks))
		sort.Slice(ob.ob.Depth.Asks, func(i, j int) bool {
			if qtyi, _ := strconv.ParseFloat(ob.ob.Depth.Asks[i][1], 64); qtyi == 0 {
				return false
			}
			if qtyj, _ := strconv.ParseFloat(ob.ob.Depth.Asks[j][1], 64); qtyj == 0 {
				return true
			}
			pricei, _ := strconv.ParseFloat(ob.ob.Depth.Asks[i][0], 64)
			pricej, _ := strconv.ParseFloat(ob.ob.Depth.Asks[j][0], 64)
			return pricei < pricej
		})
		for _, ask := range ob.ob.Depth.Asks[10:] {
			// delete(ob.ob.PriceMap, ask[0])
			ob.ob.PriceMap.Delete(ask[0])
		}
		ob.ob.Depth.Asks = ob.ob.Depth.Asks[:10]
	}
	// fmt.Println("update success")
	// fmt.Println(depthUpdate.EventTime)
	ret := &format.OrderBookStream{
		Bids: make(map[string]string),
		Asks: make(map[string]string),
	}
	for _, bid := range ob.ob.Depth.Bids {
		ret.Bids[bid[0]] = bid[1]
	}
	for _, ask := range ob.ob.Depth.Asks {
		ret.Asks[ask[0]] = ask[1]
	}
	ret.Symbol = ob.ob.Depth.Symbol
	ret.Time = depthUpdate.EventTime
	ret.Topic = depthUpdate.EventType
	return ret, nil
}

func (obs *OrderBooks) Init(symbol string) bool {
	var endpoint string
	switch obs.API.Exchange.HostType{
	case "spot":
		endpoint = "/api/v3/depth"
	case "future":
		endpoint = "/fapi/v1/depth"
	default:
		return false
	}
	req := obs.API.Request(http.MethodGet, endpoint)
	query := map[string]string{"symbol": symbol, "limit": "10"}
	req.SetQuery(query)
	data, err := req.Send()
	if err != nil {
		print(err)
		return false
	}
	time.Sleep(50 * time.Millisecond)
	var d = new(Depth)
	if err = json.Unmarshal(data, d); err != nil {
		print(err)
		return false
	}
	d.Symbol = symbol
	obj, ok := obs.Load(symbol)
	if !ok {
		print(err)
		return false
	}
	// fmt.Println(d.LastUpdateId)
	for {
		depthUpdate, ok := obj.EvQueue.TryDequeue()
		// fmt.Println(ok)
		if ok {
			// fmt.Printf("depthUpdate.FirstId: %d, depthUpdate.LastId: %d,  d.LastUpdateId: %d\n", depthUpdate.FirstId, depthUpdate.LastId, d.LastUpdateId)
			if d.LastUpdateId <= depthUpdate.FirstId {
				return false
			}
			if d.LastUpdateId <= depthUpdate.LastId {
				obj.ob = NewOrderBook(d, symbol)
				obj.LastId = depthUpdate.LastId
				return true
			}
		}else{
			break
		}
		
	}
	return false
}
