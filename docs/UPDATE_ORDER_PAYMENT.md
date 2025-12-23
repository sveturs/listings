# UpdateOrderPayment gRPC Method

## Обзор

Метод `UpdateOrderPayment` добавлен в Listings Service для обновления платёжной информации заказа. Используется Payment Service после успешной обработки платежа.

## Дата добавления

2025-12-20

## Proto Definition

```protobuf
// UpdateOrderPaymentRequest - update payment information for an order
message UpdateOrderPaymentRequest {
  int64 order_id = 1;                         // Required - order to update
  optional string payment_provider = 2;       // Payment provider (stripe, allsecure, etc.)
  optional string payment_session_id = 3;     // Checkout session ID
  optional string payment_intent_id = 4;      // Payment intent ID
  optional string payment_idempotency_key = 5; // Idempotency key for payment
  optional string payment_status = 6;         // Payment status (pending, paid, failed, etc.)
  optional string payment_transaction_id = 7; // Transaction ID from payment provider
}

message UpdateOrderPaymentResponse {
  bool success = 1;                           // True if update successful
  Order order = 2;                            // Updated order with new payment info
}

rpc UpdateOrderPayment(UpdateOrderPaymentRequest) returns (UpdateOrderPaymentResponse);
```

## Database Schema

### Новые поля в таблице `orders`

Миграция: `20251220000002_add_order_payment_fields.up.sql`

```sql
ALTER TABLE orders
    ADD COLUMN IF NOT EXISTS payment_provider VARCHAR(50);

ALTER TABLE orders
    ADD COLUMN IF NOT EXISTS payment_session_id VARCHAR(255);

ALTER TABLE orders
    ADD COLUMN IF NOT EXISTS payment_intent_id VARCHAR(255);

ALTER TABLE orders
    ADD COLUMN IF NOT EXISTS payment_idempotency_key VARCHAR(255);
```

### Индексы

- `idx_orders_payment_session` - поиск по payment_session_id (для webhook обработки)
- `idx_orders_idempotency_key` - UNIQUE constraint для предотвращения дубликатов
- `idx_orders_provider_session` - composite index для payment reconciliation
- `idx_orders_payment_intent` - поиск по payment_intent_id

## Использование

### gRPC (Go)

```go
import (
    listingspb "github.com/vondi-global/listings/api/proto/listings/v1"
)

// Обновить платёжную информацию
req := &listingspb.UpdateOrderPaymentRequest{
    OrderId:             12345,
    PaymentProvider:     strPtr("stripe"),
    PaymentSessionId:    strPtr("cs_test_abc123"),
    PaymentIntentId:     strPtr("pi_test_xyz789"),
    PaymentStatus:       strPtr("paid"),
    PaymentTransactionId: strPtr("txn_abc123"),
}

resp, err := client.UpdateOrderPayment(ctx, req)
if err != nil {
    log.Fatalf("Failed to update payment: %v", err)
}

fmt.Printf("Payment updated successfully: %v\n", resp.Success)
fmt.Printf("Updated order: %+v\n", resp.Order)
```

### Валидация

Метод требует:
1. `order_id > 0` (обязательно)
2. Хотя бы одно поле для обновления (необязательно все)

### Частичное обновление

Можно обновить только нужные поля:

```go
// Обновить только payment_status
req := &listingspb.UpdateOrderPaymentRequest{
    OrderId:       12345,
    PaymentStatus: strPtr("paid"),
}

// Обновить только payment_session_id и payment_intent_id
req := &listingspb.UpdateOrderPaymentRequest{
    OrderId:          12345,
    PaymentSessionId: strPtr("cs_test_abc123"),
    PaymentIntentId:  strPtr("pi_test_xyz789"),
}
```

## Service Layer

### Interface

```go
type OrderService interface {
    UpdatePaymentInfo(ctx context.Context, orderID int64, req *UpdatePaymentInfoRequest) (*domain.Order, error)
}

type UpdatePaymentInfoRequest struct {
    PaymentProvider       *string
    PaymentSessionID      *string
    PaymentIntentID       *string
    PaymentIdempotencyKey *string
    PaymentStatus         *string
    PaymentTransactionID  *string
}
```

## Repository Layer

### Interface

```go
type UpdatePaymentInfoParams struct {
    PaymentProvider       *string
    PaymentSessionID      *string
    PaymentIntentID       *string
    PaymentIdempotencyKey *string
    PaymentStatus         *string
    PaymentTransactionID  *string
}

type OrderRepository interface {
    UpdatePaymentInfo(ctx context.Context, orderID int64, params UpdatePaymentInfoParams) error
}
```

### SQL Query

Repository использует динамический SQL query builder для обновления только переданных полей:

```sql
UPDATE orders SET
    payment_provider = $1,
    payment_session_id = $2,
    payment_intent_id = $3
WHERE id = $4
RETURNING updated_at
```

## Идемпотентность

Для предотвращения дубликатов используйте `payment_idempotency_key`:

```go
req := &listingspb.UpdateOrderPaymentRequest{
    OrderId:               12345,
    PaymentIdempotencyKey: strPtr("pay_20251220_12345_abc123"),
    PaymentStatus:         strPtr("paid"),
}
```

Если ключ уже существует - БД вернёт ошибку UNIQUE constraint violation.

## Error Handling

Возможные ошибки:

- `codes.InvalidArgument` - order_id <= 0 или нет полей для обновления
- `codes.NotFound` - заказ не найден
- `codes.Internal` - ошибка БД или внутренняя ошибка

## Domain Model Updates

### Order struct

Добавлены новые поля:

```go
type Order struct {
    // ...existing fields...

    // Payment gateway fields
    PaymentProvider         *string `json:"payment_provider,omitempty" db:"payment_provider"`
    PaymentSessionID        *string `json:"payment_session_id,omitempty" db:"payment_session_id"`
    PaymentIntentID         *string `json:"payment_intent_id,omitempty" db:"payment_intent_id"`
    PaymentIdempotencyKey   *string `json:"payment_idempotency_key,omitempty" db:"payment_idempotency_key"`
}
```

### ToProto Method

Автоматически конвертирует новые поля в protobuf Order:

```go
func (o *Order) ToProto() *pb.Order {
    // ...
    if o.PaymentProvider != nil {
        pbOrder.PaymentProvider = o.PaymentProvider
    }
    if o.PaymentSessionID != nil {
        pbOrder.PaymentSessionId = o.PaymentSessionID
    }
    // ...
}
```

## Примеры использования

### Stripe Payment Flow

```go
// 1. После создания checkout session
req := &listingspb.UpdateOrderPaymentRequest{
    OrderId:          order.ID,
    PaymentProvider:  strPtr("stripe"),
    PaymentSessionId: strPtr(session.ID),
    PaymentStatus:    strPtr("pending"),
}
client.UpdateOrderPayment(ctx, req)

// 2. После успешного платежа (webhook)
req = &listingspb.UpdateOrderPaymentRequest{
    OrderId:          order.ID,
    PaymentIntentId:  strPtr(paymentIntent.ID),
    PaymentStatus:    strPtr("paid"),
    PaymentTransactionId: strPtr(charge.BalanceTransaction),
}
client.UpdateOrderPayment(ctx, req)
```

### AllSecure Payment Flow

```go
// После успешного платежа
req := &listingspb.UpdateOrderPaymentRequest{
    OrderId:               order.ID,
    PaymentProvider:       strPtr("allsecure"),
    PaymentIntentId:       strPtr(transaction.ID),
    PaymentStatus:         strPtr("paid"),
    PaymentIdempotencyKey: strPtr(fmt.Sprintf("allsecure_%s", transaction.ID)),
}
client.UpdateOrderPayment(ctx, req)
```

## Changelog

- **2025-12-20**: Метод добавлен для интеграции с Payment Service
  - Proto messages: `UpdateOrderPaymentRequest`, `UpdateOrderPaymentResponse`
  - Database migration: `20251220000002_add_order_payment_fields`
  - Service method: `OrderService.UpdatePaymentInfo`
  - Repository method: `OrderRepository.UpdatePaymentInfo`
  - gRPC handler: `Server.UpdateOrderPayment`

## См. также

- [Order Proto Definition](../api/proto/listings/v1/orders.proto)
- [Migration Script](../migrations/20251220000002_add_order_payment_fields.up.sql)
- [Service Implementation](../internal/service/order_service.go)
- [Repository Implementation](../internal/repository/postgres/order_repository.go)
- [gRPC Handler](../internal/transport/grpc/handlers_order_grpc.go)
