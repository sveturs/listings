package events_test

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/vondi-global/listings/internal/events"
)

// ExampleRedisOrderEventPublisher_PublishOrderConfirmed demonstrates how to publish order.confirmed event
func ExampleRedisOrderEventPublisher_PublishOrderConfirmed() {
	// Create Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer redisClient.Close()

	// Create logger
	logger := zerolog.Nop()

	// Create publisher with default warehouse ID = 1
	publisher := events.NewRedisOrderEventPublisher(
		redisClient,
		logger,
		events.OrdersStream,
		1, // default warehouse ID
	)
	defer publisher.Close()

	// Prepare order items
	items := []events.OrderItem{
		{
			ListingID:   100,
			Quantity:    2,
			WarehouseID: 1,
		},
		{
			ListingID:   101,
			Quantity:    1,
			WarehouseID: 0, // Will use default warehouse ID
		},
	}

	// Publish order confirmed event
	ctx := context.Background()
	err := publisher.PublishOrderConfirmed(ctx, 12345, 999, items)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("Order confirmed event published successfully")
}

// ExampleRedisOrderEventPublisher_PublishOrderCancelled demonstrates how to publish order.cancelled event
func ExampleRedisOrderEventPublisher_PublishOrderCancelled() {
	// Create Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer redisClient.Close()

	// Create logger
	logger := zerolog.Nop()

	// Create publisher
	publisher := events.NewRedisOrderEventPublisher(
		redisClient,
		logger,
		events.OrdersStream,
		1, // default warehouse ID
	)
	defer publisher.Close()

	// Publish order cancelled event
	ctx := context.Background()
	err := publisher.PublishOrderCancelled(ctx, 12345, "customer_request")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("Order cancelled event published successfully")
}
