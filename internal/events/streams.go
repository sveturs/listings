package events

// Stream names for inter-service communication
const (
	// OrdersStream is the Redis stream for order events
	// WMS service listens to this stream for order.confirmed and order.cancelled events
	OrdersStream = "listings:events:orders"
)

// Event types
const (
	EventTypeOrderConfirmed = "order.confirmed"
	EventTypeOrderCancelled = "order.cancelled"
)
