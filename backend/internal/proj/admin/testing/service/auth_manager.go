// Package service implements business logic for testing module
// backend/internal/proj/admin/testing/service/auth_manager.go
package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

const (
	tokenCacheDuration = 1 * time.Hour // Cache token for 1 hour
)

// TestAuthManager manages authentication for automated tests
type TestAuthManager struct {
	backendURL string // Local backend URL that proxies to auth service
	email      string
	password   string
	logger     zerolog.Logger

	mu          sync.RWMutex
	cachedToken string
	tokenExpiry time.Time
}

// NewTestAuthManager creates new test auth manager instance
func NewTestAuthManager(backendURL, email, password string, logger zerolog.Logger) *TestAuthManager {
	return &TestAuthManager{
		backendURL: backendURL,
		email:      email,
		password:   password,
		logger:     logger.With().Str("component", "test_auth_manager").Logger(),
	}
}

// LoginResponse represents auth service login response
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

// GetToken returns valid admin token (cached or fresh)
func (m *TestAuthManager) GetToken() (string, error) {
	m.mu.RLock()
	// Check if we have valid cached token
	if m.cachedToken != "" && time.Now().Before(m.tokenExpiry) {
		token := m.cachedToken
		m.mu.RUnlock()
		m.logger.Debug().Msg("Using cached admin token")
		return token, nil
	}
	m.mu.RUnlock()

	// Need to get fresh token
	m.mu.Lock()
	defer m.mu.Unlock()

	// Double-check in case another goroutine already got token
	if m.cachedToken != "" && time.Now().Before(m.tokenExpiry) {
		m.logger.Debug().Msg("Using cached admin token (after lock)")
		return m.cachedToken, nil
	}

	// Login to get new token
	m.logger.Info().Str("email", m.email).Msg("Getting fresh admin token")

	loginReq := map[string]string{
		"email":    m.email,
		"password": m.password,
	}

	reqBody, err := json.Marshal(loginReq)
	if err != nil {
		m.logger.Error().Err(err).Msg("Failed to marshal login request")
		return "", fmt.Errorf("failed to marshal login request: %w", err)
	}

	// Debug: log the actual request body (with password partially masked)
	m.logger.Debug().
		Str("email", m.email).
		Str("password_length", fmt.Sprintf("%d", len(m.password))).
		Str("request_body", string(reqBody)).
		Msg("Login request body prepared")

	loginURL := fmt.Sprintf("%s/api/v1/auth/login", m.backendURL)
	resp, err := http.Post(loginURL, "application/json", bytes.NewBuffer(reqBody)) //nolint:gosec // G107: loginURL constructed from trusted backendURL config
	if err != nil {
		m.logger.Error().Err(err).Str("url", loginURL).Msg("Failed to send login request")
		return "", fmt.Errorf("failed to send login request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		// Read response body for debugging
		bodyBytes, _ := io.ReadAll(resp.Body)
		m.logger.Error().
			Int("status", resp.StatusCode).
			Str("url", loginURL).
			Str("response_body", string(bodyBytes)).
			Msg("Login request failed")
		return "", fmt.Errorf("login request failed with status %d", resp.StatusCode)
	}

	var loginResp LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		m.logger.Error().Err(err).Msg("Failed to decode login response")
		return "", fmt.Errorf("failed to decode login response: %w", err)
	}

	if loginResp.AccessToken == "" {
		m.logger.Error().Msg("Empty access token in response")
		return "", fmt.Errorf("empty access token in response")
	}

	// Cache token
	m.cachedToken = loginResp.AccessToken
	m.tokenExpiry = time.Now().Add(tokenCacheDuration)

	m.logger.Info().
		Str("email", m.email).
		Time("expiry", m.tokenExpiry).
		Msg("Admin token cached successfully")

	return m.cachedToken, nil
}

// ClearToken clears cached token (useful for testing or after auth errors)
func (m *TestAuthManager) ClearToken() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.cachedToken = ""
	m.tokenExpiry = time.Time{}
	m.logger.Info().Msg("Admin token cache cleared")
}
