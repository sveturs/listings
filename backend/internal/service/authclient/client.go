package authclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
}

func NewClient(baseURL string) *Client {
	if baseURL == "" {
		baseURL = "http://localhost:28080"
	}

	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: baseURL,
	}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	User         User   `json:"user"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type RegisterResponse struct {
	Message string `json:"message"`
	User    User   `json:"user"`
}

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type ValidateResponse struct {
	Valid     bool   `json:"valid"`
	UserID    int    `json:"user_id,omitempty"`
	Email     string `json:"email,omitempty"`
	ExpiresAt int64  `json:"expires_at,omitempty"`
}

func (c *Client) Login(ctx context.Context, email, password string) (*LoginResponse, error) {
	req := LoginRequest{
		Email:    email,
		Password: password,
	}

	var resp LoginResponse
	if err := c.doRequest(ctx, "POST", "/api/v1/auth/login", req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) Register(ctx context.Context, email, password, name string) (*RegisterResponse, error) {
	req := RegisterRequest{
		Email:    email,
		Password: password,
		Name:     name,
	}

	var resp RegisterResponse
	if err := c.doRequest(ctx, "POST", "/api/v1/auth/register", req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) RefreshToken(ctx context.Context, refreshToken string) (*RefreshResponse, error) {
	req := RefreshRequest{
		RefreshToken: refreshToken,
	}

	var resp RefreshResponse
	if err := c.doRequest(ctx, "POST", "/api/v1/auth/refresh", req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) ValidateToken(ctx context.Context, accessToken string) (*ValidateResponse, error) {
	var resp ValidateResponse
	if err := c.doRequestWithAuth(ctx, "GET", "/api/v1/auth/validate", nil, &resp, accessToken); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) Logout(ctx context.Context, accessToken string) error {
	return c.doRequestWithAuth(ctx, "POST", "/api/v1/auth/logout", nil, nil, accessToken)
}

func (c *Client) LogoutAll(ctx context.Context, accessToken string) error {
	return c.doRequestWithAuth(ctx, "POST", "/api/v1/auth/logout-all", nil, nil, accessToken)
}

func (c *Client) GetOAuthURL(provider string) string {
	return fmt.Sprintf("%s/api/v1/auth/oauth/%s", c.baseURL, provider)
}

func (c *Client) ProxyRequest(w http.ResponseWriter, r *http.Request, path string) error {
	targetURL := c.baseURL + path

	proxyReq, err := http.NewRequestWithContext(r.Context(), r.Method, targetURL, r.Body)
	if err != nil {
		return err
	}

	for key, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(key, value)
		}
	}

	proxyReq.Header.Set("X-Forwarded-For", r.RemoteAddr)
	proxyReq.Header.Set("X-Forwarded-Host", r.Host)
	proxyReq.Header.Set("X-Forwarded-Proto", "http")

	resp, err := c.httpClient.Do(proxyReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	return err
}

func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}, response interface{}) error {
	return c.doRequestWithAuth(ctx, method, path, body, response, "")
}

func (c *Client) doRequestWithAuth(ctx context.Context, method, path string, body interface{}, response interface{}, accessToken string) error {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		var errorResp struct {
			Error   string `json:"error"`
			Message string `json:"message"`
		}
		if err := json.Unmarshal(respBody, &errorResp); err == nil {
			if errorResp.Error != "" {
				return fmt.Errorf("auth service error: %s", errorResp.Error)
			}
			if errorResp.Message != "" {
				return fmt.Errorf("auth service error: %s", errorResp.Message)
			}
		}
		return fmt.Errorf("auth service returned status %d: %s", resp.StatusCode, string(respBody))
	}

	if response != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, response); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}