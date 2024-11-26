package bybit

type Response[Res any] struct {
	RetCode int    `json:"retCode"`
	RetMsg  string `json:"retMsg"`
	Results Res    `json:"result"`
	Data    Res    `json:"data"`
}

type InterestHistoryResponce struct {
	Histories []InterestHistory `json:"list"`
}

type InterestHistory struct {
	Time             int64  `json:"timestamp"`
	HourlyBorrowRate string `json:"hourlyBorrowRate"`
}

type MarginTradeDataResponse struct {
	VipCoinList []*VipCoinList `json:"vipCoinList"`
}
type VipCoinList struct {
	VipLevel string     `json:"vipLevel"`
	List     []*VipCoin `json:"list"`
}

type VipCoin struct {
	Borrowable       bool   `json:"borrowable"`
	CollateralRatio  string `json:"collateralRatio"`
	Currency         string `json:"currency"`
	HourlyBorrowRate string `json:"hourlyBorrowRate"`
	MarginCollateral bool   `json:"marginCollateral"`
	MaxBorrowAmount  string `json:"maxBorrowingAmount"`
}

type InstrumentInfoResponse struct {
	Category       string            `json:"category"`
	NextPageCursor string            `json:"nextPageCursor"`
	List           []*InstrumentInfo `json:"list"`
}
type InstrumentInfo struct {
	Symbol           string        `json:"symbol"`
	BaseCoin         string        `json:"baseCoin"`
	QuoteCoin        string        `json:"quoteCoin"`
	Status           string        `json:"status"`
	ContractType     string        `json:"contractType,omitempty"`
	MarginTrade      string        `json:"marginTrading,omitempty"`
	LotSizeFilter    LotSizeFilter `json:"lotSizeFilter,omitempty"`
	PriceFilter      PriceFilter   `json:"priceFilter,omitempty"`
	UpperFundingRate string        `json:"upperFundingRate,omitempty"`
	LowerFundingRate string        `json:"lowerFundingRate,omitempty"`
}

type PriceFilter struct {
	TickSize string `json:"tickSize,omitempty"`
}
type LotSizeFilter struct {
	BasePrecision    string `json:"basePrecision,omitempty"`
	QuotePrecision   string `json:"quotePrecision,omitempty"`
	MinOrderQty      string `json:"minOrderQty,omitempty"`
	MaxOrderQty      string `json:"maxOrderQty,omitempty"`
	MinOrderAmt      string `json:"minOrderAmt,omitempty"`
	MaxOrderAmt      string `json:"maxOrderAmt,omitempty"`
	MaxMktOrderAmt   string `json:"maxMktOrderAmt,omitempty"`
	MaxMktOrderQty   string `json:"maxMktOrderQty,omitempty"`
	MinNotionalValue string `json:"minNotionalValue,omitempty"`
	QtyStep          string `json:"qtyStep,omitempty"`
	LimitParam       string `json:"limitParameter,omitempty"`
	MarketParam      string `json:"marketParameter,omitempty"`
}

type BorrowHistoryResponse struct {
	NextPageCursor string           `json:"nextPageCursor,omitempty"`
	History        []*BorrowHistory `json:"list,omitempty"`
}

type BorrowHistory struct {
	BorrowAmount string `json:"borrowAmount,omitempty"`
	CreatTime    int    `json:"createdTime,omitempty"`
	Currency     string `json:"currency,omitempty"`
	Size         string `json:"InterestBearingBorrowSize,omitempty"`
	BR           string `json:"hourlyBorrowRate,omitempty"`
	Cost         string `json:"borrowCost,omitempty"`
}

type PremiumIndexPriceResponse struct {
	Symbol      string     `json:"symbol"`
	Category    string     `json:"category"`
	IndexPrices [][]string `json:"list"`
}

type IndexPrice struct {
	StartTime   string
	StrartPrice string
	HighPrice   string
	LowPrice    string
	EndPrice    string
}
