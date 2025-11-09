-- Decrement Stock Integration Test Fixtures
-- This file provides comprehensive test data for DecrementStock operation testing
-- Covers: happy path, error cases, concurrency, edge cases

-- Test storefronts
INSERT INTO storefronts (id, user_id, name, slug, description, is_active, created_at, updated_at)
VALUES
    (2000, 2000, 'Decrement Stock Test Store', 'decr-stock-store', 'Store for decrement stock testing', true, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- Products for Happy Path Tests
-- ============================================================================

INSERT INTO listings (
    id, user_id, storefront_id, title, description,
    price, currency, category_id,
    sku, quantity, status,
    source_type,
    created_at, updated_at
)
VALUES
    -- 8000: Standard product with sufficient stock
    (
        8000,
        2000,
        2000,
        'Test Product - Standard Stock',
        'Product for standard decrement testing',
        100.00,
        'USD',
        2000,
        'DECR-TEST-001',
        100, -- sufficient stock
        'active',
        'b2c',
        NOW(),
        NOW()
    ),
    -- 8001: Product with moderate stock
    (
        8001,
        2000,
        2000,
        'Test Product - Moderate Stock',
        'Product with moderate inventory',
        50.00,
        'USD',
        2000,
        'DECR-TEST-002',
        50,
        'active',
        'b2c',
        NOW(),
        NOW()
    ),
    -- 8002: Product with large stock
    (
        8002,
        2000,
        2000,
        'Test Product - Large Stock',
        'Product with large inventory',
        75.00,
        'USD',
        2000,
        'DECR-TEST-003',
        200,
        'active',
        'b2c',
        NOW(),
        NOW()
    ),
    -- 8003: Product with small stock (for exact quantity test)
    (
        8003,
        2000,
        2000,
        'Test Product - Small Stock',
        'Product with small inventory',
        25.00,
        'USD',
        2000,
        'DECR-TEST-004',
        10, -- small stock
        'active',
        'b2c',
        NOW(),
        NOW()
    ),
    -- 8004: Product with variants (parent)
    (
        8004,
        2000,
        2000,
        'Test Product - With Variants',
        'Product that has size variants',
        100.00,
        'USD',
        2000,
        'DECR-TEST-005',
        0, -- no stock at product level (has variants)
        'active',
        'b2c',
        NOW(),
        NOW()
    ),
    -- ============================================================================
    -- Products for Concurrency Tests
    -- ============================================================================
    -- 8005: Product for concurrent same-product test
    (
        8005,
        2000,
        2000,
        'Test Product - Concurrent Stock',
        'Product for concurrent decrement testing',
        50.00,
        'USD',
        2000,
        'DECR-CONC-001',
        100, -- enough for concurrent requests
        'active',
        'b2c',
        NOW(),
        NOW()
    ),
    -- 8006: Product for overselling prevention test
    (
        8006,
        2000,
        2000,
        'Test Product - Overselling Test',
        'Product to test overselling prevention',
        30.00,
        'USD',
        2000,
        'DECR-OVER-001',
        20, -- limited stock to trigger overselling scenario
        'active',
        'b2c',
        NOW(),
        NOW()
    ),
    -- 8007: Product for transaction isolation test
    (
        8007,
        2000,
        2000,
        'Test Product - Transaction Isolation',
        'Product for testing transaction isolation',
        40.00,
        'USD',
        2000,
        'DECR-TX-001',
        50,
        'active',
        'b2c',
        NOW(),
        NOW()
    ),
    -- 8008: Product for high-frequency test (large stock)
    (
        8008,
        2000,
        2000,
        'Test Product - High Frequency',
        'Product with large stock for high-frequency testing',
        60.00,
        'USD',
        2000,
        'DECR-FREQ-001',
        1000, -- large stock for many operations
        'active',
        'b2c',
        NOW(),
        NOW()
    )
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- Products for Large Batch Test (8010-8059)
-- ============================================================================

-- Insert 50 products for large batch testing
INSERT INTO listings (
    id, user_id, storefront_id, title, description,
    price, currency, category_id,
    sku, quantity, status,
    source_type,
    created_at, updated_at
)
SELECT
    8010 + seq AS id,
    2000 AS user_id,
    2000 AS storefront_id,
    'Batch Test Product ' || seq AS title,
    'Product for batch testing' AS description,
    10.00 + seq AS price,
    'USD' AS currency,
    2001 AS category_id,
    'BATCH-' || LPAD(seq::text, 4, '0') AS sku,
    10 AS quantity, -- enough for batch test
    'active' AS status,
    'b2c' AS source_type,
    NOW() AS created_at,
    NOW() AS updated_at
FROM generate_series(0, 49) AS seq
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- Product Variants for Variant-Level Tests
-- ============================================================================

-- ============================================================================
-- Cleanup old test data (in case of re-runs)
-- ============================================================================

-- Clean any leftover inventory movements from previous test runs
DELETE FROM inventory_movements
WHERE listing_id >= 8000 AND listing_id < 9000;
