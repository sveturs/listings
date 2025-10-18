// Package service implements mock authentication for testing
// backend/internal/proj/admin/testing/service/mock_auth_manager.go
package service

import (
	"crypto/rand"
	"crypto/rsa"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
)

// MockAuthManager generates mock JWT tokens for testing without external auth service
type MockAuthManager struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	logger     zerolog.Logger
}

// NewMockAuthManager creates a new mock auth manager
func NewMockAuthManager(logger zerolog.Logger) (*MockAuthManager, error) {
	// Generate RSA key pair for JWT signing
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	return &MockAuthManager{
		privateKey: privateKey,
		publicKey:  &privateKey.PublicKey,
		logger:     logger,
	}, nil
}

// GetToken generates a mock admin JWT token (implements TestAuthManager interface)
func (m *MockAuthManager) GetToken() (string, error) {
	m.logger.Info().Msg("Generating mock admin token for testing")

	// Create JWT claims with admin role
	claims := jwt.MapClaims{
		"user_id": 11, // Test admin user ID
		"email":   "testadmin@test.com",
		"roles":   []string{"admin"},
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(1 * time.Hour).Unix(),
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Sign token with private key
	tokenString, err := token.SignedString(m.privateKey)
	if err != nil {
		m.logger.Error().Err(err).Msg("Failed to sign mock token")
		return "", err
	}

	m.logger.Debug().
		Str("token_preview", tokenString[:50]+"...").
		Msg("Mock admin token generated successfully")

	return tokenString, nil
}

// ClearToken is a no-op for mock auth manager (implements TestAuthManager interface)
func (m *MockAuthManager) ClearToken() {
	// No-op for mock - we generate fresh tokens each time
	m.logger.Debug().Msg("ClearToken called (no-op for mock auth)")
}

// GetPublicKey returns the public key for token verification (if needed)
func (m *MockAuthManager) GetPublicKey() *rsa.PublicKey {
	return m.publicKey
}
