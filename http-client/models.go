package httpclient

import "time"

// Request is the request model for the HTTP client
type Request struct {
	Path            string            // Path
	Body            interface{}       // Body
	Headers         map[string]string // Headers
	OverrideTimeout time.Duration     // Override Timeout of the client
}

// Response is the response model for the HTTP client
type Response struct {
	Body       interface{}
	Headers    map[string]string
	StatusCode int
}
