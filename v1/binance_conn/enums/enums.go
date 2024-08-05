package enums

// This will apply for both Rest API and WebSocket API.


// OrderSide is the type for the side of an order
type OrderSide string
const (
	Buy OrderSide = "BUY"
	Sell OrderSide = "SELL"
)

// Order types (orderTypes, type)
type OrderType string
const (
	Limit OrderType = "LIMIT"
	Market OrderType = "MARKET"
	StopLoss OrderType = "STOP_LOSS"
	StopLossLimit OrderType = "STOP_LOSS_LIMIT"
	TakeProfit OrderType = "TAKE_PROFIT"
	TakeProfitLimit OrderType = "TAKE_PROFIT_LIMIT"
	LimitMaker OrderType = "LIMIT_MAKER"
)

// Time in force (timeInForce)
// This sets how long an order will be active before expiration.
type TimeInForce string
const (
	// Good Till Cancelled. 
	// 	An order will be on the book unless the order is canceled.
	//  一定時間等待全部合約在指定的價格成交，並可以隨時靈活地取消尚未成交的合約
	GTC TimeInForce = "GTC" 

	// Immediate Or Cancel
	// 	An order will try to fill the order as much as it can before the order expires.
	// 	馬上執行部分符合價格部位，而剩下的部位會「馬上被取消」
	IOC TimeInForce = "IOC"

	// Fill Or Kill
	// 	An order will expire if the full order cannot be filled upon execution.
	// 	訂單必須立即以委託價或更佳/優的價格全部成交，否則將被完全取消，不允許部分成交
	FOK TimeInForce = "FOK" // Fill or Kill
)
// Order Response Type (newOrderRespType)
type NewOrderRespType string
const (
	ACK NewOrderRespType = "ACK" 
	Result NewOrderRespType = "RESULT"
	Full NewOrderRespType = "FULL"
)

// STP Modes (selfTradePrevention)
type SelfTradePreventionMode string
const (
	none SelfTradePreventionMode = "NONE"
	expireMaker SelfTradePreventionMode = "EXPIRE_MAKER"
	expireTaker SelfTradePreventionMode = "EXPIRE_TAKER"
	expireBoth SelfTradePreventionMode = "EXPIRE_BOTH"
) 
	
