package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPClient represents an HTTP client for making requests
type HTTPClient struct {
	client *http.Client
}

// NewHTTPClient creates a new HTTP client with default settings
func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     30 * time.Second,
			},
		},
	}
}

// Post makes a POST request to the specified URL with the given body
func (h *HTTPClient) Post(ctx context.Context, url string, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)

	return h.client.Do(req)
}

// Get makes a GET request to the specified URL
func (h *HTTPClient) Get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	return h.client.Do(req)
}

// Do executes an HTTP request
func (h *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	return h.client.Do(req)
}

// PostJSON makes a POST request with JSON body and returns response body as string
func (h *HTTPClient) PostJSON(ctx context.Context, url string, jsonBody interface{}) (string, error) {
	jsonData, err := json.Marshal(jsonBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}

	resp, err := h.Post(ctx, url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return string(body), nil
}