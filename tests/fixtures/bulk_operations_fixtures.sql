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

-- Insert test users for bulk operations
INSERT INTO users (id, email, username, created_at, updated_at)
VALUES
    (6001, 'bulk-create-user@test.com', 'bulk_create_user', NOW(), NOW()),
    (6002, 'bulk-update-user@test.com', 'bulk_update_user', NOW(), NOW()),
    (6003, 'bulk-delete-user@test.com', 'bulk_delete_user', NOW(), NOW()),
    (6004, 'bulk-mixed-user@test.com', 'bulk_mixed_user', NOW(), NOW()),
    (6005, 'bulk-perf-user@test.com', 'bulk_perf_user', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Insert test categories
INSERT INTO categories (id, name, slug, parent_id, level, is_active, sort_order, created_at, updated_at)
VALUES
    (1301, 'Bulk Test Electronics', 'bulk-electronics', NULL, 0, true, 1, NOW(), NOW()),
    (1302, 'Bulk Test Computers', 'bulk-computers', 1301, 1, true, 1, NOW(), NOW()),
    (1303, 'Bulk Test Accessories', 'bulk-accessories', 1301, 1, true, 2, NOW(), NOW()),
    (1304, 'Bulk Test Clothing', 'bulk-clothing', NULL, 0, true, 2, NOW(), NOW()),
    (1305, 'Bulk Test Home & Garden', 'bulk-home-garden', NULL, 0, true, 3, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- Storefronts for Bulk Operations Tests
-- ============================================================================

INSERT INTO b2c_storefronts (id, user_id, name, slug, description, is_active, created_at, updated_at)
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

INSERT INTO b2c_marketplace_listings (id, storefront_id, category_id, name, description, price, quantity, sku, is_active, show_on_map, has_variants, created_at, updated_at)
VALUES
    -- Happy path: Products for successful bulk update
    (20001, 6002, 1301, 'Laptop Dell XPS 13', 'High-performance ultrabook', 1299.99, 10, 'LAPTOP-DELL-001', true, false, false, NOW(), NOW()),
    (20002, 6002, 1301, 'Laptop HP Spectre', 'Premium 2-in-1 laptop', 1499.99, 8, 'LAPTOP-HP-001', true, false, false, NOW(), NOW()),
    (20003, 6002, 1302, 'Desktop Computer', 'Gaming desktop PC', 1999.99, 5, 'DESKTOP-001', true, false, false, NOW(), NOW()),
    (20004, 6002, 1303, 'Wireless Mouse', 'Ergonomic wireless mouse', 49.99, 50, 'MOUSE-001', true, false, false, NOW(), NOW()),
    (20005, 6002, 1303, 'Mechanical Keyboard', 'RGB gaming keyboard', 129.99, 30, 'KEYBOARD-001', true, false, false, NOW(), NOW()),

    -- Partial update: Only some fields will be updated
    (20006, 6002, 1301, 'Monitor 27 inch', '4K UHD display', 599.99, 15, 'MONITOR-001', true, false, false, NOW(), NOW()),
    (20007, 6002, 1301, 'Webcam HD', '1080p webcam for streaming', 89.99, 25, 'WEBCAM-001', true, false, false, NOW(), NOW()),

    -- Edge cases: Inactive products
    (20008, 6002, 1304, 'Old T-Shirt', 'Discontinued product', 19.99, 0, 'TSHIRT-OLD-001', false, false, false, NOW(), NOW()),
    (20009, 6002, 1304, 'Old Jeans', 'Out of stock product', 39.99, 0, 'JEANS-OLD-001', false, false, false, NOW(), NOW()),

    -- Duplicate SKU test
    (20010, 6002, 1303, 'USB Cable Type-C', 'Fast charging cable', 19.99, 100, 'USB-CABLE-UNIQUE', true, false, false, NOW(), NOW()),
    (20011, 6002, 1303, 'USB Cable Lightning', 'Apple lightning cable', 24.99, 80, 'USB-CABLE-LIGHTNING', true, false, false, NOW(), NOW()),

    -- Validation test: Will receive negative price
    (20012, 6002, 1305, 'Garden Tool Set', 'Complete gardening tools', 79.99, 20, 'GARDEN-TOOLS-001', true, false, false, NOW(), NOW()),

    -- Concurrency test: Will be updated by multiple threads
    (20013, 6002, 1301, 'Tablet Android', 'Android tablet 10 inch', 299.99, 40, 'TABLET-ANDROID-001', true, false, false, NOW(), NOW()),
    (20014, 6002, 1301, 'Tablet iPad', 'Apple iPad Pro', 899.99, 12, 'TABLET-IPAD-001', true, false, false, NOW(), NOW()),

    -- Products with attributes
    (20015, 6002, 1302, 'Laptop Lenovo ThinkPad', 'Business laptop', 1199.99, 15, 'LAPTOP-LENOVO-001', true, false, false, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Update product 20015 with attributes
UPDATE b2c_marketplace_listings
SET attributes = '{"brand": "Lenovo", "processor": "Intel i7", "ram": "16GB", "storage": "512GB SSD"}'::jsonb
WHERE id = 20015;

-- ============================================================================
-- Products for BulkDeleteProducts Tests (Storefront 6003)
-- ============================================================================

INSERT INTO b2c_marketplace_listings (id, storefront_id, category_id, name, description, price, quantity, sku, is_active, show_on_map, has_variants, created_at, updated_at)
VALUES
    -- Products for soft delete
    (30001, 6003, 1301, 'Product Soft Delete 1', 'Will be soft deleted', 99.99, 10, 'SOFT-DEL-001', true, false, false, NOW(), NOW()),
    (30002, 6003, 1301, 'Product Soft Delete 2', 'Will be soft deleted', 199.99, 20, 'SOFT-DEL-002', true, false, false, NOW(), NOW()),
    (30003, 6003, 1302, 'Product Soft Delete 3', 'Will be soft deleted', 299.99, 30, 'SOFT-DEL-003', true, false, false, NOW(), NOW()),
    (30004, 6003, 1302, 'Product Soft Delete 4', 'Will be soft deleted', 399.99, 40, 'SOFT-DEL-004', true, false, false, NOW(), NOW()),
    (30005, 6003, 1303, 'Product Soft Delete 5', 'Will be soft deleted', 499.99, 50, 'SOFT-DEL-005', true, false, false, NOW(), NOW()),

    -- Products for hard delete
    (30011, 6003, 1301, 'Product Hard Delete 1', 'Will be hard deleted', 99.99, 10, 'HARD-DEL-001', true, false, false, NOW(), NOW()),
    (30012, 6003, 1301, 'Product Hard Delete 2', 'Will be hard deleted', 199.99, 20, 'HARD-DEL-002', true, false, false, NOW(), NOW()),
    (30013, 6003, 1302, 'Product Hard Delete 3', 'Will be hard deleted', 299.99, 30, 'HARD-DEL-003', true, false, false, NOW(), NOW()),

    -- Products with variants (cascade delete)
    (30021, 6003, 1304, 'T-Shirt with Variants', 'Product with size/color variants', 29.99, 0, 'TSHIRT-VAR-001', true, false, true, NOW(), NOW()),
    (30022, 6003, 1304, 'Shoes with Variants', 'Product with size variants', 79.99, 0, 'SHOES-VAR-001', true, false, true, NOW(), NOW()),

    -- Products for partial success test
    (30031, 6003, 1305, 'Product Partial 1', 'Will succeed', 149.99, 15, 'PARTIAL-001', true, false, false, NOW(), NOW()),
    (30032, 6003, 1305, 'Product Partial 2', 'Will succeed', 249.99, 25, 'PARTIAL-002', true, false, false, NOW(), NOW()),
    -- 30033 doesn't exist (will fail)
    (30034, 6003, 1305, 'Product Partial 4', 'Will succeed', 449.99, 45, 'PARTIAL-004', true, false, false, NOW(), NOW()),

    -- Already deleted products (for idempotency test)
    (30041, 6003, 1301, 'Already Deleted 1', 'Already soft deleted', 99.99, 5, 'ALREADY-DEL-001', false, false, false, NOW(), NOW()),
    (30042, 6003, 1301, 'Already Deleted 2', 'Already soft deleted', 199.99, 10, 'ALREADY-DEL-002', false, false, false, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Soft delete products 30041, 30042
UPDATE b2c_marketplace_listings SET deleted_at = NOW() WHERE id IN (30041, 30042);

-- Create variants for products 30021, 30022
INSERT INTO b2c_product_variants (id, product_id, name, sku, price_adjustment, quantity, is_default, attributes, created_at, updated_at)
VALUES
    -- Variants for T-Shirt (30021)
    (40001, 30021, 'Size S / Red', 'TSHIRT-VAR-001-S-RED', 0, 10, true, '{"size": "S", "color": "Red"}'::jsonb, NOW(), NOW()),
    (40002, 30021, 'Size M / Red', 'TSHIRT-VAR-001-M-RED', 0, 15, false, '{"size": "M", "color": "Red"}'::jsonb, NOW(), NOW()),
    (40003, 30021, 'Size L / Blue', 'TSHIRT-VAR-001-L-BLUE', 2, 12, false, '{"size": "L", "color": "Blue"}'::jsonb, NOW(), NOW()),

    -- Variants for Shoes (30022)
    (40011, 30022, 'Size 42', 'SHOES-VAR-001-42', 0, 8, true, '{"size": "42"}'::jsonb, NOW(), NOW()),
    (40012, 30022, 'Size 43', 'SHOES-VAR-001-43', 0, 10, false, '{"size": "43"}'::jsonb, NOW(), NOW()),
    (40013, 30022, 'Size 44', 'SHOES-VAR-001-44', 5, 6, false, '{"size": "44"}'::jsonb, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- Products for Mixed Operations Tests (Storefront 6004)
-- ============================================================================

INSERT INTO b2c_marketplace_listings (id, storefront_id, category_id, name, description, price, quantity, sku, is_active, show_on_map, has_variants, created_at, updated_at)
VALUES
    (40001, 6004, 1301, 'Mixed Op Product 1', 'For mixed operations', 99.99, 10, 'MIXED-001', true, false, false, NOW(), NOW()),
    (40002, 6004, 1301, 'Mixed Op Product 2', 'For mixed operations', 199.99, 20, 'MIXED-002', true, false, false, NOW(), NOW()),
    (40003, 6004, 1302, 'Mixed Op Product 3', 'For mixed operations', 299.99, 30, 'MIXED-003', true, false, false, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- Performance Test Preparation (Storefront 6005)
-- ============================================================================
-- Pre-create some products for update/delete performance tests

INSERT INTO b2c_marketplace_listings (id, storefront_id, category_id, name, description, price, quantity, sku, is_active, show_on_map, has_variants, created_at, updated_at)
SELECT
    50000 + generate_series AS id,
    6005 AS storefront_id,
    1301 AS category_id,
    'Perf Test Product ' || generate_series AS name,
    'Performance test product description' AS description,
    (10 + (generate_series % 1000)) AS price,
    (generate_series % 100) AS quantity,
    'PERF-' || LPAD(generate_series::text, 6, '0') AS sku,
    true AS is_active,
    false AS show_on_map,
    false AS has_variants,
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
-- Product variants: 6 variants for cascade delete tests
