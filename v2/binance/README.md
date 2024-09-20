# Bybit SDK v2

如有需要可以參考Test case
- [http api test](./http_client/client_test.go)
- [websocket test](./ws_client/client_test.go)

## Example

* http Restful API
```Go
package main

import (
    "fmt"
    "github.com/lianyun0502/exchange_conn/v2/binance/http_client"
)

func main() {
    // create a new binance client
    client := bybit.NewAPISpotClient("YourAPIKey", "YourSecretKey")
    // restful api request body
	param := ParamMap{
		"symbol":      "BTCUSDT",
		"side":        enums.Buy,
		"type":        enums.Limit,
		"timeInForce": enums.GTC,
		"quantity":    0.0001,
		"price":       "50000",
	}
    // send a request to the exchange
	data, err := client.Request(http.MethodPost, "/api/v3/order", binance.SetSercurityType(true, true)).SetParam(param).Send()
	if err != nil {
		fmt.Println(err)
		return
	}
    // print the response data
    fmt.Println(data)
}
```
* Websocket API
```Go
    package main

    import (
        "fmt"
        "github.com/lianyun0502/exchange_conn/v2/binance/ws_client"
    )

    func main() {
        // create a new bybit websocket client
        client := bybit.NewWsAPIClient("spot", "YourAPIKey", "YourSecretKey")
        // connect to the exchange
        client.Connect()
        // start the message read loop
        go client.StartLoop()
        // print the response data
        type ParamMap map[string]string
        param := ParamMap{
            "symbol":      "BTCUSDT",
			"side":        "SELL",
			"type":        "LIMIT",
			"timeInForce": "GTC",
			"price":       "50000",
			"quantity":    "0.001",
			"timestamp":   strconv.FormatInt(time.Now().UnixMilli(), 10),
        }
        // send a order request to the exchange
        resp, err = client.Order("order.place", []ParamMap{param})
        if err != nil {
            fmt.Println(string(resp))
            return
        }
        fmt.Println(string(resp))

        // wait for the stop signal
        <- client.StopSignal
    }

``` 

* Exchange Websocket stream 
``` Go
    package main

    import (
        "fmt"
        "github.com/lianyun0502/exchange_conn/v2/binance/ws_client"
        "github.com/lianyun0502/exchange_conn/v2/consts"
    )

    func main() {
        // create a new bybit websocket client for quote
        client := bybit.NewWsQuoteClient(consts.Spot, handle)
        client.Connect()
        client.Ws_Handler = func(data []byte) {
            fmt.Println(string(data))
        }
        go client.StartLoop()
        // subscribe to the topic of the quote
        resp, err := client.Subscribe([]string{"orderbook.1.BTCUSDT", "publicTrade.BTCUSDT"})
        if err != nil {
            fmt.Println(string(resp))
            client.Stop()
            return
        }
        // wait for the stop signal, keep the program running
        <- client.StopSignal
    }
