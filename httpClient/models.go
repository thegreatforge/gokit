package httpClient

// Request is the request model for the HTTP client
type Request struct {
	Path    string
	Body    interface{}
	Headers map[string]string
}

// Response is the response model for the HTTP client
type Response struct {
	Body       interface{}
	Headers    map[string]string
	StatusCode int
}
