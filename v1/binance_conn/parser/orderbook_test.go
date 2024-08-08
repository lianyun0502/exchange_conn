package parser_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fastjson"

	// "github.com/stretchr/testify/assert"
	"github.com/lianyun0502/exchange_conn/v1/binance_conn/parser"
)




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
		quality:= string(bid.GetStringBytes("1"))
		fmt.Println(price)
		m[price] = quality
	}
	for k, v := range m {
		fmt.Println(k,",", v)
	}
	fmt.Println("==========")
	for _, ask := range asks {
		price:= ask.GetArray()[0].String()
		fmt.Println(price)
	}
}

func TestUpdateCurrentOrder(t *testing.T){
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
	parser.UpdateCurrentOrder(data.GetArray("b"), bids)
	parser.UpdateCurrentOrder(data.GetArray("a"), asks)

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


func TestOrderBook2Json(t *testing.T){
	o := parser.OrderBook{}
	o.Bids = map[string]string{ "57261.83000000": "0.00096000", "57265.01000000": "5.23762000", }
	o.Asks = map[string]string{ "57265.03000000": "0.03705000", "57266.71000000": "0.00010000", }
	o.Time = 1723025111169
	o.Symbol = "BTCUSDT"
	o.Topic = "depthUpdate"
	
	b, _ := json.Marshal(o)

	fmt.Println(string(b))
}