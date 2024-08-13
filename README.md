# exchange_conn
multiple crypto exchanges SDk


## Index
- [Environment](#environment)
- [Installation](#installation)
- [Example](#example)



## Environment
* Go version: go version go1.22.6 linux/amd64
* OS: Ubuntu Ubuntu 22.04.3 LTS

## Installation

there are two ways to install and use the package, one is to clone the repository and refer to local path and the other is to use the `go get` command set to `go.mod`.

### Git clone the repository

1. first clone the repository into your project directory

    ```bash
    git clone https://github.com/lianyun0502/exchange_conn.git
    ```

    Your directory structure should look like this:

    ```bash
    your_project/
    ├── exchange_conn/
    ├── main.go
    └── go.mod
    ```
    
2. replace the import refernce with the path of the repository in your project.

    ```bash
    go mod edit -replace=github.com/lianyun0502/exchange_conn=../exchange_conn
    ```

3. import the package in your project

    ```Go
    import (
        "github.com/lianyun0502/exchange_conn"
        "github.com/lianyun0502/exchange_conn/v1/binance_conn"
    )
    ```

### Install the package use `go get`

1. use the `go get` command to install the package

    ```bash
    go get github.com/lianyun0502/exchange_conn
    ```
2. import the package in your project

    ```Go
    import (
        "github.com/lianyun0502/exchange_conn"
        "github.com/lianyun0502/exchange_conn/v1/binance_conn"
    )
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

* Exchange Websocket stream (Binance)

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




