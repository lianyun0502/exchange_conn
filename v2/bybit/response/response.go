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
	VipCoinList []VipCoinList `json:"vipCoinList"`
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
	Category       string           `json:"category"`
	NextPageCursor string           `json:"nextPageCursor"`
	List           []InstrumentInfo `json:"list"`
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

type WalletResponse struct {
	List []*Wallet `json:"list"`
}

type Wallet struct {
	AccountType            string  `json:"accountType"`
	AccountIMRate          string  `json:"accountIMRate"`
	AccountMMRate          string  `json:"accountMMRate"`
	TotalEquity            string  `json:"totalEquity"`
	TotalBalance           string  `json:"totalWalletBalance"`
	TotalMarginBalance     string  `json:"totalMarginBalance"`
	TotalAvailableBalance  string  `json:"totalAvailableBalance"`
	TotalPerpUPL           string  `json:"totalPerpUPL"`
	TotalInitialMargin     string  `json:"totalInitialMargin"`
	TotalMaintenanceMargin string  `json:"totalMaintenanceMargin"`
	Coins                  []*Coin `json:"coin"`
}

type Coin struct {
	Coin            string `json:"coin"`
	Equity          string `json:"equity"`
	USDValue        string `json:"usdValue"`
	Walletbalance   string `json:"walletbalance"`
	BorrowAmount    string `json:"borrowAmount"`
	TotalPositionIM string `json:"totalPositionIM"`
	TotalPositionMM string `json:"totalPositionMM"`
	AccruedInterest string `json:"accruedInterest"`
	Free            string `json:"free"`
}

type MarginTrade struct {
	Leverage          string `json:"spotLeverage"`
	MarginMode        string `json:"spotMarginMode"`
	EffectiveLeverage string `json:"effectiveLeverage"`
}

type PositionResponse struct {
	Category  string      `json:"category,omitempty"`
	Positions []*Position `json:"list,omitempty"`
}
type Position struct {
	Symbol        string `json:"symbol,omitempty"`
	Side          string `json:"side,omitempty"`
	Size          string `json:"size,omitempty"`
	TradeMode     int    `json:"tradeMode,omitempty"`
	PositionValue string `json:"positionValue,omitempty"`
	Leverage      string `json:"leverage,omitempty"`
	PositionIM    string `json:"positionIM,omitempty"`
	PositionMM    string `json:"positionMM,omitempty"`
	CreateTime    string `json:"createdTime,omitempty"`
	UpdateTime    string `json:"updatedTime,omitempty"`
	SeqNum        int    `json:"seq,omitempty"`
}

type TransactionResponse struct {
	Transactions   []Transactions `json:"list,omitempty"`
	NextPageCursor string         `json:"nextPageCursor,omitempty"`
}

type Transactions struct {
	ID              string `json:"id,omitempty"`
	TradeID         string `json:"tradeId,omitempty"`
	Category        string `json:"category,omitempty"`
	Symbol          string `json:"symbol,omitempty"`
	Side            string `json:"side,omitempty"`
	OrderID         string `json:"orderId,omitempty"`
	Fee             string `json:"fee,omitempty"`
	CashFlow        string `json:"cashFlow,omitempty"`
	FeeRate         string `json:"feeRate,omitempty"`
	Funding         string `json:"funding,omitempty"`
	Type            string `json:"type,omitempty"`
	TransactionTime string `json:"transactionTime,omitempty"`
	Change          string `json:"change,omitempty"`
	Currency        string `json:"currency,omitempty"`
}

type RiskLimit struct {
	ID             int    `json:"id"`
	Symbol         string `json:"symbol"`
	RiskLimitValue string `json:"riskLimitValue"`
	MM             string `json:"maintenanceMargin"`
	IM             string `json:"initialMargin"`
	IsLowestRisk   int    `json:"isLowestRisk"`
	MaxLeverage    string `json:"maxLeverage"`
	MMDeduction    string `json:"mmDeduction"`
}

type RiskLimitResponse struct {
	Category       string      `json:"category"`
	NextPageCursor string      `json:"nextPageCursor"`
	RiskLimit      []RiskLimit `json:"list"`
}
type ReplayInfoResponse struct {
	RepayInfo []RepayInfo `json:"repayInfo"`
}
type RepayInfo struct {
	OrderID  string `json:"repayOrderId,omitempty"`
	Time     string `json:"repayTime,omitempty"`
	Token    string `json:"repayToken,omitempty"`
	Qty      string `json:"quantity,omitempty"`
	Interest string `json:"interest,omitempty"`
	Type     string `json:"businessType,omitempty"`
	Status   string `json:"status,omitempty"`
}

type Ticker struct {
	Symbol string `json:"symbol"`
	FR     string `json:"fundingRate"`
	Volume24h string `json:"volume24h"`
	Turnover24h string `json:"turnover24h"`
}

type TickerResponse struct {
	List []Ticker `json:"list"`
}

type ExecutionResponse struct {
	List []*Execution `json:"list"`
	Category string `json:"category"`
	NextPageCursor string `json:"nextPageCursor"`
}

type Execution struct {
	Category    string  `json:"category,omitempty"`
	Symbol      string  `json:"symbol,omitempty"`
	IsLeverage  string  `json:"isLeverage,omitempty"`
	Side        string  `json:"side,omitempty"`
	OrderID     string  `json:"orderId,omitempty"`
	OrderLinkID string  `json:"orderLinkId,omitempty"`
	OrderQty    float64 `json:"orderQty,omitempty,string"`
	LeavesQty   float64 `json:"leavesQty,omitempty,string"`
	OrderType   string  `json:"orderType,omitempty"`
	ExecType    string  `json:"execType,omitempty"`
	ExecPnL     string  `json:"execPnL,omitempty"`
	ExecFee     string  `json:"execFee,omitempty"`
	ExecID      string  `json:"execId,omitempty"`
	ExecPrice   float64 `json:"execPrice,omitempty,string"`
	ExecQty     float64 `json:"execQty,omitempty,string"`
	ExecTime    string  `json:"execTime,omitempty"`
	FeeRate     string  `json:"feeRate,omitempty"`
	IsMaker     bool    `json:"isMaker,omitempty"`
	SeqNum      int     `json:"seq,omitempty"`
}

type Order struct {
	Symbol       string `json:"symbol"`
	OrderID      string `json:"orderId"`
	OrderType    string `json:"orderType"`
	CancelType   string `json:"cancelType"`
	Price        string `json:"price"`
	Qty          string `json:"qty"`
	TimeInForce  string `json:"timeInForce"`
	OrderStatus  string `json:"orderStatus"`
	ReduceOnly   bool   `json:"reduceOnly"`
	Side         string `json:"side"`
	RejectReason string `json:"rejectReason"`
	CreateTime   string `json:"createdTime"`
	UpdateTime   string `json:"updatedTime"`
}

type FundingHistoryResponse struct {
	Category       string           `json:"category"`
	Histories	  []*FundingHistory `json:"list"`
}

type FundingHistory struct {
	Symbol	   string `json:"symbol"`
	Time       string `json:"fundingRateTimestamp"`
	FundingRate string `json:"fundingRate"`
}

