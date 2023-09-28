package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

var Clients map[string]*Client

// Client is a HTTP client with retry and timeout support
type Client struct {
	client           *http.Client
	defaultHeaders   map[string]string
	retries          int
	retryInterval    time.Duration
	host             string
	onRetry          OnRetryHook
	onClientResponse OnClientResponseHook
	logger           *zap.Logger
}

func init() {
	Clients = make(map[string]*Client)
}

// ClientConfig is the configuration for the HTTP client
// Host: the host of the service
// Timeout: the timeout for the HTTP request
// Retries: the number of retries for the HTTP request
// RetryInterval: the interval between retries
// Logger: the logger
// OnRetry: the hook to be called on retry
// OnClientResponse: the hook to be called on client response
type ClientConfig struct {
	Host             string
	DefaultHeaders   map[string]string
	Timeout          time.Duration
	Retries          int
	RetryInterval    time.Duration
	OnRetry          OnRetryHook
	OnClientResponse OnClientResponseHook
	Logger           *zap.Logger
}

// NewClient creates a new HTTP client with the given configuration
// config: ClientConfig
func NewClient(config ClientConfig) *Client {

	hcli := &Client{
		client: &http.Client{
			Timeout: config.Timeout,
		},
		retries:          config.Retries,
		retryInterval:    config.RetryInterval,
		host:             config.Host,
		defaultHeaders:   config.DefaultHeaders,
		onRetry:          config.OnRetry,
		onClientResponse: config.OnClientResponse,
		logger:           config.Logger,
	}
	return hcli
}

// RegisterClient registers a new HTTP client to global map
func RegisterClient(clientName string, cli *Client) error {
	_, ok := Clients[clientName]
	if ok {
		return fmt.Errorf("client %s already registered", clientName)
	}
	Clients[clientName] = cli
	return nil
}

// getRequestId returns the request id from the context
func getRequestId(ctx context.Context) string {
	if ctx == nil {
		return uuid.New().String()
	}

	ginCtx, ok := ctx.(*gin.Context)
	if ok {
		if ginCtx.GetHeader(XRequestIdHeaderKey) != "" {
			return ginCtx.GetHeader(XRequestIdHeaderKey)
		}
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		if len(md.Get(XRequestIdHeaderKey)) > 0 {
			return md.Get(XRequestIdHeaderKey)[0]
		}
	}

	return uuid.New().String()
}

// prepareRequestBody prepares the request body for the HTTP request
func prepareRequestBody(reqBody interface{}) (io.Reader, error) {
	if reqBody != nil {
		body, err := json.Marshal(reqBody)
		if err != nil {
			return nil, fmt.Errorf("error marshalling request body: %s", err)
		}
		return bytes.NewBuffer(body), nil
	}
	return nil, nil
}

// readBody reads the response body and unmarshals it into the given interface
func readBody(httpResp *http.Response, respBody interface{}) error {
	defer httpResp.Body.Close()
	if respBody != nil {
		return json.NewDecoder(httpResp.Body).Decode(respBody)
	}
	return nil
}

// overrideTimeOut overrides the timeout of the context if the timeout is not 0
func overrideTimeOut(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout != 0 {
		return context.WithTimeout(ctx, timeout)
	}
	return ctx, func() {}
}

// makeHttpRequestWithRetries makes the HTTP request with retries
func (c *Client) makeHttpRequestWithRetries(
	ctx context.Context,
	httpMethod string,
	req Request,
	resp *Response,
) error {

	httpCtx, cancel := overrideTimeOut(ctx, req.OverrideTimeout)
	defer cancel()

	requestId := getRequestId(httpCtx)

	var reqBody io.Reader
	var err error

	for i := 0; i <= c.retries; i++ {

		if httpMethod == http.MethodPut || httpMethod == http.MethodPost || httpMethod == http.MethodPatch {
			reqBody, err = prepareRequestBody(req.Body)
			if err != nil {
				return fmt.Errorf("error preparing request body: %s", err)
			}
		}

		// create the request
		httpReq, err := http.NewRequestWithContext(httpCtx, httpMethod, c.host+req.Path, reqBody)
		if err != nil {
			return fmt.Errorf("error creating request: %s", err)
		}

		// set the headers
		if c.defaultHeaders != nil {
			for k, v := range c.defaultHeaders {
				httpReq.Header.Set(k, v)
			}
		}

		httpReq.Header.Set(XRequestIdHeaderKey, requestId)
		httpReq.Header.Set("Content-Type", "application/json")
		for k, v := range req.Headers {
			httpReq.Header.Set(k, v)
		}

		startTime := time.Now()
		httpResp, err := c.client.Do(httpReq)
		if err != nil {
			c.logger.Sugar().With(XRequestIdHeaderKey, requestId).Errorf("request failed with error: %s", err)
		}
		latency := time.Since(startTime)

		// execute the onClientResponse hook
		if c.onClientResponse != nil {
			err = c.onClientResponse(req.Path, httpMethod, httpResp.StatusCode, latency)
			if err != nil {
				c.logger.Sugar().With(XRequestIdHeaderKey, requestId).Errorf("request failed with error: %s", err)
			}
		}

		if httpResp != nil {
			// success
			if httpResp.StatusCode >= 200 && httpResp.StatusCode < 400 {

				if resp != nil {
					err = readBody(httpResp, resp.Body)
					if err != nil {
						return fmt.Errorf("error reading response body: %s", err)
					}
					resp.StatusCode = httpResp.StatusCode

					if resp.Headers == nil {
						resp.Headers = make(map[string]string, len(httpResp.Header))
					}

					for k, v := range httpResp.Header {
						if len(v) > 0 {
							resp.Headers[k] = v[0]
						}
					}
				}

				return nil

			}
			// failure
			c.logger.Sugar().With(XRequestIdHeaderKey, requestId).Errorf("request failed with status %d", httpResp.StatusCode)
		}

		if i != c.retries {
			if c.onRetry != nil {
				err = c.onRetry(req.Path, httpMethod)
				if err != nil {
					c.logger.Sugar().With(XRequestIdHeaderKey, requestId).Errorf("error executing onRetry hook: %s", err)
				}
			}
			c.logger.Sugar().With(XRequestIdHeaderKey, requestId).Warnf("retrying in %s...", c.retryInterval)
			time.Sleep(c.retryInterval)
		}

		if c.retries == i {
			c.logger.Sugar().With(XRequestIdHeaderKey, requestId).Errorf("request failed after maximum %d retries", c.retries)
			if httpResp != nil && resp != nil {
				resp.StatusCode = httpResp.StatusCode
				_ = readBody(httpResp, resp.Body)

				if resp.Headers == nil {
					resp.Headers = make(map[string]string, len(httpResp.Header))
				}
				for k, v := range httpResp.Header {
					if len(v) > 0 {
						resp.Headers[k] = v[0]
					}
				}
			}
		}
	}

	return fmt.Errorf("request failed after maximum %d retries", c.retries)
}

// Close closes the idle connections of the HTTP client
func (c *Client) CloseIdleConnections() {
	c.client.CloseIdleConnections()
}

// Post makes a HTTP POST request with the given request
func (c *Client) Post(ctx context.Context, req Request, resp *Response) error {
	return c.makeHttpRequestWithRetries(ctx, http.MethodPost, req, resp)
}

// Get makes a HTTP GET request with the given request
func (c *Client) Get(ctx context.Context, req Request, resp *Response) error {
	return c.makeHttpRequestWithRetries(ctx, http.MethodGet, req, resp)
}

// Put makes a HTTP PUT request with the given request
func (c *Client) Put(ctx context.Context, req Request, resp *Response) error {
	return c.makeHttpRequestWithRetries(ctx, http.MethodPut, req, resp)
}

// Delete makes a HTTP DELETE request with the given request
func (c *Client) Delete(ctx context.Context, req Request, resp *Response) error {
	return c.makeHttpRequestWithRetries(ctx, http.MethodDelete, req, resp)
}

// Patch makes a HTTP PATCH request with the given request
func (c *Client) Patch(ctx context.Context, req Request, resp *Response) error {
	return c.makeHttpRequestWithRetries(ctx, http.MethodPatch, req, resp)
}
