package data_stream

import (
	"encoding/json"
	"fmt"
	"strconv"
	// "testing/quick"

	"github.com/lianyun0502/exchange_conn/v2/data_format"
	// "github.com/valyala/fastjson"
)

func InitMarketPrice(quote *Quote[Tickers]) (data *format.MarKetPriceStream, err error) {
	var nextFundingTime int64
	if quote.Data.NextFundingTime != "" {
		nextFundingTime, err = strconv.ParseInt(quote.Data.NextFundingTime, 10, 64)
		if err != nil {
			return nil, err
		}
	}
	data = &format.MarKetPriceStream{
		Topic:           quote.Topic,
		Time:            quote.Time,
		LastPrice:       quote.Data.LastPrice,
		Symbol:          quote.Data.Symbol,
		MarketPrice:     quote.Data.MarketPrice,
		IndexPrice:      quote.Data.IndexPrice,
		FundingRate:     quote.Data.FundingRate,
		NextFundingTime: nextFundingTime,
	}
	return data, nil
}

func UpdateMarketPrice(quote *Quote[Tickers], m map[string]*format.MarKetPriceStream) (data *format.MarKetPriceStream, err error) {
	if d, ok := m[quote.Data.Symbol]; ok{
		data = d
		data.Symbol = quote.Data.Symbol
	}else{
		return nil, fmt.Errorf("symbol is nil")
	}
	if quote.Data.NextFundingTime != "" {
		nextFundingTime, err := strconv.ParseInt(quote.Data.NextFundingTime, 10, 64)
		if err != nil {
			return nil, err
		}
		data.NextFundingTime = nextFundingTime
	}

	if quote.Data.LastPrice != "" {
		data.LastPrice = quote.Data.LastPrice
	}
	if quote.Data.MarketPrice != "" {
		data.MarketPrice = quote.Data.MarketPrice
	}
	if quote.Data.IndexPrice != "" {
		data.IndexPrice = quote.Data.IndexPrice
	}
	if quote.Data.FundingRate != "" {
		data.FundingRate = quote.Data.FundingRate
	}
	data.Topic = quote.Topic
	data.Time = quote.Time
	return data, nil
}

type MarketData struct {
	Data map[string] *format.MarKetPriceStream
}

func NewMarketData() *MarketData {
	return &MarketData{Data: make(map[string] *format.MarKetPriceStream)}

}

func (md *MarketData) Update(rawData []byte) (*format.MarKetPriceStream, error) {
	tickers := new(Quote[Tickers])
	json.Unmarshal(rawData, tickers)
	// v := fastjson.MustParseBytes(rawData)
	switch string(tickers.Type) {
	case "snapshot":
		ret, err := InitMarketPrice(tickers)
		md.Data[ret.Symbol] = ret
		return ret, err
	case "delta":
		if md.Data == nil {
			return nil, nil
		}
		return UpdateMarketPrice(tickers, md.Data)
	}
	return nil, nil
}

type Quote [T any] struct {
	Topic string `json:"topic"`
	Type string `json:"type"`
	Data T `json:"data"`
	Time int64 `json:"ts"`
}

type Tickers struct {
	Symbol string `json:"symbol"`
	MarketPrice string `json:"markPrice"`
	IndexPrice string `json:"indexPrice"`
	FundingRate string `json:"fundingRate"`
	LastPrice string `json:"lastPrice"`
	NextFundingTime string `json:"nextFundingTime"`
}