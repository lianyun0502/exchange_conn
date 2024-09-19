package binance

const (
	Name    = "binance.api.go"
	Version = "1.0.0"
	// Https
	SPOT_MAINNET = "https://api.binance.com"
	SPOT_TESTNET = "https://testnet.binance.vision"

	UFUTURE_MAINNET = "https://fapi.binance.com"
	CFUTURE_MAINNET = "https://dapi.binance.com"
	FUTURE_TESTNET  = "https://testnet.binancefuture.com"

	// WebSocket public channel - Mainnet
	WS_SPOT_MAINNET    = "wss://stream.bybit.com/v5/public/spot"
	WS_UFUTURE_MAINNET = "wss://ws-fapi.binance.com/ws-fapi/v1"

	// WebSocket public channel - Testnet
	WS_SPOT_TESTNET    = "wss://stream-testnet.bybit.com/v5/public/spot"
	WS_UFUTURE_TESTNET = "wss://testnet.binancefuture.com/ws-fapi/v1"
)
