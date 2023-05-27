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
// Path: the path of the request
// Body: the body of the request
// Headers: the headers of the request
// OverrideTimeout: override the timeout of the client, it should be less than the client timeout
type Request struct {
	Path            string
	Body            interface{}
	Headers         map[string]string
	OverrideTimeout time.Duration
}

// Response is the response model for the HTTP client
type Response struct {
	Body       interface{}
	Headers    map[string]string
	StatusCode int
}
