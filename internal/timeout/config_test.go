package timeout

import (
	"testing"
	"time"
)

func TestGetTimeout(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		expected time.Duration
	}{
		{
			name:     "GetListing endpoint",
			method:   "/listings.v1.ListingsService/GetListing",
			expected: 5 * time.Second,
		},
		{
			name:     "CreateListing endpoint",
			method:   "/listings.v1.ListingsService/CreateListing",
			expected: 10 * time.Second,
		},
		{
			name:     "BatchUpdateStock endpoint",
			method:   "/listings.v1.ListingsService/BatchUpdateStock",
			expected: 20 * time.Second,
		},
		{
			name:     "IncrementProductViews endpoint",
			method:   "/listings.v1.ListingsService/IncrementProductViews",
			expected: 3 * time.Second,
		},
		{
			name:     "unknown endpoint - returns default",
			method:   "/unknown.v1.Service/Method",
			expected: DefaultTimeout,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetTimeout(tt.method)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestSetTimeout(t *testing.T) {
	// Save original config
	original := GetTimeout("/test.Method")

	// Set new timeout
	newTimeout := 99 * time.Second
	SetTimeout("/test.Method", newTimeout)

	// Verify it was set
	result := GetTimeout("/test.Method")
	if result != newTimeout {
		t.Errorf("expected %v, got %v", newTimeout, result)
	}

	// Clean up - restore default
	if original == DefaultTimeout {
		delete(DefaultConfig, "/test.Method")
	} else {
		SetTimeout("/test.Method", original)
	}
}

func TestGetAllTimeouts(t *testing.T) {
	timeouts := GetAllTimeouts()

	// Verify it's not empty
	if len(timeouts) == 0 {
		t.Error("expected non-empty timeout map")
	}

	// Verify it contains expected endpoints
	expectedEndpoints := []string{
		"/listings.v1.ListingsService/GetListing",
		"/listings.v1.ListingsService/CreateListing",
		"/listings.v1.ListingsService/BatchUpdateStock",
	}

	for _, endpoint := range expectedEndpoints {
		if _, ok := timeouts[endpoint]; !ok {
			t.Errorf("expected endpoint %s to be in timeout map", endpoint)
		}
	}

	// Verify it's a copy (mutations don't affect original)
	timeouts["/test.Method"] = 999 * time.Second
	if _, ok := DefaultConfig["/test.Method"]; ok {
		t.Error("modifying returned map should not affect original config")
	}
}

func TestDefaultTimeoutValues(t *testing.T) {
	// Verify all configured timeouts are reasonable (between 1s and 30s)
	for method, cfg := range DefaultConfig {
		if cfg.Timeout < 1*time.Second {
			t.Errorf("timeout for %s is too short: %v", method, cfg.Timeout)
		}
		if cfg.Timeout > 30*time.Second {
			t.Errorf("timeout for %s is too long: %v", method, cfg.Timeout)
		}
	}
}

func TestTimeoutProgression(t *testing.T) {
	// Verify timeout progression makes sense:
	// Simple reads < Writes < Batch operations
	getTimeout := GetTimeout("/listings.v1.ListingsService/GetListing")
	createTimeout := GetTimeout("/listings.v1.ListingsService/CreateListing")
	batchTimeout := GetTimeout("/listings.v1.ListingsService/BatchUpdateStock")

	if getTimeout >= createTimeout {
		t.Error("GET timeout should be less than CREATE timeout")
	}

	if createTimeout >= batchTimeout {
		t.Error("CREATE timeout should be less than BATCH timeout")
	}
}
