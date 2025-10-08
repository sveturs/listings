package postexpress

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

// Client - HTTP клиент для Post Express API
type Client struct {
	config     *Config
	httpClient *http.Client
}

// NewClient создает новый HTTP клиент
func NewClient(config *Config) *Client {
	return &Client{
		config: config,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

// doRequest выполняет HTTP запрос с retry логикой
func (c *Client) doRequest(ctx context.Context, method, endpoint string, body interface{}) ([]byte, error) {
	var lastErr error

	for attempt := 0; attempt <= c.config.RetryAttempts; attempt++ {
		if attempt > 0 {
			// Exponential backoff: 1s, 2s, 4s, 8s
			// Cap the backoff to prevent overflow
			shift := attempt - 1
			if shift > 30 {
				shift = 30 // Cap to prevent overflow
			}
			backoff := time.Duration(1<<shift) * time.Second
			log.Warn().
				Int("attempt", attempt).
				Dur("backoff", backoff).
				Str("endpoint", endpoint).
				Msg("Retrying Post Express API request")

			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff):
			}
		}

		resp, err := c.doSingleRequest(ctx, method, endpoint, body)
		if err == nil {
			return resp, nil
		}

		lastErr = err

		// Не повторяем при ошибках валидации (4xx)
		if isClientError(err) {
			log.Debug().
				Err(err).
				Str("endpoint", endpoint).
				Msg("Client error, not retrying")
			return nil, err
		}
	}

	return nil, fmt.Errorf("failed after %d attempts: %w", c.config.RetryAttempts, lastErr)
}

// doSingleRequest выполняет один HTTP запрос
func (c *Client) doSingleRequest(ctx context.Context, method, endpoint string, body interface{}) ([]byte, error) {
	url := c.config.APIURL + endpoint

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)

		log.Debug().
			Str("endpoint", endpoint).
			Str("method", method).
			RawJSON("request", jsonData).
			Msg("Post Express API request")
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Basic Authentication
	req.SetBasicAuth(c.config.Username, c.config.Password)

	// Headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Svetu-Marketplace/1.0")

	startTime := time.Now()
	resp, err := c.httpClient.Do(req)
	duration := time.Since(startTime)

	if err != nil {
		log.Error().
			Err(err).
			Str("endpoint", endpoint).
			Dur("duration", duration).
			Msg("Post Express API request failed")
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Warn().Err(closeErr).Msg("Failed to close response body")
		}
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	log.Debug().
		Str("endpoint", endpoint).
		Int("status", resp.StatusCode).
		Dur("duration", duration).
		RawJSON("response", respBody).
		Msg("Post Express API response")

	// Проверяем HTTP статус
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		apiErr := &APIError{
			StatusCode: resp.StatusCode,
			Message:    string(respBody),
		}

		// Пытаемся распарсить ошибку из JSON
		var errResp struct {
			Rezultat int    `json:"Rezultat"`
			Poruka   string `json:"Poruka"`
		}
		if json.Unmarshal(respBody, &errResp) == nil && errResp.Poruka != "" {
			apiErr.Message = errResp.Poruka
			apiErr.Code = errResp.Rezultat
		}

		return nil, apiErr
	}

	return respBody, nil
}

// Post выполняет POST запрос
func (c *Client) Post(ctx context.Context, endpoint string, request interface{}, response interface{}) error {
	respBody, err := c.doRequest(ctx, http.MethodPost, endpoint, request)
	if err != nil {
		return err
	}

	if response != nil {
		if err := json.Unmarshal(respBody, response); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}

		// Проверяем результат из API
		if checker, ok := response.(ResultChecker); ok {
			if !checker.IsSuccess() {
				return &APIError{
					StatusCode: 200,
					Code:       checker.GetCode(),
					Message:    checker.GetMessage(),
				}
			}
		}
	}

	return nil
}

// Get выполняет GET запрос
func (c *Client) Get(ctx context.Context, endpoint string, response interface{}) error {
	respBody, err := c.doRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}

	if response != nil {
		if err := json.Unmarshal(respBody, response); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}

		// Проверяем результат из API
		if checker, ok := response.(ResultChecker); ok {
			if !checker.IsSuccess() {
				return &APIError{
					StatusCode: 200,
					Code:       checker.GetCode(),
					Message:    checker.GetMessage(),
				}
			}
		}
	}

	return nil
}

// APIError - ошибка API
type APIError struct {
	StatusCode int
	Code       int
	Message    string
}

func (e *APIError) Error() string {
	if e.Code != 0 {
		return fmt.Sprintf("Post Express API error (HTTP %d, Code %d): %s", e.StatusCode, e.Code, e.Message)
	}
	return fmt.Sprintf("Post Express API error (HTTP %d): %s", e.StatusCode, e.Message)
}

// IsClientError проверяет, является ли ошибка клиентской (4xx)
func (e *APIError) IsClientError() bool {
	return e.StatusCode >= 400 && e.StatusCode < 500
}

// IsServerError проверяет, является ли ошибка серверной (5xx)
func (e *APIError) IsServerError() bool {
	return e.StatusCode >= 500
}

// isClientError проверяет, является ли ошибка клиентской
func isClientError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.IsClientError()
	}
	return false
}

// ResultChecker - интерфейс для проверки результата API
type ResultChecker interface {
	IsSuccess() bool
	GetCode() int
	GetMessage() string
}

// Реализация ResultChecker для основных response типов
func (r *ManifestResponse) IsSuccess() bool {
	return r.Rezultat == 0
}

func (r *ManifestResponse) GetCode() int {
	return r.Rezultat
}

func (r *ManifestResponse) GetMessage() string {
	return r.Poruka
}

func (r *TrackingResponse) IsSuccess() bool {
	return r.Rezultat == 0
}

func (r *TrackingResponse) GetCode() int {
	return r.Rezultat
}

func (r *TrackingResponse) GetMessage() string {
	return r.Poruka
}

func (r *CancelResponse) IsSuccess() bool {
	return r.Rezultat == 0
}

func (r *CancelResponse) GetCode() int {
	return r.Rezultat
}

func (r *CancelResponse) GetMessage() string {
	return r.Poruka
}

func (r *RateResponse) IsSuccess() bool {
	return r.Rezultat == 0
}

func (r *RateResponse) GetCode() int {
	return r.Rezultat
}

func (r *RateResponse) GetMessage() string {
	return r.Poruka
}

func (r *OfficeListResponse) IsSuccess() bool {
	return r.Rezultat == 0
}

func (r *OfficeListResponse) GetCode() int {
	return r.Rezultat
}

func (r *OfficeListResponse) GetMessage() string {
	return r.Poruka
}
