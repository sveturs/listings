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

// ============================================================================
// Favorites HTTP Methods
// ============================================================================

// AddToFavorites adds a listing to user's favorites via HTTP.
func (c *HTTPClient) AddToFavorites(ctx context.Context, userID, listingID int64) error {
	url := fmt.Sprintf("%s/api/v1/favorites/%d", c.baseURL, listingID)

	// Body contains user_id for the request
	body, err := json.Marshal(map[string]int64{"user_id": userID})
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	c.addAuthHeaders(req)
	req.Header.Set("Content-Type", "application/json")

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

// RemoveFromFavorites removes a listing from user's favorites via HTTP.
func (c *HTTPClient) RemoveFromFavorites(ctx context.Context, userID, listingID int64) error {
	url := fmt.Sprintf("%s/api/v1/favorites/%d?user_id=%d", c.baseURL, listingID, userID)

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

// GetUserFavorites retrieves list of listing IDs favorited by a user via HTTP.
func (c *HTTPClient) GetUserFavorites(ctx context.Context, userID int64) ([]int64, int, error) {
	url := fmt.Sprintf("%s/api/v1/users/%d/favorites", c.baseURL, userID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	c.addAuthHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if err := c.checkHTTPError(resp); err != nil {
		return nil, 0, err
	}

	var result struct {
		Data struct {
			ListingIDs []int64 `json:"listing_ids"`
			Total      int     `json:"total"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, 0, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Data.ListingIDs, result.Data.Total, nil
}

// IsFavorite checks if a listing is in user's favorites via HTTP.
func (c *HTTPClient) IsFavorite(ctx context.Context, userID, listingID int64) (bool, error) {
	url := fmt.Sprintf("%s/api/v1/favorites/%d/is-favorite?user_id=%d", c.baseURL, listingID, userID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	c.addAuthHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if err := c.checkHTTPError(resp); err != nil {
		return false, err
	}

	var result struct {
		Data struct {
			IsFavorite bool `json:"is_favorite"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Data.IsFavorite, nil
}

// GetFavoritedUsers retrieves list of user IDs who favorited a listing via HTTP.
func (c *HTTPClient) GetFavoritedUsers(ctx context.Context, listingID int64) ([]int64, error) {
	url := fmt.Sprintf("%s/api/v1/favorites/%d/users", c.baseURL, listingID)

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
		Data struct {
			UserIDs []int64 `json:"user_ids"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Data.UserIDs, nil
}

// ============================================================================
// Image Management HTTP Methods
// ============================================================================

// DeleteListingImage removes an image from a listing via HTTP.
func (c *HTTPClient) DeleteListingImage(ctx context.Context, imageID int64) error {
	url := fmt.Sprintf("%s/api/v1/images/%d", c.baseURL, imageID)

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

// ReorderListingImages updates display order for multiple images via HTTP.
func (c *HTTPClient) ReorderListingImages(ctx context.Context, listingID int64, imageOrders []ImageOrder) error {
	url := fmt.Sprintf("%s/api/v1/listings/%d/images/reorder", c.baseURL, listingID)

	// Marshal request body
	body, err := json.Marshal(map[string]interface{}{
		"image_orders": imageOrders,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	c.addAuthHeaders(req)
	req.Header.Set("Content-Type", "application/json")

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
