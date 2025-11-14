# Order Service gRPC Integration Guide

This document describes how to integrate the Order Service gRPC handlers into the main server.

## Overview

The Order Service provides gRPC endpoints for:
- **Cart operations** (6 methods): AddToCart, UpdateCartItem, RemoveFromCart, GetCart, ClearCart, GetUserCarts
- **Order operations** (6 methods): CreateOrder, GetOrder, ListOrders, CancelOrder, UpdateOrderStatus, GetOrderStats

## Architecture

```
Client → gRPC → Middleware Chain → OrderServiceServer → Service Layer → Repository → Database
```

**Middleware Chain** (in order):
1. Recovery (catch panics)
2. Logging (log all requests)
3. Auth (extract user_id/session_id from metadata)
4. Metrics (record timing and status codes)

## Integration Steps

### 1. Initialize Services (in `cmd/server/main.go`)

Add the following service initialization after existing service initialization (around line 210):

```go
// Initialize Cart Service
cartRepo := postgres.NewCartRepository(db, zerologLogger)
cartService := service.NewCartService(
	cartRepo,
	pgRepo,
	pgRepo, // storefrontRepo
	pgRepo, // db
	zerologLogger,
)

// Initialize Reservation Repository
reservationRepo := postgres.NewReservationRepository(db, zerologLogger)

// Initialize Order Repository
orderRepo := postgres.NewOrderRepository(db, zerologLogger)

// Initialize Inventory Service
inventoryService := service.NewInventoryService(
	reservationRepo,
	pgRepo,
	orderRepo,
	pool, // pgxpool.Pool (need to add this to main.go if not exists)
	zerologLogger,
)

// Initialize Financial Config
financialConfig := service.DefaultFinancialConfig()

// Initialize Order Service
orderService := service.NewOrderService(
	orderRepo,
	cartRepo,
	reservationRepo,
	pgRepo,
	pool, // pgxpool.Pool
	financialConfig,
	zerologLogger,
)
```

### 2. Create Order Service Handler

Add after line 283 (after gRPC handler initialization):

```go
// Create Order Service gRPC handler
orderServiceHandler := grpcTransport.NewOrderServiceServer(
	cartService,
	orderService,
	inventoryService,
	zerologLogger,
)

// Create Order Service middleware
orderMiddleware := grpcTransport.NewOrderServiceMiddleware(zerologLogger, metricsInstance)

// Build interceptor chain for Order Service
// Order: Recovery → Logging → Auth → Metrics
orderInterceptors := orderMiddleware.GetDefaultInterceptors()

// Create separate gRPC server for Order Service (optional - can use same server)
// OR register to existing grpcServer if you want unified server
ordersspb.RegisterOrderServiceServer(grpcServer, orderServiceHandler)
```

### 3. Alternative: Separate gRPC Server for Orders

If you want to run Order Service on a separate port:

```go
// Create separate gRPC server for Order Service
orderGRPCServer := grpc.NewServer(
	grpc.ChainUnaryInterceptor(orderInterceptors),
)

// Register Order Service
ordersspb.RegisterOrderServiceServer(orderGRPCServer, orderServiceHandler)

// Enable reflection
reflection.Register(orderGRPCServer)

// Start Order Service gRPC server
orderGRPCListener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Server.GRPCHost, cfg.Server.OrderGRPCPort))
if err != nil {
	logger.Fatal().Err(err).Msg("failed to create Order Service gRPC listener")
}

go func() {
	logger.Info().Int("port", cfg.Server.OrderGRPCPort).Msg("Starting Order Service gRPC server")
	if err := orderGRPCServer.Serve(orderGRPCListener); err != nil {
		logger.Error().Err(err).Msg("Order Service gRPC server error")
	}
}()
```

## Configuration

Add to `.env` (if using separate server):

```bash
# Order Service gRPC port (optional, defaults to main gRPC port)
ORDER_GRPC_PORT=50052

# Financial configuration
PLATFORM_COMMISSION_RATE=0.05  # 5% commission
TAX_RATE=0.20                  # 20% VAT
ESCROW_DAYS=3                  # Hold funds for 3 days
```

## Testing

### Using grpcurl

```bash
# List services
grpcurl -plaintext localhost:50051 list

# Describe OrderService
grpcurl -plaintext localhost:50051 describe listingssvc.v1.OrderService

# Add item to cart
grpcurl -plaintext -d '{
  "user_id": 42,
  "storefront_id": 1,
  "listing_id": 100,
  "quantity": 2
}' localhost:50051 listingssvc.v1.OrderService/AddToCart

# Get cart
grpcurl -plaintext -d '{
  "user_id": 42,
  "storefront_id": 1
}' localhost:50051 listingssvc.v1.OrderService/GetCart

# Create order
grpcurl -plaintext -d '{
  "user_id": 42,
  "storefront_id": 1,
  "cart_id": 1,
  "shipping_address": {
    "street": "123 Main St",
    "city": "Belgrade",
    "postal_code": "11000",
    "country": "RS"
  },
  "shipping_method": "standard",
  "payment_method": "card"
}' localhost:50051 listingssvc.v1.OrderService/CreateOrder

# Get order
grpcurl -plaintext -d '{
  "order_id": 1,
  "user_id": 42
}' localhost:50051 listingssvc.v1.OrderService/GetOrder

# Cancel order
grpcurl -plaintext -d '{
  "order_id": 1,
  "user_id": 42,
  "reason": "Changed my mind",
  "refund": true
}' localhost:50051 listingssvc.v1.OrderService/CancelOrder
```

### Using evans (Interactive gRPC CLI)

```bash
# Start evans in REPL mode
evans -p 50051 repl

# Select service
service listingssvc.v1.OrderService

# Call methods interactively
call AddToCart
call GetCart
call CreateOrder
```

## Metrics

The Order Service exposes the following Prometheus metrics:

- `grpc_request_duration_seconds{method, code}` - Request duration histogram
- `grpc_requests_total{method, code}` - Total request counter
- `grpc_errors_total{method, code}` - Error counter

Access metrics at: `http://localhost:8080/metrics`

## Middleware Details

### Auth Interceptor

Extracts authentication information from gRPC metadata:

```go
// Client sends metadata
md := metadata.Pairs(
	"user_id", "42",
	"session_id", "session_abc123",
)
ctx := metadata.NewOutgoingContext(context.Background(), md)

// Server receives and extracts
userID := ctx.Value("user_id").(string)
sessionID := ctx.Value("session_id").(string)
```

### Logging Interceptor

Logs all RPC calls:

```
INFO: gRPC request started method=/OrderService/CreateOrder
INFO: gRPC request completed method=/OrderService/CreateOrder duration_ms=125
ERROR: gRPC request failed method=/OrderService/CreateOrder code=InvalidArgument duration_ms=5
```

### Recovery Interceptor

Catches panics and returns gRPC Internal error:

```go
// If handler panics
panic("something went wrong")

// Client receives
status.Error(codes.Internal, "internal server error: something went wrong")
```

### Metrics Interceptor

Records timing and status codes for all RPC calls.

## Error Mapping

Service layer errors are mapped to gRPC status codes:

| Service Error | gRPC Code |
|---------------|-----------|
| `ErrNotFound`, `ErrCartNotFound`, `ErrOrderNotFound`, `ErrListingNotFound` | `NotFound` |
| `ErrInvalidInput`, `ErrCartEmpty`, `ErrInvalidAddress`, `ErrInvalidPaymentMethod` | `InvalidArgument` |
| `ErrInsufficientStock`, `ErrPriceChanged`, `ErrOrderCannotCancel`, `ErrStorefrontMismatch` | `FailedPrecondition` |
| `ErrUnauthorized` | `PermissionDenied` |
| Other errors | `Internal` |

## Repository Requirements

The Order Service requires the following repositories:

1. **CartRepository** (`postgres.NewCartRepository`)
   - `GetByUserAndStorefront(ctx, userID, storefrontID) (*Cart, error)`
   - `GetBySessionAndStorefront(ctx, sessionID, storefrontID) (*Cart, error)`
   - `GetByID(ctx, cartID) (*Cart, error)`
   - `Create(ctx, cart) error`
   - `UpdateItem(ctx, item) error`
   - `AddItem(ctx, item) error`
   - `RemoveItem(ctx, cartID, itemID) error`
   - `ClearItems(ctx, cartID) error`
   - `GetUserCarts(ctx, userID) ([]*Cart, error)`

2. **OrderRepository** (`postgres.NewOrderRepository`)
   - `Create(ctx, order) error`
   - `GetByID(ctx, orderID) (*Order, error)`
   - `GetByOrderNumber(ctx, orderNumber) (*Order, error)`
   - `ListByUser(ctx, userID, limit, offset) ([]*Order, int, error)`
   - `ListByStorefront(ctx, storefrontID, limit, offset) ([]*Order, int, error)`
   - `UpdateStatus(ctx, orderID, status) error`
   - `Update(ctx, order) error`
   - `CreateItems(ctx, orderID, items) error`

3. **ReservationRepository** (`postgres.NewReservationRepository`)
   - `Create(ctx, reservation) error`
   - `GetByID(ctx, reservationID) (*Reservation, error)`
   - `GetByOrderID(ctx, orderID) ([]*Reservation, error)`
   - `Update(ctx, reservation) error`
   - `CommitReservations(ctx, orderID) error`
   - `ReleaseReservations(ctx, orderID) error`
   - `ExpireStaleReservations(ctx) (int, error)`
   - `GetActiveByListing(ctx, listingID, variantID) ([]*Reservation, error)`

## Database Migrations

Required database tables (if not exists):

```sql
-- carts table
CREATE TABLE IF NOT EXISTS carts (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT,
    session_id VARCHAR(255),
    storefront_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT carts_user_storefront_unique UNIQUE (user_id, storefront_id),
    CONSTRAINT carts_session_storefront_unique UNIQUE (session_id, storefront_id)
);

-- cart_items table
CREATE TABLE IF NOT EXISTS cart_items (
    id BIGSERIAL PRIMARY KEY,
    cart_id BIGINT NOT NULL REFERENCES carts(id) ON DELETE CASCADE,
    listing_id BIGINT NOT NULL,
    variant_id BIGINT,
    quantity INT NOT NULL CHECK (quantity > 0),
    price_snapshot DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- orders table
CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,
    order_number VARCHAR(50) UNIQUE NOT NULL,
    user_id BIGINT,
    storefront_id BIGINT NOT NULL,
    status VARCHAR(50) NOT NULL,
    payment_status VARCHAR(50) NOT NULL,
    subtotal DECIMAL(10, 2) NOT NULL,
    tax DECIMAL(10, 2) NOT NULL,
    shipping DECIMAL(10, 2) NOT NULL,
    discount DECIMAL(10, 2) DEFAULT 0,
    total DECIMAL(10, 2) NOT NULL,
    commission DECIMAL(10, 2) NOT NULL,
    seller_amount DECIMAL(10, 2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    shipping_address JSONB NOT NULL,
    billing_address JSONB,
    payment_method VARCHAR(50),
    payment_transaction_id VARCHAR(255),
    payment_completed_at TIMESTAMP,
    shipping_method VARCHAR(50),
    tracking_number VARCHAR(255),
    escrow_days INT DEFAULT 3,
    escrow_release_date TIMESTAMP,
    customer_notes TEXT,
    admin_notes TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    confirmed_at TIMESTAMP,
    shipped_at TIMESTAMP,
    delivered_at TIMESTAMP,
    cancelled_at TIMESTAMP
);

-- order_items table
CREATE TABLE IF NOT EXISTS order_items (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    listing_id BIGINT NOT NULL,
    variant_id BIGINT,
    listing_name VARCHAR(255) NOT NULL,
    quantity INT NOT NULL CHECK (quantity > 0),
    unit_price DECIMAL(10, 2) NOT NULL,
    subtotal DECIMAL(10, 2) NOT NULL,
    discount DECIMAL(10, 2) DEFAULT 0,
    total DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- inventory_reservations table
CREATE TABLE IF NOT EXISTS inventory_reservations (
    id BIGSERIAL PRIMARY KEY,
    listing_id BIGINT NOT NULL,
    variant_id BIGINT,
    order_id BIGINT NOT NULL REFERENCES orders(id),
    quantity INT NOT NULL CHECK (quantity > 0),
    status VARCHAR(50) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    committed_at TIMESTAMP,
    released_at TIMESTAMP
);
```

## Next Steps

1. ✅ Implement repository methods (CartRepository, OrderRepository, ReservationRepository)
2. ✅ Add database migrations for tables
3. ✅ Integrate Order Service into main.go
4. ⏳ Add integration tests
5. ⏳ Add end-to-end tests with real database
6. ⏳ Deploy to staging environment

## Troubleshooting

### Issue: "rpc error: code = Unimplemented desc = method AddToCart not implemented"

**Solution:** Ensure OrderServiceServer is properly registered:
```go
ordersspb.RegisterOrderServiceServer(grpcServer, orderServiceHandler)
```

### Issue: "cart not found"

**Solution:** The cart is created automatically on first AddToCart call. Make sure the service is initialized properly.

### Issue: "panic recovered: runtime error: invalid memory address"

**Solution:** Ensure all dependencies (cartService, orderService, inventoryService) are initialized before creating OrderServiceServer.

### Issue: "insufficient stock"

**Solution:** This is expected behavior. Ensure the listing has enough stock before adding to cart or creating order.

## Performance Considerations

- **Caching**: Cart data can be cached in Redis with 5-minute TTL
- **Connection Pooling**: Use `pgxpool.Pool` for efficient database connections
- **Batch Operations**: When possible, use batch repository methods
- **Transactions**: All order creation operations are wrapped in transactions
- **Deadlock Prevention**: Listings are locked in sorted order (ORDER BY id ASC)

## Security

- **User Ownership**: GetOrder and CancelOrder verify user owns the order
- **Admin Operations**: UpdateOrderStatus requires admin role (implement auth check)
- **Input Validation**: All inputs are validated before processing
- **SQL Injection**: Use parameterized queries (handled by pgx)
- **CSRF**: gRPC is not susceptible to CSRF attacks

## Monitoring

Monitor these metrics:

- `grpc_request_duration_seconds{method="CreateOrder"}` - Should be < 500ms
- `grpc_errors_total{method="CreateOrder", code="FailedPrecondition"}` - Price changes
- `grpc_errors_total{method="CreateOrder", code="NotFound"}` - Missing carts/listings
- Database connection pool metrics

Set up alerts for:
- High error rate (> 5%)
- Slow requests (p99 > 1s)
- Database connection pool exhaustion
