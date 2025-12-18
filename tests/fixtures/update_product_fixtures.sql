-- ============================================================================
-- Update Product Integration Test Fixtures
-- ============================================================================
-- Purpose: Test data for UpdateProduct and BulkUpdateProducts gRPC tests
-- Products: 10001-10010 (for UpdateProduct tests)
-- Products: 10011-10020 (for BulkUpdateProducts tests)
-- Products: 10021-10030 (for concurrency tests)
-- Storefronts: 5001-5003
-- ============================================================================

-- ============================================================================
-- Minimal Dependencies (Foreign Key References)
-- ============================================================================
-- These tables don't exist in microservice migrations but are referenced by FKs

-- Create minimal users table (if not exists)
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Note: categories table is created by migrations
-- We just insert test data into it

-- Insert test users (only if not exist)
INSERT INTO users (id, email, username, created_at, updated_at)
VALUES
    (1, 'test-user-1@test.com', 'test_user_1', NOW(), NOW()),
    (2, 'test-user-2@test.com', 'test_user_2', NOW(), NOW()),
    (3, 'test-user-3@test.com', 'test_user_3', NOW(), NOW()),
    (1001, 'update-test-1@example.com', 'update_test_user_1', NOW(), NOW()),
    (1002, 'update-test-2@example.com', 'update_test_user_2', NOW(), NOW()),
    (1003, 'update-test-3@example.com', 'update_test_user_3', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Note: Categories are auto-loaded from 00_categories_fixtures.sql
-- This fixture uses category_id 1301 (Laptops) which is provided by the auto-loaded fixture
-- No need to insert categories here

-- Cleanup existing test data
DELETE FROM listings WHERE id >= 10001 AND id <= 10030;
DELETE FROM storefronts WHERE id >= 5001 AND id <= 5003;

-- ============================================================================
-- Test Storefronts
-- ============================================================================

INSERT INTO storefronts (id, user_id, name, slug, description, is_active, created_at, updated_at)
VALUES
    (5001, 1001, 'Update Test Store 1', 'update-test-store-1', 'Store for update product tests', true, NOW(), NOW()),
    (5002, 1002, 'Update Test Store 2', 'update-test-store-2', 'Store for bulk update tests', true, NOW(), NOW()),
    (5003, 1003, 'Update Test Store 3', 'update-test-store-3', 'Store for concurrency tests', true, NOW(), NOW());

-- ============================================================================
-- UpdateProduct Test Products (10001-10010)
-- ============================================================================

-- 10001: Standard product for successful full update
INSERT INTO listings (
    id, user_id, storefront_id, title, description, price, currency, category_id,
    sku, quantity, status,
    source_type, attributes, created_at, updated_at
) VALUES (
    10001, 1001, 5001, 'Original Product Name', 'Original description', 99.99, 'USD', 1301,
    'SKU-10001', 100, 'active',
    'b2c', '{"color": "red", "size": "M"}', NOW(), NOW()
);

-- 10002: Product for partial update (name only)
INSERT INTO listings (
    id, user_id, storefront_id, title, description, price, currency, category_id,
    sku, quantity, status,
    source_type, attributes, created_at, updated_at
) VALUES (
    10002, 1001, 5001, 'Partial Update Product', 'Will only update name', 49.99, 'USD', 1301,
    'SKU-10002', 50, 'active',
    'b2c', '{"material": "cotton"}', NOW(), NOW()
);

-- 10003: Product for price update test
INSERT INTO listings (
    id, user_id, storefront_id, title, description, price, currency, category_id,
    sku, quantity, status,
    source_type, attributes, created_at, updated_at
) VALUES (
    10003, 1001, 5001, 'Price Update Product', 'Test price changes', 29.99, 'USD', 1301,
    'SKU-10003', 75, 'active',
    'b2c', '{}', NOW(), NOW()
);

-- 10004: Product for quantity/stock update test
INSERT INTO listings (
    id, user_id, storefront_id, title, description, price, currency, category_id,
    sku, quantity, status,
    source_type, attributes, created_at, updated_at
) VALUES (
    10004, 1001, 5001, 'Stock Update Product', 'Test stock quantity changes', 19.99, 'USD', 1301,
    'SKU-10004', 200, 'active',
    'b2c', '{}', NOW(), NOW()
);

-- 10005: Product with existing SKU (for duplicate SKU test)
INSERT INTO listings (
    id, user_id, storefront_id, title, description, price, currency, category_id,
    sku, quantity, status,
    source_type, attributes, created_at, updated_at
) VALUES (
    10005, 1001, 5001, 'Duplicate SKU Target', 'Product with SKU-DUPLICATE', 59.99, 'USD', 1301,
    'SKU-DUPLICATE', 30, 'active',
    'b2c', '{}', NOW(), NOW()
);

-- 10006: Product that will attempt to use duplicate SKU
INSERT INTO listings (
    id, user_id, storefront_id, title, description, price, currency, category_id,
    sku, quantity, status,
    source_type, attributes, created_at, updated_at
) VALUES (
    10006, 1001, 5001, 'Will Try Duplicate SKU', 'Product for duplicate SKU test', 39.99, 'USD', 1301,
    'SKU-10006', 40, 'active',
    'b2c', '{}', NOW(), NOW()
);

-- 10007: Product for negative price validation test
INSERT INTO listings (
    id, user_id, storefront_id, title, description, price, currency, category_id,
    sku, quantity, status,
    source_type, attributes, created_at, updated_at
) VALUES (
    10007, 1001, 5001, 'Negative Price Test', 'Will try negative price', 25.00, 'USD', 1301,
    'SKU-10007', 60, 'active',
    'b2c', '{}', NOW(), NOW()
);

-- 10008: Product for attributes update test
INSERT INTO listings (
    id, user_id, storefront_id, title, description, price, currency, category_id,
    sku, quantity, status,
    source_type, attributes, created_at, updated_at
) VALUES (
    10008, 1001, 5001, 'Attributes Update Product', 'Test JSONB updates', 79.99, 'USD', 1301,
    'SKU-10008', 80, 'active',
    'b2c', '{"brand": "TestBrand", "warranty": "1 year"}', NOW(), NOW()
);

-- 10009: Product for updated_at timestamp verification
INSERT INTO listings (
    id, user_id, storefront_id, title, description, price, currency, category_id,
    sku, quantity, status,
    source_type, attributes, created_at, updated_at
) VALUES (
    10009, 1001, 5001, 'Timestamp Test Product', 'Verify updated_at changes', 15.99, 'USD', 1301,
    'SKU-10009', 90, 'active',
    'b2c', '{}', NOW() - INTERVAL '1 hour', NOW() - INTERVAL '1 hour'
);

-- 10010: Inactive product update test
INSERT INTO listings (
    id, user_id, storefront_id, title, description, price, currency, category_id,
    sku, quantity, status,
    source_type, attributes, created_at, updated_at
) VALUES (
    10010, 1001, 5001, 'Inactive Product', 'Test updating inactive product', 45.00, 'USD', 1301,
    'SKU-10010', 20, 'inactive',
    'b2c', '{}', NOW(), NOW()
);

-- ============================================================================
-- BulkUpdateProducts Test Products (10011-10020)
-- ============================================================================

-- 10011-10013: Products for successful bulk update
INSERT INTO listings (
    id, user_id, storefront_id, title, description, price, currency, category_id,
    sku, quantity, status,
    source_type, attributes, created_at, updated_at
) VALUES
    (10011, 1002, 5002, 'Bulk Product 1', 'First bulk update product', 10.00, 'USD', 1301, 'SKU-BULK-01', 100, 'active', 'b2c', '{}', NOW(), NOW()),
    (10012, 1002, 5002, 'Bulk Product 2', 'Second bulk update product', 20.00, 'USD', 1301, 'SKU-BULK-02', 200, 'active', 'b2c', '{}', NOW(), NOW()),
    (10013, 1002, 5002, 'Bulk Product 3', 'Third bulk update product', 30.00, 'USD', 1301, 'SKU-BULK-03', 300, 'active', 'b2c', '{}', NOW(), NOW());

-- 10014-10015: Products for partial success bulk update
INSERT INTO listings (
    id, user_id, storefront_id, title, description, price, currency, category_id,
    sku, quantity, status,
    source_type, attributes, created_at, updated_at
) VALUES
    (10014, 1002, 5002, 'Partial Success 1', 'Will succeed', 40.00, 'USD', 1301, 'SKU-PARTIAL-01', 400, 'active', 'b2c', '{}', NOW(), NOW()),
    (10015, 1002, 5002, 'Partial Success 2', 'Will succeed', 50.00, 'USD', 1301, 'SKU-PARTIAL-02', 500, 'active', 'b2c', '{}', NOW(), NOW());

-- 10016-10018: Products for mixed operations bulk update
INSERT INTO listings (
    id, user_id, storefront_id, title, description, price, currency, category_id,
    sku, quantity, status,
    source_type, attributes, created_at, updated_at
) VALUES
    (10016, 1002, 5002, 'Mixed Op 1', 'Will update price', 60.00, 'USD', 1301, 'SKU-MIXED-01', 60, 'active', 'b2c', '{}', NOW(), NOW()),
    (10017, 1002, 5002, 'Mixed Op 2', 'Will update stock', 70.00, 'USD', 1301, 'SKU-MIXED-02', 70, 'active', 'b2c', '{}', NOW(), NOW()),
    (10018, 1002, 5002, 'Mixed Op 3', 'Will update attributes', 80.00, 'USD', 1301, 'SKU-MIXED-03', 80, 'active', 'b2c', '{"old": "value"}', NOW(), NOW());

-- 10019-10020: Products for transaction rollback test (duplicate SKU scenario)
INSERT INTO listings (
    id, user_id, storefront_id, title, description, price, currency, category_id,
    sku, quantity, status,
    source_type, attributes, created_at, updated_at
) VALUES
    (10019, 1002, 5002, 'Rollback Test 1', 'Has unique SKU', 90.00, 'USD', 1301, 'SKU-ROLLBACK-01', 90, 'active', 'b2c', '{}', NOW(), NOW()),
    (10020, 1002, 5002, 'Rollback Test 2', 'Target for duplicate SKU', 100.00, 'USD', 1301, 'SKU-ROLLBACK-TARGET', 100, 'active', 'b2c', '{}', NOW(), NOW());

-- ============================================================================
-- Concurrency Test Products (10021-10030)
-- ============================================================================

-- 10021: Product for concurrent update test (10 goroutines will update simultaneously)
INSERT INTO listings (
    id, user_id, storefront_id, title, description, price, currency, category_id,
    sku, quantity, status,
    source_type, attributes, created_at, updated_at
) VALUES (
    10021, 1003, 5003, 'Concurrent Update Product', 'Will be updated by 10 goroutines', 50.00, 'USD', 1301,
    'SKU-CONCURRENT-01', 100, 'active',
    'b2c', '{}', NOW(), NOW()
);

-- 10022: Product for optimistic locking test (last write wins)
INSERT INTO listings (
    id, user_id, storefront_id, title, description, price, currency, category_id,
    sku, quantity, status,
    source_type, attributes, created_at, updated_at
) VALUES (
    10022, 1003, 5003, 'Last Write Wins Product', 'Test optimistic locking behavior', 60.00, 'USD', 1301,
    'SKU-CONCURRENT-02', 150, 'active',
    'b2c', '{}', NOW(), NOW()
);

-- 10023-10030: Products for performance/load test (single item < 100ms)
INSERT INTO listings (
    id, user_id, storefront_id, title, description, price, currency, category_id,
    sku, quantity, status,
    source_type, attributes, created_at, updated_at
) VALUES
    (10023, 1003, 5003, 'Perf Test 1', 'Performance test product', 10.00, 'USD', 1301, 'SKU-PERF-01', 100, 'active', 'b2c', '{}', NOW(), NOW()),
    (10024, 1003, 5003, 'Perf Test 2', 'Performance test product', 20.00, 'USD', 1301, 'SKU-PERF-02', 100, 'active', 'b2c', '{}', NOW(), NOW()),
    (10025, 1003, 5003, 'Perf Test 3', 'Performance test product', 30.00, 'USD', 1301, 'SKU-PERF-03', 100, 'active', 'b2c', '{}', NOW(), NOW()),
    (10026, 1003, 5003, 'Perf Test 4', 'Performance test product', 40.00, 'USD', 1301, 'SKU-PERF-04', 100, 'active', 'b2c', '{}', NOW(), NOW()),
    (10027, 1003, 5003, 'Perf Test 5', 'Performance test product', 50.00, 'USD', 1301, 'SKU-PERF-05', 100, 'active', 'b2c', '{}', NOW(), NOW()),
    (10028, 1003, 5003, 'Perf Test 6', 'Performance test product', 60.00, 'USD', 1301, 'SKU-PERF-06', 100, 'active', 'b2c', '{}', NOW(), NOW()),
    (10029, 1003, 5003, 'Perf Test 7', 'Performance test product', 70.00, 'USD', 1301, 'SKU-PERF-07', 100, 'active', 'b2c', '{}', NOW(), NOW()),
    (10030, 1003, 5003, 'Perf Test 8', 'Performance test product', 80.00, 'USD', 1301, 'SKU-PERF-08', 100, 'active', 'b2c', '{}', NOW(), NOW());
