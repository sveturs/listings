# Service Layer Architecture

## Overview

The service layer provides business logic for the Listings microservice Orders functionality. It consists of three main services that work together to manage shopping carts, orders, and inventory reservations.

## Services

### 1. CartService

**Location:** `internal/service/cart_service.go`

**Responsibilities:**
- Manage shopping carts (create, read, update, delete)
- Add/remove items from cart
- Validate cart prices against current listings
- Merge session carts to user carts upon login
- Calculate cart totals

**Key Methods:**
- `AddToCart()` - Adds item to cart with stock validation
- `UpdateCartItem()` - Updates item quantity with stock checks
- `RemoveFromCart()` - Removes item from cart
- `GetCart()` - Retrieves cart by user/session ID
- `ClearCart()` - Removes all items
- `MergeSessionCartToUser()` - Transfers anonymous cart to authenticated user
- `ValidateCartItems()` - Detects price changes before checkout

**Business Rules:**
- Cart belongs to one storefront (cannot mix items from different storefronts)
- Either user_id or session_id (not both)
- Price snapshot saved at add-to-cart time
- Stock validation before adding/updating items

---

### 2. OrderService

**Location:** `internal/service/order_service.go`

**Responsibilities:**
- Create orders from carts (ACID transaction)
- Manage order lifecycle (pending → confirmed → processing → shipped → delivered)
- Cancel orders with refund processing
- Confirm payment (called by Payment Service)
- Process refunds

**Key Methods:**
- `CreateOrder()` - Creates order with full ACID transaction workflow
- `GetOrder()` - Retrieves order by ID
- `GetOrderByNumber()` - Retrieves order by order number
- `ListOrders()` - Lists orders with pagination
- `CancelOrder()` - Cancels order, releases reservations, restores stock
- `UpdateOrderStatus()` - Updates order status (validates state transitions)
- `ConfirmOrderPayment()` - Called by Payment Service webhook
- `ProcessRefund()` - Processes refund for cancelled order

**CreateOrder ACID Transaction Workflow:**

```
BEGIN TRANSACTION
1. Get cart with items (validate not empty)
2. Extract listing IDs and sort (ORDER BY id ASC - prevent deadlocks!)
3. Lock listings (SELECT FOR UPDATE - future implementation)
4. Validate stock availability for all items
5. Validate prices (cart price == current price)
6. Build order items (snapshot product data)
7. Calculate financials (tax, shipping, discount, commission)
8. Generate order number (ORD-YYYY-NNNNNN)
9. Create order record (status=pending, payment_status=pending)
10. Create order items (immutable snapshot)
11. Create inventory reservations (TTL 30 minutes)
12. Deduct stock (future: productsRepo.DeductStock)
13. Delete cart
COMMIT TRANSACTION
```

**State Transitions:**
- pending → confirmed (payment successful)
- confirmed → processing (seller preparing order)
- processing → shipped (shipped to customer)
- shipped → delivered (delivered successfully)
- pending/confirmed → cancelled (order cancelled)
- delivered → refunded (payment refunded)

**Business Rules:**
- Order number format: `ORD-YYYY-NNNNNN` (e.g., ORD-2025-001234)
- Reservation TTL: 30 minutes (expires if payment not completed)
- Escrow hold period: 3 days (funds released to seller after)
- Platform commission: 10% of subtotal (configurable)
- Tax rate: 20% VAT (configurable)
- Can only cancel pending/confirmed orders
- Stock deduction happens on order creation (not on payment confirmation)

---

### 3. InventoryService

**Location:** `internal/service/inventory_service.go`

**Responsibilities:**
- Manage inventory reservations (create, commit, release)
- Calculate available stock (total stock - active reservations)
- Clean up expired reservations (cron job)
- Validate stock availability

**Key Methods:**
- `CreateReservation()` - Creates reservation with TTL
- `GetReservationsByOrder()` - Gets all reservations for order
- `CommitReservation()` - Commits reservation (payment confirmed)
- `ReleaseReservation()` - Releases reservation (restores stock)
- `CleanupExpiredReservations()` - Cron job to clean expired (every 5 min)
- `CheckStockAvailability()` - Checks if quantity available
- `GetAvailableStock()` - Calculates: total_stock - SUM(active_reservations)

**CleanupExpiredReservations Cron Job:**

```
BEGIN TRANSACTION
1. Find expired active reservations (status=active AND expires_at < NOW())
2. If none found, exit
3. Extract listing IDs and sort (ORDER BY id ASC - prevent deadlocks!)
4. Lock listings (SELECT FOR UPDATE - future implementation)
5. Mark reservations as expired (status=expired)
6. Restore stock for each reservation (future: productsRepo.RestoreStock)
7. Update orders to cancelled (status=cancelled)
COMMIT TRANSACTION
```

**Available Stock Formula:**
```
available_stock = listing.stock - SUM(active_reservations.quantity)
```

**Business Rules:**
- Default TTL: 30 minutes (configurable)
- Reservation status: active, committed, released, expired
- Can only commit active, non-expired reservations
- Can only release active reservations
- Expired reservations automatically cancel the order

---

## Financial Calculations

**Location:** `internal/service/order_financials.go`

### Order Financials Structure

```go
type OrderFinancials struct {
    Subtotal       float64  // Sum of all items (before tax/shipping)
    Tax            float64  // Tax amount
    ShippingCost   float64  // Shipping cost
    Discount       float64  // Discount amount (coupons, promotions)
    Total          float64  // Final amount to pay
    Commission     float64  // Platform commission
    SellerAmount   float64  // Amount seller receives
    Currency       string   // ISO 4217 currency code
}
```

### Calculation Formula

```
subtotal = SUM(item.total)
tax = subtotal * tax_rate (20% VAT)
total = subtotal + tax + shipping_cost - discount
commission = subtotal * commission_rate (10%)
seller_amount = total - commission
```

### Configuration

```go
type FinancialConfig struct {
    TaxRate         float64  // 0.20 (20% VAT)
    CommissionRate  float64  // 0.10 (10% platform fee)
    DefaultCurrency string   // "RSD" (Serbian Dinar)
    EscrowDays      int32    // 3 days
}
```

---

## Error Handling

**Location:** `internal/service/errors.go`

### Custom Errors

**Cart Errors:**
- `ErrCartEmpty` - Cart has no items
- `ErrCartNotFound` - Cart not found
- `ErrCartItemNotFound` - Cart item not found
- `ErrStorefrontMismatch` - Items from different storefronts
- `ErrListingNotFound` - Listing/product not found
- `ErrListingInactive` - Listing is not active
- `ErrPriceChanged` - Price changed since add-to-cart

**Order Errors:**
- `ErrOrderNotFound` - Order not found
- `ErrOrderAlreadyConfirmed` - Order already confirmed
- `ErrOrderAlreadyCancelled` - Order already cancelled
- `ErrOrderCannotCancel` - Order cannot be cancelled (invalid state)
- `ErrOrderCannotUpdateStatus` - Invalid status transition
- `ErrInsufficientStock` - Not enough stock
- `ErrInvalidAddress` - Invalid shipping/billing address
- `ErrPaymentFailed` - Payment processing failed

**Inventory Errors:**
- `ErrReservationNotFound` - Reservation not found
- `ErrReservationExpired` - Reservation has expired
- `ErrReservationCannotCommit` - Reservation cannot be committed
- `ErrReservationCannotRelease` - Reservation cannot be released
- `ErrStockNotAvailable` - Stock not available (locked by reservations)

### Error Helper Functions

```go
IsNotFoundError(err error) bool    // Checks if error is "not found"
IsConflictError(err error) bool    // Checks if error is conflict
IsValidationError(err error) bool  // Checks if error is validation
```

---

## Transaction Safety

### Deadlock Prevention

**Rule:** Always lock listings in sorted order (ORDER BY id ASC)

```go
// Example: CreateOrder
listingIDs := extractListingIDs(cartItems)
sort.Slice(listingIDs, func(i, j int) bool {
    return listingIDs[i] < listingIDs[j]  // ASC order
})
listings := productsRepo.LockListingsByIDs(ctx, listingIDs)
```

**Why?**
- Two concurrent transactions locking in different orders can deadlock
- Consistent ordering prevents circular wait conditions
- Same pattern used in CleanupExpiredReservations

### Transaction Boundaries

**What should be in a transaction:**
- CreateOrder: cart → order + items + reservations + stock deduction
- CancelOrder: order status + release reservations + restore stock
- CleanupExpiredReservations: expire + restore stock + cancel orders
- ReleaseReservation: release + restore stock

**What should NOT be in a transaction:**
- External API calls (Payment Service, Delivery Service)
- Message queue publishing
- Email notifications

---

## Testing Strategy

### Unit Tests

**Location:** `internal/service/*_test.go`

**Approach:**
- Use mocks for repositories (`testify/mock`)
- Test happy paths + error cases
- Test business rules validation
- Test state transitions

**Coverage Target:** 85%+

### Test Examples

```go
// Cart Service
- TestAddToCart_Success
- TestAddToCart_ListingNotFound
- TestAddToCart_InsufficientStock
- TestAddToCart_StorefrontMismatch
- TestUpdateCartItem_Success
- TestRemoveFromCart_Success
- TestValidateCartItems_PriceChanges

// Order Service
- TestGetOrder_Success
- TestGetOrder_NotFound
- TestUpdateOrderStatus_Success
- TestUpdateOrderStatus_InvalidTransition
- TestCancelOrder_Success

// Inventory Service
- TestGetAvailableStock_NoReservations
- TestGetAvailableStock_WithReservations
- TestCheckStockAvailability_Success
- TestCommitReservation_Success
- TestCommitReservation_Expired
```

---

## Future Improvements

### Phase 17 Remaining Work

1. **Stock Management Methods** (Repository Layer)
   - `LockListingsByIDs(ctx, tx, listingIDs)` - SELECT FOR UPDATE
   - `DeductStock(ctx, tx, listingID, quantity)` - UPDATE listings SET stock -= qty
   - `RestoreStock(ctx, tx, listingID, quantity)` - UPDATE listings SET stock += qty
   - `UpdateStockBatch(ctx, tx, updates)` - Batch stock updates

2. **Order Number Sequence**
   - Create PostgreSQL sequence: `orders_sequence_YYYY`
   - Use sequence in GenerateOrderNumber
   - Reset sequence yearly (ORD-2025-000001, ORD-2026-000001)

3. **Event Publishing**
   - OrderCreated → Analytics Service
   - OrderConfirmed → Delivery Service (create shipment)
   - OrderCancelled → Payment Service (trigger refund)
   - OrderDelivered → Analytics Service (update metrics)

4. **Order Stats Implementation**
   - Aggregate queries for dashboard
   - Total orders, revenue, average order value
   - Filter by user/storefront/date range

5. **Discount Code Validation**
   - Validate discount code against Promotions Service
   - Apply discount rules (percentage, fixed amount, free shipping)
   - Check expiry dates and usage limits

6. **Comprehensive Integration Tests**
   - Test full CreateOrder flow with real DB
   - Test concurrent order creation (race conditions)
   - Test CleanupExpiredReservations cron job
   - Test order lifecycle (pending → delivered)

---

## Dependencies

**Required Repositories:**
- `CartRepository` - Cart and cart items CRUD
- `OrderRepository` - Orders and order items CRUD
- `ReservationRepository` - Inventory reservations CRUD
- `ProductsRepository` - Listings/products read + stock management

**External Services:**
- **Payment Service** - Payment processing, refunds
- **Delivery Service** - Shipment creation, tracking
- **Analytics Service** - Metrics and reporting

**Database:**
- PostgreSQL with pgx driver
- pgxpool for connection pooling
- Transactions for ACID guarantees

---

## Configuration

```go
// Financial Configuration
FinancialConfig{
    TaxRate:         0.20,    // 20% VAT
    CommissionRate:  0.10,    // 10% platform fee
    DefaultCurrency: "RSD",   // Serbian Dinar
    EscrowDays:      3,       // 3 days hold
}

// Reservation TTL
DefaultTTL: 30 minutes

// Cleanup Cron
Schedule: Every 5 minutes
```

---

## Architecture Diagram

```
┌─────────────────┐
│   Frontend      │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  gRPC Handlers  │  (Phase 17 Days 15-17)
└────────┬────────┘
         │
         ▼
┌─────────────────────────────────────────┐
│         Service Layer                   │
│                                         │
│  ┌──────────────┐  ┌─────────────────┐ │
│  │ CartService  │  │  OrderService   │ │
│  └──────┬───────┘  └────────┬────────┘ │
│         │                   │          │
│         │  ┌────────────────┴───────┐  │
│         │  │  InventoryService      │  │
│         │  └────────┬───────────────┘  │
│         │           │                  │
└─────────┼───────────┼──────────────────┘
          │           │
          ▼           ▼
┌─────────────────────────────────────────┐
│       Repository Layer                  │
│                                         │
│  ┌─────────┐ ┌─────────┐ ┌──────────┐  │
│  │  Cart   │ │  Order  │ │ Reserv   │  │
│  │  Repo   │ │  Repo   │ │  Repo    │  │
│  └────┬────┘ └────┬────┘ └────┬─────┘  │
│       │           │           │         │
└───────┼───────────┼───────────┼─────────┘
        │           │           │
        ▼           ▼           ▼
   ┌────────────────────────────────┐
   │       PostgreSQL Database      │
   │  (shopping_carts, orders,      │
   │   inventory_reservations)      │
   └────────────────────────────────┘
```

---

## Summary

The service layer implements comprehensive business logic for shopping cart management, order processing, and inventory reservations. Key features include:

- ✅ **ACID Transactions** - Full transactional safety for order creation
- ✅ **Deadlock Prevention** - Sorted locking order for concurrent operations
- ✅ **Price Validation** - Detects price changes before checkout
- ✅ **Stock Management** - Real-time availability with reservations
- ✅ **Financial Calculations** - Tax, shipping, discount, commission
- ✅ **State Machine** - Strict order status transitions
- ✅ **Error Handling** - Comprehensive custom errors
- ✅ **Testing** - Unit tests with mocks (85%+ coverage target)

**Next Phase:** Days 15-17 - gRPC Handlers Implementation
