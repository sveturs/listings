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
    (1303, 'Test Category 1303', 'test-category-1303', NULL, 0, true, 3, NOW(), NOW()),
    (9000, 'Test Electronics', 'test-electronics', NULL, 0, true, 1, NOW(), NOW()),
    (9001, 'Test Clothing', 'test-clothing', NULL, 0, true, 2, NOW(), NOW()),
    (9002, 'Test Books', 'test-books', NULL, 0, true, 3, NOW(), NOW())
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
INSERT INTO b2c_products (
    id, storefront_id, name, description, price, currency, category_id,
    sku, barcode, stock_quantity, stock_status, is_active,
    attributes, view_count, sold_count, has_variants,
    created_at, updated_at
)
VALUES
    (
        9000, 9000, 'Test Laptop', 'High-performance laptop for testing',
        1200.00, 'USD', 9000, 'TEST-LAPTOP-001', '1234567890000',
        50, 'in_stock', true,
        '{"brand": "TestBrand", "model": "X1", "year": 2024}'::jsonb,
        100, 25, false,
        NOW(), NOW()
    ),
    -- 9001: Product with variants
    (
        9001, 9000, 'Test T-Shirt', 'Cotton t-shirt with multiple sizes',
        29.99, 'USD', 9001, 'TEST-TSHIRT-001', '1234567890001',
        100, 'in_stock', true,
        '{"material": "cotton", "brand": "TestClothing"}'::jsonb,
        250, 80, true,
        NOW(), NOW()
    ),
    -- 9002: Product with images
    (
        9002, 9000, 'Test Headphones', 'Wireless bluetooth headphones',
        199.99, 'USD', 9000, 'TEST-HEADPHONES-001', '1234567890002',
        30, 'in_stock', true,
        '{"wireless": true, "bluetooth": "5.0"}'::jsonb,
        500, 120, false,
        NOW(), NOW()
    ),
    -- 9003: Out of stock product
    (
        9003, 9000, 'Test Mouse', 'Gaming mouse',
        59.99, 'USD', 9000, 'TEST-MOUSE-001', '1234567890003',
        0, 'out_of_stock', true,
        '{"dpi": 16000, "buttons": 8}'::jsonb,
        75, 0, false,
        NOW(), NOW()
    ),
    -- 9004: Inactive product
    (
        9004, 9000, 'Test Keyboard', 'Mechanical keyboard',
        149.99, 'USD', 9000, 'TEST-KEYBOARD-001', '1234567890004',
        20, 'in_stock', false,
        '{"switches": "cherry_mx_blue"}'::jsonb,
        45, 10, false,
        NOW(), NOW()
    ),
    -- 9005-9009: Products for batch GetProductsByIDs (5 products)
    (9005, 9000, 'Test Product 5', 'Product for batch test', 10.00, 'USD', 9002, 'TEST-BATCH-005', '1234567890005', 100, 'in_stock', true, '{}'::jsonb, 0, 0, false, NOW(), NOW()),
    (9006, 9000, 'Test Product 6', 'Product for batch test', 15.00, 'USD', 9002, 'TEST-BATCH-006', '1234567890006', 100, 'in_stock', true, '{}'::jsonb, 0, 0, false, NOW(), NOW()),
    (9007, 9000, 'Test Product 7', 'Product for batch test', 20.00, 'USD', 9002, 'TEST-BATCH-007', '1234567890007', 100, 'in_stock', true, '{}'::jsonb, 0, 0, false, NOW(), NOW()),
    (9008, 9000, 'Test Product 8', 'Product for batch test', 25.00, 'USD', 9002, 'TEST-BATCH-008', '1234567890008', 100, 'in_stock', true, '{}'::jsonb, 0, 0, false, NOW(), NOW()),
    (9009, 9000, 'Test Product 9', 'Product for batch test', 30.00, 'USD', 9002, 'TEST-BATCH-009', '1234567890009', 100, 'in_stock', true, '{}'::jsonb, 0, 0, false, NOW(), NOW()),
    -- 9010-9019: Products for GetProductsBySKUs (10 products)
    (9010, 9000, 'SKU Test 10', 'Product with unique SKU', 35.00, 'USD', 9000, 'SKU-BATCH-010', '1234567890010', 50, 'in_stock', true, '{}'::jsonb, 0, 0, false, NOW(), NOW()),
    (9011, 9000, 'SKU Test 11', 'Product with unique SKU', 40.00, 'USD', 9000, 'SKU-BATCH-011', '1234567890011', 50, 'in_stock', true, '{}'::jsonb, 0, 0, false, NOW(), NOW()),
    (9012, 9000, 'SKU Test 12', 'Product with unique SKU', 45.00, 'USD', 9000, 'SKU-BATCH-012', '1234567890012', 50, 'in_stock', true, '{}'::jsonb, 0, 0, false, NOW(), NOW()),
    (9013, 9000, 'SKU Test 13', 'Product with unique SKU', 50.00, 'USD', 9000, 'SKU-BATCH-013', '1234567890013', 50, 'in_stock', true, '{}'::jsonb, 0, 0, false, NOW(), NOW()),
    (9014, 9000, 'SKU Test 14', 'Product with unique SKU', 55.00, 'USD', 9000, 'SKU-BATCH-014', '1234567890014', 50, 'in_stock', true, '{}'::jsonb, 0, 0, false, NOW(), NOW()),
    (9015, 9000, 'SKU Test 15', 'Product with unique SKU', 60.00, 'USD', 9001, 'SKU-BATCH-015', '1234567890015', 50, 'in_stock', true, '{}'::jsonb, 0, 0, false, NOW(), NOW()),
    (9016, 9000, 'SKU Test 16', 'Product with unique SKU', 65.00, 'USD', 9001, 'SKU-BATCH-016', '1234567890016', 50, 'in_stock', true, '{}'::jsonb, 0, 0, false, NOW(), NOW()),
    (9017, 9000, 'SKU Test 17', 'Product with unique SKU', 70.00, 'USD', 9001, 'SKU-BATCH-017', '1234567890017', 50, 'in_stock', true, '{}'::jsonb, 0, 0, false, NOW(), NOW()),
    (9018, 9000, 'SKU Test 18', 'Product with unique SKU', 75.00, 'USD', 9001, 'SKU-BATCH-018', '1234567890018', 50, 'in_stock', true, '{}'::jsonb, 0, 0, false, NOW(), NOW()),
    (9019, 9000, 'SKU Test 19', 'Product with unique SKU', 80.00, 'USD', 9001, 'SKU-BATCH-019', '1234567890019', 50, 'in_stock', true, '{}'::jsonb, 0, 0, false, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Variants for product 9001 (Test T-Shirt)
INSERT INTO b2c_product_variants (
    id, product_id, sku, barcode, price, stock_quantity, stock_status,
    variant_attributes, is_active, is_default, view_count, sold_count,
    created_at, updated_at
)
VALUES
    (19000, 9001, 'TEST-TSHIRT-001-S', '1234567890001-S', 29.99, 30, 'in_stock', '{"size": "S", "color": "black"}'::jsonb, true, false, 50, 20, NOW(), NOW()),
    (19001, 9001, 'TEST-TSHIRT-001-M', '1234567890001-M', 29.99, 40, 'in_stock', '{"size": "M", "color": "black"}'::jsonb, true, true, 100, 35, NOW(), NOW()),
    (19002, 9001, 'TEST-TSHIRT-001-L', '1234567890001-L', 29.99, 30, 'in_stock', '{"size": "L", "color": "black"}'::jsonb, true, false, 80, 25, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Product images for product 9002 (Test Headphones)
INSERT INTO b2c_product_images (
    id, product_id, url, storage_path, thumbnail_url, display_order, is_primary,
    width, height, file_size, mime_type, created_at, updated_at
)
VALUES
    (29000, 9002, 'https://test.example.com/headphones-1.jpg', '/test/headphones-1.jpg', 'https://test.example.com/headphones-1-thumb.jpg', 1, true, 1920, 1080, 524288, 'image/jpeg', NOW(), NOW()),
    (29001, 9002, 'https://test.example.com/headphones-2.jpg', '/test/headphones-2.jpg', 'https://test.example.com/headphones-2-thumb.jpg', 2, false, 1920, 1080, 498432, 'image/jpeg', NOW(), NOW()),
    (29002, 9002, 'https://test.example.com/headphones-3.jpg', '/test/headphones-3.jpg', 'https://test.example.com/headphones-3-thumb.jpg', 3, false, 1920, 1080, 512000, 'image/jpeg', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- Soft Deleted Products (Product IDs 9020-9024) - For GetProduct NOT FOUND tests
-- ============================================================================

INSERT INTO b2c_products (
    id, storefront_id, name, description, price, currency, category_id,
    sku, barcode, stock_quantity, stock_status, is_active,
    attributes, view_count, sold_count, has_variants,
    created_at, updated_at, deleted_at
)
VALUES
    (
        9020, 9000, 'Soft Deleted Product 1', 'This product is soft deleted',
        100.00, 'USD', 9000, 'SOFT-DEL-001', '1234567890020',
        10, 'in_stock', false,
        '{}'::jsonb, 0, 0, false,
        NOW() - INTERVAL '10 days', NOW() - INTERVAL '5 days', NOW() - INTERVAL '5 days'
    ),
    (
        9021, 9000, 'Soft Deleted Product 2', 'This product is soft deleted',
        200.00, 'USD', 9001, 'SOFT-DEL-002', '1234567890021',
        20, 'in_stock', false,
        '{}'::jsonb, 0, 0, false,
        NOW() - INTERVAL '20 days', NOW() - INTERVAL '10 days', NOW() - INTERVAL '10 days'
    ),
    (
        9022, 9000, 'Soft Deleted Product 3', 'This product is soft deleted',
        300.00, 'USD', 9002, 'SOFT-DEL-003', '1234567890022',
        30, 'in_stock', false,
        '{}'::jsonb, 0, 0, false,
        NOW() - INTERVAL '30 days', NOW() - INTERVAL '15 days', NOW() - INTERVAL '15 days'
    ),
    (
        9023, 9000, 'Soft Deleted Product 4', 'This product is soft deleted',
        400.00, 'USD', 9000, 'SOFT-DEL-004', '1234567890023',
        40, 'in_stock', false,
        '{}'::jsonb, 0, 0, false,
        NOW() - INTERVAL '40 days', NOW() - INTERVAL '20 days', NOW() - INTERVAL '20 days'
    ),
    (
        9024, 9000, 'Soft Deleted Product 5', 'This product is soft deleted with variants',
        500.00, 'USD', 9001, 'SOFT-DEL-005', '1234567890024',
        50, 'in_stock', false,
        '{}'::jsonb, 0, 0, true,
        NOW() - INTERVAL '50 days', NOW() - INTERVAL '25 days', NOW() - INTERVAL '25 days'
    )
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- DeleteProduct Test Data (Product IDs 9100-9199)
-- ============================================================================

-- 9100-9102: Products for delete tests (active products)
INSERT INTO b2c_products (
    id, storefront_id, name, description, price, currency, category_id,
    sku, barcode, stock_quantity, stock_status, is_active,
    attributes, view_count, sold_count, has_variants,
    created_at, updated_at
)
VALUES
    (9100, 9001, 'Delete Test Product 1', 'Product for hard delete test', 50.00, 'USD', 9000, 'DEL-TEST-001', '1234567890100', 10, 'in_stock', true, '{}'::jsonb, 0, 0, false, NOW(), NOW()),
    -- 9101: Product for soft delete
    (9101, 9001, 'Delete Test Product 2', 'Product for soft delete test', 60.00, 'USD', 9000, 'DEL-TEST-002', '1234567890101', 15, 'in_stock', true, '{}'::jsonb, 0, 0, false, NOW(), NOW()),
    -- 9102: Product with variants (CASCADE test)
    (9102, 9001, 'Delete Test Product 3', 'Product with variants for cascade delete', 70.00, 'USD', 9001, 'DEL-TEST-003', '1234567890102', 20, 'in_stock', true, '{}'::jsonb, 0, 0, true, NOW(), NOW()),
    -- 9104-9108: Products for BulkDeleteProducts (5 products)
    (9104, 9001, 'Bulk Delete 1', 'Product for bulk delete', 90.00, 'USD', 9000, 'BULK-DEL-001', '1234567890104', 30, 'in_stock', true, '{}'::jsonb, 0, 0, false, NOW(), NOW()),
    (9105, 9001, 'Bulk Delete 2', 'Product for bulk delete', 100.00, 'USD', 9001, 'BULK-DEL-002', '1234567890105', 35, 'in_stock', true, '{}'::jsonb, 0, 0, false, NOW(), NOW()),
    (9106, 9001, 'Bulk Delete 3', 'Product for bulk delete', 110.00, 'USD', 9002, 'BULK-DEL-003', '1234567890106', 40, 'in_stock', true, '{}'::jsonb, 0, 0, false, NOW(), NOW()),
    (9107, 9001, 'Bulk Delete 4', 'Product for bulk delete', 120.00, 'USD', 9000, 'BULK-DEL-004', '1234567890107', 45, 'in_stock', true, '{}'::jsonb, 0, 0, false, NOW(), NOW()),
    (9108, 9001, 'Bulk Delete 5', 'Product for bulk delete', 130.00, 'USD', 9001, 'BULK-DEL-005', '1234567890108', 50, 'in_stock', true, '{}'::jsonb, 0, 0, false, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- 9103: Already soft-deleted product (separate INSERT with deleted_at)
INSERT INTO b2c_products (
    id, storefront_id, name, description, price, currency, category_id,
    sku, barcode, stock_quantity, stock_status, is_active,
    attributes, view_count, sold_count, has_variants,
    created_at, updated_at, deleted_at
)
VALUES
    (9103, 9001, 'Delete Test Product 4', 'Already soft-deleted product', 80.00, 'USD', 9000, 'DEL-TEST-004', '1234567890103', 25, 'in_stock', false, '{}'::jsonb, 0, 0, false, NOW() - INTERVAL '7 days', NOW() - INTERVAL '3 days', NOW() - INTERVAL '3 days')
ON CONFLICT (id) DO NOTHING;

-- Variants for product 9102 (cascade delete test)
INSERT INTO b2c_product_variants (
    id, product_id, sku, barcode, price, stock_quantity, stock_status,
    variant_attributes, is_active, is_default, view_count, sold_count,
    created_at, updated_at
)
VALUES
    (19100, 9102, 'DEL-TEST-003-S', '1234567890102-S', 70.00, 10, 'in_stock', '{"size": "S"}'::jsonb, true, false, 0, 0, NOW(), NOW()),
    (19101, 9102, 'DEL-TEST-003-M', '1234567890102-M', 70.00, 10, 'in_stock', '{"size": "M"}'::jsonb, true, true, 0, 0, NOW(), NOW()),
    (19102, 9102, 'DEL-TEST-003-L', '1234567890102-L', 70.00, 10, 'in_stock', '{"size": "L"}'::jsonb, true, false, 0, 0, NOW(), NOW()),
    (19103, 9102, 'DEL-TEST-003-XL', '1234567890102-XL', 70.00, 10, 'in_stock', '{"size": "XL"}'::jsonb, true, false, 0, 0, NOW(), NOW()),
    (19104, 9102, 'DEL-TEST-003-XXL', '1234567890102-XXL', 70.00, 10, 'in_stock', '{"size": "XXL"}'::jsonb, true, false, 0, 0, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- Large Batch Test Data (Product IDs 9150-9249) - 100 products
-- ============================================================================

INSERT INTO b2c_products (
    id, storefront_id, name, description, price, currency, category_id,
    sku, barcode, stock_quantity, stock_status, is_active,
    attributes, view_count, sold_count, has_variants,
    created_at, updated_at
)
SELECT
    9150 + i,
    9001,
    'Bulk Product ' || i,
    'Product for bulk operations testing',
    (50.00 + (i * 1.5))::numeric(10,2),
    'USD',
    9000 + (i % 3),
    'BULK-PROD-' || LPAD(i::text, 3, '0'),
    '1234567890' || LPAD((150 + i)::text, 3, '0'),
    50 + i,
    CASE
        WHEN i % 10 = 0 THEN 'out_of_stock'
        WHEN i % 5 = 0 THEN 'low_stock'
        ELSE 'in_stock'
    END,
    true,
    '{}'::jsonb,
    i * 2,
    i,
    false,
    NOW(),
    NOW()
FROM generate_series(1, 100) AS i
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- Performance & E2E Test Data (Product IDs 9300-9310)
-- ============================================================================

INSERT INTO b2c_products (
    id, storefront_id, name, description, price, currency, category_id,
    sku, barcode, stock_quantity, stock_status, is_active,
    attributes, view_count, sold_count, has_variants,
    created_at, updated_at
)
VALUES
    (9300, 9000, 'E2E Test Product', 'Product for end-to-end workflow', 999.99, 'USD', 9000, 'E2E-PROD-001', '1234567890300', 100, 'in_stock', true, '{"test": "e2e"}'::jsonb, 0, 0, true, NOW(), NOW()),
    (9301, 9000, 'Performance Test 1', 'Product for performance testing', 10.00, 'USD', 9000, 'PERF-001', '1234567890301', 1000, 'in_stock', true, '{}'::jsonb, 0, 0, false, NOW(), NOW()),
    (9302, 9000, 'Performance Test 2', 'Product for performance testing', 20.00, 'USD', 9001, 'PERF-002', '1234567890302', 2000, 'in_stock', true, '{}'::jsonb, 0, 0, false, NOW(), NOW()),
    (9303, 9000, 'Performance Test 3', 'Product for performance testing', 30.00, 'USD', 9002, 'PERF-003', '1234567890303', 3000, 'in_stock', true, '{}'::jsonb, 0, 0, false, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Variants for E2E test product 9300
INSERT INTO b2c_product_variants (
    id, product_id, sku, barcode, price, stock_quantity, stock_status,
    variant_attributes, is_active, is_default, view_count, sold_count,
    created_at, updated_at
)
VALUES
    (19300, 9300, 'E2E-PROD-001-V1', '1234567890300-V1', 999.99, 50, 'in_stock', '{"variant": "V1"}'::jsonb, true, true, 0, 0, NOW(), NOW()),
    (19301, 9300, 'E2E-PROD-001-V2', '1234567890300-V2', 999.99, 50, 'in_stock', '{"variant": "V2"}'::jsonb, true, false, 0, 0, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- Verification Queries (for manual testing - commented out)
-- ============================================================================

-- SELECT COUNT(*) FROM b2c_products WHERE id BETWEEN 9000 AND 9099; -- Should be ~20
-- SELECT COUNT(*) FROM b2c_products WHERE id BETWEEN 9020 AND 9024 AND deleted_at IS NOT NULL; -- Should be 5
-- SELECT COUNT(*) FROM b2c_products WHERE id BETWEEN 9100 AND 9199; -- Should be ~15
-- SELECT COUNT(*) FROM b2c_products WHERE id BETWEEN 9150 AND 9249; -- Should be 100
-- SELECT COUNT(*) FROM b2c_product_variants WHERE product_id IN (9001, 9102, 9300); -- Should be ~10
-- SELECT COUNT(*) FROM b2c_product_images WHERE product_id = 9002; -- Should be 3
