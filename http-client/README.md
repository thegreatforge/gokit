# http-client Package

http-client package is http client wrapper which will reduces the boiler plate needed to
marshall/un-marshall request/response bodies, retry the requests with timeouts and send headers like
`x-request-id` and `service` for tracing via logs.

## Usage


```go
package main

import (
	"context"
	"fmt"

	httpclient "github.com/thegreatforge/gokit/http-client"
	"go.uber.org/zap"
)

type ResponseBody struct {
	Name string `json:"name"`
	Job  string `json:"job"`
}

type RequestBody struct {
	Name string `json:"name"`
	Job  string `json:"job"`
}

func main() {
	config := ClientConfig{
		Host:             "https://reqres.in",
		Service:          "example",
		RemoteService:    "example-service",
		Timeout:          1 * time.Second,
		Retries:          3,
		RetryInterval:    5 * time.Second,
		Logger:           zap.NewNop(),
		OnRetry:          nil,
		OnClientResponse: nil,
	}

	// create new client and register it to global variable
	cli := httpclient.NewClient(config)
	httpclient.RegisterClient(cli)

	resp := &httpclient.Response{
		Body: &ResponseBody{},
	}
	err := httpclient.Clients["example"].Post(context.Background(), httpclient.Request{
		Path: "/api/users/2",
		Body: &RequestBody{
			Name: "morpheus",
			Job:  "leader",
		},
	}, resp)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(resp.Body.(*ResponseBody).Name)
	for k, v := range resp.Headers {
		fmt.Println(k, v)
	}
}

```