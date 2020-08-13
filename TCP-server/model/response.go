package model

import (
	"fmt"
	"strings"
)

//Builder for Response model
type Builder interface {
	Protocol(string) Builder
	Status(string) Builder
	Body(string) Builder
	Headers(...string) Builder
	Build() Response
}

func NewBuilder() Builder {
	return new(Response)
}

//Response model for http response
type Response struct {
	protocol, status, body string
	headers                map[string]string
}

func (res *Response) Protocol(protocol string) Builder {
	res.protocol = protocol
	return res
}

func (res *Response) Status(status string) Builder {
	res.status = status
	return res
}

func (res *Response) Body(body string) Builder {
	res.body = body
	return res
}

func (res *Response) Headers(headers ...string) Builder {
	if res.headers == nil {
		res.headers = make(map[string]string)
	}
	for i := 0; i < len(headers); i += 2 {
		res.headers[headers[i]] = headers[i+1]
	}
	return res
}

func (res *Response) Build() Response {
	return *res
}

func (res *Response) String() string {
	var str strings.Builder
	fmt.Fprintf(&str, "%s %s\r\n", res.protocol, res.status)
	for key, value := range res.headers {
		fmt.Fprintf(&str, "%s: %s\r\n", key, value)
	}
	fmt.Fprintf(&str, "\r\n")
	fmt.Fprintf(&str, "%s", res.body)
	return str.String()
}
