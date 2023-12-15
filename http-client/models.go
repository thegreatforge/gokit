package httpclient

import (
	"time"
)

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
