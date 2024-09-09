package bybit

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/lianyun0502/exchange_conn/v2/common"
	"github.com/lianyun0502/exchange_conn/v2/http_client"
	"github.com/lianyun0502/exchange_conn/v2/http_client/consts"
	"github.com/sirupsen/logrus"
)


const (
	Signed = 0b01
	ApiKey = 0b10
)

type SecurityT int
const (
	None        SecurityT = 0               // all public access
	Trade                 = Signed | ApiKey // API-key and Singnature required
	UserData              = Signed | ApiKey // API-key and Singnature required
	UserStream            = ApiKey          // API-key required
	MARKET_DATA           = ApiKey          // API-key required
)

type Request struct {
	Method   string
	Endpoint string
	body  []byte
	query url.Values

	SecType int
	recvWindow string
}

func (r *Request) SetBody(body []byte) {
	r.body = body
}

func (r *Request) SetQuery(query url.Values) {
	r.query = query
}

func (r *Request) Body() []byte {
	return r.body
}

func (r *Request) Query() url.Values {
	return r.query
}

func NewRequest(method, endpoint string, reqOpts ...func(*Request)) *Request {
	req := &Request{
		Method: method,
		Endpoint: endpoint,
		recvWindow: "5000",
	}
	for _, opt := range reqOpts {
		opt(req)
	}
	return req
}
		
func SetRequest(r *Request, c*exchange_conn.HttpClient[*Request]) (req *http.Request, err error) {
	fullURL := fmt.Sprintf("%s%s", c.BaseURL, r.Endpoint)

	queryString := r.Query().Encode()
	if queryString != "" {
		fullURL = fmt.Sprintf("%s?%s", fullURL, queryString)
	}

	req, err = http.NewRequest(r.Method, fullURL, bytes.NewBuffer(r.Body()))
	if err != nil {
		return nil, err
	}
	c.Log.WithFields(logrus.Fields{"url": fullURL}).Debug("Compose URL")
	c.Log.WithFields(logrus.Fields{"body": r.Body()}).Debug("Body Ready")

	req.Header.Set("User-Agent", fmt.Sprintf("%s/%s", "bybit_connect", "v1"))
	if r.Body() != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	if (r.SecType & Signed) != 0x00 {
		timeStamp := common.GetCurrentTime()
		req.Header.Set(signTypeKey, "2")
		req.Header.Set(apiRequestKey, c.Exchange.APIKey)
		req.Header.Set(timestampKey, strconv.FormatInt(timeStamp, 10))
		req.Header.Set(recvWindowKey, r.recvWindow)

		var signatureBase string
		if r.Method == "POST" {
			req.Header.Set("Content-Type", "application/json")
			signatureBase = strconv.FormatInt(timeStamp, 10) + c.Exchange.APIKey + r.recvWindow + string(r.Body())
		} else {
			signatureBase = strconv.FormatInt(timeStamp, 10) + c.Exchange.APIKey + r.recvWindow + queryString
		}
		c.Log.WithFields(logrus.Fields{"signatureBase": signatureBase}).Debug("Signature Base")
		signature := common.GetSignature(c.Exchange.SecretKey, signatureBase)
		req.Header.Set(signatureKey, signature)
	}
	return req, err
}

type ByBitSpotClient struct{
	*exchange_conn.HttpClient[*Request]
}


func NewSpotClient(apiKey, secretKey string) *ByBitSpotClient {
	return &ByBitSpotClient{
		HttpClient: &exchange_conn.HttpClient[*Request]{
			Client: http.DefaultClient,
			BaseURL: "https://api-testnet.bybit.com",
			Exchange: &exchange_conn.ExchangeApi{
				Name: consts.Bybit,
				HostType: consts.Spot,
				APIKey: apiKey,
				SecretKey: secretKey,
			},
			NewRequest: NewRequest,
			SetRequest: SetRequest,
			Log: logrus.New(),
		},
	}
}

func test() {
	client := NewSpotClient("apiKey	", "secretKey")
	
	resp, _ := client.Request(http.MethodGet, "/api/v2/ping")
	fmt.Println(resp)
}