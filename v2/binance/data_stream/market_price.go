package data_stream

import (
	"github.com/lianyun0502/exchange_conn/v2/data_format"
	"github.com/valyala/fastjson"
)

func UpdateMarketPrice(rawData []byte) (data *format.MarKetPriceStream, err error) {
	v := fastjson.MustParseBytes(rawData)
	data = &format.MarKetPriceStream{
		Topic:           string(v.GetStringBytes("e")),
		Time:            v.GetInt64("E"),
		Symbol:          string(v.GetStringBytes("s")),
		MarketPrice:     string(v.GetStringBytes("p")),
		IndexPrice:      string(v.GetStringBytes("i")),
		FundingRate:     string(v.GetStringBytes("r")),
		NextFundingTime: v.GetInt64("T"),
	}
	return data, nil
}
