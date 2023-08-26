# HTTP Client Go Library

The `httpclient` package provides an HTTP client library for making HTTP requests with retry and timeout support. This library is designed to simplify making HTTP requests in Go applications while providing features like request retries and customizable timeouts.

## Installation

To use this library in your Go application, you can import it using the following import statement:

```go
import "github.com/thegreatforge/gokit/http-client"
```

## Usage

### Creating a Client

To create an instance of the HTTP client, you need to configure it using the `ClientConfig` struct and pass it to the `NewClient` function. Here's an example:

```go
import (
    "time"
    httpclient "github.com/thegreatforge/gokit/http-client"
)

func main() {
    config := httpclient.ClientConfig{
        Host:             "https://api.example.com",
        DefaultHeaders:   map[string]string{"Authorization": "Bearer YOUR_TOKEN"},
        Timeout:          10 * time.Second,
        Retries:          3,
        RetryInterval:    2 * time.Second,
        Logger:           yourLogger, // Replace with your zap logger instance
        OnRetry:          yourRetryHook, // For retry metrics
        OnClientResponse: yourResponseHook, // For every client response metrics
    }

    client := httpclient.NewClient(config)

    // Register the client for later use
    err := httpclient.RegisterClient("exampleClient", client)
    if err != nil {
        // Handle registration error
    }
}
```

### Making Requests

The `Client` instance provides methods for making different types of HTTP requests:

- `Post(ctx context.Context, req Request, resp *Response) error`
- `Get(ctx context.Context, req Request, resp *Response) error`
- `Put(ctx context.Context, req Request, resp *Response) error`
- `Delete(ctx context.Context, req Request, resp *Response) error`
- `Patch(ctx context.Context, req Request, resp *Response) error`

Here's how you can use these methods to make requests:

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
		DefaultHeaders:   map[string]string{"Authorization": "Bearer YOUR_TOKEN"} 
		Timeout:          1 * time.Second,
		Retries:          3,
		RetryInterval:    5 * time.Second,
		Logger:           zap.NewNop(),
		OnRetry:          nil,
		OnClientResponse: nil,
	}

	// create new client and register it to global variable
	cli := httpclient.NewClient(config)
	_ = httpclient.RegisterClient("example", cli)

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

## Closing Thoughts

The `httpclient` package simplifies the process of making HTTP requests with retry and timeout support. By providing an intuitive API and customizable hooks, it aims to improve the reliability of HTTP communication in your Go applications.