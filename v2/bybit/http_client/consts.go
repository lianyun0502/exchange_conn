package bybit
const (
	Name    = "bybit.api.go"
	Version = "1.0.0"
	// Https
	MAINNET       = "https://api.bybit.com"
	MAINNET_BACKT = "https://api.bytick.com"
	TESTNET       = "https://api-testnet.bybit.com"

	// V3
	V3_CONTRACT_PRIVATE = "wss://stream.bybit.com/contract/private/v3"
	V3_UNIFIED_PRIVATE  = "wss://stream.bybit.com/unified/private/v3"
	V3_SPOT_PRIVATE     = "wss://stream.bybit.com/spot/private/v3"

	// Globals
	timestampKey  = "X-BAPI-TIMESTAMP"
	signatureKey  = "X-BAPI-SIGN"
	apiRequestKey = "X-BAPI-API-KEY"
	recvWindowKey = "X-BAPI-RECV-WINDOW"
	signTypeKey   = "X-BAPI-SIGN-TYPE"
)