package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

// HTTPClient provides HTTP REST client for listings service.
type HTTPClient struct {
	baseURL    string
	authToken  string
	httpClient *http.Client
	logger     zerolog.Logger
}

// NewHTTPClient creates a new HTTP client for listings service.
func NewHTTPClient(baseURL, authToken string, timeout time.Duration, logger zerolog.Logger) (*HTTPClient, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("baseURL cannot be empty")
	}

	return &HTTPClient{
		baseURL:   baseURL,
		authToken: authToken,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		logger: logger,
	}, nil
}

// GetListing retrieves a listing via HTTP.
func (c *HTTPClient) GetListing(ctx context.Context, id int64) (*Listing, error) {
	url := fmt.Sprintf("%s/api/v1/listings/%d", c.baseURL, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.addAuthHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if err := c.checkHTTPError(resp); err != nil {
		return nil, err
	}

	var result struct {
		Data *Listing `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Data, nil
}

// CreateListing creates a listing via HTTP.
func (c *HTTPClient) CreateListing(ctx context.Context, req *CreateListingRequest) (*Listing, error) {
	url := fmt.Sprintf("%s/api/v1/listings", c.baseURL)

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.addAuthHeaders(httpReq)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if err := c.checkHTTPError(resp); err != nil {
		return nil, err
	}

	var result struct {
		Data *Listing `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Data, nil
}

// UpdateListing updates a listing via HTTP.
func (c *HTTPClient) UpdateListing(ctx context.Context, id int64, req *UpdateListingRequest) (*Listing, error) {
	url := fmt.Sprintf("%s/api/v1/listings/%d", c.baseURL, id)

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.addAuthHeaders(httpReq)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if err := c.checkHTTPError(resp); err != nil {
		return nil, err
	}

	var result struct {
		Data *Listing `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Data, nil
}

// DeleteListing deletes a listing via HTTP.
func (c *HTTPClient) DeleteListing(ctx context.Context, id int64) error {
	url := fmt.Sprintf("%s/api/v1/listings/%d", c.baseURL, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	c.addAuthHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if err := c.checkHTTPError(resp); err != nil {
		return err
	}

	return nil
}

// SearchListings searches listings via HTTP.
func (c *HTTPClient) SearchListings(ctx context.Context, req *SearchListingsRequest) (*SearchListingsResponse, error) {
	url := fmt.Sprintf("%s/api/v1/listings/search?query=%s&limit=%d&offset=%d",
		c.baseURL, req.Query, req.Limit, req.Offset)

	if req.CategoryID != nil {
		url += fmt.Sprintf("&category_id=%d", *req.CategoryID)
	}

	if req.MinPrice != nil {
		url += fmt.Sprintf("&min_price=%.2f", *req.MinPrice)
	}

	if req.MaxPrice != nil {
		url += fmt.Sprintf("&max_price=%.2f", *req.MaxPrice)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.addAuthHeaders(httpReq)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if err := c.checkHTTPError(resp); err != nil {
		return nil, err
	}

	var result struct {
		Data struct {
			Listings []*Listing `json:"listings"`
			Total    int32      `json:"total"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &SearchListingsResponse{
		Listings: result.Data.Listings,
		Total:    result.Data.Total,
	}, nil
}

// ListListings lists listings via HTTP.
func (c *HTTPClient) ListListings(ctx context.Context, req *ListListingsRequest) (*ListListingsResponse, error) {
	url := fmt.Sprintf("%s/api/v1/listings?limit=%d&offset=%d",
		c.baseURL, req.Limit, req.Offset)

	if req.UserID != nil {
		url += fmt.Sprintf("&user_id=%d", *req.UserID)
	}

	if req.StorefrontID != nil {
		url += fmt.Sprintf("&storefront_id=%d", *req.StorefrontID)
	}

	if req.CategoryID != nil {
		url += fmt.Sprintf("&category_id=%d", *req.CategoryID)
	}

	if req.Status != nil {
		url += fmt.Sprintf("&status=%s", *req.Status)
	}

	if req.MinPrice != nil {
		url += fmt.Sprintf("&min_price=%.2f", *req.MinPrice)
	}

	if req.MaxPrice != nil {
		url += fmt.Sprintf("&max_price=%.2f", *req.MaxPrice)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.addAuthHeaders(httpReq)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if err := c.checkHTTPError(resp); err != nil {
		return nil, err
	}

	var result struct {
		Data struct {
			Listings []*Listing `json:"listings"`
			Total    int32      `json:"total"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &ListListingsResponse{
		Listings: result.Data.Listings,
		Total:    result.Data.Total,
	}, nil
}

// addAuthHeaders adds authentication headers to the request.
func (c *HTTPClient) addAuthHeaders(req *http.Request) {
	if c.authToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.authToken)
	}
	req.Header.Set("User-Agent", "listings-service-client/1.0")
}

// checkHTTPError converts HTTP status codes to appropriate errors.
func (c *HTTPClient) checkHTTPError(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	// Try to read error message from response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("HTTP %d: failed to read error response", resp.StatusCode)
	}

	var errorResp struct {
		Error   string `json:"error"`
		Message string `json:"message"`
	}

	if err := json.Unmarshal(bodyBytes, &errorResp); err != nil {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(bodyBytes))
	}

	errorMsg := errorResp.Error
	if errorResp.Message != "" {
		errorMsg = errorResp.Message
	}

	// Map HTTP status codes to our error types
	switch resp.StatusCode {
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusBadRequest:
		return fmt.Errorf("%w: %s", ErrInvalidInput, errorMsg)
	case http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return ErrUnavailable
	default:
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, errorMsg)
	}
}
