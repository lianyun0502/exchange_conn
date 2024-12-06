package bybit

import (
	"encoding/json"
	"errors"
	"net/http"
	."github.com/lianyun0502/exchange_conn/v2/bybit/response"
	"github.com/lianyun0502/exchange_conn/v2/http_client"
	"github.com/sirupsen/logrus"
)

func WithQuery(query map[string]string) func(queryMap map[string]string) {
	return func(queryMap map[string]string) {
		for k, v := range query {
			queryMap[k] = v
		}
	}
}

func (api *ByBitClient) InsLoan_RepaidHistory(opts ...func(map[string]string)) ([]RepayInfo, error) {
	req := api.Request(http.MethodGet, "/v5/ins-loan/repaid-history", SetSercurityType(true, true))
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
	ret := new(Response[*ReplayInfoResponse])
	json.Unmarshal(resp, ret)
	if ret.RetCode != 0 {
		api.Log.WithFields(logrus.Fields{
			"retCode": ret.RetCode,
			"retMsg":  ret.RetMsg,
		}).Warning("InsLoan_RepaidHistory")
		return nil, errors.New(ret.RetMsg)
	}
	return ret.Results.RepayInfo, nil
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
	req.SetQuery(query)
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

// https://bybit-exchange.github.io/docs/zh-TW/v5/spot-margin-uta/status
func (api *ByBitClient) MarginTrade_State() (*MarginTrade, error) {
	req := api.Request(http.MethodGet, "/v5/spot-margin-trade/state", SetSercurityType(true, true))
	resp, err := req.Send()
	if err != nil {
		api.Log.Error(err)
		return nil, err
	}
	ret := new(Response[*MarginTrade])
	json.Unmarshal(resp, ret)
	if ret.RetCode != 0 {
		api.Log.WithFields(logrus.Fields{
			"retCode": ret.RetCode,
			"retMsg":  ret.RetMsg,
		}).Warning("MarginTrade_State")
		return nil, errors.New(ret.RetMsg)
	}
	return ret.Results, nil
}

/*https://bybit-exchange.github.io/docs/zh-TW/v5/spot-margin-uta/vip-margin*/
func (api *ByBitClient) MarginTrade_Data(opts ...func(map[string]string)) ([]VipCoinList, error) {
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
	return ret.Results.VipCoinList, nil
}

// https://bybit-exchange.github.io/docs/zh-TW/v5/market/tickers
func (api *ByBitClient) Market_Tickers(category string, opts ...func(map[string]string))([]Ticker, error){
	req := api.Request(http.MethodGet, "/v5/market/tickers")
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
	ret := new(Response[*TickerResponse])
	json.Unmarshal(resp, ret)
	if ret.RetCode != 0 {
		api.Log.WithFields(logrus.Fields{
			"retCode": ret.RetCode,
			"retMsg":  ret.RetMsg,
		}).Warning("Market_Tickers")
		return nil, errors.New(ret.RetMsg)
	}
	return ret.Results.List, nil
}

// https://bybit-exchange.github.io/docs/zh-TW/v5/market/instrument
func (api *ByBitClient) Market_InstrumentsInfo(category string, opts ...func(map[string]string)) ([]InstrumentInfo, error) {
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
	ret := new(Response[*InstrumentInfoResponse])
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
		r := new(Response[*InstrumentInfoResponse])
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
	return ret.Results.List, nil
}

// https://bybit-exchange.github.io/docs/zh-TW/v5/market/premium-index-kline
func (api *ByBitClient) Market_PremiumIndexPrice(symbol, interval string, opts ...func(map[string]string)) (*PremiumIndexPriceResponse, error) {
	req := api.Request(http.MethodGet, "/v5/market/premium-index-price-kline")
	query := httpClient.QueryMap{
		"symbol":   symbol,
		"category": "linear",
		"interval": interval,
	}
	req.SetQuery(query)
	resp, err := req.Send()
	if err != nil {
		api.Log.Error(err)
		return nil, err
	}
	ret := new(Response[*PremiumIndexPriceResponse])
	json.Unmarshal(resp, ret)
	if ret.RetCode != 0 {
		api.Log.WithFields(logrus.Fields{
			"retCode": ret.RetCode,
			"retMsg":  ret.RetMsg,
		}).Warning("Market_PremiumIndexPrice")
		return nil, errors.New(ret.RetMsg)
	}
	return ret.Results, nil
}
// https://bybit-exchange.github.io/docs/zh-TW/v5/market/risk-limit
func (api *ByBitClient) Market_RiskLimit(category string, opts ...func(map[string]string)) ([]RiskLimit, error) {
	req := api.Request(http.MethodGet, "/v5/market/risk-limit", SetSercurityType(true, true))
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
	ret := new(Response[*RiskLimitResponse])
	json.Unmarshal(resp, ret)
	if ret.RetCode != 0 {
		api.Log.WithFields(logrus.Fields{
			"retCode": ret.RetCode,
			"retMsg":  ret.RetMsg,
		}).Warning("Market_RiskLimit")
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
		r := new(Response[*RiskLimitResponse])
		json.Unmarshal(resp, r)
		if r.RetCode != 0 {
			api.Log.WithFields(logrus.Fields{
				"retCode": r.RetCode,
				"retMsg":  r.RetMsg,
			}).Warning("Market_InstrumentsInfo")
			return nil, errors.New(r.RetMsg)
		}
		ret.Results.RiskLimit = append(ret.Results.RiskLimit, r.Results.RiskLimit...)
	}
	return ret.Results.RiskLimit, nil
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
	ret := new(Response[*BorrowHistoryResponse])
	json.Unmarshal(resp, ret)
	if ret.RetCode != 0 {
		api.Log.WithFields(logrus.Fields{
			"retCode": ret.RetCode,
			"retMsg":  ret.RetMsg,
		}).Warning("Account_BorrowHistory")
		return nil, errors.New(ret.RetMsg)
	}
	return ret.Results, nil
}

// https://bybit-exchange.github.io/docs/zh-TW/v5/account/wallet-balance
func (api *ByBitClient) Account_WalletBalance(accountType string, opts ...func(map[string]string)) (*WalletResponse, error) {
	req := api.Request(http.MethodGet, "/v5/account/wallet-balance", SetSercurityType(true, true))
	query := httpClient.QueryMap{
		"accountType": accountType,
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
	ret := new(Response[*WalletResponse])
	json.Unmarshal(resp, ret)
	if ret.RetCode != 0 {
		api.Log.WithFields(logrus.Fields{
			"retCode": ret.RetCode,
			"retMsg":  ret.RetMsg,
		}).Warning("Account_WalletBalance")
		return nil, errors.New(ret.RetMsg)
	}
	return ret.Results, nil
}
// https://bybit-exchange.github.io/docs/zh-TW/v5/account/transaction-log
func (api *ByBitClient) Account_TransactionLog(category string, opts ...func(map[string]string)) ([]Transactions, error) {
	req := api.Request(http.MethodGet, "/v5/account/transaction-log", SetSercurityType(true, true))
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
	ret := new(Response[*TransactionResponse])
	json.Unmarshal(resp, ret)
	if ret.RetCode != 0 {
		api.Log.WithFields(logrus.Fields{
			"retCode": ret.RetCode,
			"retMsg":  ret.RetMsg,
		}).Warning("Account_TransactionLog")
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
		r := new(Response[*TransactionResponse])
		json.Unmarshal(resp, r)
		if r.RetCode != 0 {
			api.Log.WithFields(logrus.Fields{
				"retCode": r.RetCode,
				"retMsg":  r.RetMsg,
			}).Warning("Market_InstrumentsInfo")
			return nil, errors.New(r.RetMsg)
		}
		ret.Results.Transactions = append(ret.Results.Transactions, r.Results.Transactions...)
	}
	return ret.Results.Transactions, nil
}

// https://bybit-exchange.github.io/docs/zh-TW/v5/position
func (api *ByBitClient) Position_List (category string, opts ...func(map[string]string)) (*PositionResponse, error) {
	req := api.Request(http.MethodGet, "/v5/position/list", SetSercurityType(true, true))
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
	ret := new(Response[*PositionResponse])
	json.Unmarshal(resp, ret)
	if ret.RetCode != 0 {
		api.Log.WithFields(logrus.Fields{
			"retCode": ret.RetCode,
			"retMsg":  ret.RetMsg,
		}).Warning("Position_List")
		return nil, errors.New(ret.RetMsg)
	}
	return ret.Results, nil
}

// https://bybit-exchange.github.io/docs/zh-TW/v5/order/execution
func (api *ByBitClient) Trade_ExecutionList(category string, opts ...func(map[string]string)) ([]*Execution, error) {
	req := api.Request(http.MethodGet, "/v5/execution/list", SetSercurityType(true, true))
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
	ret := new(Response[ExecutionResponse])
	json.Unmarshal(resp, ret)
	if ret.RetCode != 0 {
		api.Log.WithFields(logrus.Fields{
			"retCode": ret.RetCode,
			"retMsg":  ret.RetMsg,
		}).Warning("Trade_ExecutionList")
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
		r := new(Response[ExecutionResponse])
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
	return ret.Results.List, nil
}