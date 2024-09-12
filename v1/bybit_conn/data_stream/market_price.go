package data_stream

import (
	"strconv"

	"github.com/lianyun0502/exchange_conn/v1"
	"github.com/valyala/fastjson"
)

func InitMarketPrice(v *fastjson.Value) (data *exchange_conn.MarKetPriceStream, err error) {
	nft := v.GetStringBytes("data", "nextFundingTime")
	var nextFundingTime int64
	if nft != nil{
		nextFundingTime, err = strconv.ParseInt(string(nft), 10, 64)
		if err != nil {
			return nil, err
		}
	}
	data = &exchange_conn.MarKetPriceStream{
		Topic:           string(v.GetStringBytes("topic")),
		Time:            v.GetInt64("ts"),
		Symbol:          string(v.GetStringBytes("data", "symbol")),
		MarketPrice:     string(v.GetStringBytes("data", "markPrice")),
		IndexPrice:      string(v.GetStringBytes("data", "indexPrice")),
		FundingRate:     string(v.GetStringBytes("data", "fundingRate")),
		NextFundingTime: nextFundingTime,
	}
	return data, nil
}

func UpdateMarketPrice(v *fastjson.Value, d *exchange_conn.MarKetPriceStream) (data *exchange_conn.MarKetPriceStream, err error) {
	nft := v.GetStringBytes("data", "nextFundingTime")
	if nft != nil{
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
	symbol := v.GetStringBytes("data", "symbol")
	if symbol != nil {
		d.Symbol = string(symbol)
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

type MarketData struct{
	Data *exchange_conn.MarKetPriceStream
}

func NewMarketData() *MarketData {
	return new(MarketData)

}

func (md *MarketData) Update(rawData []byte) (*exchange_conn.MarKetPriceStream, error) {
	v := fastjson.MustParseBytes(rawData)
	switch string(v.GetStringBytes("type")){
	case "snapshot":
		ret, err := InitMarketPrice(v)
		md.Data = ret
		return ret, err
	case "delta":
		if md.Data == nil {
			return nil, nil
		}
		return UpdateMarketPrice(v, md.Data)
	}
	return nil, nil
}
