# exchange_conn
Multiple Exchanges SDK

## Index
- [Log](#log)
- [Environment](#environment)
- [Installation](#installation)
- [Example](#example)

## Log
- **2024-9/19 : release Bybit SDk v2 version**


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
        "github.com/lianyun0502/exchange_conn/v2"
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
        "github.com/lianyun0502/exchange_conn/v2"
    )
    ```

## Example

* [v1 Example](v1/README.md)
* [Bybit v2 Example](v2/Bybit/README.md)





