package data_stream

import (
	"fmt"
	"strconv"

	"github.com/lianyun0502/exchange_conn/v2/data_format"
	"github.com/valyala/fastjson"
)

func InitMarketPrice(v *fastjson.Value) (data *format.MarKetPriceStream, err error) {
	var nextFundingTime int64
	if v.GetStringBytes("data", "nextFundingTime") != nil {
		nextFundingTime, err = strconv.ParseInt(string(v.GetStringBytes("data", "nextFundingTime")), 10, 64)
		if err != nil {
			return nil, err
		}
	}
	data = &format.MarKetPriceStream{
		Topic:           string(v.GetStringBytes("topic")),
		Time:            v.GetInt64("ts"),
		LastPrice:       string(v.GetStringBytes("data", "lastPrice")),
		Symbol:          string(v.GetStringBytes("data", "symbol")),
		MarketPrice:     string(v.GetStringBytes("data", "markPrice")),
		IndexPrice:      string(v.GetStringBytes("data", "indexPrice")),
		FundingRate:     string(v.GetStringBytes("data", "fundingRate")),
		NextFundingTime: nextFundingTime,
	}
	return data, nil
}

func UpdateMarketPrice(v *fastjson.Value, m map[string]*format.MarKetPriceStream) (data *format.MarKetPriceStream, err error) {
	var d *format.MarKetPriceStream
	symbol := v.GetStringBytes("data", "symbol")
	if symbol != nil {
		d = m[string(symbol)]
		d.Symbol = string(symbol)
	}else{
		return nil, fmt.Errorf("symbol is nil")
	}
	nft := v.GetStringBytes("data", "nextFundingTime")
	if nft != nil {
		nextFundingTime, err := strconv.ParseInt(string(nft), 10, 64)
		if err != nil {
			return nil, err
		}
		d.NextFundingTime = nextFundingTime
	}

	topic := v.GetStringBytes("topic")
	if topic != nil {
		d.Topic = string(topic)
	}
	ts := v.GetInt64("ts")
	if ts != 0 {
		d.Time = ts
	}
	lastPrice := v.GetStringBytes("data", "lastPrice")
	if lastPrice != nil {
		d.LastPrice = string(lastPrice)
	}
	markPrice := v.GetStringBytes("data", "markPrice")
	if markPrice != nil {
		d.MarketPrice = string(markPrice)
	}
	indexPrice := v.GetStringBytes("data", "indexPrice")
	if indexPrice != nil {
		d.IndexPrice = string(indexPrice)
	}
	fundingRate := v.GetStringBytes("data", "fundingRate")
	if fundingRate != nil {
		d.FundingRate = string(fundingRate)
	}
	return d, nil
}

type MarketData struct {
	Data map[string] *format.MarKetPriceStream
}

func NewMarketData() *MarketData {
	return &MarketData{Data: make(map[string] *format.MarKetPriceStream)}

}

func (md *MarketData) Update(rawData []byte) (*format.MarKetPriceStream, error) {
	v := fastjson.MustParseBytes(rawData)
	switch string(v.GetStringBytes("type")) {
	case "snapshot":
		ret, err := InitMarketPrice(v)
		md.Data[ret.Symbol] = ret
		return ret, err
	case "delta":
		if md.Data == nil {
			return nil, nil
		}
		return UpdateMarketPrice(v, md.Data)
	}
	return nil, nil
}
