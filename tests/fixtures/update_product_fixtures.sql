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

-- Create minimal categories table (if not exists)
CREATE TABLE IF NOT EXISTS categories (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    parent_id BIGINT REFERENCES categories(id),
    level INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

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

-- Insert test categories (only if not exist)
INSERT INTO categories (id, name, slug, parent_id, level, is_active, sort_order, created_at, updated_at)
VALUES
    (1, 'Electronics', 'electronics', NULL, 0, true, 1, NOW(), NOW()),
    (2, 'Laptops', 'laptops', 1, 1, true, 1, NOW(), NOW()),
    (3, 'Phones', 'phones', 1, 1, true, 2, NOW(), NOW()),
    (4, 'Clothing', 'clothing', NULL, 0, true, 2, NOW(), NOW()),
    (5, 'Shoes', 'shoes', 4, 1, true, 1, NOW(), NOW()),
    (1301, 'Test Category 1301', 'test-category-1301', NULL, 0, true, 1, NOW(), NOW()),
    (1302, 'Test Category 1302', 'test-category-1302', NULL, 0, true, 2, NOW(), NOW()),
    (1303, 'Test Category 1303', 'test-category-1303', NULL, 0, true, 3, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Cleanup existing test data
DELETE FROM b2c_products WHERE id >= 10001 AND id <= 10030;
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
INSERT INTO b2c_products (
    id, storefront_id, name, description, price, currency, category_id,
    sku, barcode, stock_quantity, stock_status,
    is_active, attributes, created_at, updated_at
) VALUES (
    10001, 5001, 'Original Product Name', 'Original description', 99.99, 'USD', 1301,
    'SKU-10001', 'BAR-10001', 100, 'in_stock',
    true, '{"color": "red", "size": "M"}', NOW(), NOW()
);

-- 10002: Product for partial update (name only)
INSERT INTO b2c_products (
    id, storefront_id, name, description, price, currency, category_id,
    sku, barcode, stock_quantity, stock_status,
    is_active, attributes, created_at, updated_at
) VALUES (
    10002, 5001, 'Partial Update Product', 'Will only update name', 49.99, 'USD', 1301,
    'SKU-10002', 'BAR-10002', 50, 'in_stock',
    true, '{"material": "cotton"}', NOW(), NOW()
);

-- 10003: Product for price update test
INSERT INTO b2c_products (
    id, storefront_id, name, description, price, currency, category_id,
    sku, barcode, stock_quantity, stock_status,
    is_active, attributes, created_at, updated_at
) VALUES (
    10003, 5001, 'Price Update Product', 'Test price changes', 29.99, 'USD', 1301,
    'SKU-10003', 'BAR-10003', 75, 'in_stock',
    true, '{}', NOW(), NOW()
);

-- 10004: Product for quantity/stock update test
INSERT INTO b2c_products (
    id, storefront_id, name, description, price, currency, category_id,
    sku, barcode, stock_quantity, stock_status,
    is_active, attributes, created_at, updated_at
) VALUES (
    10004, 5001, 'Stock Update Product', 'Test stock quantity changes', 19.99, 'USD', 1301,
    'SKU-10004', 'BAR-10004', 200, 'in_stock',
    true, '{}', NOW(), NOW()
);

-- 10005: Product with existing SKU (for duplicate SKU test)
INSERT INTO b2c_products (
    id, storefront_id, name, description, price, currency, category_id,
    sku, barcode, stock_quantity, stock_status,
    is_active, attributes, created_at, updated_at
) VALUES (
    10005, 5001, 'Duplicate SKU Target', 'Product with SKU-DUPLICATE', 59.99, 'USD', 1301,
    'SKU-DUPLICATE', 'BAR-10005', 30, 'in_stock',
    true, '{}', NOW(), NOW()
);

-- 10006: Product that will attempt to use duplicate SKU
INSERT INTO b2c_products (
    id, storefront_id, name, description, price, currency, category_id,
    sku, barcode, stock_quantity, stock_status,
    is_active, attributes, created_at, updated_at
) VALUES (
    10006, 5001, 'Will Try Duplicate SKU', 'Product for duplicate SKU test', 39.99, 'USD', 1301,
    'SKU-10006', 'BAR-10006', 40, 'in_stock',
    true, '{}', NOW(), NOW()
);

-- 10007: Product for negative price validation test
INSERT INTO b2c_products (
    id, storefront_id, name, description, price, currency, category_id,
    sku, barcode, stock_quantity, stock_status,
    is_active, attributes, created_at, updated_at
) VALUES (
    10007, 5001, 'Negative Price Test', 'Will try negative price', 25.00, 'USD', 1301,
    'SKU-10007', 'BAR-10007', 60, 'in_stock',
    true, '{}', NOW(), NOW()
);

-- 10008: Product for attributes update test
INSERT INTO b2c_products (
    id, storefront_id, name, description, price, currency, category_id,
    sku, barcode, stock_quantity, stock_status,
    is_active, attributes, created_at, updated_at
) VALUES (
    10008, 5001, 'Attributes Update Product', 'Test JSONB updates', 79.99, 'USD', 1301,
    'SKU-10008', 'BAR-10008', 80, 'in_stock',
    true, '{"brand": "TestBrand", "warranty": "1 year"}', NOW(), NOW()
);

-- 10009: Product for updated_at timestamp verification
INSERT INTO b2c_products (
    id, storefront_id, name, description, price, currency, category_id,
    sku, barcode, stock_quantity, stock_status,
    is_active, attributes, created_at, updated_at
) VALUES (
    10009, 5001, 'Timestamp Test Product', 'Verify updated_at changes', 15.99, 'USD', 1301,
    'SKU-10009', 'BAR-10009', 90, 'in_stock',
    true, '{}', NOW() - INTERVAL '1 hour', NOW() - INTERVAL '1 hour'
);

-- 10010: Inactive product update test
INSERT INTO b2c_products (
    id, storefront_id, name, description, price, currency, category_id,
    sku, barcode, stock_quantity, stock_status,
    is_active, attributes, created_at, updated_at
) VALUES (
    10010, 5001, 'Inactive Product', 'Test updating inactive product', 45.00, 'USD', 1301,
    'SKU-10010', 'BAR-10010', 20, 'in_stock',
    false, '{}', NOW(), NOW()
);

-- ============================================================================
-- BulkUpdateProducts Test Products (10011-10020)
-- ============================================================================

-- 10011-10013: Products for successful bulk update
INSERT INTO b2c_products (
    id, storefront_id, name, description, price, currency, category_id,
    sku, barcode, stock_quantity, stock_status,
    is_active, attributes, created_at, updated_at
) VALUES
    (10011, 5002, 'Bulk Product 1', 'First bulk update product', 10.00, 'USD', 1301, 'SKU-BULK-01', 'BAR-BULK-01', 100, 'in_stock', true, '{}', NOW(), NOW()),
    (10012, 5002, 'Bulk Product 2', 'Second bulk update product', 20.00, 'USD', 1301, 'SKU-BULK-02', 'BAR-BULK-02', 200, 'in_stock', true, '{}', NOW(), NOW()),
    (10013, 5002, 'Bulk Product 3', 'Third bulk update product', 30.00, 'USD', 1301, 'SKU-BULK-03', 'BAR-BULK-03', 300, 'in_stock', true, '{}', NOW(), NOW());

-- 10014-10015: Products for partial success bulk update
INSERT INTO b2c_products (
    id, storefront_id, name, description, price, currency, category_id,
    sku, barcode, stock_quantity, stock_status,
    is_active, attributes, created_at, updated_at
) VALUES
    (10014, 5002, 'Partial Success 1', 'Will succeed', 40.00, 'USD', 1301, 'SKU-PARTIAL-01', 'BAR-PARTIAL-01', 400, 'in_stock', true, '{}', NOW(), NOW()),
    (10015, 5002, 'Partial Success 2', 'Will succeed', 50.00, 'USD', 1301, 'SKU-PARTIAL-02', 'BAR-PARTIAL-02', 500, 'in_stock', true, '{}', NOW(), NOW());

-- 10016-10018: Products for mixed operations bulk update
INSERT INTO b2c_products (
    id, storefront_id, name, description, price, currency, category_id,
    sku, barcode, stock_quantity, stock_status,
    is_active, attributes, created_at, updated_at
) VALUES
    (10016, 5002, 'Mixed Op 1', 'Will update price', 60.00, 'USD', 1301, 'SKU-MIXED-01', 'BAR-MIXED-01', 60, 'in_stock', true, '{}', NOW(), NOW()),
    (10017, 5002, 'Mixed Op 2', 'Will update stock', 70.00, 'USD', 1301, 'SKU-MIXED-02', 'BAR-MIXED-02', 70, 'in_stock', true, '{}', NOW(), NOW()),
    (10018, 5002, 'Mixed Op 3', 'Will update attributes', 80.00, 'USD', 1301, 'SKU-MIXED-03', 'BAR-MIXED-03', 80, 'in_stock', true, '{"old": "value"}', NOW(), NOW());

-- 10019-10020: Products for transaction rollback test (duplicate SKU scenario)
INSERT INTO b2c_products (
    id, storefront_id, name, description, price, currency, category_id,
    sku, barcode, stock_quantity, stock_status,
    is_active, attributes, created_at, updated_at
) VALUES
    (10019, 5002, 'Rollback Test 1', 'Has unique SKU', 90.00, 'USD', 1301, 'SKU-ROLLBACK-01', 'BAR-ROLLBACK-01', 90, 'in_stock', true, '{}', NOW(), NOW()),
    (10020, 5002, 'Rollback Test 2', 'Target for duplicate SKU', 100.00, 'USD', 1301, 'SKU-ROLLBACK-TARGET', 'BAR-ROLLBACK-02', 100, 'in_stock', true, '{}', NOW(), NOW());

-- ============================================================================
-- Concurrency Test Products (10021-10030)
-- ============================================================================

-- 10021: Product for concurrent update test (10 goroutines will update simultaneously)
INSERT INTO b2c_products (
    id, storefront_id, name, description, price, currency, category_id,
    sku, barcode, stock_quantity, stock_status,
    is_active, attributes, created_at, updated_at
) VALUES (
    10021, 5003, 'Concurrent Update Product', 'Will be updated by 10 goroutines', 50.00, 'USD', 1301,
    'SKU-CONCURRENT-01', 'BAR-CONCURRENT-01', 100, 'in_stock',
    true, '{}', NOW(), NOW()
);

-- 10022: Product for optimistic locking test (last write wins)
INSERT INTO b2c_products (
    id, storefront_id, name, description, price, currency, category_id,
    sku, barcode, stock_quantity, stock_status,
    is_active, attributes, created_at, updated_at
) VALUES (
    10022, 5003, 'Last Write Wins Product', 'Test optimistic locking behavior', 60.00, 'USD', 1301,
    'SKU-CONCURRENT-02', 'BAR-CONCURRENT-02', 150, 'in_stock',
    true, '{}', NOW(), NOW()
);

-- 10023-10030: Products for performance/load test (single item < 100ms)
INSERT INTO b2c_products (
    id, storefront_id, name, description, price, currency, category_id,
    sku, barcode, stock_quantity, stock_status,
    is_active, attributes, created_at, updated_at
) VALUES
    (10023, 5003, 'Perf Test 1', 'Performance test product', 10.00, 'USD', 1301, 'SKU-PERF-01', 'BAR-PERF-01', 100, 'in_stock', true, '{}', NOW(), NOW()),
    (10024, 5003, 'Perf Test 2', 'Performance test product', 20.00, 'USD', 1301, 'SKU-PERF-02', 'BAR-PERF-02', 100, 'in_stock', true, '{}', NOW(), NOW()),
    (10025, 5003, 'Perf Test 3', 'Performance test product', 30.00, 'USD', 1301, 'SKU-PERF-03', 'BAR-PERF-03', 100, 'in_stock', true, '{}', NOW(), NOW()),
    (10026, 5003, 'Perf Test 4', 'Performance test product', 40.00, 'USD', 1301, 'SKU-PERF-04', 'BAR-PERF-04', 100, 'in_stock', true, '{}', NOW(), NOW()),
    (10027, 5003, 'Perf Test 5', 'Performance test product', 50.00, 'USD', 1301, 'SKU-PERF-05', 'BAR-PERF-05', 100, 'in_stock', true, '{}', NOW(), NOW()),
    (10028, 5003, 'Perf Test 6', 'Performance test product', 60.00, 'USD', 1301, 'SKU-PERF-06', 'BAR-PERF-06', 100, 'in_stock', true, '{}', NOW(), NOW()),
    (10029, 5003, 'Perf Test 7', 'Performance test product', 70.00, 'USD', 1301, 'SKU-PERF-07', 'BAR-PERF-07', 100, 'in_stock', true, '{}', NOW(), NOW()),
    (10030, 5003, 'Perf Test 8', 'Performance test product', 80.00, 'USD', 1301, 'SKU-PERF-08', 'BAR-PERF-08', 100, 'in_stock', true, '{}', NOW(), NOW());
