package data_stream

import (
	"fmt"
	"regexp"
	"github.com/valyala/fastjson"
)


type IDataParser interface {
	Update(rawdata []byte) (any, error)
}


type DataParser struct {
	OrderBook *OrderBook
	Trade *Trade
}

func NewDataParser() *DataParser {
	return &DataParser{
		OrderBook: NewOrderBook(),
		Trade: NewTrade(),
	}
}

func ByBitSymbolToTopic(symbol string) string {
	switch {
	case regexp.MustCompile("publicTrade").MatchString(symbol):
		return "publicTrade"
	case regexp.MustCompile("orderbook").MatchString(symbol):
		return "orderbook"
	case regexp.MustCompile("tickers").MatchString(symbol):
		return "tickers"
	default:
		return ""
	}
}

func (dp *DataParser) Parse(rawData []byte) (any, error) {
	v := fastjson.MustParseBytes(rawData)
	symbol := string(v.GetStringBytes("topic"))
	topic := ByBitSymbolToTopic(symbol)
	switch topic{
	case "publicTrade":
		return dp.Trade.Update(rawData)
	case "orderbook":
		return dp.OrderBook.Update(rawData)
	// case "tickers":
	// 	return data_stream.MarketPrice()(rawData)
	default:
		return nil, fmt.Errorf("topic %s not found", topic)
	} 
}
