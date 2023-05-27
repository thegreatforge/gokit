package httpclient

import (
	"time"
)

// OnRetryHook is a function to be called on retry
// string: remote service name
// string: http request path
// string: http method
type OnRetryHook func(string, string, string) error

// OnClientResponseHook is a function to be called on http client response
// string: remote service name
// string: http request path
// string: http method
// int: http status code
// time.Duration: duration of the request
type OnClientResponseHook func(string, string, string, int, time.Duration) error

// Request is the request model for the HTTP client
type Request struct {
	Path            string            // Path
	Body            interface{}       // Body
	Headers         map[string]string // Headers
	OverrideTimeout time.Duration     // Override Timeout of the client in seconds
}

// Response is the response model for the HTTP client
type Response struct {
	Body       interface{}
	Headers    map[string]string
	StatusCode int
}
