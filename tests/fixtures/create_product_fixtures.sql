-- CreateProduct Integration Test Fixtures
-- Comprehensive test data for CreateProduct and BulkCreateProducts APIs
-- Phase 9.7.2 - Product CRUD Integration Tests

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
    (1100, 'create-test@example.com', 'create_test_user', NOW(), NOW()),
    (1101, 'bulk-test@example.com', 'bulk_test_user', NOW(), NOW()),
    (1102, 'edge-test@example.com', 'edge_test_user', NOW(), NOW())
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

-- ============================================================================
-- Test Storefronts
-- ============================================================================
-- Storefront for basic product tests
INSERT INTO storefronts (id, user_id, name, slug, description, is_active, created_at, updated_at)
VALUES
    (1100, 1100, 'CreateProduct Test Store', 'create-product-test', 'Storefront for CreateProduct tests', true, NOW(), NOW()),
    (1101, 1101, 'Bulk Import Store', 'bulk-import-store', 'Storefront for bulk operations', true, NOW(), NOW()),
    (1102, 1102, 'Edge Case Store', 'edge-case-store', 'Storefront for edge case testing', true, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- Additional Test Categories (CreateProduct specific)
-- ============================================================================
-- Top-level categories
INSERT INTO categories (id, name, slug, parent_id, level, is_active, sort_order, created_at, updated_at)
VALUES
    (2100, 'Electronics Test', 'electronics-test', NULL, 0, true, 1, NOW(), NOW()),
    (2101, 'Clothing Test', 'clothing-test', NULL, 0, true, 2, NOW(), NOW()),
    (2102, 'Books Test', 'books-test', NULL, 0, true, 3, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Sub-categories
INSERT INTO categories (id, name, slug, parent_id, level, is_active, sort_order, created_at, updated_at)
VALUES
    (2110, 'Smartphones', 'smartphones-test', 2100, 1, true, 1, NOW(), NOW()),
    (2111, 'Laptops', 'laptops-test', 2100, 1, true, 2, NOW(), NOW()),
    (2120, 'T-Shirts', 'tshirts-test', 2101, 1, true, 1, NOW(), NOW()),
    (2121, 'Jeans', 'jeans-test', 2101, 1, true, 2, NOW(), NOW()),
    (2130, 'Fiction', 'fiction-test', 2102, 1, true, 1, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- Existing Products for SKU Uniqueness Tests
-- ============================================================================
-- These products are created to test duplicate SKU validation
INSERT INTO listings (
    id, user_id, storefront_id, title, description,
    price, currency, category_id,
    sku, quantity, status,
    source_type,
    created_at, updated_at
)
VALUES
    -- Product with existing SKU for duplicate test
    (
        7000,
        1100,
        1100,
        'Existing Product 1',
        'Product with SKU TEST-SKU-001 already in database',
        50.00,
        'USD',
        2110,
        'TEST-SKU-001', -- This SKU will be used in duplicate test
        100,
        'active',
        'b2c',
        NOW(),
        NOW()
    ),
    -- Product with another existing SKU
    (
        7001,
        1100,
        1100,
        'Existing Product 2',
        'Product with SKU TEST-SKU-002 already in database',
        75.00,
        'USD',
        2111,
        'TEST-SKU-002', -- This SKU will be used in duplicate test
        50,
        'active',
        'b2c',
        NOW(),
        NOW()
    ),
    -- Product in different storefront with unique SKU
    (
        7002,
        1101,
        1101,
        'Different Storefront Product',
        'Product in different storefront (SKU must be globally unique)',
        60.00,
        'USD',
        2110,
        'TEST-SKU-003', -- Changed to unique SKU (global constraint)
        25,
        'active',
        'b2c',
        NOW(),
        NOW()
    )
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- Products for Performance Testing (Bulk Create)
-- ============================================================================
-- Pre-create storefront capacity check
-- This ensures storefront 1101 can handle bulk imports
COMMENT ON TABLE storefronts IS 'Test storefronts with capacity for bulk product creation';

-- ============================================================================
-- Cleanup Helper Comments
-- ============================================================================
-- To clean up test data, run:
-- DELETE FROM listing_variants WHERE listing_id >= 7000 AND listing_id < 8000;
-- DELETE FROM listings WHERE id >= 7000 AND id < 8000;
-- DELETE FROM categories WHERE id >= 2100 AND id < 2200;
-- DELETE FROM storefronts WHERE id >= 1100 AND id < 1200;
-- DELETE FROM users WHERE id >= 1100 AND id < 1200;

-- ============================================================================
-- Verification Queries (for debugging)
-- ============================================================================
-- SELECT COUNT(*) FROM storefronts WHERE id >= 1100 AND id < 1200; -- Should return 3
-- SELECT COUNT(*) FROM categories WHERE id >= 2100 AND id < 2200; -- Should return 8
-- SELECT COUNT(*) FROM listings WHERE id >= 7000 AND id < 8000; -- Should return 3
-- SELECT sku, COUNT(*) FROM listings WHERE id >= 7000 GROUP BY sku HAVING COUNT(*) > 1; -- Should return 0 (no duplicates within storefront)
