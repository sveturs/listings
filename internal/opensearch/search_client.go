package opensearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	opensearch "github.com/opensearch-project/opensearch-go/v2"
	"github.com/rs/zerolog"
)

// SearchClient handles OpenSearch search operations with circuit breaker and retry logic
type SearchClient struct {
	client         *opensearch.Client
	index          string
	circuitBreaker *CircuitBreaker
	maxRetries     int
	retryDelay     time.Duration
	timeout        time.Duration
	logger         zerolog.Logger
}

// SearchConfig contains configuration for SearchClient
type SearchConfig struct {
	Addresses    []string
	Username     string
	Password     string
	Index        string
	MaxRetries   int           // Default: 3
	RetryDelay   time.Duration // Default: 100ms
	Timeout      time.Duration // Default: 5s
	MaxFailures  int           // Circuit breaker max failures (default: 5)
	ResetTimeout time.Duration // Circuit breaker reset timeout (default: 60s)
}

// NewSearchClient creates a new OpenSearch search client with circuit breaker
func NewSearchClient(cfg *SearchConfig, logger zerolog.Logger) (*SearchClient, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is required")
	}

	// Set defaults
	if cfg.MaxRetries <= 0 {
		cfg.MaxRetries = 3
	}
	if cfg.RetryDelay <= 0 {
		cfg.RetryDelay = 100 * time.Millisecond
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 5 * time.Second
	}
	if cfg.MaxFailures <= 0 {
		cfg.MaxFailures = 5
	}
	if cfg.ResetTimeout <= 0 {
		cfg.ResetTimeout = 60 * time.Second
	}
	if cfg.Index == "" {
		cfg.Index = "marketplace_listings"
	}

	// Create OpenSearch client
	osCfg := opensearch.Config{
		Addresses: cfg.Addresses,
		Username:  cfg.Username,
		Password:  cfg.Password,
	}

	client, err := opensearch.NewClient(osCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create OpenSearch client: %w", err)
	}

	// Test connection
	res, err := client.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to OpenSearch: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("OpenSearch returned error: %s", res.Status())
	}

	logger.Info().
		Strs("addresses", cfg.Addresses).
		Str("index", cfg.Index).
		Int("max_retries", cfg.MaxRetries).
		Dur("timeout", cfg.Timeout).
		Msg("OpenSearch search client initialized")

	// Create circuit breaker
	circuitBreaker := NewCircuitBreaker(cfg.MaxFailures, cfg.ResetTimeout, logger)

	return &SearchClient{
		client:         client,
		index:          cfg.Index,
		circuitBreaker: circuitBreaker,
		maxRetries:     cfg.MaxRetries,
		retryDelay:     cfg.RetryDelay,
		timeout:        cfg.Timeout,
		logger:         logger.With().Str("component", "opensearch_search_client").Logger(),
	}, nil
}

// SearchResponse represents OpenSearch search response
type SearchResponse struct {
	Took int64 `json:"took"` // Time in milliseconds
	Hits struct {
		Total struct {
			Value    int64  `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		Hits []struct {
			ID     string                 `json:"_id"`
			Score  float64                `json:"_score"`
			Source map[string]interface{} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

// Search executes a search query with retry logic and circuit breaker
func (sc *SearchClient) Search(ctx context.Context, query map[string]interface{}) (*SearchResponse, error) {
	start := time.Now()

	// Apply timeout to context
	ctx, cancel := context.WithTimeout(ctx, sc.timeout)
	defer cancel()

	var lastErr error
	var response *SearchResponse

	// Execute with circuit breaker
	err := sc.circuitBreaker.Execute(func() error {
		// Retry loop
		for attempt := 0; attempt <= sc.maxRetries; attempt++ {
			if attempt > 0 {
				// Exponential backoff
				backoff := sc.retryDelay * time.Duration(1<<uint(attempt-1))
				sc.logger.Debug().
					Int("attempt", attempt).
					Dur("backoff", backoff).
					Msg("retrying search request")

				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(backoff):
				}
			}

			// Execute search
			resp, err := sc.executeSearch(ctx, query)
			if err != nil {
				lastErr = err
				sc.logger.Warn().
					Err(err).
					Int("attempt", attempt+1).
					Int("max_retries", sc.maxRetries).
					Msg("search request failed")
				continue
			}

			response = resp
			return nil
		}

		return fmt.Errorf("search failed after %d retries: %w", sc.maxRetries, lastErr)
	})

	duration := time.Since(start)

	if err != nil {
		sc.logger.Error().
			Err(err).
			Dur("duration", duration).
			Str("circuit_state", sc.circuitBreaker.GetState().String()).
			Msg("search request failed")
		return nil, err
	}

	sc.logger.Debug().
		Dur("duration", duration).
		Int64("took_ms", response.Took).
		Int64("total", response.Hits.Total.Value).
		Msg("search completed successfully")

	return response, nil
}

// executeSearch executes a single search request
func (sc *SearchClient) executeSearch(ctx context.Context, query map[string]interface{}) (*SearchResponse, error) {
	// Marshal query to JSON
	queryJSON, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query: %w", err)
	}

	// Execute search
	res, err := sc.client.Search(
		sc.client.Search.WithContext(ctx),
		sc.client.Search.WithIndex(sc.index),
		sc.client.Search.WithBody(bytes.NewReader(queryJSON)),
		sc.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, fmt.Errorf("search request failed: %w", err)
	}
	defer res.Body.Close()

	// Check for HTTP errors
	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("search returned error [%s]: %s", res.Status(), string(bodyBytes))
	}

	// Parse response
	var searchResp SearchResponse
	if err := json.NewDecoder(res.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &searchResp, nil
}

// HealthCheck performs a health check on OpenSearch connection
func (sc *SearchClient) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	res, err := sc.client.Ping(sc.client.Ping.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("opensearch ping failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("opensearch health check failed: %s", res.Status())
	}

	return nil
}

// GetCircuitBreakerStats returns circuit breaker statistics
func (sc *SearchClient) GetCircuitBreakerStats() map[string]interface{} {
	return sc.circuitBreaker.GetStats()
}

// ResetCircuitBreaker manually resets the circuit breaker
func (sc *SearchClient) ResetCircuitBreaker() {
	sc.circuitBreaker.Reset()
}

// Close closes the search client (no-op for OpenSearch, kept for interface compatibility)
func (sc *SearchClient) Close() error {
	sc.logger.Info().Msg("search client closed")
	return nil
}
