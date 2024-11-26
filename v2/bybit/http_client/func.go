package bybit

import (
	"net/http"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/lianyun0502/exchange_conn/v2/http_client"
)

func WithQuery(query map[string]string) func(queryMap map[string]string) {
	return func(queryMap map[string]string) {
		for k, v := range query {
			queryMap[k] = v
		}
	}
}
/*
https://bybit-exchange.github.io/docs/zh-TW/v5/spot-margin-uta/historical-interest
*/
func (api *ByBitClient) MarginTrade_InterestRateHistory(currency string, opts ...func(map[string]string)) (*InterestHistoryResponce, error) {
	req := api.Request(http.MethodGet, "/v5/spot-margin-trade/interest-rate-history", SetSercurityType(true, true))
	query := httpClient.QueryMap{
		"currency": currency,
	}
	for _, opt := range opts {
		opt(query)
	}
	resp, err := req.Send()
	if err != nil {
		api.Log.Error(err)
		return nil, err
	}
	ret := new(Response[InterestHistoryResponce])
	json.Unmarshal(resp, ret)
	if ret.RetCode != 0 {
		api.Log.WithFields(logrus.Fields{
			"retCode": ret.RetCode,
			"retMsg":  ret.RetMsg,
		}).Warning("MarginTrade_InterestRateHistory")
		return nil, errors.New(ret.RetMsg)
	}
	return &ret.Results, nil
}

/*https://bybit-exchange.github.io/docs/zh-TW/v5/spot-margin-uta/vip-margin*/
func (api *ByBitClient) MarginTrade_Data(opts ...func(map[string]string)) (*MarginTradeDataResponse, error) {
	req := api.Request(http.MethodGet, "/v5/spot-margin-trade/data")
	query := make(httpClient.QueryMap)
	for _, opt := range opts {
		opt(query)
	}
	req.SetQuery(query)
	resp, err := req.Send()
	if err != nil {
		api.Log.Error(err)
		return nil, err
	}
	ret := new(Response[MarginTradeDataResponse])
	json.Unmarshal(resp, ret)
	if ret.RetCode != 0 {
		api.Log.WithFields(logrus.Fields{
			"retCode": ret.RetCode,
			"retMsg":  ret.RetMsg,
		}).Warning("MarginTrade_Data")
		return nil, errors.New(ret.RetMsg)
	}
	return &ret.Results, nil
}

// https://bybit-exchange.github.io/docs/zh-TW/v5/market/instrument
func (api *ByBitClient) Market_InstrumentsInfo(category string, opts ...func(map[string]string)) (*InstrumentInfoResponse, error) {
	req := api.Request(http.MethodGet, "/v5/market/instruments-info")
	query := httpClient.QueryMap{
		"category": category,
	}
	for _, opt := range opts {
		opt(query)
	}
	req.SetQuery(query)
	resp, err := req.Send()
	if err != nil {
		api.Log.Error(err)
		return nil, err
	}
	ret := new(Response[InstrumentInfoResponse])
	json.Unmarshal(resp, ret)
	if ret.RetCode != 0 {
		api.Log.WithFields(logrus.Fields{
			"retCode": ret.RetCode,
			"retMsg":  ret.RetMsg,
		}).Warning("Market_InstrumentsInfo")
		return nil, errors.New(ret.RetMsg)
	}
	if ret.Results.NextPageCursor != "" {
		query["cursor"] = ret.Results.NextPageCursor
		req.SetQuery(query)
		resp, err := req.Send()
		if err != nil {
			api.Log.Error(err)
			return nil, err
		}
		r := new(Response[InstrumentInfoResponse])
		json.Unmarshal(resp, r)
		if r.RetCode != 0 {
			api.Log.WithFields(logrus.Fields{
				"retCode": r.RetCode,
				"retMsg":  r.RetMsg,
			}).Warning("Market_InstrumentsInfo")
			return nil, errors.New(r.RetMsg)
		}
		ret.Results.List = append(ret.Results.List, r.Results.List...)
	}
	return &ret.Results, nil
}

// https://bybit-exchange.github.io/docs/zh-TW/v5/account/borrow-history
func (api *ByBitClient) Account_BorrowHistory(opts ...func(map[string]string)) (*BorrowHistoryResponse, error) {
	req := api.Request(http.MethodGet, "/v5/account/borrow-history", SetSercurityType(true, true))
	query := make(httpClient.QueryMap)
	for _, opt := range opts {
		opt(query)
	}
	req.SetQuery(query)
	resp, err := req.Send()
	if err != nil {
		api.Log.Error(err)
		return nil, err
	}
	ret := new(Response[BorrowHistoryResponse])
	json.Unmarshal(resp, ret)
	if ret.RetCode != 0 {
		api.Log.WithFields(logrus.Fields{
			"retCode": ret.RetCode,
			"retMsg":  ret.RetMsg,
		}).Warning("Account_BorrowHistory")
		return nil, errors.New(ret.RetMsg)
	}
	return &ret.Results, nil
}
// https://bybit-exchange.github.io/docs/zh-TW/v5/market/premium-index-kline
func (api *ByBitClient) Market_PremiumIndexPrice(symbol, interval string, opts ...func(map[string]string)) (*PremiumIndexPriceResponse, error) {
	req := api.Request(http.MethodGet, "/v5/market/premium-index-price-kline")
	query := httpClient.QueryMap{
		"symbol": symbol,
		"category": "linear",
		"interval": interval,
	}
	req.SetQuery(query)
	resp, err := req.Send()
	if err != nil {
		api.Log.Error(err)
		return nil, err
	}
	ret := new(Response[PremiumIndexPriceResponse])
	json.Unmarshal(resp, ret)
	if ret.RetCode != 0 {
		api.Log.WithFields(logrus.Fields{
			"retCode": ret.RetCode,
			"retMsg":  ret.RetMsg,
		}).Warning("Market_PremiumIndexPrice")
		return nil, errors.New(ret.RetMsg)
	}
	return &ret.Results, nil
}