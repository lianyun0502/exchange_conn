# exchange_conn
multiple crypto exchanges SDk


## Index
- [Environment](#environment)
- [Installation](#installation)
- [Example](#example)



## Environment
* go version: go version go1.22.6 linux/amd64
* OS: Ubuntu Ubuntu 22.04.3 LTS

## Installation
* git clone the repository

```bash
git clone https://github.com/lianyun0502/exchange_conn.git
```

* install the package

```bash
go get github.com/lianyun0502/exchange_conn
```

## Example

* Exchange Restful API (Binance)
    - simple request
    ``` Go   
    package main

    import (
        "fmt"
        "github.com/lianyun0502/exchange_conn"
        "github.com/lianyun0502/exchange_conn/v1/binance_conn"
    )

    func main() {
        // create a new agent
        // agent is and adpater for every exchange client, set the client instance to the agent.  
        agent := exchange_conn.NewAgent(binance_conn.NewClient("YourAPIKey", "YourSecretKey", "https://api.binance.com"))

        // send a request to the exchange
        data, err := agent.Request(http.MethodGet, "/api/v3/ping", false, false).Send()
        if err != nil {
            t.Error(err)
            return
        }

        fmt.Println(data)
    }
    ```
    - request with query parameters or body parameters

        there are two ways to set the query/body parameters, one is to use the `SetQuery()`/`SetParam()` method for sigle param, the other is to use the `SetQueries`/`SetParams()` method for mutiple params.
    ``` Go
    package main

    import (
        "fmt"
        "github.com/lianyun0502/exchange_conn"
        "github.com/lianyun0502/exchange_conn/v1/binance_conn"
    )


    func main() {
        // create a new agent
        // agent is and adpater for every exchange client, set the client instance to the agent.  
        agent := exchange_conn.NewAgent(binance_conn.NewClient("YourAPIKey", "YourSecretKey", "https://api.binance.com"))

        // send a request to the exchange with query parameters
        data, err := agent.Request(http.MethodGet, "/api/v3/depth", false, false).SetQuery("symbol", "BTCUSDT").SetQuery("limit", "10").Send()
        // set the query parameters with a map
        // params := map[string]string{"symbol": "BTCUSDT", "limit": "10"}
        // data, err := agent.Request(http.MethodGet, "/api/v3/depth", false, false).SetQueries(params).Send()
	    if err != nil {
		t.Error(err)
		return
	    }

        fmt.Println(data)
    }
    ```

* Exchange Websocket API (Binance)

    ``` Go
    package main

    import (
        "fmt"
        "time"
        "github.com/lianyun0502/exchange_conn"
        "github.com/lianyun0502/exchange_conn/v1/binance_conn"
    )

    // define the handler for the data stream
    func wsHandler(data []byte) {
        fmt.Println(string(data))
    }

    // define the error handler
    func errorHandler(err error) {
        fmt.Println(err)
    }

    func main() {
        url := "wss://stream.binance.com:443/ws"

        // create a new agent for exchange websocket
        agent := exchange_conn.NewWebSocketAgent(binance_conn.NewWsClient(wsHandler, errorHandler, 10))

        // connect to the exchange
        agent.Connect(url)

        // start the loop to receive the data stream
        go agent.StartLoop()

        // send a request to the exchange to ask for the data
        agent.SendString(`{"method": "SUBSCRIBE","params": ["btcusdt@trade", "btcusdt@aggTrade", "btcusdt@depth@100ms"],"id": 1}`)

        // goroutine to stop the program after 10 seconds
        go func() {
            time.Sleep(10 * time.Second)
            agent.Stop()
        }()

        // block the program until the agent is stopped
        <-agent.Client.DoneSignal
    }
    ```




