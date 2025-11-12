package ratelimit

import (
	"time"
)

// EndpointConfig defines rate limiting configuration for a specific endpoint
type EndpointConfig struct {
	// Limit is the maximum number of requests allowed
	Limit int

	// Window is the time window for the limit
	Window time.Duration

	// Identifier defines how to identify the client (IP, UserID, or both)
	Identifier IdentifierType

	// Enabled allows disabling rate limiting for specific endpoints
	Enabled bool
}

// Config holds rate limiting configuration for all endpoints
type Config struct {
	// Endpoints maps full gRPC method names to their rate limit configs
	Endpoints map[string]EndpointConfig

	// DefaultConfig is used when no specific config is found for an endpoint
	DefaultConfig EndpointConfig
}

// NewDefaultConfig creates a production-ready rate limiting configuration
func NewDefaultConfig() *Config {
	return &Config{
		Endpoints: map[string]EndpointConfig{
			// Listings endpoints
			"/listings.v1.ListingsService/GetListing": {
				Limit:      200,
				Window:     time.Minute,
				Identifier: ByIP,
				Enabled:    true,
			},
			"/listings.v1.ListingsService/CreateListing": {
				Limit:      50,
				Window:     time.Minute,
				Identifier: ByUserID,
				Enabled:    true,
			},
			"/listings.v1.ListingsService/UpdateListing": {
				Limit:      50,
				Window:     time.Minute,
				Identifier: ByUserID,
				Enabled:    true,
			},
			"/listings.v1.ListingsService/DeleteListing": {
				Limit:      20,
				Window:     time.Minute,
				Identifier: ByUserID,
				Enabled:    true,
			},
			"/listings.v1.ListingsService/SearchListings": {
				Limit:      300,
				Window:     time.Minute,
				Identifier: ByIP,
				Enabled:    true,
			},
			"/listings.v1.ListingsService/ListListings": {
				Limit:      200,
				Window:     time.Minute,
				Identifier: ByIP,
				Enabled:    true,
			},

			// Inventory endpoints (from the spec)
			"/inventory.InventoryService/IncrementProductViews": {
				Limit:      100,
				Window:     time.Minute,
				Identifier: ByIP,
				Enabled:    true,
			},
			"/inventory.InventoryService/GetProductStats": {
				Limit:      300,
				Window:     time.Minute,
				Identifier: ByIP,
				Enabled:    true,
			},
			"/inventory.InventoryService/RecordInventoryMovement": {
				Limit:      50,
				Window:     time.Minute,
				Identifier: ByUserID,
				Enabled:    true,
			},
			"/inventory.InventoryService/BatchUpdateStock": {
				Limit:      20,
				Window:     time.Minute,
				Identifier: ByUserID,
				Enabled:    true,
			},
			"/inventory.InventoryService/GetInventoryStatus": {
				Limit:      200,
				Window:     time.Minute,
				Identifier: ByIP,
				Enabled:    true,
			},
		},
		DefaultConfig: EndpointConfig{
			Limit:      100,
			Window:     time.Minute,
			Identifier: ByIP,
			Enabled:    true,
		},
	}
}

// GetEndpointConfig returns the rate limit configuration for a specific endpoint
func (c *Config) GetEndpointConfig(method string) EndpointConfig {
	if config, ok := c.Endpoints[method]; ok {
		return config
	}
	return c.DefaultConfig
}

// IsEnabled checks if rate limiting is enabled for a specific endpoint
func (c *Config) IsEnabled(method string) bool {
	config := c.GetEndpointConfig(method)
	return config.Enabled
}

// DisableEndpoint disables rate limiting for a specific endpoint
func (c *Config) DisableEndpoint(method string) {
	if config, ok := c.Endpoints[method]; ok {
		config.Enabled = false
		c.Endpoints[method] = config
	}
}

// EnableEndpoint enables rate limiting for a specific endpoint
func (c *Config) EnableEndpoint(method string) {
	if config, ok := c.Endpoints[method]; ok {
		config.Enabled = true
		c.Endpoints[method] = config
	}
}
