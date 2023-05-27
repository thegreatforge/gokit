package httpclient

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestNewClient(t *testing.T) {

	logger := zap.NewNop()

	host := "http://example.com"
	service := "example"
	remoteService := "remote"
	timeout := time.Second
	retries := 3
	retryInterval := time.Millisecond

	// Create a new client
	client := NewClient(host, service, remoteService, timeout, retries, retryInterval, logger, nil, nil)

	// Check if the client's properties were set correctly
	if client.host != host {
		t.Errorf("Expected client's host to be '%s', but got '%s'", host, client.host)
	}

	if client.remoteService != remoteService {
		t.Errorf("Expected client's remoteService to be '%s', but got '%s'", remoteService, client.remoteService)
	}

	if client.retries != retries {
		t.Errorf("Expected client's retries to be %d, but got %d", retries, client.retries)
	}

	if client.retryInterval != retryInterval {
		t.Errorf("Expected client's retryInterval to be %s, but got %s", retryInterval*time.Second, client.retryInterval)
	}

}

func TestRegisterClient(t *testing.T) {
	logger := zap.NewNop()

	host := "http://example.com"
	service := "example"
	remoteService := "remote"
	timeout := time.Second
	retries := 3
	retryInterval := time.Millisecond

	// Create a new client
	client := NewClient(host, service, remoteService, timeout, retries, retryInterval, logger, nil, nil)

	RegisterClient(client)

	if _, exists := Clients[service]; !exists {
		t.Errorf("Expected client's service to be '%s', but got '%s'", service, client.service)
	}
	Clients = make(map[string]*Client)
}

func TestGet(t *testing.T) {

	logger := zap.NewNop()

	service := "example"
	remoteService := "remote"
	timeout := 5 * time.Second
	retries := 3
	retryInterval := time.Millisecond

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
		if r.Header.Get(RemoteServiceHeaderKey) != remoteService {
			t.Errorf("Expected request header '%s' to be '%s', but got '%s'", RemoteServiceHeaderKey, remoteService, r.Header.Get(RemoteServiceHeaderKey))
		}

		// Respond with a sample JSON response
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"success"}`))
	}))
	defer server.Close()

	// Create a new client
	client := NewClient(server.URL, service, remoteService, timeout, retries, retryInterval, logger, nil, nil)

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
	logger := zap.NewNop()

	service := "example"
	remoteService := "remote"
	timeout := 1 * time.Second
	retries := 3
	retryInterval := time.Millisecond

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		time.Sleep(2 * time.Second)
		// Respond with a sample JSON response
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"success"}`))
	}))
	defer server.Close()

	// Create a new client
	client := NewClient(server.URL, service, remoteService, timeout, retries, retryInterval, logger, nil, nil)

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
			OverrideTimeout: 5 * time.Second,
		},
		resp)

	if err != nil {
		t.Errorf("Expected no error, but got '%s'", err.Error())
	}
}

func TestGetWithRetries(t *testing.T) {

	logger := zap.NewNop()

	service := "example"
	remoteService := "remote"
	timeout := 5 * time.Second
	retries := 3
	retryInterval := time.Millisecond

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Respond with a sample JSON response
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"message":"error"}`))
	}))
	defer server.Close()

	// Create a new client
	client := NewClient(server.URL, service, remoteService, timeout, retries, retryInterval, logger, nil, nil)

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

	logger := zap.NewNop()

	service := "example"
	remoteService := "remote"
	timeout := 5 * time.Second
	retries := 3
	retryInterval := time.Millisecond

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
		if r.Header.Get(RemoteServiceHeaderKey) != remoteService {
			t.Errorf("Expected request header '%s' to be '%s', but got '%s'", RemoteServiceHeaderKey, remoteService, r.Header.Get(RemoteServiceHeaderKey))
		}
		if r.Body == nil {
			t.Errorf("Expected request body to be set")
		}

		// Respond with a sample JSON response
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"success"}`))
	}))
	defer server.Close()

	// Create a new client
	client := NewClient(server.URL, service, remoteService, timeout, retries, retryInterval, logger, nil, nil)

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

	logger := zap.NewNop()

	service := "example"
	remoteService := "remote"
	timeout := 5 * time.Second
	retries := 3
	retryInterval := time.Millisecond

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
		if r.Header.Get(RemoteServiceHeaderKey) != remoteService {
			t.Errorf("Expected request header '%s' to be '%s', but got '%s'", RemoteServiceHeaderKey, remoteService, r.Header.Get(RemoteServiceHeaderKey))
		}
		if r.Body == nil {
			t.Errorf("Expected request body to be set")
		}

		// Respond with a sample JSON response
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"success"}`))
	}))
	defer server.Close()

	// Create a new client
	client := NewClient(server.URL, service, remoteService, timeout, retries, retryInterval, logger, nil, nil)

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

	logger := zap.NewNop()

	service := "example"
	remoteService := "remote"
	timeout := 5 * time.Second
	retries := 3
	retryInterval := time.Millisecond

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
		if r.Header.Get(RemoteServiceHeaderKey) != remoteService {
			t.Errorf("Expected request header '%s' to be '%s', but got '%s'", RemoteServiceHeaderKey, remoteService, r.Header.Get(RemoteServiceHeaderKey))
		}
		if r.Body == nil {
			t.Errorf("Expected request body to be set")
		}

		// Respond with a sample JSON response
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"success"}`))
	}))
	defer server.Close()

	// Create a new client
	client := NewClient(server.URL, service, remoteService, timeout, retries, retryInterval, logger, nil, nil)

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

	logger := zap.NewNop()

	service := "example"
	remoteService := "remote"
	timeout := 5 * time.Second
	retries := 3
	retryInterval := time.Millisecond

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
		if r.Header.Get(RemoteServiceHeaderKey) != remoteService {
			t.Errorf("Expected request header '%s' to be '%s', but got '%s'", RemoteServiceHeaderKey, remoteService, r.Header.Get(RemoteServiceHeaderKey))
		}
		if r.Body == nil {
			t.Errorf("Expected request body to be set")
		}

		// Respond with a sample JSON response
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"success"}`))
	}))
	defer server.Close()

	// Create a new client
	client := NewClient(server.URL, service, remoteService, timeout, retries, retryInterval, logger, nil, nil)

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
