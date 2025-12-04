package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

// OrderEventPublisher defines interface for publishing order events
type OrderEventPublisher interface {
	PublishOrderConfirmed(ctx context.Context, orderID int64, storefrontID int64, items []OrderItem) error
	PublishOrderCancelled(ctx context.Context, orderID int64, reason string) error
	Close() error
}

// OrderItem represents item data for order events
type OrderItem struct {
	ListingID   int64 `json:"listing_id"`
	Quantity    int32 `json:"quantity"`
	WarehouseID int64 `json:"warehouse_id"`
}

// RedisOrderEventPublisher publishes order events to Redis Streams
type RedisOrderEventPublisher struct {
	redis       *redis.Client
	logger      zerolog.Logger
	streamName  string
	warehouseID int64 // Default warehouse ID
}

// NewRedisOrderEventPublisher creates a new publisher
func NewRedisOrderEventPublisher(redisClient *redis.Client, logger zerolog.Logger, streamName string, defaultWarehouseID int64) *RedisOrderEventPublisher {
	return &RedisOrderEventPublisher{
		redis:       redisClient,
		logger:      logger.With().Str("component", "order_event_publisher").Logger(),
		streamName:  streamName,
		warehouseID: defaultWarehouseID,
	}
}

// PublishOrderConfirmed publishes order.confirmed event to Redis Stream
func (p *RedisOrderEventPublisher) PublishOrderConfirmed(ctx context.Context, orderID int64, storefrontID int64, items []OrderItem) error {
	// Set default warehouse if not specified
	for i := range items {
		if items[i].WarehouseID == 0 {
			items[i].WarehouseID = p.warehouseID
		}
	}

	itemsJSON, err := json.Marshal(items)
	if err != nil {
		return fmt.Errorf("failed to marshal items: %w", err)
	}

	values := map[string]interface{}{
		"type":          "order.confirmed",
		"order_id":      fmt.Sprintf("%d", orderID),
		"storefront_id": fmt.Sprintf("%d", storefrontID),
		"items":         string(itemsJSON),
		"timestamp":     fmt.Sprintf("%d", time.Now().Unix()),
	}

	result, err := p.redis.XAdd(ctx, &redis.XAddArgs{
		Stream: p.streamName,
		Values: values,
	}).Result()

	if err != nil {
		p.logger.Error().Err(err).
			Int64("order_id", orderID).
			Str("stream", p.streamName).
			Msg("failed to publish order.confirmed event")
		return fmt.Errorf("failed to publish order.confirmed event: %w", err)
	}

	p.logger.Info().
		Str("message_id", result).
		Int64("order_id", orderID).
		Int64("storefront_id", storefrontID).
		Int("items_count", len(items)).
		Str("stream", p.streamName).
		Msg("published order.confirmed event")

	return nil
}

// PublishOrderCancelled publishes order.cancelled event to Redis Stream
func (p *RedisOrderEventPublisher) PublishOrderCancelled(ctx context.Context, orderID int64, reason string) error {
	values := map[string]interface{}{
		"type":      "order.cancelled",
		"order_id":  fmt.Sprintf("%d", orderID),
		"reason":    reason,
		"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
	}

	result, err := p.redis.XAdd(ctx, &redis.XAddArgs{
		Stream: p.streamName,
		Values: values,
	}).Result()

	if err != nil {
		p.logger.Error().Err(err).
			Int64("order_id", orderID).
			Str("stream", p.streamName).
			Msg("failed to publish order.cancelled event")
		return fmt.Errorf("failed to publish order.cancelled event: %w", err)
	}

	p.logger.Info().
		Str("message_id", result).
		Int64("order_id", orderID).
		Str("reason", reason).
		Str("stream", p.streamName).
		Msg("published order.cancelled event")

	return nil
}

// Close gracefully closes the publisher
func (p *RedisOrderEventPublisher) Close() error {
	// No resources to cleanup as redis client is shared
	return nil
}
