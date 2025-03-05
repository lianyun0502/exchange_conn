package data_stream_test

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fastjson"

	// "github.com/stretchr/testify/assert"
	ws "github.com/lianyun0502/exchange_conn/v2/binance/ws_client"
	// "github.com/lianyun0502/exchange_conn/v2/common"
	"github.com/lianyun0502/exchange_conn/v2/consts"
	"github.com/lianyun0502/exchange_conn/v2/binance/data_stream"
	"github.com/lianyun0502/exchange_conn/v2/data_format"
	"github.com/sirupsen/logrus"
	
)

func TestFastJson(t *testing.T) {
	v := fastjson.MustParse("[]")

	a, _ := v.Array()
	v.SetArrayItem(len(a), fastjson.MustParse("1"))
	a, _ = v.Array()
	v.SetArrayItem(len(a), fastjson.MustParse("2"))
	fmt.Println(v.GetInt("0"))
}

func TestDepthToSlice(t *testing.T) {
	rawData := `{"e":"depthUpdate","E":1723025111169,"s":"BTCUSDT","U":49981777515,"u":49981777539,"b":[["57261.83000000","0.00096000"],["57265.01000000","5.23762000"],["57261.66000000","0.00091000"],["57251.30000000","0.03018000"],["57249.66000000","0.04961000"],["57249.36000000","0.00000000"],["57248.79000000","0.00000000"],["57245.17000000","0.43654000"],["57197.85000000","8.11363000"],["57196.99000000","0.00000000"],["57185.39000000","0.00000000"],["57165.01000000","0.01000000"],["57158.97000000","0.00000000"],["33000.00000000","20.95952000"]],"a":[["57265.03000000","0.03705000"],["57266.71000000","0.00010000"],["57279.41000000","0.43654000"],["57282.15000000","0.74283000"],["57282.45000000","0.11479000"],["57284.50000000","0.13960000"],["57291.80000000","0.03024000"]]}`

	var p fastjson.Parser

	v, _ := p.Parse(rawData)
	bids := v.GetArray("b")
	asks := v.GetArray("a")
	m := map[string]string{}
	fmt.Println("==========")
	for _, bid := range bids {
		price := string(bid.GetStringBytes("0"))
		quality := string(bid.GetStringBytes("1"))
		fmt.Println(price)
		m[price] = quality
	}
	for k, v := range m {
		fmt.Println(k, ",", v)
	}
	fmt.Println("==========")
	for _, ask := range asks {
		price := ask.GetArray()[0].String()
		fmt.Println(price)
	}
	assert.Equal(t, 49981777515, v.GetInt("U"))
}

func TestUpdateCurrentOrder(t *testing.T) {
	rawData :=
	`
	{"e":"depthUpdate","E":1723025111169,"s":"BTCUSDT","U":49981777515,"u":49981777539,
	"b":[["57261.83000000","0.00096000"],["57265.01000000","5.23762000"],["57261.66000000","0.00091000"],
	["57251.30000000","0.03018000"],["57249.66000000","0.04961000"],["57249.36000000","0.00000000"],
	["57248.79000000","0.00000000"],["57245.17000000","0.43654000"],["57197.85000000","8.11363000"],
	["57196.99000000","0.00000000"],["57185.39000000","0.00000000"],["57165.01000000","0.01000000"],
	["57158.97000000","0.00000000"],["33000.00000000","20.95952000"]],
	"a":[["57265.03000000","0.03705000"],["57266.71000000","0.00010000"],["57279.41000000","0.43654000"],
	["57282.15000000","0.74283000"],["57282.45000000","0.11479000"],["57284.50000000","0.13960000"],
	["57291.80000000","0.03024000"]]}
	`

	bids := map[string]string{
		"57249.36000000": "1.00000000",
		"57248.79000000": "1.00000000",
		"57196.99000000": "1.00000000",
		"57185.39000000": "1.00000000",
		"57158.97000000": "1.00000000",
	}
	asks := map[string]string{}

	var p fastjson.Parser

	data, _ := p.Parse(rawData)
	data_stream.UpdateCurrentOrder(data.GetArray("b"), bids)
	data_stream.UpdateCurrentOrder(data.GetArray("a"), asks)

	assert.Equal(t, "0.00096000", bids["57261.83000000"])
	assert.Equal(t, "5.23762000", bids["57265.01000000"])
	assert.Equal(t, "0.00091000", bids["57261.66000000"])
	assert.Equal(t, "0.03018000", bids["57251.30000000"])
	assert.Equal(t, "0.04961000", bids["57249.66000000"])
	assert.Equal(t, "", bids["57249.36000000"])
	assert.Equal(t, "", bids["57248.79000000"])
	assert.Equal(t, "0.43654000", bids["57245.17000000"])
	assert.Equal(t, "8.11363000", bids["57197.85000000"])
	assert.Equal(t, "", bids["57196.99000000"])
	assert.Equal(t, "", bids["57185.39000000"])
	assert.Equal(t, "0.01000000", bids["57165.01000000"])
	assert.Equal(t, "", bids["57158.97000000"])
	assert.Equal(t, "20.95952000", bids["33000.00000000"])

	assert.Equal(t, "0.03705000", asks["57265.03000000"])
	assert.Equal(t, "0.00010000", asks["57266.71000000"])
	assert.Equal(t, "0.43654000", asks["57279.41000000"])
	assert.Equal(t, "0.74283000", asks["57282.15000000"])
	assert.Equal(t, "0.11479000", asks["57282.45000000"])
	assert.Equal(t, "0.13960000", asks["57284.50000000"])
	assert.Equal(t, "0.03024000", asks["57291.80000000"])

}

func TestOrderBook2Json(t *testing.T) {
	o := format.OrderBookStream{}
	o.Bids = map[string]string{"57261.83000000": "0.00096000", "57265.01000000": "5.23762000"}
	o.Asks = map[string]string{"57265.03000000": "0.03705000", "57266.71000000": "0.00010000"}
	o.Time = 1723025111169
	o.Symbol = "BTCUSDT"
	o.Topic = "depthUpdate"

	b, _ := json.Marshal(o)

	fmt.Println(string(b))
}

func TestPartialOrderBookUpdate(t *testing.T) {
	rawData := `{"lastUpdateId":49981777515,"bids":[["57261.83000000","0.00096000"],["57265.01000000","5.23762000"]],"asks":[["57265.03000000","0.03705000"],["57266.71000000","0.00010000"]]}`
	parse := data_stream.NewPartialOrderBook(2)
	ob, _ := parse.Update([]byte(rawData))

	assert.Equal(t, ob.Bids["57261.83000000"], "0.00096000")
	assert.Equal(t, ob.Bids["57265.01000000"], "5.23762000")
	assert.Equal(t, ob.Asks["57265.03000000"], "0.03705000")
	assert.Equal(t, ob.Asks["57266.71000000"], "0.00010000")

}


func TestOrderBookUpdate(t *testing.T) {
	snapShot := []byte(`{"lastUpdateId":49981777514,"bids":[["57261.83000000","0.00096000"],["57265.01000000","5.23762000"],["57261.66000000","0.00091000"],["57251.30000000","0.03018000"],["57249.66000000","0.04961000"],["57249.36000000","0.00000000"],["57248.79000000","0.00000000"],["57245.17000000","0.43654000"],["57197.85000000","8.11363000"],["57196.99000000","0.00000000"],["57185.39000000","0.00000000"],["57165.01000000","0.01000000"],["57158.97000000","0.00000000"],["33000.00000000","20.95952000"]],"asks":[["57265.03000000","0.03705000"],["57266.71000000","0.00010000"],["57279.41000000","0.43654000"],["57282.15000000","0.74283000"],["57282.45000000","0.11479000"],["57284.50000000","0.13960000"],["57291.80000000","0.03024000"]]}`)
	parse := data_stream.NewOrderBookParser(2)
	update := []byte(`{"e":"depthUpdate","E":1723025111269,"s":"BTCUSDT","U":49981777515,"u":49981777539,"b":[["57261.83000000","0.00096000"],["57265.01000000","0"],["57261.66000000","1.00091000"],["57251.30000000","1.03018000"],["57250.66000000","0.04961000"],["57249.36000000","0.00000000"]],"a":[["57265.03000000","0.03705000"],["57266.71000000","0.00000000"],["57279.41000000","0.000"]]}`)
	update2 := []byte(`{"e":"depthUpdate","E":1723025111369,"s":"BTCUSDT","U":49981777540,"u":49981777541,"b":[],"a":[["57265.03000000","1.03705000"]]}`)
	parse.Update(update)
	parse.SetSnapshot(snapShot)
	parse.Update(update2)
	assert.Equal(t, parse.Bids["57261.83000000"], "0.00096000")
	assert.Equal(t, parse.Bids["57265.01000000"], "")
	assert.Equal(t, parse.Bids["57261.66000000"], "1.00091000")
	assert.Equal(t, parse.Bids["57251.30000000"], "0.03018000")
	assert.Equal(t, parse.Bids["57250.66000000"], "")
	assert.Equal(t, parse.Bids["57249.36000000"], "")

	assert.Equal(t, parse.Asks["57265.03000000"], "1.03705000")
	assert.Equal(t, parse.Asks["57266.71000000"], "0.00010000")
	assert.Equal(t, parse.Asks["57279.41000000"], "0.43654000")

}

func TestSort(t *testing.T)	{
	bids := [][]string{
		{"57261.83000000","0.00096000"},
		{"57265.01000000","5.23762000"},	
		{"57261.66000000","0.00091000"},
		{"57251.30000000","0.00000000"},
		{"57249.66000000","0.04961000"},
		{"57249.36000000","0.00000000"},
		{"57248.79000000","0.00000000"},
		{"57245.17000000","0.43654000"},
		{"57197.85000000","8.11363000"},
	}
	sort.Slice(bids, func(i, j int) bool {
		if price, _ := strconv.ParseFloat(bids[j][1], 64); price == 0 {
			return true
		}
		return bids[i][0] > bids[j][0]
	})
	for _, bid := range bids[:5] {
		fmt.Println(bid)
	}

	fmt.Println("==========")

	asks := [][]string{
		{"57265.03000000","0.03705000"},
		{"57266.71000000","0.00010000"},
		{"57279.41000000","0.43654000"},
		{"57282.15000000","0.74283000"},
		{"57282.45000000","0.00000000"},
		{"57284.50000000","0.13960000"},
		{"57291.80000000","0.00000000"},
		{"57291.90000000","0.03024000"},
	}
	sort.Slice(asks, func(i, j int) bool {
		if price, _ := strconv.ParseFloat(asks[j][1], 64); price == 0 {
			return true
		}
		return asks[i][0] < asks[j][0]
	})
	for _, ask := range asks[:5] {
		fmt.Println(ask)
	}
}


func TestOrderBookUpdateRealtime(t *testing.T) {
	logger := logrus.New()

	updater, _ := data_stream.NewOrderBookMap("spot")
	handle := func(data []byte) {
		// logger.Infof(`%s`, string(data))
		ob, err := updater.Update(data)
		if err != nil {
			logger.Error(err)
			return
		}
		if ob == nil {
			return
		}
		logger.Infof("Bids: %v", ob.Bids)
		logger.Infof("Asks: %v", ob.Asks)
	}
	client, _ := ws.NewWsQuoteClient(consts.Spot, handle)
	client.Logger = logger
	client.Logger.SetLevel(logrus.DebugLevel)

	resp, err := client.Connect()
	if err != nil {
		logger.Println(resp)
		t.Error(err)
		return
	}

	go func() {
		for range client.StartSignal {
			client.Subscribe([]string{"btcusdt@depth@100ms"})
		}
	}()

	client.StartLoop()


	time.Sleep(30 * time.Second)
	client.Stop()


	<-client.StopSignal
}