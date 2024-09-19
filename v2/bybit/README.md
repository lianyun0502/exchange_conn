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
    "github.com/lianyun0502/exchange_conn/v2/bybit/http_client"
)

func main() {
    // create a new bybit client
    client := bybit.NewSpotClient("YourAPIKey", "YourSecretKey")
    // restful api request body
	param := ParamMap{
		"category":  "spot",
		"symbol":    "BTCUSDT",
		"side":      "Buy",
		"orderType": "Limit",
		"qty":       "0.001",
		"price":     "50000",
	}
    // send a request to the exchange
	data, err := client.Request(http.MethodPost, "/v5/order/create", bybit.SetSercurityType(true, true)).SetParam(param).Send()
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
        "github.com/lianyun0502/exchange_conn/v2/bybit/ws_client"
    )

    func main() {
        // create a new bybit websocket client
        client := bybit.NewWsTradeClient("YourAPIKey", "YourSecretKey")
        // connect to the exchange
        client.Connect()
        // start the message read loop
        go client.Conn.ReadLoop()
        // get the signature from the exchange
        resp, err := client.GetSignature()
        if err != nil {
            t.Log(string(resp))
            t.Error(err)
            return
        }
        // print the response data
        type ParamMap map[string]string
        param := ParamMap{
            "category":  "spot",
            "symbol":    "BTCUSDT",
            "side":      "Buy",
            "orderType": "Limit",
            "qty":       "0.001",
            "price":     "50000",
        }
        // send a order request to the exchange
        resp, err = client.Order("order.create", []ParamMap{param})
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
        "github.com/lianyun0502/exchange_conn/v2/bybit/ws_client"
    )

    func main() {
        // create a new bybit websocket client for quote
        client := bybit.NewWsSpotClient()
        client.Connect()
        client.Ws_Handler = func(data []byte) {
            fmt.Println(string(data))
        }
        go client.Conn.ReadLoop()
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
