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

// 接收自交易所 orderbook data 的結構體(websocket)
type UpdateDepth struct {
	EventType string     `json:"e"`
	EventTime int64      `json:"E"`
	Symbol    string     `json:"s"`
	FirstId   int64      `json:"U"`
	LastId    int64      `json:"u"`
	PreID     int64      `json:"pu"`
	Bids      [][]string `json:"b"`
	Asks      [][]string `json:"a"`
}


// 接收自交易所 更新 orderbook data 的結構體 (API)
type SnapshotDepth struct {
	LastUpdateId int64 `json:"lastUpdateId"`
	Symbol       string `json:"symbol"`
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`
}

// OB組合物件
type OrderBook struct {
	API *http_client.BinanceClient
	Depth *SnapshotDepth
	Symbol string
	bidPriceMap *xsync.MapOf[string, []string]
	askPriceMap *xsync.MapOf[string, []string]

	DataQueue *xsync.SPSCQueueOf[*UpdateDepth]

	LastId  int64
	IsInit  bool
}

func NewOrderBook(symbol string, host_type string) (*OrderBook, error) {
	api, err := http_client.NewAPIClient(host_type, "", "")
	if err != nil {
		return nil, err
	}
	return &OrderBook{
		API: api,
		Symbol: symbol,
		DataQueue: xsync.NewSPSCQueueOf[*UpdateDepth](100),
		LastId:  0,
		IsInit:  false,
		bidPriceMap: xsync.NewMapOf[string, []string](),
		askPriceMap: xsync.NewMapOf[string, []string](),
	}, nil
}

func (ob *OrderBook) UpdateSnapshot(snapshot SnapshotDepth) *OrderBook {
	ob.Depth = &snapshot
	ob.Depth.Symbol = ob.Symbol
	ob.askPriceMap.Clear()
	ob.bidPriceMap.Clear()

	for _, bid := range ob.Depth.Bids {
		ob.bidPriceMap.Store(bid[0], bid)
	}
	for _, ask := range ob.Depth.Asks {
		ob.askPriceMap.Store(ask[0], ask)
	}
	return ob
}

func (ob *OrderBook) IsFrameLose(depthUpdate *UpdateDepth, hostType string) bool{
	switch hostType {
	case "spot":
		return depthUpdate.FirstId-ob.LastId > 20 
	case "future":
		return depthUpdate.PreID != ob.LastId
	default:
		return false
	}

}

func (ob *OrderBook) ComposeDepth() (bool, error) {
	var endpoint string
	switch ob.API.Exchange.HostType {
	case "spot":
		endpoint = "/api/v3/depth"
	case "future":
		endpoint = "/fapi/v1/depth"
	default:
		return false , fmt.Errorf("host type error")	
	}
	req := ob.API.Request(http.MethodGet, endpoint)
	query := map[string]string{"symbol": ob.Symbol, "limit": "10"}
	req.SetQuery(query)
	data, err := req.SendWithTimeout(5)
	if err != nil {
		return false, err
	}
	var d = new(SnapshotDepth)
	if err = json.Unmarshal(data, d); err != nil {
		return false, err
	}
	for {
		depthUpdate, ok := ob.DataQueue.TryDequeue()
		// fmt.Println(ok)
		if ok {
			ob.API.Log.Debugf("depthUpdate.FirstId: %d, depthUpdate.LastId: %d,  d.LastUpdateId: %d\n", depthUpdate.FirstId, depthUpdate.LastId, d.LastUpdateId)
			if d.LastUpdateId <= depthUpdate.FirstId {
				return false, fmt.Errorf("out of orderbook frame")
			}
			if d.LastUpdateId <= depthUpdate.LastId {
				ob.UpdateSnapshot(*d)
				ob.LastId = depthUpdate.LastId
				ob.IsInit = true
				return true, nil
			}
		} else {
			break
		}

	}
	return false, fmt.Errorf("no valid depth update")

}

type OrderBookManager struct {
	Host_type string
	*xsync.MapOf[string, *OrderBook] // map[coin]* OB

	InitQueue chan *OrderBook
}

func NewOrderBookManager(host_type string) (*OrderBookManager, error) {
	// api, err := http_client.NewAPIClient(host_type, "", "")
	// if err != nil {
	// 	return nil, err
	// }
	ob := &OrderBookManager{
		Host_type: host_type,
		MapOf: xsync.NewMapOf[string, *OrderBook](),
		InitQueue: make(chan *OrderBook, 100),
	}
	go ob.StartInitQueue()
	return ob, nil
}


// 判斷掉frame情況
func (obs *OrderBookManager) IsFrameLose(depthUpdate *UpdateDepth, ob *OrderBook) bool{
	return ob.IsFrameLose(depthUpdate, ob.API.Exchange.HostType)
}

func (obs *OrderBookManager) Update(rawData []byte, opts ...func(*OrderBookManager)) (data *format.OrderBookStream, err error) {
	for _, opt := range opts {
		opt(obs)
	}
	// fmt.Println(string(rawData))
	var depthUpdate = new(UpdateDepth)
	json.Unmarshal(rawData, depthUpdate)
	orderBook, exist := obs.Load(depthUpdate.Symbol)
	if exist {
		if orderBook.DataQueue.TryEnqueue(depthUpdate) {
			// fmt.Println("load success")
			// fmt.Println(depthUpdate)
		}
	} else {
		orderBook, err := NewOrderBook(depthUpdate.Symbol, obs.Host_type)
		if err != nil {
			return nil, err
		}
		if orderBook.DataQueue.TryEnqueue(depthUpdate) {
			// fmt.Println("store success")
		}
		fmt.Printf("store new orderbook %s\n", depthUpdate.Symbol)
		obs.Store(depthUpdate.Symbol, orderBook)
		obs.InitQueue <- orderBook

		return nil, nil
	}
	if orderBook == nil {
		return nil, fmt.Errorf("%s orderbook is nil", depthUpdate.Symbol)
	}

	if !orderBook.IsInit {
		return nil, fmt.Errorf("%s orderbook is not init", depthUpdate.Symbol)
	}

	for {
		if depth, ok := orderBook.DataQueue.TryDequeue(); !ok {
			break
		} else {
			depthUpdate = depth
			if obs.IsFrameLose(depthUpdate, orderBook) {
				obs.Delete(depthUpdate.Symbol)
				return nil, fmt.Errorf("%s frame loss", depthUpdate.Symbol)
			}
			orderBook.LastId = depthUpdate.LastId
		}

		for _, bid := range depthUpdate.Bids {
			if b, ok := orderBook.bidPriceMap.Load(bid[0]); ok {
				// fmt.Printf("b[1]: %s, bid[1]: %s\n", b[1], bid[1])
				b[1] = bid[1]
			} else {
				orderBook.Depth.Bids = append(orderBook.Depth.Bids, bid)
				orderBook.bidPriceMap.Store(bid[0], bid)
			}

		}
		// fmt.Printf("ob.ob.Depth.Bids len : %d\n", len(ob.ob.Depth.Bids))
		sort.Slice(orderBook.Depth.Bids, func(i, j int) bool {
			if qtyj, _ := strconv.ParseFloat(orderBook.Depth.Bids[j][1], 64); qtyj == 0 {
				return true
			}
			if qtyi, _ := strconv.ParseFloat(orderBook.Depth.Bids[i][1], 64); qtyi == 0 {
				// fmt.Printf("ob.ob.Depth.Bids[i][1]: %s\n", ob.ob.Depth.Bids[i][1])
				return false
			}
			pricei, _ := strconv.ParseFloat(orderBook.Depth.Bids[i][0], 64)
			pricej, _ := strconv.ParseFloat(orderBook.Depth.Bids[j][0], 64)
			// fmt.Printf("ob.ob.Depth.Bids[i][0]: %d, ob.ob.Depth.Bids[j][0]: %d\n", i, j)
			return pricei > pricej
		})
		if len(orderBook.Depth.Bids) > 100 {
			for _, bid := range orderBook.Depth.Bids[100:] {
				// delete(ob.ob.PriceMap, bid[0])
				orderBook.bidPriceMap.Delete(bid[0])
			}
			orderBook.Depth.Bids = orderBook.Depth.Bids[:100]
		}
		

		for _, ask := range depthUpdate.Asks {
			if a, ok := orderBook.askPriceMap.Load(ask[0]); ok {
				a[1] = ask[1]
			} else {
				orderBook.Depth.Asks = append(orderBook.Depth.Asks, ask)
				orderBook.askPriceMap.Store(ask[0], ask)
			}
		}
		// fmt.Printf("ob.ob.Depth.Asks len: %d\n", len(ob.ob.Depth.Asks))
		sort.Slice(orderBook.Depth.Asks, func(i, j int) bool {
			if qtyi, _ := strconv.ParseFloat(orderBook.Depth.Asks[i][1], 64); qtyi == 0 {
				return false
			}
			if qtyj, _ := strconv.ParseFloat(orderBook.Depth.Asks[j][1], 64); qtyj == 0 {
				return true
			}
			pricei, _ := strconv.ParseFloat(orderBook.Depth.Asks[i][0], 64)
			pricej, _ := strconv.ParseFloat(orderBook.Depth.Asks[j][0], 64)
			return pricei < pricej
		})
		if len(orderBook.Depth.Asks) > 100 {
			for _, ask := range orderBook.Depth.Asks[100:] {
				// delete(ob.ob.PriceMap, ask[0])
				orderBook.askPriceMap.Delete(ask[0])
			}
			orderBook.Depth.Asks = orderBook.Depth.Asks[:100]
		}
	}
	// fmt.Println("update success")
	// fmt.Println(depthUpdate.EventTime)
	ret := &format.OrderBookStream{
		Bids: make(map[string]string),
		Asks: make(map[string]string),
	}
	if orderBook.Depth.Bids[0][1] == "0.000" {
		fmt.Println(depthUpdate)
		fmt.Println(orderBook.Depth.Bids[0])
	}
	for _, bid := range orderBook.Depth.Bids[:10] {
		ret.Bids[bid[0]] = bid[1]
	}
	for _, ask := range orderBook.Depth.Asks[:10] {
		ret.Asks[ask[0]] = ask[1]
	}
	ret.Symbol = orderBook.Depth.Symbol
	ret.Time = depthUpdate.EventTime
	ret.Topic = depthUpdate.EventType
	return ret, nil
}

func (obs *OrderBookManager) StartInitQueue() {
	for orderbook := range obs.InitQueue {
		if orderbook.IsInit {
			continue
		}
		time.Sleep(10 * time.Millisecond)
		if ok, err := orderbook.ComposeDepth(); !ok {
			orderbook.API.Log.Error(err)
			if (err.Error() == "request timeout") || 
			(err.Error() == "no valid depth update") ||
			(err.Error() == "out of orderbook frame") {
				obs.Delete(orderbook.Symbol)
			}
		}
	}

}


