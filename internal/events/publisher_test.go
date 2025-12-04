package events

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRedis(t *testing.T) (*redis.Client, *miniredis.Miniredis) {
	mr := miniredis.RunT(t)

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return client, mr
}

func TestPublishOrderConfirmed(t *testing.T) {
	client, mr := setupTestRedis(t)
	defer mr.Close()

	logger := zerolog.Nop()
	publisher := NewRedisOrderEventPublisher(client, logger, OrdersStream, 1)

	ctx := context.Background()
	items := []OrderItem{
		{ListingID: 100, Quantity: 2, WarehouseID: 0}, // Should use default warehouse
		{ListingID: 101, Quantity: 1, WarehouseID: 5}, // Should keep specified warehouse
	}

	err := publisher.PublishOrderConfirmed(ctx, 12345, 999, items)
	require.NoError(t, err)

	// Verify stream exists and has messages
	entries, err := client.XRange(ctx, OrdersStream, "-", "+").Result()
	require.NoError(t, err)
	assert.Len(t, entries, 1)

	// Verify event type
	assert.Equal(t, EventTypeOrderConfirmed, entries[0].Values["type"])

	// Verify order ID
	assert.Equal(t, "12345", entries[0].Values["order_id"])

	// Verify storefront ID
	assert.Equal(t, "999", entries[0].Values["storefront_id"])
}

func TestPublishOrderCancelled(t *testing.T) {
	client, mr := setupTestRedis(t)
	defer mr.Close()

	logger := zerolog.Nop()
	publisher := NewRedisOrderEventPublisher(client, logger, OrdersStream, 1)

	ctx := context.Background()
	err := publisher.PublishOrderCancelled(ctx, 12345, "customer_request")
	require.NoError(t, err)

	// Verify stream exists and has messages
	entries, err := client.XRange(ctx, OrdersStream, "-", "+").Result()
	require.NoError(t, err)
	assert.Len(t, entries, 1)

	// Verify event type
	assert.Equal(t, EventTypeOrderCancelled, entries[0].Values["type"])

	// Verify order ID
	assert.Equal(t, "12345", entries[0].Values["order_id"])

	// Verify reason
	assert.Equal(t, "customer_request", entries[0].Values["reason"])
}

func TestDefaultWarehouseID(t *testing.T) {
	client, mr := setupTestRedis(t)
	defer mr.Close()

	logger := zerolog.Nop()
	defaultWarehouse := int64(42)
	publisher := NewRedisOrderEventPublisher(client, logger, OrdersStream, defaultWarehouse)

	ctx := context.Background()
	items := []OrderItem{
		{ListingID: 100, Quantity: 2, WarehouseID: 0}, // Should get default
	}

	err := publisher.PublishOrderConfirmed(ctx, 12345, 999, items)
	require.NoError(t, err)

	// Read the stream entry
	entries, err := client.XRange(ctx, OrdersStream, "-", "+").Result()
	require.NoError(t, err)
	require.Len(t, entries, 1)

	// Verify the items JSON contains the default warehouse
	itemsJSON := entries[0].Values["items"].(string)

	var parsedItems []OrderItem
	err = json.Unmarshal([]byte(itemsJSON), &parsedItems)
	require.NoError(t, err)
	require.Len(t, parsedItems, 1)

	// Verify default warehouse was set
	assert.Equal(t, defaultWarehouse, parsedItems[0].WarehouseID)
}

func TestMultipleEvents(t *testing.T) {
	client, mr := setupTestRedis(t)
	defer mr.Close()

	logger := zerolog.Nop()
	publisher := NewRedisOrderEventPublisher(client, logger, OrdersStream, 1)

	ctx := context.Background()

	// Publish multiple events
	items := []OrderItem{{ListingID: 100, Quantity: 1, WarehouseID: 1}}

	err := publisher.PublishOrderConfirmed(ctx, 1, 999, items)
	require.NoError(t, err)

	time.Sleep(10 * time.Millisecond) // Ensure different timestamps

	err = publisher.PublishOrderCancelled(ctx, 1, "test")
	require.NoError(t, err)

	// Verify both messages
	entries, err := client.XRange(ctx, OrdersStream, "-", "+").Result()
	require.NoError(t, err)
	require.Len(t, entries, 2)

	// Verify event types
	assert.Equal(t, EventTypeOrderConfirmed, entries[0].Values["type"])
	assert.Equal(t, EventTypeOrderCancelled, entries[1].Values["type"])
}

func TestKeepSpecifiedWarehouse(t *testing.T) {
	client, mr := setupTestRedis(t)
	defer mr.Close()

	logger := zerolog.Nop()
	defaultWarehouse := int64(1)
	publisher := NewRedisOrderEventPublisher(client, logger, OrdersStream, defaultWarehouse)

	ctx := context.Background()
	items := []OrderItem{
		{ListingID: 100, Quantity: 1, WarehouseID: 5}, // Should keep 5
	}

	err := publisher.PublishOrderConfirmed(ctx, 12345, 999, items)
	require.NoError(t, err)

	// Read the stream entry
	entries, err := client.XRange(ctx, OrdersStream, "-", "+").Result()
	require.NoError(t, err)
	require.Len(t, entries, 1)

	// Parse items JSON
	itemsJSON := entries[0].Values["items"].(string)
	var parsedItems []OrderItem
	err = json.Unmarshal([]byte(itemsJSON), &parsedItems)
	require.NoError(t, err)

	// Verify warehouse was NOT changed to default
	assert.Equal(t, int64(5), parsedItems[0].WarehouseID)
}
