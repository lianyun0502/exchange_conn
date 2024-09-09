package exchange_conn

import (
	"net/url"
)


type Request struct {
	Method   string
	Endpoint string
	body  []byte
	query url.Values
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