package httpclient

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.uber.org/zap"
)

func getClientConfig(host string, timeout, retryInterval time.Duration) ClientConfig {
	if host == "" {
		host = "http://example.com"
	}
	return ClientConfig{
		Host:             host,
		Service:          "example",
		RemoteService:    "remote",
		Timeout:          timeout,
		Retries:          3,
		RetryInterval:    retryInterval,
		Logger:           zap.NewNop(),
		OnRetry:          nil,
		OnClientResponse: nil,
	}
}

func TestNewClient(t *testing.T) {
	config := getClientConfig("", time.Second, time.Second)
	// Create a new client
	client := NewClient(config)

	// Check if the client's properties were set correctly
	if client.host != config.Host {
		t.Errorf("Expected client's host to be '%s', but got '%s'", config.Host, client.host)
	}

	if client.remoteService != config.RemoteService {
		t.Errorf("Expected client's remoteService to be '%s', but got '%s'", config.RemoteService, client.remoteService)
	}

	if client.retries != config.Retries {
		t.Errorf("Expected client's retries to be %d, but got %d", config.Retries, client.retries)
	}

	if client.retryInterval != config.RetryInterval {
		t.Errorf("Expected client's retryInterval to be %s, but got %s", config.RetryInterval, client.retryInterval)
	}

}

func TestRegisterClient(t *testing.T) {
	config := getClientConfig("", time.Second, time.Second)
	// Create a new client
	client := NewClient(config)

	RegisterClient(client)

	if _, exists := Clients[config.Service]; !exists {
		t.Errorf("Expected client's service to be '%s', but got '%s'", config.Service, client.service)
	}
	Clients = make(map[string]*Client)
}

func TestGet(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// test the request
		if r.URL.Path != "/api/resource" {
			t.Errorf("Expected request path to be '%s', but got '%s'", "/api/resource", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("Expected request method to be '%s', but got '%s'", http.MethodGet, r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected request content-type to be '%s', but got '%s'", "application/json", r.Header.Get("Content-Type"))
		}
		if r.Header.Get(XRequestIdHeaderKey) == "" {
			t.Errorf("Expected request header '%s' to be set", XRequestIdHeaderKey)
		}

		// Respond with a sample JSON response
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"success"}`))
	}))
	defer server.Close()

	config := getClientConfig(server.URL, time.Second, time.Second)
	// Create a new client
	client := NewClient(config)

	type responseBody struct {
		Message string `json:"message"`
	}

	resp := &Response{
		Body: &responseBody{},
	}

	err := client.Get(
		context.Background(),
		Request{
			Path: "/api/resource",
		},
		resp)

	if err != nil {
		t.Errorf("Expected no error, but got '%s'", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected response status code to be '%d', but got '%d'", http.StatusOK, resp.StatusCode)
	}

	if resp.Body.(*responseBody).Message != "success" {
		t.Errorf("Expected response body to be '%s', but got '%s'", "success", resp.Body.(*responseBody).Message)
	}

}

func TestGetTimeoutOveride(t *testing.T) {
	timeout := 5 * time.Second
	retryInterval := 10 * time.Millisecond

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		time.Sleep(2 * time.Second)
		// Respond with a sample JSON response
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"error"}`))
	}))
	defer server.Close()

	config := getClientConfig(server.URL, timeout, retryInterval)
	// Create a new client
	client := NewClient(config)

	type responseBody struct {
		Message string `json:"message"`
	}

	resp := &Response{
		Body: &responseBody{},
	}

	err := client.Get(
		context.Background(),
		Request{
			Path: "/api/resource",
			// Override the timeout
			OverrideTimeout: 1 * time.Second,
		},
		resp)

	if err == nil {
		t.Errorf("Expected error, but got nil")
	}
}

func TestGetWithRetries(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Respond with a sample JSON response
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"message":"error"}`))
	}))
	defer server.Close()

	config := getClientConfig(server.URL, time.Second, time.Second)
	// Create a new client
	client := NewClient(config)

	type responseBody struct {
		Message string `json:"message"`
	}

	resp := &Response{
		Body: &responseBody{},
	}

	err := client.Get(
		context.Background(),
		Request{
			Path: "/api/resource",
		},
		resp)

	if err == nil {
		t.Errorf("Expected error, but got nil")
	}

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected response status code to be '%d', but got '%d'", http.StatusInternalServerError, resp.StatusCode)
	}

	if resp.Body.(*responseBody).Message != "error" {
		t.Errorf("Expected response body to be '%s', but got '%s'", "error", resp.Body.(*responseBody).Message)
	}
}

func TestPut(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// test the request
		if r.URL.Path != "/api/resource" {
			t.Errorf("Expected request path to be '%s', but got '%s'", "/api/resource", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("Expected request method to be '%s', but got '%s'", http.MethodPut, r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected request content-type to be '%s', but got '%s'", "application/json", r.Header.Get("Content-Type"))
		}
		if r.Header.Get(XRequestIdHeaderKey) == "" {
			t.Errorf("Expected request header '%s' to be set", XRequestIdHeaderKey)
		}
		if r.Body == nil {
			t.Errorf("Expected request body to be set")
		}

		// Respond with a sample JSON response
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"success"}`))
	}))
	defer server.Close()

	config := getClientConfig(server.URL, time.Second, time.Second)
	// Create a new client
	client := NewClient(config)

	type requestBody struct {
		Message string `json:"message"`
	}

	type responseBody struct {
		Message string `json:"message"`
	}

	resp := &Response{
		Body: &responseBody{},
	}

	err := client.Put(
		context.Background(),
		Request{
			Path: "/api/resource",
			Body: &requestBody{
				Message: "hello",
			},
		},
		resp)

	if err != nil {
		t.Errorf("Expected no error, but got '%s'", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected response status code to be '%d', but got '%d'", http.StatusOK, resp.StatusCode)
	}

	if resp.Body.(*responseBody).Message != "success" {
		t.Errorf("Expected response body to be '%s', but got '%s'", "success", resp.Body.(*responseBody).Message)
	}
}

func TestPost(t *testing.T) {

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// test the request
		if r.URL.Path != "/api/resource" {
			t.Errorf("Expected request path to be '%s', but got '%s'", "/api/resource", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("Expected request method to be '%s', but got '%s'", http.MethodPost, r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected request content-type to be '%s', but got '%s'", "application/json", r.Header.Get("Content-Type"))
		}
		if r.Header.Get(XRequestIdHeaderKey) == "" {
			t.Errorf("Expected request header '%s' to be set", XRequestIdHeaderKey)
		}
		if r.Body == nil {
			t.Errorf("Expected request body to be set")
		}

		// Respond with a sample JSON response
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"success"}`))
	}))
	defer server.Close()

	config := getClientConfig(server.URL, time.Second, time.Second)
	// Create a new client
	client := NewClient(config)

	type requestBody struct {
		Message string `json:"message"`
	}

	type responseBody struct {
		Message string `json:"message"`
	}

	resp := &Response{
		Body: &responseBody{},
	}

	err := client.Post(
		context.Background(),
		Request{
			Path: "/api/resource",
			Body: &requestBody{
				Message: "hello",
			},
		},
		resp)

	if err != nil {
		t.Errorf("Expected no error, but got '%s'", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected response status code to be '%d', but got '%d'", http.StatusOK, resp.StatusCode)
	}

	if resp.Body.(*responseBody).Message != "success" {
		t.Errorf("Expected response body to be '%s', but got '%s'", "success", resp.Body.(*responseBody).Message)
	}
}

func TestDelete(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// test the request
		if r.URL.Path != "/api/resource" {
			t.Errorf("Expected request path to be '%s', but got '%s'", "/api/resource", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("Expected request method to be '%s', but got '%s'", http.MethodDelete, r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected request content-type to be '%s', but got '%s'", "application/json", r.Header.Get("Content-Type"))
		}
		if r.Header.Get(XRequestIdHeaderKey) == "" {
			t.Errorf("Expected request header '%s' to be set", XRequestIdHeaderKey)
		}
		if r.Body == nil {
			t.Errorf("Expected request body to be set")
		}

		// Respond with a sample JSON response
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"success"}`))
	}))
	defer server.Close()

	config := getClientConfig(server.URL, time.Second, time.Second)
	// Create a new client
	client := NewClient(config)

	type requestBody struct {
		Message string `json:"message"`
	}

	type responseBody struct {
		Message string `json:"message"`
	}

	resp := &Response{
		Body: &responseBody{},
	}

	err := client.Delete(
		context.Background(),
		Request{
			Path: "/api/resource",
			Body: &requestBody{
				Message: "hello",
			},
		},
		resp)

	if err != nil {
		t.Errorf("Expected no error, but got '%s'", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected response status code to be '%d', but got '%d'", http.StatusOK, resp.StatusCode)
	}
}

func TestPatch(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// test the request
		if r.URL.Path != "/api/resource" {
			t.Errorf("Expected request path to be '%s', but got '%s'", "/api/resource", r.URL.Path)
		}
		if r.Method != http.MethodPatch {
			t.Errorf("Expected request method to be '%s', but got '%s'", http.MethodPatch, r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected request content-type to be '%s', but got '%s'", "application/json", r.Header.Get("Content-Type"))
		}
		if r.Header.Get(XRequestIdHeaderKey) == "" {
			t.Errorf("Expected request header '%s' to be set", XRequestIdHeaderKey)
		}
		if r.Body == nil {
			t.Errorf("Expected request body to be set")
		}

		// Respond with a sample JSON response
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"success"}`))
	}))
	defer server.Close()

	config := getClientConfig(server.URL, time.Second, time.Second)
	// Create a new client
	client := NewClient(config)

	type requestBody struct {
		Message string `json:"message"`
	}

	type responseBody struct {
		Message string `json:"message"`
	}

	resp := &Response{
		Body: &responseBody{},
	}

	err := client.Patch(
		context.Background(),
		Request{
			Path: "/api/resource",
			Body: &requestBody{
				Message: "hello",
			},
		},
		resp)

	if err != nil {
		t.Errorf("Expected no error, but got '%s'", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected response status code to be '%d', but got '%d'", http.StatusOK, resp.StatusCode)
	}

	if resp.Body.(*responseBody).Message != "success" {
		t.Errorf("Expected response body to be '%s', but got '%s'", "success", resp.Body.(*responseBody).Message)
	}
}
