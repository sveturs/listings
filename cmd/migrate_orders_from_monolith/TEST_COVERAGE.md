# Migration Test Coverage Report

## Summary

| Metric | Value |
|--------|-------|
| **Total Tests** | 10 |
| **Helper Functions** | 12 |
| **Lines of Code** | 767 |
| **Compilation Status** | ✅ Success |
| **Ready to Run** | 🔄 Waiting for migration script |

## Test Scenarios Covered

### 1. Core Migration Tests (6 tests)

| Test | Description | Coverage |
|------|-------------|----------|
| `TestMigrateShoppingCarts` | Active carts migration, date filtering | Carts created < 30 days |
| `TestMigrateCartItems` | Cart items with FK constraints | All items from active carts |
| `TestMigrateOrders` | Recent orders migration, date filtering | Orders created < 12 months |
| `TestMigrateOrderItems` | Order items with snapshots | All items from recent orders |
| `TestMigrateInventoryReservations` | Active reservations with TTL | Reservations with expires_at > NOW() |
| `TestMigrationDryRun` | Dry-run mode validation | No DB modifications |

### 2. Quality Assurance Tests (3 tests)

| Test | Description | What It Validates |
|------|-------------|-------------------|
| `TestIdempotency` | Multiple runs safety | No duplicates on re-run |
| `TestFKIntegrity` | Foreign key constraints | All FKs valid after migration |
| `TestRollback` | Transaction rollback | Data consistency on error |

### 3. Integration Test (1 test)

| Test | Description | Dataset Size |
|------|-------------|--------------|
| `TestFullE2EMigration` | End-to-end comprehensive test | 10 carts, 50 cart items, 20 orders, 100+ order items, 5 reservations |

## Helper Functions

### Test Setup (3 functions)
1. `setupMigrationTest()` - Creates dual-DB test environment
2. `setupMonolithSchema()` - Creates source schema
3. `cleanup()` - Tears down test DBs

### Data Creation (5 functions)
4. `createTestCart()` - Creates shopping cart
5. `createTestCartItem()` - Creates cart item with price snapshot
6. `createTestOrder()` - Creates order with financials
7. `createTestOrderItem()` - Creates order item
8. `createTestReservation()` - Creates inventory reservation

### Validation (4 functions)
9. `checkFKIntegrity()` - Validates foreign keys
10. `compareCartRecord()` - Compares cart data
11. `compareOrderRecord()` - Compares order data
12. `countRows()` - Counts table rows

## Test Data Scenarios

### Shopping Carts
- ✅ Active carts (0-29 days old)
- ❌ Old carts (30+ days old)
- ✅ Carts with multiple items
- ✅ Carts with variant data

### Orders
- ✅ Recent orders (0-11 months old)
- ❌ Old orders (12+ months old)
- ✅ Orders with multiple items
- ✅ Orders with financial data
- ✅ Orders with shipping/billing addresses

### Inventory Reservations
- ✅ Active reservations (expires_at > NOW())
- ❌ Expired reservations (expires_at < NOW())
- ✅ Reservations linked to orders
- ✅ Reservations with quantity

## Foreign Key Constraints Validated

### Cart Items
```sql
cart_items.cart_id → shopping_carts.id
```

### Order Items
```sql
order_items.order_id → orders.id
```

### Orphan Checks
- No orphaned cart items (without parent cart)
- No orphaned order items (without parent order)

## Edge Cases Covered

### Date Filtering
- ✅ Exactly 30 days old (boundary)
- ✅ 29 days old (should migrate)
- ✅ 31 days old (should skip)
- ✅ 365 days old (boundary for orders)

### Data Integrity
- ✅ NULL user_id (guest carts/orders)
- ✅ JSONB fields (shipping_address, variant_data)
- ✅ Decimal precision (financial amounts)
- ✅ VARCHAR length limits

### Concurrency
- ✅ Multiple carts for same user
- ✅ Multiple orders in same storefront
- ✅ Multiple reservations for same listing

## E2E Test Dataset Details

### Created Records
```
Monolith DB:
  - 10 shopping_carts (5 active + 5 old)
  - 55 cart_items (50 from active + 5 from old)
  - 20 orders (15 recent + 5 old)
  - 120+ order_items (varied per order)
  - 5 inventory_reservations (3 active + 2 expired)

Expected Migration:
  - 5 shopping_carts → microservice
  - 50 cart_items → microservice
  - 15 orders → microservice
  - 100+ order_items → microservice
  - 3 inventory_reservations → microservice
```

## Test Execution Time (Estimated)

| Test | Estimated Time |
|------|----------------|
| Setup (per test) | ~5s (Docker container) |
| Data creation | ~1s |
| Migration (when implemented) | ~2-5s |
| Validation | ~1s |
| Teardown | ~2s |

**Total per test**: ~10-15 seconds
**Full suite**: ~2-3 minutes

## Coverage Gaps (Intentional)

These scenarios are NOT covered (by design):

1. **Concurrent migrations** - Single-threaded migration assumed
2. **Network failures** - Local Docker containers used
3. **Schema version mismatches** - Assumes migrations are run
4. **Production data volume** - Test dataset is small (10-20 records)
5. **Cross-DB transactions** - PostgreSQL doesn't support this

## Ready for Production?

### ✅ Ready
- Test structure is complete
- All helper functions implemented
- Compilation verified
- Edge cases identified

### 🔄 Waiting
- Migration script from Agent #1
- Uncomment migration calls
- Uncomment assertions
- Run full test suite

## Test Execution Commands

```bash
# Run all tests
go test -tags=migration -v ./cmd/migrate_orders_from_monolith/

# Run specific test
go test -tags=migration -v ./cmd/migrate_orders_from_monolith/ -run TestMigrateShoppingCarts

# Run with coverage
go test -tags=migration -v -coverprofile=coverage.out ./cmd/migrate_orders_from_monolith/
go tool cover -html=coverage.out -o coverage.html

# Run E2E test only
go test -tags=migration -v ./cmd/migrate_orders_from_monolith/ -run TestFullE2EMigration
```

## Maintenance Notes

### When to Update Tests

1. **Schema changes** → Update `setupMonolithSchema()`
2. **Migration rules change** → Update test assertions
3. **New tables added** → Add new test cases
4. **FK constraints change** → Update `checkFKIntegrity()`

### Test Maintenance Checklist

- [ ] Run tests after any schema migration
- [ ] Update test data when business rules change
- [ ] Add new tests for new migration features
- [ ] Keep README.md in sync with tests
- [ ] Document any new edge cases

## Related Files

- **Test file**: `/p/github.com/sveturs/listings/cmd/migrate_orders_from_monolith/migration_test.go`
- **README**: `/p/github.com/sveturs/listings/cmd/migrate_orders_from_monolith/README.md`
- **Migrations**: `/p/github.com/sveturs/listings/migrations/`
- **Test helpers**: `/p/github.com/sveturs/listings/tests/testing.go`
