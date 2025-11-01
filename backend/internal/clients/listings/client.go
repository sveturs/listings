package listings

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/rs/zerolog"
	pb "github.com/sveturs/listings/api/proto/listings/v1"
)

const (
	defaultTimeout    = 30 * time.Second
	maxRetries        = 3
	initialBackoff    = 100 * time.Millisecond
	maxBackoff        = 2 * time.Second
	backoffMultiplier = 2.0
	// Circuit breaker: открывается после 5 последовательных ошибок
	circuitBreakerThreshold = 5
	// Circuit breaker: переходит в half-open после 30 секунд
	circuitBreakerTimeout = 30 * time.Second
)

// Client представляет gRPC клиент для listings микросервиса
type Client struct {
	conn            *grpc.ClientConn
	client          pb.ListingsServiceClient
	logger          zerolog.Logger
	failureCount    int
	lastFailureTime time.Time
	isOpen          bool
}

// NewClient создает новый gRPC клиент для listings микросервиса
func NewClient(serverURL string, logger zerolog.Logger) (*Client, error) {
	logger.Info().
		Str("url", serverURL).
		Msg("Connecting to listings gRPC service")

	// grpc.NewClient заменяет deprecated grpc.DialContext
	// Соединение устанавливается лениво при первом вызове
	conn, err := grpc.NewClient(
		serverURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Error().
			Err(err).
			Str("url", serverURL).
			Msg("Failed to create listings gRPC client")
		return nil, fmt.Errorf("failed to create listings gRPC client: %w", err)
	}

	logger.Info().Msg("Successfully created listings gRPC client")

	return &Client{
		conn:   conn,
		client: pb.NewListingsServiceClient(conn),
		logger: logger.With().Str("component", "listings_grpc_client").Logger(),
	}, nil
}

// Close закрывает соединение с gRPC сервером
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// GetListing получает listing по ID
func (c *Client) GetListing(ctx context.Context, req *pb.GetListingRequest) (*pb.GetListingResponse, error) {
	if c.isCircuitBreakerOpen() {
		c.logger.Warn().Msg("Circuit breaker is open, rejecting GetListing request")
		return nil, ErrServiceUnavailable
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var resp *pb.GetListingResponse
	var err error

	backoff := initialBackoff
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			c.logger.Info().Int("attempt", attempt+1).Msg("Retrying GetListing")
			time.Sleep(backoff)
			backoff = time.Duration(float64(backoff) * backoffMultiplier)
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}

		resp, err = c.client.GetListing(ctx, req)
		if err == nil {
			c.onSuccess()
			c.logger.Debug().Int64("listing_id", req.Id).Msg("GetListing successful")
			return resp, nil
		}

		if !c.shouldRetry(err) {
			c.onFailure()
			c.logger.Error().Err(err).Msg("Non-retryable error in GetListing")
			return nil, MapGRPCError(err)
		}

		c.logger.Warn().Err(err).Int("attempt", attempt+1).Msg("Retryable error in GetListing")
	}

	c.onFailure()
	c.logger.Error().Err(err).Msg("GetListing failed after all retries")
	return nil, fmt.Errorf("failed to get listing after %d attempts: %w", maxRetries, MapGRPCError(err))
}

// CreateListing создает новый listing
func (c *Client) CreateListing(ctx context.Context, req *pb.CreateListingRequest) (*pb.CreateListingResponse, error) {
	if c.isCircuitBreakerOpen() {
		c.logger.Warn().Msg("Circuit breaker is open, rejecting CreateListing request")
		return nil, ErrServiceUnavailable
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var resp *pb.CreateListingResponse
	var err error

	backoff := initialBackoff
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			c.logger.Info().Int("attempt", attempt+1).Msg("Retrying CreateListing")
			time.Sleep(backoff)
			backoff = time.Duration(float64(backoff) * backoffMultiplier)
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}

		resp, err = c.client.CreateListing(ctx, req)
		if err == nil {
			c.onSuccess()
			c.logger.Info().
				Int64("listing_id", resp.Listing.Id).
				Str("title", req.Title).
				Msg("Listing created successfully")
			return resp, nil
		}

		if !c.shouldRetry(err) {
			c.onFailure()
			c.logger.Error().Err(err).Msg("Non-retryable error in CreateListing")
			return nil, MapGRPCError(err)
		}

		c.logger.Warn().Err(err).Int("attempt", attempt+1).Msg("Retryable error in CreateListing")
	}

	c.onFailure()
	c.logger.Error().Err(err).Msg("CreateListing failed after all retries")
	return nil, fmt.Errorf("failed to create listing after %d attempts: %w", maxRetries, MapGRPCError(err))
}

// UpdateListing обновляет существующий listing
func (c *Client) UpdateListing(ctx context.Context, req *pb.UpdateListingRequest) (*pb.UpdateListingResponse, error) {
	if c.isCircuitBreakerOpen() {
		c.logger.Warn().Msg("Circuit breaker is open, rejecting UpdateListing request")
		return nil, ErrServiceUnavailable
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var resp *pb.UpdateListingResponse
	var err error

	backoff := initialBackoff
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			c.logger.Info().Int("attempt", attempt+1).Msg("Retrying UpdateListing")
			time.Sleep(backoff)
			backoff = time.Duration(float64(backoff) * backoffMultiplier)
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}

		resp, err = c.client.UpdateListing(ctx, req)
		if err == nil {
			c.onSuccess()
			c.logger.Info().Int64("listing_id", req.Id).Msg("Listing updated successfully")
			return resp, nil
		}

		if !c.shouldRetry(err) {
			c.onFailure()
			c.logger.Error().Err(err).Msg("Non-retryable error in UpdateListing")
			return nil, MapGRPCError(err)
		}

		c.logger.Warn().Err(err).Int("attempt", attempt+1).Msg("Retryable error in UpdateListing")
	}

	c.onFailure()
	c.logger.Error().Err(err).Msg("UpdateListing failed after all retries")
	return nil, fmt.Errorf("failed to update listing after %d attempts: %w", maxRetries, MapGRPCError(err))
}

// DeleteListing удаляет listing (soft delete)
func (c *Client) DeleteListing(ctx context.Context, req *pb.DeleteListingRequest) (*pb.DeleteListingResponse, error) {
	if c.isCircuitBreakerOpen() {
		c.logger.Warn().Msg("Circuit breaker is open, rejecting DeleteListing request")
		return nil, ErrServiceUnavailable
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var resp *pb.DeleteListingResponse
	var err error

	backoff := initialBackoff
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			c.logger.Info().Int("attempt", attempt+1).Msg("Retrying DeleteListing")
			time.Sleep(backoff)
			backoff = time.Duration(float64(backoff) * backoffMultiplier)
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}

		resp, err = c.client.DeleteListing(ctx, req)
		if err == nil {
			c.onSuccess()
			c.logger.Info().Int64("listing_id", req.Id).Msg("Listing deleted successfully")
			return resp, nil
		}

		if !c.shouldRetry(err) {
			c.onFailure()
			c.logger.Error().Err(err).Msg("Non-retryable error in DeleteListing")
			return nil, MapGRPCError(err)
		}

		c.logger.Warn().Err(err).Int("attempt", attempt+1).Msg("Retryable error in DeleteListing")
	}

	c.onFailure()
	c.logger.Error().Err(err).Msg("DeleteListing failed after all retries")
	return nil, fmt.Errorf("failed to delete listing after %d attempts: %w", maxRetries, MapGRPCError(err))
}

// SearchListings выполняет поиск по listings
func (c *Client) SearchListings(ctx context.Context, req *pb.SearchListingsRequest) (*pb.SearchListingsResponse, error) {
	if c.isCircuitBreakerOpen() {
		c.logger.Warn().Msg("Circuit breaker is open, rejecting SearchListings request")
		return nil, ErrServiceUnavailable
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var resp *pb.SearchListingsResponse
	var err error

	backoff := initialBackoff
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			c.logger.Info().Int("attempt", attempt+1).Msg("Retrying SearchListings")
			time.Sleep(backoff)
			backoff = time.Duration(float64(backoff) * backoffMultiplier)
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}

		resp, err = c.client.SearchListings(ctx, req)
		if err == nil {
			c.onSuccess()
			c.logger.Debug().
				Int("count", len(resp.Listings)).
				Str("query", req.Query).
				Msg("SearchListings successful")
			return resp, nil
		}

		if !c.shouldRetry(err) {
			c.onFailure()
			c.logger.Error().Err(err).Msg("Non-retryable error in SearchListings")
			return nil, MapGRPCError(err)
		}

		c.logger.Warn().Err(err).Int("attempt", attempt+1).Msg("Retryable error in SearchListings")
	}

	c.onFailure()
	c.logger.Error().Err(err).Msg("SearchListings failed after all retries")
	return nil, fmt.Errorf("failed to search listings after %d attempts: %w", maxRetries, MapGRPCError(err))
}

// ListListings получает список listings с фильтрацией
func (c *Client) ListListings(ctx context.Context, req *pb.ListListingsRequest) (*pb.ListListingsResponse, error) {
	if c.isCircuitBreakerOpen() {
		c.logger.Warn().Msg("Circuit breaker is open, rejecting ListListings request")
		return nil, ErrServiceUnavailable
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var resp *pb.ListListingsResponse
	var err error

	backoff := initialBackoff
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			c.logger.Info().Int("attempt", attempt+1).Msg("Retrying ListListings")
			time.Sleep(backoff)
			backoff = time.Duration(float64(backoff) * backoffMultiplier)
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}

		resp, err = c.client.ListListings(ctx, req)
		if err == nil {
			c.onSuccess()
			c.logger.Debug().Int("count", len(resp.Listings)).Msg("ListListings successful")
			return resp, nil
		}

		if !c.shouldRetry(err) {
			c.onFailure()
			c.logger.Error().Err(err).Msg("Non-retryable error in ListListings")
			return nil, MapGRPCError(err)
		}

		c.logger.Warn().Err(err).Int("attempt", attempt+1).Msg("Retryable error in ListListings")
	}

	c.onFailure()
	c.logger.Error().Err(err).Msg("ListListings failed after all retries")
	return nil, fmt.Errorf("failed to list listings after %d attempts: %w", maxRetries, MapGRPCError(err))
}

// shouldRetry определяет, нужно ли повторять запрос при данной ошибке
func (c *Client) shouldRetry(err error) bool {
	st, ok := status.FromError(err)
	if !ok {
		// Неизвестная ошибка - пробуем повторить
		return true
	}

	code := st.Code()
	switch code {
	// Non-retryable errors - ошибки валидации и бизнес-логики
	case codes.InvalidArgument, codes.NotFound, codes.AlreadyExists,
		codes.PermissionDenied, codes.Unauthenticated, codes.FailedPrecondition,
		codes.OutOfRange, codes.Unimplemented:
		return false
	// Retryable errors - временные сетевые проблемы
	case codes.Unavailable, codes.DeadlineExceeded, codes.ResourceExhausted,
		codes.Aborted, codes.Canceled:
		return true
	// Success and other codes - не повторяем
	case codes.OK, codes.Unknown, codes.Internal, codes.DataLoss:
		return false
	default:
		return true
	}
}

// isCircuitBreakerOpen проверяет, открыт ли circuit breaker
func (c *Client) isCircuitBreakerOpen() bool {
	if !c.isOpen {
		return false
	}

	// Переход в half-open state после таймаута
	if time.Since(c.lastFailureTime) > circuitBreakerTimeout {
		c.logger.Info().Msg("Circuit breaker transitioning to half-open state")
		c.isOpen = false
		c.failureCount = 0
		return false
	}

	return true
}

// onSuccess обрабатывает успешный ответ
func (c *Client) onSuccess() {
	if c.failureCount > 0 || c.isOpen {
		c.logger.Info().Msg("Circuit breaker closed after successful request")
	}
	c.failureCount = 0
	c.isOpen = false
}

// onFailure обрабатывает ошибку
func (c *Client) onFailure() {
	c.failureCount++
	c.lastFailureTime = time.Now()

	if c.failureCount >= circuitBreakerThreshold {
		c.isOpen = true
		c.logger.Warn().
			Int("failure_count", c.failureCount).
			Msg("Circuit breaker opened due to consecutive failures")
	}
}
