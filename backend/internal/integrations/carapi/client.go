package carapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"backend/internal/logger"
	"github.com/go-redis/redis/v8"
)

const (
	defaultTimeout = 30 * time.Second
	cacheTTL       = 24 * time.Hour
)

// Client represents CarAPI client
type Client struct {
	token      string
	baseURL    string
	httpClient *http.Client
	cache      *redis.Client
	rateLimit  int
}

// NewClient creates new CarAPI client
func NewClient(token string, cache *redis.Client) *Client {
	return &Client{
		token:   token,
		baseURL: "https://carapi.app/api",
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
		cache:     cache,
		rateLimit: 1500, // Base plan limit
	}
}

// Make represents a car make
type Make struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Model represents a car model
type Model struct {
	ID     int    `json:"id"`
	MakeID int    `json:"make_id"`
	Name   string `json:"name"`
}

// Trim represents a car trim/configuration
type Trim struct {
	ID          int    `json:"id"`
	ModelID     int    `json:"model_id"`
	Year        int    `json:"year"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// VINDecodeResponse represents VIN decode response
type VINDecodeResponse struct {
	Make        string `json:"make"`
	Model       string `json:"model"`
	ModelYear   int    `json:"model_year"`
	Trim        string `json:"trim"`
	BodyType    string `json:"body_type"`
	Engine      string `json:"engine"`
	Drivetrain  string `json:"drivetrain"`
	Transmission string `json:"transmission"`
}

// GetMakes fetches all car makes
func (c *Client) GetMakes(ctx context.Context) ([]Make, error) {
	// Try cache first
	cacheKey := "carapi:makes:all"
	cached, err := c.cache.Get(ctx, cacheKey).Result()
	if err == nil {
		var makes []Make
		if err := json.Unmarshal([]byte(cached), &makes); err == nil {
			logger.Info().Msg("CarAPI: returning makes from cache")
			return makes, nil
		}
	}

	// Fetch from API
	resp, err := c.doRequest(ctx, "GET", "/makes", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch makes: %w", err)
	}
	defer resp.Body.Close()

	var makes []Make
	if err := json.NewDecoder(resp.Body).Decode(&makes); err != nil {
		return nil, fmt.Errorf("failed to decode makes: %w", err)
	}

	// Cache the result
	data, _ := json.Marshal(makes)
	c.cache.Set(ctx, cacheKey, data, cacheTTL)

	logger.Info().Int("count", len(makes)).Msg("CarAPI: fetched makes from API")
	return makes, nil
}

// GetModels fetches models for a specific make
func (c *Client) GetModels(ctx context.Context, makeID int) ([]Model, error) {
	cacheKey := fmt.Sprintf("carapi:models:make:%d", makeID)
	cached, err := c.cache.Get(ctx, cacheKey).Result()
	if err == nil {
		var models []Model
		if err := json.Unmarshal([]byte(cached), &models); err == nil {
			return models, nil
		}
	}

	endpoint := fmt.Sprintf("/models?make_id=%d", makeID)
	resp, err := c.doRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch models: %w", err)
	}
	defer resp.Body.Close()

	var models []Model
	if err := json.NewDecoder(resp.Body).Decode(&models); err != nil {
		return nil, fmt.Errorf("failed to decode models: %w", err)
	}

	// Cache the result
	data, _ := json.Marshal(models)
	c.cache.Set(ctx, cacheKey, data, cacheTTL)

	return models, nil
}

// GetTrims fetches trims for a specific model and year
func (c *Client) GetTrims(ctx context.Context, modelID int, year int) ([]Trim, error) {
	cacheKey := fmt.Sprintf("carapi:trims:model:%d:year:%d", modelID, year)
	cached, err := c.cache.Get(ctx, cacheKey).Result()
	if err == nil {
		var trims []Trim
		if err := json.Unmarshal([]byte(cached), &trims); err == nil {
			return trims, nil
		}
	}

	endpoint := fmt.Sprintf("/trims?model_id=%d&year=%d", modelID, year)
	resp, err := c.doRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch trims: %w", err)
	}
	defer resp.Body.Close()

	var trims []Trim
	if err := json.NewDecoder(resp.Body).Decode(&trims); err != nil {
		return nil, fmt.Errorf("failed to decode trims: %w", err)
	}

	// Cache the result
	data, _ := json.Marshal(trims)
	c.cache.Set(ctx, cacheKey, data, 7*24*time.Hour) // Cache trims for 7 days

	return trims, nil
}

// DecodeVIN decodes a VIN number
func (c *Client) DecodeVIN(ctx context.Context, vin string) (*VINDecodeResponse, error) {
	// VIN decode results are cached for 30 days
	cacheKey := fmt.Sprintf("carapi:vin:%s", vin)
	cached, err := c.cache.Get(ctx, cacheKey).Result()
	if err == nil {
		var result VINDecodeResponse
		if err := json.Unmarshal([]byte(cached), &result); err == nil {
			return &result, nil
		}
	}

	endpoint := fmt.Sprintf("/vin/decode?vin=%s", vin)
	resp, err := c.doRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decode VIN: %w", err)
	}
	defer resp.Body.Close()

	var result VINDecodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode VIN response: %w", err)
	}

	// Cache the result for 30 days
	data, _ := json.Marshal(result)
	c.cache.Set(ctx, cacheKey, data, 30*24*time.Hour)

	return &result, nil
}

// doRequest performs HTTP request to CarAPI
func (c *Client) doRequest(ctx context.Context, method, endpoint string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+endpoint, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, fmt.Errorf("rate limit exceeded")
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("API error: status=%d, body=%s", resp.StatusCode, string(body))
	}

	return resp, nil
}