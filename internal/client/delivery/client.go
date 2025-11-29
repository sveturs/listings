// Package delivery provides a gRPC client for interacting with the Delivery microservice.
// It implements retry logic, timeout handling, and proper error conversion.
package delivery

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	deliveryv1 "github.com/vondi-global/delivery/gen/go/delivery/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// Domain errors for delivery operations
var (
	ErrDeliveryUnavailable   = errors.New("delivery service unavailable")
	ErrShipmentNotFound      = errors.New("shipment not found")
	ErrInvalidAddress        = errors.New("invalid address")
	ErrInvalidPackage        = errors.New("invalid package dimensions")
	ErrProviderUnavailable   = errors.New("delivery provider unavailable")
	ErrRateCalculationFailed = errors.New("rate calculation failed")
	ErrShipmentCancelled     = errors.New("shipment already cancelled")
	ErrTimeout               = errors.New("delivery service timeout")
)

// Config holds delivery client configuration
type Config struct {
	Address    string
	Timeout    time.Duration
	MaxRetries int
	RetryDelay time.Duration
	PoolSize   int
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		Address:    "localhost:50052",
		Timeout:    10 * time.Second,
		MaxRetries: 3,
		RetryDelay: 100 * time.Millisecond,
		PoolSize:   1,
	}
}

// Client provides methods to interact with Delivery microservice
type Client struct {
	conn   *grpc.ClientConn
	client deliveryv1.DeliveryServiceClient
	config *Config
	logger zerolog.Logger
}

// NewClient creates a new delivery gRPC client
func NewClient(cfg *Config, logger zerolog.Logger) (*Client, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	// Create gRPC connection with interceptors
	conn, err := grpc.NewClient(
		cfg.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			loggingInterceptor(logger),
			retryInterceptor(cfg.MaxRetries, cfg.RetryDelay, logger),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to delivery service: %w", err)
	}

	return &Client{
		conn:   conn,
		client: deliveryv1.NewDeliveryServiceClient(conn),
		config: cfg,
		logger: logger.With().Str("component", "delivery_client").Logger(),
	}, nil
}

// Close closes the gRPC connection
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// CreateShipment creates a new shipment
func (c *Client) CreateShipment(ctx context.Context, req *CreateShipmentRequest) (*Shipment, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	grpcReq := &deliveryv1.CreateShipmentRequest{
		Provider:    mapProviderToProto(req.Provider),
		FromAddress: mapAddressToProto(req.FromAddress),
		ToAddress:   mapAddressToProto(req.ToAddress),
		Package:     mapPackageToProto(req.Package),
		UserId:      req.UserID,
	}

	resp, err := c.client.CreateShipment(ctx, grpcReq)
	if err != nil {
		return nil, c.mapError(err, "CreateShipment")
	}

	return mapShipmentFromProto(resp.Shipment), nil
}

// GetShipment retrieves shipment by ID
func (c *Client) GetShipment(ctx context.Context, shipmentID string) (*Shipment, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	resp, err := c.client.GetShipment(ctx, &deliveryv1.GetShipmentRequest{
		Id: shipmentID,
	})
	if err != nil {
		return nil, c.mapError(err, "GetShipment")
	}

	return mapShipmentFromProto(resp.Shipment), nil
}

// TrackShipment tracks shipment by tracking number
func (c *Client) TrackShipment(ctx context.Context, trackingNumber string) (*TrackingInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	resp, err := c.client.TrackShipment(ctx, &deliveryv1.TrackShipmentRequest{
		TrackingNumber: trackingNumber,
	})
	if err != nil {
		return nil, c.mapError(err, "TrackShipment")
	}

	return &TrackingInfo{
		Shipment: mapShipmentFromProto(resp.Shipment),
		Events:   mapTrackingEventsFromProto(resp.Events),
	}, nil
}

// CancelShipment cancels a shipment
func (c *Client) CancelShipment(ctx context.Context, shipmentID, reason string) (*Shipment, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	resp, err := c.client.CancelShipment(ctx, &deliveryv1.CancelShipmentRequest{
		Id:     shipmentID,
		Reason: reason,
	})
	if err != nil {
		return nil, c.mapError(err, "CancelShipment")
	}

	return mapShipmentFromProto(resp.Shipment), nil
}

// CalculateRate calculates delivery cost
func (c *Client) CalculateRate(ctx context.Context, req *CalculateRateRequest) (*RateInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	grpcReq := &deliveryv1.CalculateRateRequest{
		Provider:    mapProviderToProto(req.Provider),
		FromAddress: mapAddressToProto(req.FromAddress),
		ToAddress:   mapAddressToProto(req.ToAddress),
		Package:     mapPackageToProto(req.Package),
	}

	resp, err := c.client.CalculateRate(ctx, grpcReq)
	if err != nil {
		return nil, c.mapError(err, "CalculateRate")
	}

	return &RateInfo{
		Cost:              resp.Cost,
		Currency:          resp.Currency,
		EstimatedDelivery: resp.EstimatedDelivery.AsTime(),
	}, nil
}

// GetSettlements returns list of available settlements
func (c *Client) GetSettlements(ctx context.Context, provider DeliveryProvider, country, searchQuery string) ([]Settlement, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	resp, err := c.client.GetSettlements(ctx, &deliveryv1.GetSettlementsRequest{
		Provider:    mapProviderToProto(provider),
		Country:     country,
		SearchQuery: searchQuery,
	})
	if err != nil {
		return nil, c.mapError(err, "GetSettlements")
	}

	return mapSettlementsFromProto(resp.Settlements), nil
}

// GetStreets returns list of streets for a settlement
func (c *Client) GetStreets(ctx context.Context, provider DeliveryProvider, settlementName, searchQuery string) ([]Street, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	resp, err := c.client.GetStreets(ctx, &deliveryv1.GetStreetsRequest{
		Provider:       mapProviderToProto(provider),
		SettlementName: settlementName,
		SearchQuery:    searchQuery,
	})
	if err != nil {
		return nil, c.mapError(err, "GetStreets")
	}

	return mapStreetsFromProto(resp.Streets), nil
}

// GetParcelLockers returns list of available parcel lockers
func (c *Client) GetParcelLockers(ctx context.Context, provider DeliveryProvider, city, searchQuery string) ([]ParcelLocker, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	resp, err := c.client.GetParcelLockers(ctx, &deliveryv1.GetParcelLockersRequest{
		Provider:    mapProviderToProto(provider),
		City:        city,
		SearchQuery: searchQuery,
	})
	if err != nil {
		return nil, c.mapError(err, "GetParcelLockers")
	}

	return mapParcelLockersFromProto(resp.ParcelLockers), nil
}

// mapError converts gRPC errors to domain errors
func (c *Client) mapError(err error, method string) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		c.logger.Error().Err(err).Str("method", method).Msg("non-gRPC error")
		return fmt.Errorf("%w: %v", ErrDeliveryUnavailable, err)
	}

	c.logger.Warn().
		Str("method", method).
		Str("code", st.Code().String()).
		Str("message", st.Message()).
		Msg("gRPC error")

	switch st.Code() {
	case codes.NotFound:
		return ErrShipmentNotFound
	case codes.InvalidArgument:
		if contains(st.Message(), "address") {
			return ErrInvalidAddress
		}
		if contains(st.Message(), "package") {
			return ErrInvalidPackage
		}
		return fmt.Errorf("invalid argument: %s", st.Message())
	case codes.Unavailable:
		return ErrDeliveryUnavailable
	case codes.DeadlineExceeded:
		return ErrTimeout
	case codes.FailedPrecondition:
		if contains(st.Message(), "cancelled") {
			return ErrShipmentCancelled
		}
		return ErrProviderUnavailable
	case codes.Internal:
		return fmt.Errorf("%w: %s", ErrDeliveryUnavailable, st.Message())
	default:
		return fmt.Errorf("delivery error [%s]: %s", st.Code(), st.Message())
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
