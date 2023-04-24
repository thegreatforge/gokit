package httpClient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

var Clients map[string]*Client

// Client is a HTTP client with retry and timeout support
type Client struct {
	client        *http.Client
	retries       int
	retryInterval time.Duration
	host          string
	service       string
	remoteService string
	logger        *zap.Logger
}

func init() {
	Clients = make(map[string]*Client)
}

// RegisterNewClient registers a new HTTP client with the given configuration
func RegisterNewClient(host string, service string, remoteService string,
	timeout time.Duration, retries int, retryInterval time.Duration,
	logger *zap.Logger) {
	hcli := &Client{
		client: &http.Client{
			Timeout: timeout * time.Second,
		},
		retries:       retries,
		retryInterval: retryInterval * time.Second,
		host:          host,
		service:       service,
		remoteService: remoteService,
		logger:        logger,
	}
	Clients[service] = hcli
	logger.Sugar().Infof("registered new http client for service %s", service)
}

// getRequestId returns the request id from the context
func getRequestId(ctx context.Context) string {
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

// Post makes a HTTP POST request with the given request
func (c *Client) Post(ctx context.Context, req Request, resp *Response) error {
	reqBody, err := prepareRequestBody(req.Body)
	if err != nil {
		return fmt.Errorf("error preparing request body: %s", err)
	}

	requestId := getRequestId(ctx)
	for i := 0; i <= c.retries; i++ {

		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.host+req.Path, reqBody)
		if err != nil {
			return fmt.Errorf("error creating request: %s", err)
		}
		httpReq.Header.Set(XRequestIdHeaderKey, requestId)
		httpReq.Header.Set(RemoteServiceHeaderKey, c.remoteService)
		httpReq.Header.Set("Content-Type", "application/json")
		for k, v := range req.Headers {
			httpReq.Header.Set(k, v)
		}

		httpResp, err := c.client.Do(httpReq)
		if err != nil {
			c.logger.Sugar().With(XRequestIdHeaderKey, requestId).Errorf("request failed with error: %s", err)
		}

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
					resp.Headers[k] = v[0]
				}
				return nil
			}
			return nil
		} else {
			c.logger.Sugar().With(XRequestIdHeaderKey, requestId).Errorf("request failed with status %d", httpResp.StatusCode)
		}

		if i != c.retries {
			c.logger.Sugar().With(XRequestIdHeaderKey, requestId).Warnf("Retrying in %s...", c.retryInterval)
			time.Sleep(c.retryInterval)
		}
	}

	return fmt.Errorf("request failed after maximum retries")
}

// Get makes a HTTP GET request with the given request
func (c *Client) Get(ctx context.Context, req Request, resp *Response) error {
	requestId := getRequestId(ctx)
	for i := 0; i <= c.retries; i++ {

		httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, c.host+req.Path, nil)
		if err != nil {
			return fmt.Errorf("error creating request: %s", err)
		}
		httpReq.Header.Set(XRequestIdHeaderKey, requestId)
		httpReq.Header.Set(RemoteServiceHeaderKey, c.remoteService)
		httpReq.Header.Set("Content-Type", "application/json")
		for k, v := range req.Headers {
			httpReq.Header.Set(k, v)
		}

		httpResp, err := c.client.Do(httpReq)
		if err != nil {
			c.logger.Sugar().With(XRequestIdHeaderKey, requestId).Errorf("request failed with error: %s", err)
		}

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
					resp.Headers[k] = v[0]
				}
				return nil
			}
			return nil
		} else {
			c.logger.Sugar().With(XRequestIdHeaderKey, requestId).Errorf("request failed with status %d", httpResp.StatusCode)
		}

		if i != c.retries {
			c.logger.Sugar().With(XRequestIdHeaderKey, requestId).Warnf("Retrying in %s...", c.retryInterval)
			time.Sleep(c.retryInterval)
		}
	}

	return fmt.Errorf("request failed after maximum retries")
}

// Put makes a HTTP PUT request with the given request
func (c *Client) Put(ctx context.Context, req Request, resp *Response) error {
	reqBody, err := prepareRequestBody(req.Body)
	if err != nil {
		return fmt.Errorf("error preparing request body: %s", err)
	}

	requestId := getRequestId(ctx)
	for i := 0; i <= c.retries; i++ {

		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPut, c.host+req.Path, reqBody)
		if err != nil {
			return fmt.Errorf("error creating request: %s", err)
		}
		httpReq.Header.Set(XRequestIdHeaderKey, requestId)
		httpReq.Header.Set(RemoteServiceHeaderKey, c.remoteService)
		httpReq.Header.Set("Content-Type", "application/json")
		for k, v := range req.Headers {
			httpReq.Header.Set(k, v)
		}

		httpResp, err := c.client.Do(httpReq)
		if err != nil {
			c.logger.Sugar().With(XRequestIdHeaderKey, requestId).Errorf("request failed with error: %s", err)
		}

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
					resp.Headers[k] = v[0]
				}
				return nil
			}
			return nil
		} else {
			c.logger.Sugar().With(XRequestIdHeaderKey, requestId).Errorf("request failed with status %d", httpResp.StatusCode)
		}

		if i != c.retries {
			c.logger.Sugar().With(XRequestIdHeaderKey, requestId).Warnf("Retrying in %s...", c.retryInterval)
			time.Sleep(c.retryInterval)
		}
	}

	return fmt.Errorf("request failed after maximum retries")
}

// Delete makes a HTTP DELETE request with the given request
func (c *Client) Delete(ctx context.Context, req Request, resp *Response) error {
	requestId := getRequestId(ctx)
	for i := 0; i <= c.retries; i++ {

		httpReq, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.host+req.Path, nil)
		if err != nil {
			return fmt.Errorf("error creating request: %s", err)
		}
		httpReq.Header.Set(XRequestIdHeaderKey, requestId)
		httpReq.Header.Set(RemoteServiceHeaderKey, c.remoteService)
		httpReq.Header.Set("Content-Type", "application/json")
		for k, v := range req.Headers {
			httpReq.Header.Set(k, v)
		}

		httpResp, err := c.client.Do(httpReq)
		if err != nil {
			c.logger.Sugar().With(XRequestIdHeaderKey, requestId).Errorf("request failed with error: %s", err)
		}

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
					resp.Headers[k] = v[0]
				}
				return nil
			}
			return nil
		} else {
			c.logger.Sugar().With(XRequestIdHeaderKey, requestId).Errorf("request failed with status %d", httpResp.StatusCode)
		}

		if i != c.retries {
			c.logger.Sugar().With(XRequestIdHeaderKey, requestId).Warnf("Retrying in %s...", c.retryInterval)
			time.Sleep(c.retryInterval)
		}
	}

	return fmt.Errorf("request failed after maximum retries")
}

// Patch makes a HTTP PATCH request with the given request
func (c *Client) Patch(ctx context.Context, req Request, resp *Response) error {
	reqBody, err := prepareRequestBody(req.Body)
	if err != nil {
		return fmt.Errorf("error preparing request body: %s", err)
	}

	requestId := getRequestId(ctx)
	for i := 0; i <= c.retries; i++ {

		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPatch, c.host+req.Path, reqBody)
		if err != nil {
			return fmt.Errorf("error creating request: %s", err)
		}
		httpReq.Header.Set(XRequestIdHeaderKey, requestId)
		httpReq.Header.Set(RemoteServiceHeaderKey, c.remoteService)
		httpReq.Header.Set("Content-Type", "application/json")
		for k, v := range req.Headers {
			httpReq.Header.Set(k, v)
		}

		httpResp, err := c.client.Do(httpReq)
		if err != nil {
			c.logger.Sugar().With(XRequestIdHeaderKey, requestId).Errorf("request failed with error: %s", err)
		}

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
					resp.Headers[k] = v[0]
				}
				return nil
			}
			return nil
		} else {
			c.logger.Sugar().With(XRequestIdHeaderKey, requestId).Errorf("request failed with status %d", httpResp.StatusCode)
		}

		if i != c.retries {
			c.logger.Sugar().With(XRequestIdHeaderKey, requestId).Warnf("Retrying in %s...", c.retryInterval)
			time.Sleep(c.retryInterval)
		}
	}

	return fmt.Errorf("request failed after maximum retries")
}
