-- ============================================================================
-- GetProduct & DeleteProduct Integration Test Fixtures
-- ============================================================================
-- Purpose: Test data for GetProduct and DeleteProduct gRPC API integration tests
-- Phase: 9.7.2 - Product CRUD Integration Tests
-- Created: 2025-11-05
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

-- Note: c2c_categories table is created by migrations
-- We just insert test data into it

-- Create b2c_product_images table (if not exists) - for GetProduct tests with images
CREATE TABLE IF NOT EXISTS b2c_product_images (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL,
    url VARCHAR(500) NOT NULL,
    storage_path VARCHAR(500),
    thumbnail_url VARCHAR(500),
    display_order INTEGER DEFAULT 0,
    is_primary BOOLEAN DEFAULT false,
    width INTEGER,
    height INTEGER,
    file_size BIGINT,
    mime_type VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Insert test users (only if not exist)
INSERT INTO users (id, email, username, created_at, updated_at)
VALUES
    (1, 'test-user-1@test.com', 'test_user_1', NOW(), NOW()),
    (2, 'test-user-2@test.com', 'test_user_2', NOW(), NOW()),
    (3, 'test-user-3@test.com', 'test_user_3', NOW(), NOW()),
    (9000, 'get-test@example.com', 'get_test_user', NOW(), NOW()),
    (9001, 'delete-test@example.com', 'delete_test_user', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Insert test categories (only ones actually used by products in this fixture)
-- Note: 00_categories_fixtures.sql is auto-loaded and provides common categories
-- We only insert specific test categories with high IDs to avoid conflicts
INSERT INTO c2c_categories (id, name, slug, parent_id, level, is_active, sort_order, created_at)
VALUES
    (9000, 'Test Electronics', 'test-electronics-9000', NULL, 0, true, 1, NOW()),
    (9001, 'Test Clothing', 'test-clothing-9001', NULL, 0, true, 2, NOW()),
    (9002, 'Test Books', 'test-books-9002', NULL, 0, true, 3, NOW())
ON CONFLICT (id) DO NOTHING;

-- Test storefronts (unique IDs to avoid conflicts)
INSERT INTO storefronts (id, user_id, name, slug, description, is_active, created_at, updated_at)
VALUES
    (9000, 9000, 'GetProduct Test Store', 'get-product-test-store', 'Storefront for GetProduct tests', true, NOW(), NOW()),
    (9001, 9001, 'DeleteProduct Test Store', 'delete-product-test-store', 'Storefront for DeleteProduct tests', true, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- GetProduct Test Data (Product IDs 9000-9099)
-- ============================================================================

-- 9000: Basic product (success case)
INSERT INTO listings (
    id, user_id, storefront_id, title, description,
    price, currency, category_id,
    sku, quantity, status,
    source_type,
    created_at, updated_at
)
VALUES
    (
        9000, 9000, 9000, 'Test Laptop', 'High-performance laptop for testing',
        1200.00, 'USD', 9000,
        'TEST-LAPTOP-001', 50, 'active',
        'b2c', NOW(), NOW()
    ),
    (
        9001, 9000, 9000, 'Test T-Shirt', 'Cotton t-shirt with multiple sizes',
        29.99, 'USD', 9001,
        'TEST-TSHIRT-001', 100, 'active',
        'b2c', NOW(), NOW()
    ),
    (
        9002, 9000, 9000, 'Test Headphones', 'Wireless bluetooth headphones',
        199.99, 'USD', 9000,
        'TEST-HEADPHONES-001', 30, 'active',
        'b2c', NOW(), NOW()
    ),
    (
        9003, 9000, 9000, 'Test Mouse', 'Gaming mouse',
        59.99, 'USD', 9000,
        'TEST-MOUSE-001', 0, 'inactive',
        'b2c', NOW(), NOW()
    ),
    (
        9004, 9000, 9000, 'Test Keyboard', 'Mechanical keyboard',
        149.99, 'USD', 9000,
        'TEST-KEYBOARD-001', 20, 'inactive',
        'b2c', NOW(), NOW()
    ),
    (9005, 9000, 9000, 'Test Product 5', 'Product for batch test', 10.00, 'USD', 9002, 'TEST-BATCH-005', 100, 'active', 'b2c', NOW(), NOW()),
    (9006, 9000, 9000, 'Test Product 6', 'Product for batch test', 15.00, 'USD', 9002, 'TEST-BATCH-006', 100, 'active', 'b2c', NOW(), NOW()),
    (9007, 9000, 9000, 'Test Product 7', 'Product for batch test', 20.00, 'USD', 9002, 'TEST-BATCH-007', 100, 'active', 'b2c', NOW(), NOW()),
    (9008, 9000, 9000, 'Test Product 8', 'Product for batch test', 25.00, 'USD', 9002, 'TEST-BATCH-008', 100, 'active', 'b2c', NOW(), NOW()),
    (9009, 9000, 9000, 'Test Product 9', 'Product for batch test', 30.00, 'USD', 9002, 'TEST-BATCH-009', 100, 'active', 'b2c', NOW(), NOW()),
    (9010, 9000, 9000, 'SKU Test 10', 'Product with unique SKU', 35.00, 'USD', 9000, 'SKU-BATCH-010', 50, 'active', 'b2c', NOW(), NOW()),
    (9011, 9000, 9000, 'SKU Test 11', 'Product with unique SKU', 40.00, 'USD', 9000, 'SKU-BATCH-011', 50, 'active', 'b2c', NOW(), NOW()),
    (9012, 9000, 9000, 'SKU Test 12', 'Product with unique SKU', 45.00, 'USD', 9000, 'SKU-BATCH-012', 50, 'active', 'b2c', NOW(), NOW()),
    (9013, 9000, 9000, 'SKU Test 13', 'Product with unique SKU', 50.00, 'USD', 9000, 'SKU-BATCH-013', 50, 'active', 'b2c', NOW(), NOW()),
    (9014, 9000, 9000, 'SKU Test 14', 'Product with unique SKU', 55.00, 'USD', 9000, 'SKU-BATCH-014', 50, 'active', 'b2c', NOW(), NOW()),
    (9015, 9000, 9000, 'SKU Test 15', 'Product with unique SKU', 60.00, 'USD', 9001, 'SKU-BATCH-015', 50, 'active', 'b2c', NOW(), NOW()),
    (9016, 9000, 9000, 'SKU Test 16', 'Product with unique SKU', 65.00, 'USD', 9001, 'SKU-BATCH-016', 50, 'active', 'b2c', NOW(), NOW()),
    (9017, 9000, 9000, 'SKU Test 17', 'Product with unique SKU', 70.00, 'USD', 9001, 'SKU-BATCH-017', 50, 'active', 'b2c', NOW(), NOW()),
    (9018, 9000, 9000, 'SKU Test 18', 'Product with unique SKU', 75.00, 'USD', 9001, 'SKU-BATCH-018', 50, 'active', 'b2c', NOW(), NOW()),
    (9019, 9000, 9000, 'SKU Test 19', 'Product with unique SKU', 80.00, 'USD', 9001, 'SKU-BATCH-019', 50, 'active', 'b2c', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Product images for product 9002 (Test Headphones)
INSERT INTO listing_images (
    id, listing_id, url, display_order, created_at, updated_at
)
VALUES
    (29000, 9002, 'https://test.example.com/headphones-1.jpg', 1, NOW(), NOW()),
    (29001, 9002, 'https://test.example.com/headphones-2.jpg', 2, NOW(), NOW()),
    (29002, 9002, 'https://test.example.com/headphones-3.jpg', 3, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- Soft Deleted Products (Product IDs 9020-9024) - For GetProduct NOT FOUND tests
-- ============================================================================

INSERT INTO listings (
    id, user_id, storefront_id, title, description,
    price, currency, category_id,
    sku, quantity, status,
    source_type,
    created_at, updated_at, deleted_at
)
VALUES
    (
        9020, 9000, 9000, 'Soft Deleted Product 1', 'This product is soft deleted',
        100.00, 'USD', 9000,
        'SOFT-DEL-001', 10, 'inactive',
        'b2c',
        NOW() - INTERVAL '10 days', NOW() - INTERVAL '5 days', NOW() - INTERVAL '5 days'
    ),
    (
        9021, 9000, 9000, 'Soft Deleted Product 2', 'This product is soft deleted',
        200.00, 'USD', 9001,
        'SOFT-DEL-002', 20, 'inactive',
        'b2c',
        NOW() - INTERVAL '20 days', NOW() - INTERVAL '10 days', NOW() - INTERVAL '10 days'
    ),
    (
        9022, 9000, 9000, 'Soft Deleted Product 3', 'This product is soft deleted',
        300.00, 'USD', 9002,
        'SOFT-DEL-003', 30, 'inactive',
        'b2c',
        NOW() - INTERVAL '30 days', NOW() - INTERVAL '15 days', NOW() - INTERVAL '15 days'
    ),
    (
        9023, 9000, 9000, 'Soft Deleted Product 4', 'This product is soft deleted',
        400.00, 'USD', 9000,
        'SOFT-DEL-004', 40, 'inactive',
        'b2c',
        NOW() - INTERVAL '40 days', NOW() - INTERVAL '20 days', NOW() - INTERVAL '20 days'
    ),
    (
        9024, 9000, 9000, 'Soft Deleted Product 5', 'This product is soft deleted',
        500.00, 'USD', 9001,
        'SOFT-DEL-005', 50, 'inactive',
        'b2c',
        NOW() - INTERVAL '50 days', NOW() - INTERVAL '25 days', NOW() - INTERVAL '25 days'
    )
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- DeleteProduct Test Data (Product IDs 9100-9199)
-- ============================================================================

-- 9100-9102: Products for delete tests (active products)
INSERT INTO listings (
    id, user_id, storefront_id, title, description,
    price, currency, category_id,
    sku, quantity, status,
    source_type,
    created_at, updated_at
)
VALUES
    (9100, 9001, 9001, 'Delete Test Product 1', 'Product for hard delete test', 50.00, 'USD', 9000, 'DEL-TEST-001', 10, 'active', 'b2c', NOW(), NOW()),
    (9101, 9001, 9001, 'Delete Test Product 2', 'Product for soft delete test', 60.00, 'USD', 9000, 'DEL-TEST-002', 15, 'active', 'b2c', NOW(), NOW()),
    (9102, 9001, 9001, 'Delete Test Product 3', 'Product for cascade delete', 70.00, 'USD', 9001, 'DEL-TEST-003', 20, 'active', 'b2c', NOW(), NOW()),
    (9104, 9001, 9001, 'Bulk Delete 1', 'Product for bulk delete', 90.00, 'USD', 9000, 'BULK-DEL-001', 30, 'active', 'b2c', NOW(), NOW()),
    (9105, 9001, 9001, 'Bulk Delete 2', 'Product for bulk delete', 100.00, 'USD', 9001, 'BULK-DEL-002', 35, 'active', 'b2c', NOW(), NOW()),
    (9106, 9001, 9001, 'Bulk Delete 3', 'Product for bulk delete', 110.00, 'USD', 9002, 'BULK-DEL-003', 40, 'active', 'b2c', NOW(), NOW()),
    (9107, 9001, 9001, 'Bulk Delete 4', 'Product for bulk delete', 120.00, 'USD', 9000, 'BULK-DEL-004', 45, 'active', 'b2c', NOW(), NOW()),
    (9108, 9001, 9001, 'Bulk Delete 5', 'Product for bulk delete', 130.00, 'USD', 9001, 'BULK-DEL-005', 50, 'active', 'b2c', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- 9103: Already soft-deleted product (separate INSERT with deleted_at)
INSERT INTO listings (
    id, user_id, storefront_id, title, description,
    price, currency, category_id,
    sku, quantity, status,
    source_type,
    created_at, updated_at, deleted_at
)
VALUES
    (9103, 9001, 9001, 'Delete Test Product 4', 'Already soft-deleted product', 80.00, 'USD', 9000, 'DEL-TEST-004', 25, 'inactive', 'b2c', NOW() - INTERVAL '7 days', NOW() - INTERVAL '3 days', NOW() - INTERVAL '3 days')
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- Large Batch Test Data (Product IDs 9150-9249) - 100 products
-- ============================================================================

INSERT INTO listings (
    id, user_id, storefront_id, title, description,
    price, currency, category_id,
    sku, quantity, status,
    source_type,
    created_at, updated_at
)
SELECT
    9150 + i,
    9001,
    9001,
    'Bulk Product ' || i,
    'Product for bulk operations testing',
    (50.00 + (i * 1.5))::numeric(10,2),
    'USD',
    9000 + (i % 3),
    'BULK-PROD-' || LPAD(i::text, 3, '0'),
    50 + i,
    CASE
        WHEN i % 10 = 0 THEN 'inactive'
        ELSE 'active'
    END,
    'b2c',
    NOW(),
    NOW()
FROM generate_series(1, 100) AS i
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- Performance & E2E Test Data (Product IDs 9300-9310)
-- ============================================================================

INSERT INTO listings (
    id, user_id, storefront_id, title, description,
    price, currency, category_id,
    sku, quantity, status,
    source_type,
    created_at, updated_at
)
VALUES
    (9300, 9000, 9000, 'E2E Test Product', 'Product for end-to-end workflow', 999.99, 'USD', 9000, 'E2E-PROD-001', 100, 'active', 'b2c', NOW(), NOW()),
    (9301, 9000, 9000, 'Performance Test 1', 'Product for performance testing', 10.00, 'USD', 9000, 'PERF-001', 1000, 'active', 'b2c', NOW(), NOW()),
    (9302, 9000, 9000, 'Performance Test 2', 'Product for performance testing', 20.00, 'USD', 9001, 'PERF-002', 2000, 'active', 'b2c', NOW(), NOW()),
    (9303, 9000, 9000, 'Performance Test 3', 'Product for performance testing', 30.00, 'USD', 9002, 'PERF-003', 3000, 'active', 'b2c', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- Product Variants for GetProduct/DeleteProduct Tests
-- ============================================================================

-- Variants for Product 9001 (Test T-Shirt with sizes)
INSERT INTO b2c_product_variants (
    id, product_id, sku, barcode,
    price, compare_at_price, cost_price,
    stock_quantity, stock_status, low_stock_threshold,
    variant_attributes, weight, dimensions,
    is_active, is_default, view_count, sold_count,
    created_at, updated_at
)
VALUES
    -- Size S
    (
        9101, 9001, 'TEST-TSHIRT-S', 'BAR-TSHIRT-S',
        29.99, 35.00, 15.00,
        30, 'in_stock', 5,
        '{"size": "S", "color": "white"}'::jsonb, 0.2, '{"length": 30, "width": 20, "height": 2}'::jsonb,
        true, false, 0, 0,
        NOW(), NOW()
    ),
    -- Size M
    (
        9102, 9001, 'TEST-TSHIRT-M', 'BAR-TSHIRT-M',
        29.99, 35.00, 15.00,
        40, 'in_stock', 5,
        '{"size": "M", "color": "white"}'::jsonb, 0.2, '{"length": 30, "width": 20, "height": 2}'::jsonb,
        true, true, 0, 0,
        NOW(), NOW()
    ),
    -- Size L
    (
        9103, 9001, 'TEST-TSHIRT-L', 'BAR-TSHIRT-L',
        29.99, 35.00, 15.00,
        30, 'in_stock', 5,
        '{"size": "L", "color": "white"}'::jsonb, 0.2, '{"length": 30, "width": 20, "height": 2}'::jsonb,
        true, false, 0, 0,
        NOW(), NOW()
    )
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- Verification Queries (for manual testing - commented out)
-- ============================================================================

-- SELECT COUNT(*) FROM listings WHERE id BETWEEN 9000 AND 9099; -- Should be ~20
-- SELECT COUNT(*) FROM listings WHERE id BETWEEN 9020 AND 9024 AND deleted_at IS NOT NULL; -- Should be 5
-- SELECT COUNT(*) FROM listings WHERE id BETWEEN 9100 AND 9199; -- Should be ~15
-- SELECT COUNT(*) FROM listings WHERE id BETWEEN 9150 AND 9249; -- Should be 100
-- SELECT COUNT(*) FROM listing_images WHERE listing_id = 9002; -- Should be 3
