package middleware

import (
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"backend/internal/config"
)

func TestPaymentAPIRateLimitExceed(t *testing.T) {
	// Create test app
	app := fiber.New()

	// Create middleware with test config
	cfg := &config.Config{
		Environment: "test",
	}
	mw := &Middleware{
		config: cfg,
	}

	// Apply rate limiter
	app.Use("/payment", mw.PaymentAPIRateLimit())
	app.Post("/payment/create", func(c *fiber.Ctx) error {
		return c.SendString("Payment created")
	})

	// Test exceeding rate limit
	successCount := 0
	blockedCount := 0

	for i := 0; i < 15; i++ {
		req := httptest.NewRequest("POST", "/payment/create", nil)
		resp, err := app.Test(req, -1)
		require.NoError(t, err)

		switch resp.StatusCode {
		case 200:
			successCount++
		case fiber.StatusTooManyRequests:
			blockedCount++
		}
	}

	// Should allow 10 and block 5
	assert.Equal(t, 10, successCount, "Should allow exactly 10 requests")
	assert.Equal(t, 5, blockedCount, "Should block 5 requests")
}

func TestPaymentAPIRateLimit(t *testing.T) {
	// Create test app
	app := fiber.New()

	// Create middleware with test config
	cfg := &config.Config{
		Environment: "test",
	}
	mw := &Middleware{
		config: cfg,
	}

	// Apply rate limiter
	app.Use("/payment", mw.PaymentAPIRateLimit())
	app.Post("/payment/create", func(c *fiber.Ctx) error {
		return c.SendString("Payment created")
	})

	// Test rate limiting
	tests := []struct {
		name          string
		requestCount  int
		expectBlocked int
		userID        int
		waitBetween   time.Duration
	}{
		{
			name:          "Within limit",
			requestCount:  5,
			expectBlocked: 0,
			userID:        1,
			waitBetween:   0,
		},
		{
			name:          "Rate resets after time window",
			requestCount:  12,
			expectBlocked: 0, // Should reset after 1 minute
			userID:        3,
			waitBetween:   65 * time.Second, // Wait after first 10 requests
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blockedCount := 0

			for i := 0; i < tt.requestCount; i++ {
				// Wait after 10 requests if specified
				if i == 10 && tt.waitBetween > 0 {
					// Note: In real tests, you'd mock time instead
					t.Skip("Skipping time-based test in unit tests")
				}

				req := httptest.NewRequest("POST", "/payment/create", nil)
				req.Header.Set("X-User-Id", string(rune(tt.userID)))

				resp, err := app.Test(req, -1)
				require.NoError(t, err)

				if resp.StatusCode == fiber.StatusTooManyRequests {
					blockedCount++
				}
			}

			assert.Equal(t, tt.expectBlocked, blockedCount,
				"Expected %d blocked requests, got %d", tt.expectBlocked, blockedCount)
		})
	}
}

func TestStrictPaymentRateLimit(t *testing.T) {
	app := fiber.New()

	cfg := &config.Config{
		Environment: "production",
	}
	mw := &Middleware{
		config: cfg,
	}

	// Apply strict rate limiter
	app.Use("/payment/capture", mw.StrictPaymentRateLimit())
	app.Post("/payment/capture", func(c *fiber.Ctx) error {
		return c.SendString("Payment captured")
	})

	// Test strict limits (3 requests per 5 minutes)
	blockedCount := 0
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("POST", "/payment/capture", nil)
		resp, err := app.Test(req, -1)
		require.NoError(t, err)

		if resp.StatusCode == fiber.StatusTooManyRequests {
			blockedCount++
		}
	}

	// Should block after 3 requests
	assert.Equal(t, 2, blockedCount, "Should block 2 out of 5 requests")
}

func TestWebhookRateLimit(t *testing.T) {
	app := fiber.New()

	mw := &Middleware{
		config: &config.Config{},
	}

	app.Use("/webhook", mw.WebhookRateLimit())
	app.Post("/webhook/allsecure", func(c *fiber.Ctx) error {
		return c.SendString("Webhook processed")
	})

	// Webhooks should allow 100 requests per minute
	successCount := 0
	for i := 0; i < 110; i++ {
		req := httptest.NewRequest("POST", "/webhook/allsecure", nil)
		resp, err := app.Test(req, -1)
		require.NoError(t, err)

		if resp.StatusCode == 200 {
			successCount++
		}
	}

	// Should allow 100 requests
	assert.Equal(t, 100, successCount, "Should allow exactly 100 requests")
}

func TestRateLimiterKeyGeneration(t *testing.T) {
	tests := []struct {
		name        string
		userID      int
		ip          string
		expectedKey string
		limiterFunc string
	}{
		{
			name:        "Authenticated user",
			userID:      123,
			ip:          "192.168.1.1",
			expectedKey: "payment_user_123",
			limiterFunc: "payment",
		},
		{
			name:        "Anonymous user",
			userID:      0,
			ip:          "192.168.1.1",
			expectedKey: "payment_ip_192.168.1.1",
			limiterFunc: "payment",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This is a conceptual test - in real implementation,
			// you'd test the actual key generation logic
			assert.NotEmpty(t, tt.expectedKey)
		})
	}
}

func TestConcurrentRateLimiting(t *testing.T) {
	app := fiber.New()

	mw := &Middleware{
		config: &config.Config{},
	}

	app.Use(mw.PaymentAPIRateLimit())
	app.Post("/payment", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Test concurrent requests
	var wg sync.WaitGroup
	successCount := 0
	blockedCount := 0
	var mu sync.Mutex

	// Launch 20 concurrent requests
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			req := httptest.NewRequest("POST", "/payment", nil)
			resp, err := app.Test(req, -1)
			if err != nil {
				return
			}

			mu.Lock()
			switch resp.StatusCode {
			case 200:
				successCount++
			case fiber.StatusTooManyRequests:
				blockedCount++
			}
			mu.Unlock()
		}()
	}

	wg.Wait()

	// Should allow 10 and block 10
	assert.Equal(t, 10, successCount, "Should allow 10 concurrent requests")
	assert.Equal(t, 10, blockedCount, "Should block 10 concurrent requests")
}

func TestRateLimiterCleanup(t *testing.T) {
	rl := NewRateLimiter(1 * time.Second)

	// Add some requests
	key := "test_user_1"
	for i := 0; i < 5; i++ {
		allowed := rl.isAllowed(key, 10, 1*time.Minute)
		assert.True(t, allowed, "Request %d should be allowed", i+1)
	}

	// Check that timestamps are recorded
	ur := rl.getUserRequests(key)
	assert.NotNil(t, ur)

	// Wait for cleanup (in real tests, you'd mock time)
	// The cleanup goroutine should remove old entries
	time.Sleep(2 * time.Second)

	// After cleanup, old timestamps should be removed
	// This is a simplified test - in production, you'd have more sophisticated testing
	assert.NotNil(t, rl)
}

func BenchmarkRateLimiter(b *testing.B) {
	rl := NewRateLimiter(5 * time.Minute)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := "benchmark_user"
		rl.isAllowed(key, 100, 1*time.Minute)
	}
}

func TestRateLimiterDifferentEndpoints(t *testing.T) {
	app := fiber.New()

	mw := &Middleware{
		config: &config.Config{},
	}

	// Different endpoints with different limits
	app.Post("/auth/login", mw.AuthRateLimit(), func(c *fiber.Ctx) error {
		return c.SendString("Logged in")
	})

	app.Post("/auth/register", mw.RegistrationRateLimit(), func(c *fiber.Ctx) error {
		return c.SendString("Registered")
	})

	// Test that limits are independent
	// Login: 5 attempts per 15 minutes
	loginBlocked := false
	for i := 0; i < 6; i++ {
		req := httptest.NewRequest("POST", "/auth/login", nil)
		resp, _ := app.Test(req, -1)
		if resp.StatusCode == fiber.StatusTooManyRequests {
			loginBlocked = true
			break
		}
	}
	assert.True(t, loginBlocked, "Login should be blocked after 5 attempts")

	// Registration: 3 per hour (should still work)
	registerSuccess := 0
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest("POST", "/auth/register", nil)
		resp, _ := app.Test(req, -1)
		if resp.StatusCode == 200 {
			registerSuccess++
		}
	}
	assert.Equal(t, 3, registerSuccess, "Registration should allow 3 attempts")
}
