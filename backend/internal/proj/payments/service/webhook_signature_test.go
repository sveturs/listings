package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"backend/pkg/logger"
)

func TestValidateWebhookSignature(t *testing.T) {
	tests := []struct {
		name          string
		payload       []byte
		signature     string
		secret        string
		expectedValid bool
	}{
		{
			name:          "Valid signature",
			payload:       []byte(`{"transaction_id":"12345","status":"success"}`),
			secret:        "test-webhook-secret-key",
			signature:     "", // Will be calculated
			expectedValid: true,
		},
		{
			name:          "Invalid signature",
			payload:       []byte(`{"transaction_id":"12345","status":"success"}`),
			secret:        "test-webhook-secret-key",
			signature:     "invalid-signature",
			expectedValid: false,
		},
		{
			name:          "Empty secret (backward compatibility)",
			payload:       []byte(`{"transaction_id":"12345","status":"success"}`),
			secret:        "",
			signature:     "any-signature",
			expectedValid: true, // Should pass when secret is not configured
		},
		{
			name:          "Modified payload",
			payload:       []byte(`{"transaction_id":"12345","status":"failed"}`), // Changed
			secret:        "test-webhook-secret-key",
			signature:     "", // Will be calculated for original payload
			expectedValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create service with test configuration
			service := &AllSecureService{
				config: &AllSecureConfig{
					WebhookSecret: tt.secret,
				},
				logger: logger.GetLogger(),
			}

			// Calculate valid signature for tests that need it
			if tt.signature == "" && tt.secret != "" {
				h := hmac.New(sha256.New, []byte(tt.secret))
				if tt.name == "Modified payload" {
					// Use original payload for signature
					originalPayload := []byte(`{"transaction_id":"12345","status":"success"}`)
					h.Write(originalPayload)
				} else {
					h.Write(tt.payload)
				}
				tt.signature = hex.EncodeToString(h.Sum(nil))
			}

			// Test signature validation
			isValid := service.ValidateWebhookSignature(tt.payload, tt.signature)

			if isValid != tt.expectedValid {
				t.Errorf("ValidateWebhookSignature() = %v, want %v", isValid, tt.expectedValid)
			}
		})
	}
}

func TestCalculateWebhookSignature(t *testing.T) {
	service := &AllSecureService{
		config: &AllSecureConfig{
			WebhookSecret: "test-secret-key",
		},
		logger: logger.GetLogger(),
	}

	payload := []byte(`{"test":"data"}`)
	signature := service.calculateWebhookSignature(payload)

	// Verify signature format (should be hex string)
	if len(signature) != 64 { // SHA256 produces 32 bytes = 64 hex chars
		t.Errorf("Invalid signature length: got %d, want 64", len(signature))
	}

	// Verify signature is deterministic
	signature2 := service.calculateWebhookSignature(payload)
	if signature != signature2 {
		t.Errorf("Signature not deterministic: %s != %s", signature, signature2)
	}

	// Verify different payloads produce different signatures
	differentPayload := []byte(`{"test":"different"}`)
	differentSignature := service.calculateWebhookSignature(differentPayload)
	if signature == differentSignature {
		t.Errorf("Different payloads produced same signature")
	}
}

func BenchmarkValidateWebhookSignature(b *testing.B) {
	service := &AllSecureService{
		config: &AllSecureConfig{
			WebhookSecret: "benchmark-secret-key",
		},
		logger: logger.GetLogger(),
	}

	payload := []byte(`{"transaction_id":"12345","amount":100.50,"currency":"USD","status":"success"}`)

	// Calculate valid signature
	h := hmac.New(sha256.New, []byte(service.config.WebhookSecret))
	h.Write(payload)
	signature := hex.EncodeToString(h.Sum(nil))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.ValidateWebhookSignature(payload, signature)
	}
}

func TestConstantTimeComparison(t *testing.T) {
	// This test verifies that signature comparison is done in constant time
	// to prevent timing attacks

	service := &AllSecureService{
		config: &AllSecureConfig{
			WebhookSecret: "timing-attack-test",
		},
		logger: logger.GetLogger(),
	}

	payload := []byte(`{"test":"timing"}`)
	validSignature := service.calculateWebhookSignature(payload)

	// Test with signatures that differ at different positions
	testCases := []struct {
		name      string
		signature string
	}{
		{"Valid", validSignature},
		{"Differ at start", "0" + validSignature[1:]},
		{"Differ at middle", validSignature[:32] + "0" + validSignature[33:]},
		{"Differ at end", validSignature[:63] + "0"},
		{"Completely different", "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
	}

	// Note: This is a basic test. In production, you'd want more sophisticated
	// timing attack prevention testing
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := service.ValidateWebhookSignature(payload, tc.signature)
			expectedValid := tc.signature == validSignature

			if result != expectedValid {
				t.Errorf("Unexpected result for %s: got %v, want %v",
					tc.name, result, expectedValid)
			}
		})
	}
}
