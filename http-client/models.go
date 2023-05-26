package httpclient

import "time"

// Request is the request model for the HTTP client
type Request struct {
	Path      string            // Path
	Body      interface{}       // Body
	Headers   map[string]string // Headers
	Timeout   time.Duration     // Override Timeout
	RequestId string            // Unique identifier for the request, if not set, it will be searched in context or generated if not found
}

// Response is the response model for the HTTP client
type Response struct {
	Body       interface{}
	Headers    map[string]string
	StatusCode int
}
