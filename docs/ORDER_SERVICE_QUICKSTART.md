# Order Service - Quick Start Guide

This is a **quick reference** for using the Order Service gRPC API.

For complete documentation, see: `ORDER_SERVICE_INTEGRATION.md`

## Prerequisites

1. gRPC server running on `localhost:50051`
2. Install grpcurl: `go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest`

## Testing Flow

### Step 1: Add Item to Cart

```bash
grpcurl -plaintext -d '{
  "user_id": 42,
  "storefront_id": 1,
  "listing_id": 100,
  "quantity": 2
}' localhost:50051 listingssvc.v1.OrderService/AddToCart
```

Save the `cart.id` from response (e.g., 5).

### Step 2: View Cart

```bash
grpcurl -plaintext -d '{
  "user_id": 42,
  "storefront_id": 1
}' localhost:50051 listingssvc.v1.OrderService/GetCart
```

### Step 3: Create Order

```bash
grpcurl -plaintext -d '{
  "user_id": 42,
  "storefront_id": 1,
  "cart_id": 5,
  "shipping_address": {
    "street": "123 Main St",
    "city": "Belgrade",
    "postal_code": "11000",
    "country": "RS"
  },
  "shipping_method": "standard",
  "payment_method": "card"
}' localhost:50051 listingssvc.v1.OrderService/CreateOrder
```

Save the `order.id` from response (e.g., 1).

### Step 4: View Order

```bash
grpcurl -plaintext -d '{
  "order_id": 1,
  "user_id": 42
}' localhost:50051 listingssvc.v1.OrderService/GetOrder
```

### Step 5: Cancel Order (if needed)

```bash
grpcurl -plaintext -d '{
  "order_id": 1,
  "user_id": 42,
  "reason": "Changed my mind",
  "refund": true
}' localhost:50051 listingssvc.v1.OrderService/CancelOrder
```

## All Available Methods

### Cart Operations

#### AddToCart
```bash
grpcurl -plaintext -d '{
  "user_id": 42,
  "storefront_id": 1,
  "listing_id": 100,
  "quantity": 2
}' localhost:50051 listingssvc.v1.OrderService/AddToCart
```

#### GetCart
```bash
grpcurl -plaintext -d '{
  "user_id": 42,
  "storefront_id": 1
}' localhost:50051 listingssvc.v1.OrderService/GetCart
```

#### ClearCart
```bash
grpcurl -plaintext -d '{
  "user_id": 42,
  "storefront_id": 1
}' localhost:50051 listingssvc.v1.OrderService/ClearCart
```

#### GetUserCarts (all storefronts)
```bash
grpcurl -plaintext -d '{
  "user_id": 42
}' localhost:50051 listingssvc.v1.OrderService/GetUserCarts
```

### Order Operations

#### CreateOrder
```bash
grpcurl -plaintext -d '{
  "user_id": 42,
  "storefront_id": 1,
  "cart_id": 5,
  "shipping_address": {
    "street": "123 Main St",
    "city": "Belgrade",
    "postal_code": "11000",
    "country": "RS"
  },
  "shipping_method": "standard",
  "payment_method": "card"
}' localhost:50051 listingssvc.v1.OrderService/CreateOrder
```

#### GetOrder
```bash
grpcurl -plaintext -d '{
  "order_id": 1,
  "user_id": 42
}' localhost:50051 listingssvc.v1.OrderService/GetOrder
```

#### ListOrders
```bash
grpcurl -plaintext -d '{
  "user_id": 42,
  "page": 1,
  "page_size": 10
}' localhost:50051 listingssvc.v1.OrderService/ListOrders
```

#### CancelOrder
```bash
grpcurl -plaintext -d '{
  "order_id": 1,
  "user_id": 42,
  "reason": "Changed my mind",
  "refund": true
}' localhost:50051 listingssvc.v1.OrderService/CancelOrder
```

#### UpdateOrderStatus (Admin)
```bash
grpcurl -plaintext -d '{
  "order_id": 1,
  "new_status": "ORDER_STATUS_CONFIRMED"
}' localhost:50051 listingssvc.v1.OrderService/UpdateOrderStatus
```

#### GetOrderStats (Admin)
```bash
grpcurl -plaintext -d '{
  "storefront_id": 1
}' localhost:50051 listingssvc.v1.OrderService/GetOrderStats
```

## Order Status Values

Use these values for `new_status` in UpdateOrderStatus:

- `ORDER_STATUS_PENDING` - Order created, awaiting payment
- `ORDER_STATUS_CONFIRMED` - Payment successful
- `ORDER_STATUS_PROCESSING` - Order being prepared
- `ORDER_STATUS_SHIPPED` - Order shipped
- `ORDER_STATUS_DELIVERED` - Order delivered
- `ORDER_STATUS_CANCELLED` - Order cancelled
- `ORDER_STATUS_REFUNDED` - Payment refunded
- `ORDER_STATUS_FAILED` - Order failed

## Anonymous Cart (Session-Based)

For anonymous users, use `session_id` instead of `user_id`:

```bash
grpcurl -plaintext -d '{
  "session_id": "session_abc123",
  "storefront_id": 1,
  "listing_id": 100,
  "quantity": 1
}' localhost:50051 listingssvc.v1.OrderService/AddToCart
```

## Common Error Codes

| Error Code | Meaning | Solution |
|------------|---------|----------|
| `InvalidArgument` | Validation error | Check required fields |
| `NotFound` | Cart/Order/Listing not found | Verify IDs |
| `FailedPrecondition` | Insufficient stock / Price changed | Check inventory |
| `PermissionDenied` | Unauthorized access | Verify user_id |
| `Internal` | Server error | Check logs |

## Debugging

### List all available services
```bash
grpcurl -plaintext localhost:50051 list
```

### Describe OrderService
```bash
grpcurl -plaintext localhost:50051 describe listingssvc.v1.OrderService
```

### Describe specific method
```bash
grpcurl -plaintext localhost:50051 describe listingssvc.v1.OrderService.AddToCart
```

### View server reflection
```bash
grpcurl -plaintext localhost:50051 list listingssvc.v1.OrderService
```

## Interactive Mode (evans)

Install evans: `go install github.com/ktr0731/evans@latest`

```bash
# Start evans
evans -p 50051 repl

# Inside evans:
service listingssvc.v1.OrderService
call AddToCart
# (enter values interactively)
```

## Metrics

View Prometheus metrics:
```bash
curl http://localhost:8080/metrics | grep grpc_request
```

Monitor specific method:
```bash
curl http://localhost:8080/metrics | grep 'grpc_request_duration_seconds{method="CreateOrder"}'
```

## Logs

Logs are structured JSON (zerolog):

```json
{
  "level": "info",
  "component": "grpc_orders_handler",
  "method": "/OrderService/CreateOrder",
  "user_id": 42,
  "cart_id": 5,
  "message": "CreateOrder called",
  "time": "2025-11-14T01:30:00Z"
}
```

Filter logs by method:
```bash
journalctl -u listings-service | grep "CreateOrder"
```

## Testing with Real Data

### 1. Populate test listings
```bash
# Use existing listings service to create products
grpcurl -plaintext -d '{
  "storefront_id": 1,
  "name": "Test Product",
  "price": 29.99,
  "stock_quantity": 100
}' localhost:50051 listingssvc.v1.ListingsService.CreateProduct
```

### 2. Add to cart and create order
Follow the "Testing Flow" section above

### 3. Verify in database
```sql
-- Check cart
SELECT * FROM carts WHERE user_id = 42;
SELECT * FROM cart_items WHERE cart_id = 5;

-- Check order
SELECT * FROM orders WHERE user_id = 42;
SELECT * FROM order_items WHERE order_id = 1;

-- Check reservations
SELECT * FROM inventory_reservations WHERE order_id = 1;
```

## Performance Testing

### Benchmark with ghz

Install: `go install github.com/bojand/ghz/cmd/ghz@latest`

```bash
# Benchmark AddToCart
ghz --insecure \
  --proto /p/github.com/sveturs/listings/api/proto/listings/v1/orders.proto \
  --call listingssvc.v1.OrderService/AddToCart \
  -d '{"user_id": 42, "storefront_id": 1, "listing_id": 100, "quantity": 1}' \
  -c 10 -n 1000 \
  localhost:50051
```

Expected results:
- Latency p99: < 100ms
- Requests/sec: > 500
- Error rate: < 1%

## Security

### Authentication (when enabled)

Add metadata for authentication:

```bash
grpcurl -plaintext \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -d '{"user_id": 42, ...}' \
  localhost:50051 listingssvc.v1.OrderService/AddToCart
```

### Rate Limiting

Default limits:
- 100 requests per minute per user
- 1000 requests per minute global

Exceeded rate limit returns:
```json
{
  "code": "RESOURCE_EXHAUSTED",
  "message": "rate limit exceeded"
}
```

## Troubleshooting

### "rpc error: code = Unimplemented"
**Solution**: OrderService not registered. Check main.go integration.

### "cart not found"
**Solution**: Cart created on first AddToCart. Verify storefront_id matches.

### "insufficient stock"
**Solution**: Check product stock: `SELECT stock_quantity FROM products WHERE id = 100`

### "price changed"
**Solution**: Expected behavior. Product price changed since adding to cart.

### "unauthorized"
**Solution**: user_id doesn't match order owner. Verify user_id parameter.

## Next Steps

1. ✅ Read full documentation: `ORDER_SERVICE_INTEGRATION.md`
2. ⏳ Implement repository layer
3. ⏳ Create database migrations
4. ⏳ Run integration tests
5. ⏳ Deploy to staging

---

**Quick Start Complete** - You now have working Order Service gRPC API!
