package binance

const (
	Name    = "binance.api.go"
	Version = "1.0.0"

	SPOT_QUOTE_MAINNET = "wss://stream.binance.com:9443/ws"
	SPOT_QUOTE_TESTNET = "wss://testnet.binance.vision/ws"
	USD_QUOTE_MAINNET  = "wss://fstream.binance.com/ws"
	USD_QUOTE_TESTNET  = "wss://fstream.binancefuture.com/ws"
	COIN_QUOTE_MAINNET = "wss://dstream.binance.com"
	COIN_QUOTE_TESTNET = "wss://dstream.binancefuture.com"


	// WebSocket public channel - Mainnet
	SPOT_MAINNET    = "wss://ws-api.binance.com/ws-api/v3"
	USD_MAINNET  = "wss://ws-fapi.binance.com/ws-fapi/v1"
	COIN_MAINNET = "wss://ws-dapi.binance.com/ws-fapi/v1"
	OPTION_MAINNET  = "wss://stream.bybit.com/v5/public/option"

	// WebSocket public channel - Testnet
	SPOT_TESTNET    = "wss://testnet.binance.vision/ws-api/v3"
	USD_TESTNET  = "wss://testnet.binancefuture.com/ws-fapi/v1"
	COIN_TESTNET = "wss://testnet.binancefuture.com/ws-dapi/v1"
	OPTION_TESTNET  = "wss://stream-testnet.bybit.com/v5/public/option"

	// WebSocket private channel
	WEBSOCKET_PRIVATE_MAINNET = "wss://stream.bybit.com/v5/private"
	WEBSOCKET_PRIVATE_TESTNET = "wss://stream-testnet.bybit.com/v5/private"

	// WebSocket Trade channel
	WEBSOCKET_TRADE_MAINNET = "wss://stream.bybit.com/v5/trade"
	WEBSOCKET_TRADE_TESTNET = "wss://stream-testnet.bybit.com/v5/trade"

	// V3
	V3_CONTRACT_PRIVATE = "wss://stream.bybit.com/contract/private/v3"
	V3_UNIFIED_PRIVATE  = "wss://stream.bybit.com/unified/private/v3"
	V3_SPOT_PRIVATE     = "wss://stream.bybit.com/spot/private/v3"
)