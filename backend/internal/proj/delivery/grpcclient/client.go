package grpcclient

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	pb "backend/pkg/grpc/delivery/v1"
	"backend/pkg/logger"
)

const (
	defaultTimeout     = 30 * time.Second
	maxRetries         = 3
	initialBackoff     = 100 * time.Millisecond
	maxBackoff         = 2 * time.Second
	backoffMultiplier  = 2.0
	circuitBreakerOpen = 5
)

type Client struct {
	conn            *grpc.ClientConn
	client          pb.DeliveryServiceClient
	logger          *logger.Logger
	failureCount    int
	lastFailureTime time.Time
	isOpen          bool
}

func NewClient(serverURL string, log *logger.Logger) (*Client, error) {
	log.Info("Connecting to delivery gRPC service: %s", serverURL)

	// grpc.NewClient заменяет deprecated grpc.DialContext
	// Соединение устанавливается лениво при первом вызове
	conn, err := grpc.NewClient(
		serverURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("Failed to create delivery gRPC client: %s, error: %v", serverURL, err)
		return nil, fmt.Errorf("failed to create delivery client: %w", err)
	}

	log.Info("Successfully created delivery gRPC client")

	return &Client{
		conn:   conn,
		client: pb.NewDeliveryServiceClient(conn),
		logger: log,
	}, nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *Client) CreateShipment(ctx context.Context, req *pb.CreateShipmentRequest) (*pb.CreateShipmentResponse, error) {
	if c.isCircuitBreakerOpen() {
		c.logger.Info("Circuit breaker is open, rejecting CreateShipment request")
		return nil, fmt.Errorf("delivery service temporarily unavailable")
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var resp *pb.CreateShipmentResponse
	var err error

	backoff := initialBackoff
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			c.logger.Info("Retrying CreateShipment, attempt: %d", attempt+1)
			time.Sleep(backoff)
			backoff = time.Duration(float64(backoff) * backoffMultiplier)
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}

		resp, err = c.client.CreateShipment(ctx, req)
		if err == nil {
			c.onSuccess()
			c.logger.Info("Shipment created successfully: %s", resp.Shipment.Id)
			return resp, nil
		}

		if !c.shouldRetry(err) {
			c.onFailure()
			c.logger.Error("Non-retryable error in CreateShipment: %v", err)
			return nil, err
		}

		c.logger.Info("Retryable error in CreateShipment, attempt %d: %v", attempt+1, err)
	}

	c.onFailure()
	c.logger.Error("CreateShipment failed after all retries: %v", err)
	return nil, fmt.Errorf("failed to create shipment after %d attempts: %w", maxRetries, err)
}

func (c *Client) GetShipment(ctx context.Context, req *pb.GetShipmentRequest) (*pb.GetShipmentResponse, error) {
	if c.isCircuitBreakerOpen() {
		return nil, fmt.Errorf("delivery service temporarily unavailable")
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var resp *pb.GetShipmentResponse
	var err error

	backoff := initialBackoff
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(backoff)
			backoff = time.Duration(float64(backoff) * backoffMultiplier)
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}

		resp, err = c.client.GetShipment(ctx, req)
		if err == nil {
			c.onSuccess()
			return resp, nil
		}

		if !c.shouldRetry(err) {
			c.onFailure()
			return nil, err
		}
	}

	c.onFailure()
	return nil, fmt.Errorf("failed to get shipment after %d attempts: %w", maxRetries, err)
}

func (c *Client) TrackShipment(ctx context.Context, req *pb.TrackShipmentRequest) (*pb.TrackShipmentResponse, error) {
	if c.isCircuitBreakerOpen() {
		return nil, fmt.Errorf("delivery service temporarily unavailable")
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var resp *pb.TrackShipmentResponse
	var err error

	backoff := initialBackoff
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(backoff)
			backoff = time.Duration(float64(backoff) * backoffMultiplier)
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}

		resp, err = c.client.TrackShipment(ctx, req)
		if err == nil {
			c.onSuccess()
			return resp, nil
		}

		if !c.shouldRetry(err) {
			c.onFailure()
			return nil, err
		}
	}

	c.onFailure()
	return nil, fmt.Errorf("failed to track shipment after %d attempts: %w", maxRetries, err)
}

func (c *Client) CancelShipment(ctx context.Context, req *pb.CancelShipmentRequest) (*pb.CancelShipmentResponse, error) {
	if c.isCircuitBreakerOpen() {
		return nil, fmt.Errorf("delivery service temporarily unavailable")
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var resp *pb.CancelShipmentResponse
	var err error

	backoff := initialBackoff
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(backoff)
			backoff = time.Duration(float64(backoff) * backoffMultiplier)
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}

		resp, err = c.client.CancelShipment(ctx, req)
		if err == nil {
			c.onSuccess()
			return resp, nil
		}

		if !c.shouldRetry(err) {
			c.onFailure()
			return nil, err
		}
	}

	c.onFailure()
	return nil, fmt.Errorf("failed to cancel shipment after %d attempts: %w", maxRetries, err)
}

func (c *Client) CalculateRate(ctx context.Context, req *pb.CalculateRateRequest) (*pb.CalculateRateResponse, error) {
	if c.isCircuitBreakerOpen() {
		return nil, fmt.Errorf("delivery service temporarily unavailable")
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var resp *pb.CalculateRateResponse
	var err error

	backoff := initialBackoff
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(backoff)
			backoff = time.Duration(float64(backoff) * backoffMultiplier)
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}

		resp, err = c.client.CalculateRate(ctx, req)
		if err == nil {
			c.onSuccess()
			return resp, nil
		}

		if !c.shouldRetry(err) {
			c.onFailure()
			return nil, err
		}
	}

	c.onFailure()
	return nil, fmt.Errorf("failed to calculate rate after %d attempts: %w", maxRetries, err)
}

func (c *Client) shouldRetry(err error) bool {
	st, ok := status.FromError(err)
	if !ok {
		return true
	}

	code := st.Code()
	switch code {
	// Non-retryable errors
	case codes.InvalidArgument, codes.NotFound, codes.AlreadyExists,
		codes.PermissionDenied, codes.Unauthenticated, codes.FailedPrecondition,
		codes.OutOfRange, codes.Unimplemented:
		return false
	// Retryable errors
	case codes.Unavailable, codes.DeadlineExceeded, codes.ResourceExhausted,
		codes.Aborted, codes.Canceled:
		return true
	// Success and other codes - don't retry
	case codes.OK, codes.Unknown, codes.Internal, codes.DataLoss:
		return false
	default:
		return true
	}
}

func (c *Client) isCircuitBreakerOpen() bool {
	if !c.isOpen {
		return false
	}

	if time.Since(c.lastFailureTime) > 30*time.Second {
		c.logger.Info("Circuit breaker transitioning to half-open state")
		c.isOpen = false
		c.failureCount = 0
		return false
	}

	return true
}

func (c *Client) onSuccess() {
	if c.failureCount > 0 || c.isOpen {
		c.logger.Info("Circuit breaker closed after successful request")
	}
	c.failureCount = 0
	c.isOpen = false
}

func (c *Client) onFailure() {
	c.failureCount++
	c.lastFailureTime = time.Now()

	if c.failureCount >= circuitBreakerOpen {
		c.isOpen = true
		c.logger.Info("Circuit breaker opened due to consecutive failures: %d", c.failureCount)
	}
}

// GetSettlements получает список населенных пунктов
func (c *Client) GetSettlements(ctx context.Context, req *pb.GetSettlementsRequest) (*pb.GetSettlementsResponse, error) {
	if c.isCircuitBreakerOpen() {
		c.logger.Info("Circuit breaker is open, rejecting GetSettlements request")
		return nil, fmt.Errorf("delivery service temporarily unavailable")
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var resp *pb.GetSettlementsResponse
	var err error

	backoff := initialBackoff
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			c.logger.Info("Retrying GetSettlements, attempt: %d", attempt+1)
			time.Sleep(backoff)
			backoff = time.Duration(float64(backoff) * backoffMultiplier)
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}

		resp, err = c.client.GetSettlements(ctx, req)
		if err == nil {
			c.onSuccess()
			c.logger.Info("GetSettlements successful: %d settlements found", len(resp.Settlements))
			return resp, nil
		}

		if !c.shouldRetry(err) {
			c.onFailure()
			c.logger.Error("Non-retryable error in GetSettlements: %v", err)
			return nil, err
		}

		c.logger.Info("Retryable error in GetSettlements, attempt %d: %v", attempt+1, err)
	}

	c.onFailure()
	c.logger.Error("GetSettlements failed after all retries: %v", err)
	return nil, fmt.Errorf("failed to get settlements after %d attempts: %w", maxRetries, err)
}

// GetStreets получает список улиц для населенного пункта
func (c *Client) GetStreets(ctx context.Context, req *pb.GetStreetsRequest) (*pb.GetStreetsResponse, error) {
	if c.isCircuitBreakerOpen() {
		c.logger.Info("Circuit breaker is open, rejecting GetStreets request")
		return nil, fmt.Errorf("delivery service temporarily unavailable")
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var resp *pb.GetStreetsResponse
	var err error

	backoff := initialBackoff
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			c.logger.Info("Retrying GetStreets, attempt: %d", attempt+1)
			time.Sleep(backoff)
			backoff = time.Duration(float64(backoff) * backoffMultiplier)
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}

		resp, err = c.client.GetStreets(ctx, req)
		if err == nil {
			c.onSuccess()
			c.logger.Info("GetStreets successful: %d streets found", len(resp.Streets))
			return resp, nil
		}

		if !c.shouldRetry(err) {
			c.onFailure()
			c.logger.Error("Non-retryable error in GetStreets: %v", err)
			return nil, err
		}

		c.logger.Info("Retryable error in GetStreets, attempt %d: %v", attempt+1, err)
	}

	c.onFailure()
	c.logger.Error("GetStreets failed after all retries: %v", err)
	return nil, fmt.Errorf("failed to get streets after %d attempts: %w", maxRetries, err)
}

// GetParcelLockers получает список паккетоматов
func (c *Client) GetParcelLockers(ctx context.Context, req *pb.GetParcelLockersRequest) (*pb.GetParcelLockersResponse, error) {
	if c.isCircuitBreakerOpen() {
		c.logger.Info("Circuit breaker is open, rejecting GetParcelLockers request")
		return nil, fmt.Errorf("delivery service temporarily unavailable")
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var resp *pb.GetParcelLockersResponse
	var err error

	backoff := initialBackoff
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			c.logger.Info("Retrying GetParcelLockers, attempt: %d", attempt+1)
			time.Sleep(backoff)
			backoff = time.Duration(float64(backoff) * backoffMultiplier)
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}

		resp, err = c.client.GetParcelLockers(ctx, req)
		if err == nil {
			c.onSuccess()
			c.logger.Info("GetParcelLockers successful: %d parcel lockers found", len(resp.ParcelLockers))
			return resp, nil
		}

		if !c.shouldRetry(err) {
			c.onFailure()
			c.logger.Error("Non-retryable error in GetParcelLockers: %v", err)
			return nil, err
		}

		c.logger.Info("Retryable error in GetParcelLockers, attempt %d: %v", attempt+1, err)
	}

	c.onFailure()
	c.logger.Error("GetParcelLockers failed after all retries: %v", err)
	return nil, fmt.Errorf("failed to get parcel lockers after %d attempts: %w", maxRetries, err)
}
