package exchange_conn_test

import (
	"encoding/json"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lianyun0502/exchange_conn/v1/binance_conn"
	"github.com/lianyun0502/exchange_conn/v1"
	"github.com/lianyun0502/exchange_conn/v1/common"
)

func TestExchangePing(t *testing.T) {
	client := binance_conn.NewClient(
		"",
		"",
		"https://api.binance.com",
	)
	agent := exchange_conn.NewAgent(client)

	data, err := agent.Request(http.MethodGet, "/api/v3/ping", false, false).Send()
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, string(data), "{}")
}

func TestExchangeGetInfo(t *testing.T) {
	client := binance_conn.NewClient(
		"",
		"",
		"https://api.binance.com",
	)
	agent := exchange_conn.NewAgent(client)

	data, err := agent.Request(http.MethodGet, "/api/v3/exchangeInfo", false, false).Send()
	if err != nil {
		t.Error(err)
		return
	}

	assert.NotEqual(t, string(data), "{}")
	j := new(interface{})
	json.Unmarshal(data, &j)
	log.Println(common.PrettyPrint(j))
}
