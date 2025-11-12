package timeout

import (
	"time"
)

// EndpointConfig contains timeout configuration for a specific endpoint.
type EndpointConfig struct {
	Timeout time.Duration
}

// DefaultConfig defines timeouts for all gRPC endpoints.
// Timeouts are based on SLA requirements (p95 < 100ms for most endpoints)
// and expected operation complexity.
var DefaultConfig = map[string]EndpointConfig{
	// Listing CRUD operations
	"/listings.v1.ListingsService/GetListing": {
		Timeout: 5 * time.Second, // Simple DB query
	},
	"/listings.v1.ListingsService/CreateListing": {
		Timeout: 10 * time.Second, // DB write + validation
	},
	"/listings.v1.ListingsService/UpdateListing": {
		Timeout: 10 * time.Second, // DB write + potential cascade
	},
	"/listings.v1.ListingsService/DeleteListing": {
		Timeout: 15 * time.Second, // Cascade deletes
	},

	// Search and list operations
	"/listings.v1.ListingsService/SearchListings": {
		Timeout: 8 * time.Second, // OpenSearch query
	},
	"/listings.v1.ListingsService/ListListings": {
		Timeout: 5 * time.Second, // Paginated query
	},

	// Inventory operations
	"/listings.v1.ListingsService/IncrementProductViews": {
		Timeout: 3 * time.Second, // Simple counter update
	},
	"/listings.v1.ListingsService/GetProductStats": {
		Timeout: 5 * time.Second, // Aggregation query
	},
	"/listings.v1.ListingsService/RecordInventoryMovement": {
		Timeout: 8 * time.Second, // DB write + audit
	},
	"/listings.v1.ListingsService/BatchUpdateStock": {
		Timeout: 20 * time.Second, // Bulk operation
	},
	"/listings.v1.ListingsService/GetInventoryStatus": {
		Timeout: 5 * time.Second, // Read query
	},
	"/listings.v1.ListingsService/UpdateStock": {
		Timeout: 5 * time.Second, // Single stock update
	},
	"/listings.v1.ListingsService/GetStock": {
		Timeout: 3 * time.Second, // Simple read
	},

	// Additional endpoints
	"/listings.v1.ListingsService/GetListingsByStorefront": {
		Timeout: 5 * time.Second, // Paginated query
	},
	"/listings.v1.ListingsService/GetListingsByIDs": {
		Timeout: 5 * time.Second, // Batch read
	},
	"/listings.v1.ListingsService/UpdateListingStatus": {
		Timeout: 5 * time.Second, // Simple update
	},
}

// DefaultTimeout is used when no specific timeout is configured for an endpoint.
const DefaultTimeout = 30 * time.Second

// GetTimeout returns the configured timeout for a specific gRPC method.
// If no specific configuration exists, returns the default timeout.
func GetTimeout(method string) time.Duration {
	if cfg, ok := DefaultConfig[method]; ok {
		return cfg.Timeout
	}
	return DefaultTimeout
}

// SetTimeout allows runtime configuration of endpoint timeouts.
// This is useful for testing or dynamic tuning.
func SetTimeout(method string, timeout time.Duration) {
	DefaultConfig[method] = EndpointConfig{Timeout: timeout}
}

// GetAllTimeouts returns a copy of all configured timeouts.
// Useful for diagnostics and monitoring.
func GetAllTimeouts() map[string]time.Duration {
	result := make(map[string]time.Duration, len(DefaultConfig))
	for method, cfg := range DefaultConfig {
		result[method] = cfg.Timeout
	}
	return result
}
