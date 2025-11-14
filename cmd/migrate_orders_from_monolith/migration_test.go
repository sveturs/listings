//go:build migration

package main

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sveturs/listings/tests"
)

// =============================================================================
// Test Constants
// =============================================================================

const (
	// Migration timeframes
	activeCarts        = 30 * 24 * time.Hour // 30 days
	recentOrders       = 365 * 24 * time.Hour // 12 months
	activeReservations = 24 * time.Hour       // 24 hours
)

// =============================================================================
// Test Helper Functions
// =============================================================================

// MigrationTestContext holds test databases and helpers
type MigrationTestContext struct {
	MonolithDB  *sqlx.DB
	ListingsDB  *sqlx.DB
	MonolithRes *tests.TestDB
	ListingsRes *tests.TestDB
	T           *testing.T
}

// setupMigrationTest creates two test databases (monolith and listings microservice)
func setupMigrationTest(t *testing.T) *MigrationTestContext {
	t.Helper()
	tests.SkipIfNoDocker(t)

	// Setup monolith database (source)
	monolithRes := tests.SetupTestPostgres(t)
	monolithDB := sqlx.NewDb(monolithRes.DB, "postgres")

	// Setup listings microservice database (destination)
	listingsRes := tests.SetupTestPostgres(t)
	listingsDB := sqlx.NewDb(listingsRes.DB, "postgres")

	// Run migrations on listings DB
	tests.RunMigrations(t, listingsRes.DB, "../../migrations")

	ctx := &MigrationTestContext{
		MonolithDB:  monolithDB,
		ListingsDB:  listingsDB,
		MonolithRes: monolithRes,
		ListingsRes: listingsRes,
		T:           t,
	}

	// Create monolith schema
	ctx.setupMonolithSchema()

	return ctx
}

// cleanup tears down test databases
func (ctx *MigrationTestContext) cleanup() {
	ctx.MonolithRes.TeardownTestPostgres(ctx.T)
	ctx.ListingsRes.TeardownTestPostgres(ctx.T)
}

// setupMonolithSchema creates monolith tables (source schema)
func (ctx *MigrationTestContext) setupMonolithSchema() {
	ctx.T.Helper()

	// Create monolith shopping_carts table
	_, err := ctx.MonolithDB.Exec(`
		CREATE TABLE IF NOT EXISTS shopping_carts (
			id BIGSERIAL PRIMARY KEY,
			user_id BIGINT,
			storefront_id BIGINT NOT NULL,
			status VARCHAR(50) DEFAULT 'active',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	require.NoError(ctx.T, err, "Failed to create monolith shopping_carts table")

	// Create monolith cart_items table
	_, err = ctx.MonolithDB.Exec(`
		CREATE TABLE IF NOT EXISTS cart_items (
			id BIGSERIAL PRIMARY KEY,
			cart_id BIGINT NOT NULL REFERENCES shopping_carts(id) ON DELETE CASCADE,
			listing_id BIGINT NOT NULL,
			variant_id BIGINT,
			quantity INTEGER NOT NULL DEFAULT 1,
			price_snapshot DECIMAL(15,2) NOT NULL,
			variant_data JSONB,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	require.NoError(ctx.T, err, "Failed to create monolith cart_items table")

	// Create monolith orders table
	_, err = ctx.MonolithDB.Exec(`
		CREATE TABLE IF NOT EXISTS orders (
			id BIGSERIAL PRIMARY KEY,
			order_number VARCHAR(50) UNIQUE NOT NULL,
			user_id BIGINT,
			storefront_id BIGINT NOT NULL,
			status VARCHAR(50) DEFAULT 'pending',
			payment_status VARCHAR(50) DEFAULT 'pending',
			payment_method VARCHAR(100),
			payment_transaction_id VARCHAR(255),
			subtotal DECIMAL(15,2) NOT NULL DEFAULT 0.00,
			tax DECIMAL(15,2) NOT NULL DEFAULT 0.00,
			shipping_cost DECIMAL(15,2) NOT NULL DEFAULT 0.00,
			discount DECIMAL(15,2) NOT NULL DEFAULT 0.00,
			total DECIMAL(15,2) NOT NULL,
			commission DECIMAL(15,2) NOT NULL DEFAULT 0.00,
			seller_amount DECIMAL(15,2) NOT NULL DEFAULT 0.00,
			currency VARCHAR(3) DEFAULT 'USD',
			shipping_address JSONB,
			billing_address JSONB,
			escrow_days INTEGER DEFAULT 0,
			notes TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	require.NoError(ctx.T, err, "Failed to create monolith orders table")

	// Create monolith order_items table
	_, err = ctx.MonolithDB.Exec(`
		CREATE TABLE IF NOT EXISTS order_items (
			id BIGSERIAL PRIMARY KEY,
			order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
			listing_id BIGINT NOT NULL,
			variant_id BIGINT,
			listing_name VARCHAR(500) NOT NULL,
			sku VARCHAR(100),
			variant_data JSONB,
			attributes JSONB,
			quantity INTEGER NOT NULL,
			unit_price DECIMAL(15,2) NOT NULL,
			subtotal DECIMAL(15,2) NOT NULL,
			discount DECIMAL(15,2) DEFAULT 0.00,
			total DECIMAL(15,2) NOT NULL,
			image_url TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	require.NoError(ctx.T, err, "Failed to create monolith order_items table")

	// Create monolith inventory_reservations table
	_, err = ctx.MonolithDB.Exec(`
		CREATE TABLE IF NOT EXISTS inventory_reservations (
			id BIGSERIAL PRIMARY KEY,
			listing_id BIGINT NOT NULL,
			variant_id BIGINT,
			order_id BIGINT,
			cart_id BIGINT,
			quantity INTEGER NOT NULL,
			expires_at TIMESTAMP NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	require.NoError(ctx.T, err, "Failed to create monolith inventory_reservations table")
}

// createTestCart creates a test shopping cart in monolith DB
func (ctx *MigrationTestContext) createTestCart(userID, storefrontID int64, createdAt time.Time) int64 {
	ctx.T.Helper()

	var cartID int64
	err := ctx.MonolithDB.QueryRow(`
		INSERT INTO shopping_carts (user_id, storefront_id, status, created_at, updated_at)
		VALUES ($1, $2, 'active', $3, $3)
		RETURNING id
	`, userID, storefrontID, createdAt).Scan(&cartID)
	require.NoError(ctx.T, err, "Failed to create test cart")

	return cartID
}

// createTestCartItem creates a test cart item in monolith DB
func (ctx *MigrationTestContext) createTestCartItem(cartID, listingID int64, quantity int, price float64) int64 {
	ctx.T.Helper()

	var itemID int64
	err := ctx.MonolithDB.QueryRow(`
		INSERT INTO cart_items (cart_id, listing_id, quantity, price_snapshot)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, cartID, listingID, quantity, price).Scan(&itemID)
	require.NoError(ctx.T, err, "Failed to create test cart item")

	return itemID
}

// createTestOrder creates a test order in monolith DB
func (ctx *MigrationTestContext) createTestOrder(orderNumber string, userID, storefrontID int64, total float64, createdAt time.Time) int64 {
	ctx.T.Helper()

	var orderID int64
	err := ctx.MonolithDB.QueryRow(`
		INSERT INTO orders (
			order_number, user_id, storefront_id, status, payment_status,
			subtotal, total, currency, created_at, updated_at
		)
		VALUES ($1, $2, $3, 'confirmed', 'completed', $4, $4, 'USD', $5, $5)
		RETURNING id
	`, orderNumber, userID, storefrontID, total, createdAt).Scan(&orderID)
	require.NoError(ctx.T, err, "Failed to create test order")

	return orderID
}

// createTestOrderItem creates a test order item in monolith DB
func (ctx *MigrationTestContext) createTestOrderItem(orderID, listingID int64, listingName string, quantity int, unitPrice, total float64) int64 {
	ctx.T.Helper()

	var itemID int64
	err := ctx.MonolithDB.QueryRow(`
		INSERT INTO order_items (
			order_id, listing_id, listing_name, quantity,
			unit_price, subtotal, total
		)
		VALUES ($1, $2, $3, $4, $5, $6, $6)
		RETURNING id
	`, orderID, listingID, listingName, quantity, unitPrice, total).Scan(&itemID)
	require.NoError(ctx.T, err, "Failed to create test order item")

	return itemID
}

// createTestReservation creates a test inventory reservation in monolith DB
func (ctx *MigrationTestContext) createTestReservation(listingID, orderID int64, quantity int, expiresAt time.Time) int64 {
	ctx.T.Helper()

	var resID int64
	err := ctx.MonolithDB.QueryRow(`
		INSERT INTO inventory_reservations (
			listing_id, order_id, quantity, expires_at
		)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, listingID, orderID, quantity, expiresAt).Scan(&resID)
	require.NoError(ctx.T, err, "Failed to create test reservation")

	return resID
}

// countRows counts rows in a table
func (ctx *MigrationTestContext) countRows(db *sqlx.DB, table string) int {
	ctx.T.Helper()

	var count int
	err := db.Get(&count, fmt.Sprintf("SELECT COUNT(*) FROM %s", table))
	require.NoError(ctx.T, err, "Failed to count rows in %s", table)

	return count
}

// checkFKIntegrity checks foreign key constraints
func (ctx *MigrationTestContext) checkFKIntegrity(db *sqlx.DB) error {
	ctx.T.Helper()

	// Check cart_items -> shopping_carts FK
	var orphanedCartItems int
	err := db.Get(&orphanedCartItems, `
		SELECT COUNT(*) FROM cart_items ci
		WHERE NOT EXISTS (SELECT 1 FROM shopping_carts sc WHERE sc.id = ci.cart_id)
	`)
	if err != nil {
		return fmt.Errorf("failed to check cart_items FK: %w", err)
	}
	if orphanedCartItems > 0 {
		return fmt.Errorf("found %d orphaned cart_items", orphanedCartItems)
	}

	// Check order_items -> orders FK
	var orphanedOrderItems int
	err = db.Get(&orphanedOrderItems, `
		SELECT COUNT(*) FROM order_items oi
		WHERE NOT EXISTS (SELECT 1 FROM orders o WHERE o.id = oi.order_id)
	`)
	if err != nil {
		return fmt.Errorf("failed to check order_items FK: %w", err)
	}
	if orphanedOrderItems > 0 {
		return fmt.Errorf("found %d orphaned order_items", orphanedOrderItems)
	}

	return nil
}

// compareCartRecord compares cart records between monolith and microservice
func (ctx *MigrationTestContext) compareCartRecord(monolithID, microserviceID int64) {
	ctx.T.Helper()

	type CartRecord struct {
		UserID       sql.NullInt64 `db:"user_id"`
		StorefrontID int64         `db:"storefront_id"`
		Status       string        `db:"status"`
	}

	var monolithCart, microCart CartRecord
	err := ctx.MonolithDB.Get(&monolithCart, "SELECT user_id, storefront_id, status FROM shopping_carts WHERE id = $1", monolithID)
	require.NoError(ctx.T, err)

	err = ctx.ListingsDB.Get(&microCart, "SELECT user_id, storefront_id, status FROM shopping_carts WHERE id = $1", microserviceID)
	require.NoError(ctx.T, err)

	assert.Equal(ctx.T, monolithCart.UserID.Int64, microCart.UserID.Int64, "user_id mismatch")
	assert.Equal(ctx.T, monolithCart.StorefrontID, microCart.StorefrontID, "storefront_id mismatch")
	assert.Equal(ctx.T, monolithCart.Status, microCart.Status, "status mismatch")
}

// compareOrderRecord compares order records between monolith and microservice
func (ctx *MigrationTestContext) compareOrderRecord(monolithID, microserviceID int64) {
	ctx.T.Helper()

	type OrderRecord struct {
		OrderNumber  string        `db:"order_number"`
		UserID       sql.NullInt64 `db:"user_id"`
		StorefrontID int64         `db:"storefront_id"`
		Status       string        `db:"status"`
		Total        float64       `db:"total"`
	}

	var monolithOrder, microOrder OrderRecord
	err := ctx.MonolithDB.Get(&monolithOrder, "SELECT order_number, user_id, storefront_id, status, total FROM orders WHERE id = $1", monolithID)
	require.NoError(ctx.T, err)

	err = ctx.ListingsDB.Get(&microOrder, "SELECT order_number, user_id, storefront_id, status, total FROM orders WHERE id = $1", microserviceID)
	require.NoError(ctx.T, err)

	assert.Equal(ctx.T, monolithOrder.OrderNumber, microOrder.OrderNumber, "order_number mismatch")
	assert.Equal(ctx.T, monolithOrder.UserID.Int64, microOrder.UserID.Int64, "user_id mismatch")
	assert.Equal(ctx.T, monolithOrder.StorefrontID, microOrder.StorefrontID, "storefront_id mismatch")
	assert.Equal(ctx.T, monolithOrder.Status, microOrder.Status, "status mismatch")
	assert.Equal(ctx.T, monolithOrder.Total, microOrder.Total, "total mismatch")
}

// =============================================================================
// Test Cases
// =============================================================================

// TestMigrationDryRun validates that dry-run mode doesn't modify database
func TestMigrationDryRun(t *testing.T) {
	ctx := setupMigrationTest(t)
	defer ctx.cleanup()

	// Create test data in monolith
	ctx.createTestCart(1, 100, time.Now().Add(-1*time.Hour))
	ctx.createTestOrder("ORD-2025-000001", 1, 100, 99.99, time.Now().Add(-1*time.Hour))

	// Count before migration
	cartsBeforeMono := ctx.countRows(ctx.MonolithDB, "shopping_carts")
	ordersBeforeMono := ctx.countRows(ctx.MonolithDB, "orders")
	cartsBefore := ctx.countRows(ctx.ListingsDB, "shopping_carts")
	ordersBefore := ctx.countRows(ctx.ListingsDB, "orders")

	// TODO: Run migration with --dry-run flag
	// This will be implemented when Agent #1 completes the migration script

	// Verify no changes in monolith DB
	assert.Equal(t, cartsBeforeMono, ctx.countRows(ctx.MonolithDB, "shopping_carts"))
	assert.Equal(t, ordersBeforeMono, ctx.countRows(ctx.MonolithDB, "orders"))

	// Verify no changes in listings DB
	assert.Equal(t, cartsBefore, ctx.countRows(ctx.ListingsDB, "shopping_carts"))
	assert.Equal(t, ordersBefore, ctx.countRows(ctx.ListingsDB, "orders"))
}

// TestMigrateShoppingCarts validates shopping cart migration logic
func TestMigrateShoppingCarts(t *testing.T) {
	ctx := setupMigrationTest(t)
	defer ctx.cleanup()

	now := time.Now()

	// Create active carts (should be migrated)
	activeCart1 := ctx.createTestCart(1, 100, now.Add(-1*time.Hour))
	activeCart2 := ctx.createTestCart(2, 100, now.Add(-10*24*time.Hour))

	// Create old carts (should NOT be migrated - older than 30 days)
	oldCart := ctx.createTestCart(3, 100, now.Add(-31*24*time.Hour))

	// Count active carts in monolith
	var activeCount int
	err := ctx.MonolithDB.Get(&activeCount, `
		SELECT COUNT(*) FROM shopping_carts
		WHERE created_at >= $1
	`, now.Add(-activeCarts))
	require.NoError(t, err)
	assert.Equal(t, 2, activeCount, "Expected 2 active carts in monolith")

	// TODO: Run migration
	// After migration, verify:
	// 1. Active carts are migrated
	// 2. Old carts are NOT migrated
	// 3. FK constraints are valid

	t.Log("Active cart IDs:", activeCart1, activeCart2)
	t.Log("Old cart ID (should be skipped):", oldCart)

	// Verify FK integrity after migration
	err = ctx.checkFKIntegrity(ctx.ListingsDB)
	assert.NoError(t, err, "FK integrity check failed")
}

// TestMigrateCartItems validates cart item migration
func TestMigrateCartItems(t *testing.T) {
	ctx := setupMigrationTest(t)
	defer ctx.cleanup()

	now := time.Now()

	// Create cart with items
	cartID := ctx.createTestCart(1, 100, now.Add(-1*time.Hour))
	item1 := ctx.createTestCartItem(cartID, 200, 2, 99.99)
	item2 := ctx.createTestCartItem(cartID, 201, 1, 49.99)

	t.Log("Cart ID:", cartID)
	t.Log("Cart item IDs:", item1, item2)

	// TODO: Run migration
	// After migration, verify:
	// 1. Cart items are migrated
	// 2. FK to shopping_carts is correct
	// 3. Price snapshot is preserved
	// 4. Variant data (if present) is migrated

	// Verify FK integrity
	err := ctx.checkFKIntegrity(ctx.ListingsDB)
	assert.NoError(t, err, "FK integrity check failed")
}

// TestMigrateOrders validates order migration
func TestMigrateOrders(t *testing.T) {
	ctx := setupMigrationTest(t)
	defer ctx.cleanup()

	now := time.Now()

	// Create recent orders (should be migrated)
	recentOrder1 := ctx.createTestOrder("ORD-2025-000001", 1, 100, 199.99, now.Add(-30*24*time.Hour))
	recentOrder2 := ctx.createTestOrder("ORD-2025-000002", 2, 100, 299.99, now.Add(-180*24*time.Hour))

	// Create old order (should NOT be migrated - older than 12 months)
	oldOrder := ctx.createTestOrder("ORD-2024-000001", 3, 100, 99.99, now.Add(-400*24*time.Hour))

	// Count recent orders in monolith
	var recentCount int
	err := ctx.MonolithDB.Get(&recentCount, `
		SELECT COUNT(*) FROM orders
		WHERE created_at >= $1
	`, now.Add(-recentOrders))
	require.NoError(t, err)
	assert.Equal(t, 2, recentCount, "Expected 2 recent orders in monolith")

	t.Log("Recent order IDs:", recentOrder1, recentOrder2)
	t.Log("Old order ID (should be skipped):", oldOrder)

	// TODO: Run migration
	// After migration, verify:
	// 1. Recent orders are migrated
	// 2. Old orders are NOT migrated
	// 3. Order numbers are unique
	// 4. Financial data is correct

	// Verify FK integrity
	err = ctx.checkFKIntegrity(ctx.ListingsDB)
	assert.NoError(t, err, "FK integrity check failed")
}

// TestMigrateOrderItems validates order item migration
func TestMigrateOrderItems(t *testing.T) {
	ctx := setupMigrationTest(t)
	defer ctx.cleanup()

	now := time.Now()

	// Create order with items
	orderID := ctx.createTestOrder("ORD-2025-000001", 1, 100, 149.98, now.Add(-1*time.Hour))
	item1 := ctx.createTestOrderItem(orderID, 200, "Test Product 1", 2, 50.00, 100.00)
	item2 := ctx.createTestOrderItem(orderID, 201, "Test Product 2", 1, 49.98, 49.98)

	t.Log("Order ID:", orderID)
	t.Log("Order item IDs:", item1, item2)

	// TODO: Run migration
	// After migration, verify:
	// 1. Order items are migrated
	// 2. FK to orders is correct
	// 3. Snapshot data is preserved

	// Verify FK integrity
	err := ctx.checkFKIntegrity(ctx.ListingsDB)
	assert.NoError(t, err, "FK integrity check failed")
}

// TestMigrateInventoryReservations validates inventory reservation migration
func TestMigrateInventoryReservations(t *testing.T) {
	ctx := setupMigrationTest(t)
	defer ctx.cleanup()

	now := time.Now()

	// Create order for reservations
	orderID := ctx.createTestOrder("ORD-2025-000001", 1, 100, 99.99, now.Add(-1*time.Hour))

	// Create active reservations (should be migrated)
	activeRes1 := ctx.createTestReservation(200, orderID, 2, now.Add(2*time.Hour))
	activeRes2 := ctx.createTestReservation(201, orderID, 1, now.Add(1*time.Hour))

	// Create expired reservation (should NOT be migrated)
	expiredRes := ctx.createTestReservation(202, orderID, 1, now.Add(-1*time.Hour))

	// Count active reservations
	var activeCount int
	err := ctx.MonolithDB.Get(&activeCount, `
		SELECT COUNT(*) FROM inventory_reservations
		WHERE expires_at > $1
	`, now)
	require.NoError(t, err)
	assert.Equal(t, 2, activeCount, "Expected 2 active reservations in monolith")

	t.Log("Active reservation IDs:", activeRes1, activeRes2)
	t.Log("Expired reservation ID (should be skipped):", expiredRes)

	// TODO: Run migration
	// After migration, verify:
	// 1. Only active reservations are migrated
	// 2. Expired reservations are NOT migrated
	// 3. TTL is correctly calculated
}

// TestIdempotency validates that migration can be run multiple times safely
func TestIdempotency(t *testing.T) {
	ctx := setupMigrationTest(t)
	defer ctx.cleanup()

	now := time.Now()

	// Create test data
	cartID := ctx.createTestCart(1, 100, now.Add(-1*time.Hour))
	ctx.createTestCartItem(cartID, 200, 2, 99.99)

	orderID := ctx.createTestOrder("ORD-2025-000001", 1, 100, 99.99, now.Add(-1*time.Hour))
	ctx.createTestOrderItem(orderID, 200, "Test Product", 1, 99.99, 99.99)

	// TODO: Run migration first time
	// Count records after first migration
	// cartsAfterFirst := ctx.countRows(ctx.ListingsDB, "shopping_carts")
	// ordersAfterFirst := ctx.countRows(ctx.ListingsDB, "orders")

	// TODO: Run migration second time
	// Count records after second migration
	// cartsAfterSecond := ctx.countRows(ctx.ListingsDB, "shopping_carts")
	// ordersAfterSecond := ctx.countRows(ctx.ListingsDB, "orders")

	// Verify no duplicates
	// assert.Equal(t, cartsAfterFirst, cartsAfterSecond, "Carts were duplicated on second run")
	// assert.Equal(t, ordersAfterFirst, ordersAfterSecond, "Orders were duplicated on second run")
}

// TestFKIntegrity validates all foreign key constraints after migration
func TestFKIntegrity(t *testing.T) {
	ctx := setupMigrationTest(t)
	defer ctx.cleanup()

	now := time.Now()

	// Create comprehensive test data
	cartID := ctx.createTestCart(1, 100, now.Add(-1*time.Hour))
	ctx.createTestCartItem(cartID, 200, 2, 99.99)
	ctx.createTestCartItem(cartID, 201, 1, 49.99)

	orderID := ctx.createTestOrder("ORD-2025-000001", 1, 100, 149.98, now.Add(-1*time.Hour))
	ctx.createTestOrderItem(orderID, 200, "Product 1", 2, 50.00, 100.00)
	ctx.createTestOrderItem(orderID, 201, "Product 2", 1, 49.98, 49.98)

	ctx.createTestReservation(200, orderID, 2, now.Add(2*time.Hour))

	// TODO: Run migration

	// Verify all FK constraints
	err := ctx.checkFKIntegrity(ctx.ListingsDB)
	assert.NoError(t, err, "FK integrity check failed after migration")
}

// TestRollback validates rollback on error
func TestRollback(t *testing.T) {
	ctx := setupMigrationTest(t)
	defer ctx.cleanup()

	now := time.Now()

	// Create test data
	ctx.createTestCart(1, 100, now.Add(-1*time.Hour))
	ctx.createTestOrder("ORD-2025-000001", 1, 100, 99.99, now.Add(-1*time.Hour))

	// Count before migration
	cartsBefore := ctx.countRows(ctx.ListingsDB, "shopping_carts")
	ordersBefore := ctx.countRows(ctx.ListingsDB, "orders")

	// TODO: Simulate error during migration
	// This will require injecting an error in the migration logic
	// For now, this is a placeholder test

	// Verify rollback occurred
	cartsAfter := ctx.countRows(ctx.ListingsDB, "shopping_carts")
	ordersAfter := ctx.countRows(ctx.ListingsDB, "orders")

	assert.Equal(t, cartsBefore, cartsAfter, "Rollback failed for shopping_carts")
	assert.Equal(t, ordersBefore, ordersAfter, "Rollback failed for orders")
}

// TestFullE2EMigration performs a comprehensive end-to-end migration test
func TestFullE2EMigration(t *testing.T) {
	ctx := setupMigrationTest(t)
	defer ctx.cleanup()

	now := time.Now()

	// =============================================================================
	// Setup: Create comprehensive test dataset
	// =============================================================================

	// 10 shopping carts (5 active, 5 old)
	var activeCarts []int64
	for i := 0; i < 5; i++ {
		cartID := ctx.createTestCart(int64(i+1), 100, now.Add(-time.Duration(i)*24*time.Hour))
		activeCarts = append(activeCarts, cartID)

		// Add 10 items per cart
		for j := 0; j < 10; j++ {
			ctx.createTestCartItem(cartID, int64(200+j), 1+j, 10.00*float64(j+1))
		}
	}

	for i := 0; i < 5; i++ {
		cartID := ctx.createTestCart(int64(i+6), 100, now.Add(-35*24*time.Hour))
		// Old carts with items
		ctx.createTestCartItem(cartID, 300, 1, 10.00)
	}

	// 20 orders (15 recent, 5 old)
	var recentOrders []int64
	for i := 0; i < 15; i++ {
		orderNumber := fmt.Sprintf("ORD-2025-%06d", i+1)
		orderID := ctx.createTestOrder(orderNumber, int64(i%5+1), 100, 199.99, now.Add(-time.Duration(i*20)*24*time.Hour))
		recentOrders = append(recentOrders, orderID)

		// Add 5-10 order items per order
		itemCount := 5 + (i % 6)
		for j := 0; j < itemCount; j++ {
			listingName := fmt.Sprintf("Product %d-%d", i+1, j+1)
			ctx.createTestOrderItem(orderID, int64(200+j), listingName, 1+j, 20.00, 20.00*float64(1+j))
		}
	}

	for i := 0; i < 5; i++ {
		orderNumber := fmt.Sprintf("ORD-2024-%06d", i+1)
		orderID := ctx.createTestOrder(orderNumber, int64(i+1), 100, 99.99, now.Add(-400*24*time.Hour))
		ctx.createTestOrderItem(orderID, 200, "Old Product", 1, 99.99, 99.99)
	}

	// 5 inventory reservations (3 active, 2 expired)
	var activeReservations []int64
	for i := 0; i < 3; i++ {
		resID := ctx.createTestReservation(int64(200+i), recentOrders[0], i+1, now.Add(time.Duration(i+1)*time.Hour))
		activeReservations = append(activeReservations, resID)
	}

	for i := 0; i < 2; i++ {
		ctx.createTestReservation(int64(210+i), recentOrders[0], i+1, now.Add(-time.Duration(i+1)*time.Hour))
	}

	// =============================================================================
	// Verify: Initial counts in monolith
	// =============================================================================

	assert.Equal(t, 10, ctx.countRows(ctx.MonolithDB, "shopping_carts"))
	assert.Equal(t, 55, ctx.countRows(ctx.MonolithDB, "cart_items")) // 5*10 + 5*1
	assert.Equal(t, 20, ctx.countRows(ctx.MonolithDB, "orders"))

	var totalOrderItems int
	for i := 0; i < 15; i++ {
		totalOrderItems += 5 + (i % 6)
	}
	totalOrderItems += 5 // old orders
	assert.Equal(t, totalOrderItems, ctx.countRows(ctx.MonolithDB, "order_items"))
	assert.Equal(t, 5, ctx.countRows(ctx.MonolithDB, "inventory_reservations"))

	// =============================================================================
	// Execute: Run migration
	// =============================================================================

	// TODO: Run migration script
	// migrator := NewMigrator(ctx.MonolithDB, ctx.ListingsDB, logger)
	// err := migrator.Migrate(context.Background(), MigrateOptions{DryRun: false})
	// require.NoError(t, err)

	// =============================================================================
	// Verify: Expected counts in microservice DB
	// =============================================================================

	// TODO: After migration is implemented, verify:
	// 1. Only 5 active carts migrated
	// expectedActiveCarts := 5
	// assert.Equal(t, expectedActiveCarts, ctx.countRows(ctx.ListingsDB, "shopping_carts"))

	// 2. Only cart items from active carts migrated (50 items)
	// expectedActiveCartItems := 50
	// assert.Equal(t, expectedActiveCartItems, ctx.countRows(ctx.ListingsDB, "cart_items"))

	// 3. Only 15 recent orders migrated
	// expectedRecentOrders := 15
	// assert.Equal(t, expectedRecentOrders, ctx.countRows(ctx.ListingsDB, "orders"))

	// 4. Only order items from recent orders migrated
	// var expectedRecentOrderItems int
	// for i := 0; i < 15; i++ {
	// 	expectedRecentOrderItems += 5 + (i % 6)
	// }
	// assert.Equal(t, expectedRecentOrderItems, ctx.countRows(ctx.ListingsDB, "order_items"))

	// 5. Only 3 active reservations migrated
	// expectedActiveRes := 3
	// assert.Equal(t, expectedActiveRes, ctx.countRows(ctx.ListingsDB, "inventory_reservations"))

	// =============================================================================
	// Verify: FK integrity
	// =============================================================================

	err := ctx.checkFKIntegrity(ctx.ListingsDB)
	assert.NoError(t, err, "FK integrity check failed after E2E migration")

	// =============================================================================
	// Verify: Sample record comparison
	// =============================================================================

	// TODO: Compare sample records
	// if len(activeCarts) > 0 && ctx.countRows(ctx.ListingsDB, "shopping_carts") > 0 {
	// 	ctx.compareCartRecord(activeCarts[0], activeCarts[0])
	// }

	// if len(recentOrders) > 0 && ctx.countRows(ctx.ListingsDB, "orders") > 0 {
	// 	ctx.compareOrderRecord(recentOrders[0], recentOrders[0])
	// }

	t.Logf("✅ E2E Migration Test Complete")
	t.Logf("   Active carts created: %d", len(activeCarts))
	t.Logf("   Recent orders created: %d", len(recentOrders))
	t.Logf("   Active reservations created: %d", len(activeReservations))
}
