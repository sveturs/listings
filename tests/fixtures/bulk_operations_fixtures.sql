-- Bulk Operations Integration Test Fixtures
-- Comprehensive test data for BulkCreateProducts, BulkUpdateProducts, BulkDeleteProducts
-- Phase 9.7.3 - Bulk Operations Integration Tests

-- ============================================================================
-- Minimal Dependencies (Foreign Key References)
-- ============================================================================

-- Create minimal users table (if not exists)
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Categories are auto-loaded from 00_categories_fixtures.sql (c2c_categories table)
-- We don't need to create the table here

-- Insert test users for bulk operations
INSERT INTO users (id, email, username, created_at, updated_at)
VALUES
    (6001, 'bulk-create-user@test.com', 'bulk_create_user', NOW(), NOW()),
    (6002, 'bulk-update-user@test.com', 'bulk_update_user', NOW(), NOW()),
    (6003, 'bulk-delete-user@test.com', 'bulk_delete_user', NOW(), NOW()),
    (6004, 'bulk-mixed-user@test.com', 'bulk_mixed_user', NOW(), NOW()),
    (6005, 'bulk-perf-user@test.com', 'bulk_perf_user', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Categories 1301-1305 are already loaded from 00_categories_fixtures.sql
-- No need to insert them again - use existing categories from that file

-- ============================================================================
-- Storefronts for Bulk Operations Tests
-- ============================================================================

INSERT INTO storefronts (id, user_id, name, slug, description, is_active, created_at, updated_at)
VALUES
    -- BulkCreateProducts tests (storefront 6001)
    (6001, 6001, 'Bulk Create Test Store', 'bulk-create-store', 'Store for bulk create operations', true, NOW(), NOW()),

    -- BulkUpdateProducts tests (storefront 6002)
    (6002, 6002, 'Bulk Update Test Store', 'bulk-update-store', 'Store for bulk update operations', true, NOW(), NOW()),

    -- BulkDeleteProducts tests (storefront 6003)
    (6003, 6003, 'Bulk Delete Test Store', 'bulk-delete-store', 'Store for bulk delete operations', true, NOW(), NOW()),

    -- Mixed operations tests (storefront 6004)
    (6004, 6004, 'Bulk Mixed Operations Store', 'bulk-mixed-store', 'Store for mixed bulk operations', true, NOW(), NOW()),

    -- Performance tests (storefront 6005)
    (6005, 6005, 'Bulk Performance Test Store', 'bulk-perf-store', 'Store for performance benchmarking', true, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- Products for BulkUpdateProducts Tests (Storefront 6002)
-- ============================================================================

INSERT INTO listings (id, user_id, storefront_id, category_id, title, description, price, currency, quantity, sku, status, source_type, created_at, updated_at)
VALUES
    -- Happy path: Products for successful bulk update
    (20001, 6002, 6002, 1301, 'Laptop Dell XPS 13', 'High-performance ultrabook', 1299.99, 'RSD', 10, 'LAPTOP-DELL-001', 'active', 'b2c', NOW(), NOW()),
    (20002, 6002, 6002, 1301, 'Laptop HP Spectre', 'Premium 2-in-1 laptop', 1499.99, 'RSD', 8, 'LAPTOP-HP-001', 'active', 'b2c', NOW(), NOW()),
    (20003, 6002, 6002, 1302, 'Desktop Computer', 'Gaming desktop PC', 1999.99, 'RSD', 5, 'DESKTOP-001', 'active', 'b2c', NOW(), NOW()),
    (20004, 6002, 6002, 1303, 'Wireless Mouse', 'Ergonomic wireless mouse', 49.99, 'RSD', 50, 'MOUSE-001', 'active', 'b2c', NOW(), NOW()),
    (20005, 6002, 6002, 1303, 'Mechanical Keyboard', 'RGB gaming keyboard', 129.99, 'RSD', 30, 'KEYBOARD-001', 'active', 'b2c', NOW(), NOW()),

    -- Partial update: Only some fields will be updated
    (20006, 6002, 6002, 1301, 'Monitor 27 inch', '4K UHD display', 599.99, 'RSD', 15, 'MONITOR-001', 'active', 'b2c', NOW(), NOW()),
    (20007, 6002, 6002, 1301, 'Webcam HD', '1080p webcam for streaming', 89.99, 'RSD', 25, 'WEBCAM-001', 'active', 'b2c', NOW(), NOW()),

    -- Edge cases: Inactive products
    (20008, 6002, 6002, 1304, 'Old T-Shirt', 'Discontinued product', 19.99, 'RSD', 0, 'TSHIRT-OLD-001', 'inactive', 'b2c', NOW(), NOW()),
    (20009, 6002, 6002, 1304, 'Old Jeans', 'Out of stock product', 39.99, 'RSD', 0, 'JEANS-OLD-001', 'inactive', 'b2c', NOW(), NOW()),

    -- Duplicate SKU test
    (20010, 6002, 6002, 1303, 'USB Cable Type-C', 'Fast charging cable', 19.99, 'RSD', 100, 'USB-CABLE-UNIQUE', 'active', 'b2c', NOW(), NOW()),
    (20011, 6002, 6002, 1303, 'USB Cable Lightning', 'Apple lightning cable', 24.99, 'RSD', 80, 'USB-CABLE-LIGHTNING', 'active', 'b2c', NOW(), NOW()),

    -- Validation test: Will receive negative price
    (20012, 6002, 6002, 1305, 'Garden Tool Set', 'Complete gardening tools', 79.99, 'RSD', 20, 'GARDEN-TOOLS-001', 'active', 'b2c', NOW(), NOW()),

    -- Concurrency test: Will be updated by multiple threads
    (20013, 6002, 6002, 1301, 'Tablet Android', 'Android tablet 10 inch', 299.99, 'RSD', 40, 'TABLET-ANDROID-001', 'active', 'b2c', NOW(), NOW()),
    (20014, 6002, 6002, 1301, 'Tablet iPad', 'Apple iPad Pro', 899.99, 'RSD', 12, 'TABLET-IPAD-001', 'active', 'b2c', NOW(), NOW()),

    -- Products with attributes (stored in listing_attributes table)
    (20015, 6002, 6002, 1302, 'Laptop Lenovo ThinkPad', 'Business laptop', 1199.99, 'RSD', 15, 'LAPTOP-LENOVO-001', 'active', 'b2c', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- Products for BulkDeleteProducts Tests (Storefront 6003)
-- ============================================================================

INSERT INTO listings (id, user_id, storefront_id, category_id, title, description, price, currency, quantity, sku, status, source_type, created_at, updated_at)
VALUES
    -- Products for soft delete
    (30001, 6003, 6003, 1301, 'Product Soft Delete 1', 'Will be soft deleted', 99.99, 'RSD', 10, 'SOFT-DEL-001', 'active', 'b2c', NOW(), NOW()),
    (30002, 6003, 6003, 1301, 'Product Soft Delete 2', 'Will be soft deleted', 199.99, 'RSD', 20, 'SOFT-DEL-002', 'active', 'b2c', NOW(), NOW()),
    (30003, 6003, 6003, 1302, 'Product Soft Delete 3', 'Will be soft deleted', 299.99, 'RSD', 30, 'SOFT-DEL-003', 'active', 'b2c', NOW(), NOW()),
    (30004, 6003, 6003, 1302, 'Product Soft Delete 4', 'Will be soft deleted', 399.99, 'RSD', 40, 'SOFT-DEL-004', 'active', 'b2c', NOW(), NOW()),
    (30005, 6003, 6003, 1303, 'Product Soft Delete 5', 'Will be soft deleted', 499.99, 'RSD', 50, 'SOFT-DEL-005', 'active', 'b2c', NOW(), NOW()),

    -- Products for hard delete
    (30011, 6003, 6003, 1301, 'Product Hard Delete 1', 'Will be hard deleted', 99.99, 'RSD', 10, 'HARD-DEL-001', 'active', 'b2c', NOW(), NOW()),
    (30012, 6003, 6003, 1301, 'Product Hard Delete 2', 'Will be hard deleted', 199.99, 'RSD', 20, 'HARD-DEL-002', 'active', 'b2c', NOW(), NOW()),
    (30013, 6003, 6003, 1302, 'Product Hard Delete 3', 'Will be hard deleted', 299.99, 'RSD', 30, 'HARD-DEL-003', 'active', 'b2c', NOW(), NOW()),

    -- Products for partial success test
    (30031, 6003, 6003, 1305, 'Product Partial 1', 'Will succeed', 149.99, 'RSD', 15, 'PARTIAL-001', 'active', 'b2c', NOW(), NOW()),
    (30032, 6003, 6003, 1305, 'Product Partial 2', 'Will succeed', 249.99, 'RSD', 25, 'PARTIAL-002', 'active', 'b2c', NOW(), NOW()),
    -- 30033 doesn't exist (will fail)
    (30034, 6003, 6003, 1305, 'Product Partial 4', 'Will succeed', 449.99, 'RSD', 45, 'PARTIAL-004', 'active', 'b2c', NOW(), NOW()),

    -- Already deleted products (for idempotency test)
    (30041, 6003, 6003, 1301, 'Already Deleted 1', 'Already soft deleted', 99.99, 'RSD', 5, 'ALREADY-DEL-001', 'inactive', 'b2c', NOW(), NOW()),
    (30042, 6003, 6003, 1301, 'Already Deleted 2', 'Already soft deleted', 199.99, 'RSD', 10, 'ALREADY-DEL-002', 'inactive', 'b2c', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Soft delete products 30041, 30042
UPDATE listings SET deleted_at = NOW() WHERE id IN (30041, 30042);

-- ============================================================================
-- Products for Mixed Operations Tests (Storefront 6004)
-- ============================================================================

INSERT INTO listings (id, user_id, storefront_id, category_id, title, description, price, currency, quantity, sku, status, source_type, created_at, updated_at)
VALUES
    (40001, 6004, 6004, 1301, 'Mixed Op Product 1', 'For mixed operations', 99.99, 'RSD', 10, 'MIXED-001', 'active', 'b2c', NOW(), NOW()),
    (40002, 6004, 6004, 1301, 'Mixed Op Product 2', 'For mixed operations', 199.99, 'RSD', 20, 'MIXED-002', 'active', 'b2c', NOW(), NOW()),
    (40003, 6004, 6004, 1302, 'Mixed Op Product 3', 'For mixed operations', 299.99, 'RSD', 30, 'MIXED-003', 'active', 'b2c', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- Performance Test Preparation (Storefront 6005)
-- ============================================================================
-- Pre-create some products for update/delete performance tests

INSERT INTO listings (id, user_id, storefront_id, category_id, title, description, price, currency, quantity, sku, status, source_type, created_at, updated_at)
SELECT
    50000 + generate_series AS id,
    6005 AS user_id,
    6005 AS storefront_id,
    1301 AS category_id,
    'Perf Test Product ' || generate_series AS title,
    'Performance test product description' AS description,
    (10 + (generate_series % 1000)) AS price,
    'RSD' AS currency,
    (generate_series % 100) AS quantity,
    'PERF-' || LPAD(generate_series::text, 6, '0') AS sku,
    'active' AS status,
    'b2c' AS source_type,
    NOW() AS created_at,
    NOW() AS updated_at
FROM generate_series(1, 200)
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- Test Summary
-- ============================================================================
-- Storefronts created:
--   6001: Bulk Create Test Store (for BulkCreateProducts)
--   6002: Bulk Update Test Store (15 products for BulkUpdateProducts)
--   6003: Bulk Delete Test Store (20+ products for BulkDeleteProducts)
--   6004: Mixed Operations Store (3 products for mixed operations)
--   6005: Performance Test Store (200 products for benchmarks)
--
-- Total products: 15 (storefront 6002) + 24 (storefront 6003) + 3 (storefront 6004) + 200 (storefront 6005) = 242 products
--
-- Categories: 5 test categories (1301-1305)
-- Users: 5 test users (6001-6005)
