package data_stream

import (
	"strconv"
	"strings"

	"github.com/lianyun0502/exchange_conn/v1"
	"github.com/valyala/fastjson"
)

func InitMarketPrice(v *fastjson.Value) (data *exchange_conn.MarKetPriceStream, err error) {
	nextFundingTime, err := strconv.ParseInt(string(v.GetStringBytes("data", "nextFundingTime")), 10, 64)
	if err != nil {
		return nil, err
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

func MarketPrice() func([]byte) (*exchange_conn.MarKetPriceStream, error) {
	data := new(exchange_conn.MarKetPriceStream)
	return func (rawData []byte) (*exchange_conn.MarKetPriceStream, error) {
		v := fastjson.MustParseBytes(rawData)
		topic := string(v.GetStringBytes("topic"))
		if strings.Split(topic, ".")[0] != "tickers" {
			return nil, nil
		}
		switch string(v.GetStringBytes("type")){
		// case "snapshot":
		// 	ret, err := InitMarketPrice(v)
		// 	data = ret
		// 	return ret, err
		case "snapshot":
			if data == nil {
				return nil, nil
			}
			return UpdateMarketPrice(v, data)
		}
		return nil, nil
	}
}

