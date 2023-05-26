# http-client Package

http-client package is http client wrapper which will reduces the boiler plate needed to
marshall/un-marshall request/response bodies, retry the requests with timeouts and send headers like
`x-request-id` and `remote-service` for tracing.

## Usage


```go

import httpclient "github.com/thegreatforge/gokit/http-client"

type ResponseBody struct {
        Name string `json:"name"`
        Job  string `json:"job"`
}

type RequestBody struct {
        Name string `json:"name"`
        Job  string `json:"job"`
}

func Example() {
        // create new client and register it to global variable
        cli := httpclient.NewClient("https://reqres.in", "example", "example-service", 5, 3, 1, zap.New())
        httpclient.RegisterNewClient(cli)

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