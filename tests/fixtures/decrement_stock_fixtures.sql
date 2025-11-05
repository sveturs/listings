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

INSERT INTO b2c_products (
    id, storefront_id, name, description,
    price, currency, category_id,
    sku, barcode, stock_quantity, stock_status,
    is_active, view_count, sold_count,
    created_at, updated_at
)
VALUES
    -- 8000: Standard product with sufficient stock
    (
        8000,
        2000,
        'Test Product - Standard Stock',
        'Product for standard decrement testing',
        100.00,
        'USD',
        2000,
        'DECR-TEST-001',
        '8000000000001',
        100, -- sufficient stock
        'in_stock',
        true,
        0,
        0,
        NOW(),
        NOW()
    ),
    -- 8001: Product with moderate stock
    (
        8001,
        2000,
        'Test Product - Moderate Stock',
        'Product with moderate inventory',
        50.00,
        'USD',
        2000,
        'DECR-TEST-002',
        '8000000000002',
        50,
        'in_stock',
        true,
        0,
        0,
        NOW(),
        NOW()
    ),
    -- 8002: Product with large stock
    (
        8002,
        2000,
        'Test Product - Large Stock',
        'Product with large inventory',
        75.00,
        'USD',
        2000,
        'DECR-TEST-003',
        '8000000000003',
        200,
        'in_stock',
        true,
        0,
        0,
        NOW(),
        NOW()
    ),
    -- 8003: Product with small stock (for exact quantity test)
    (
        8003,
        2000,
        'Test Product - Small Stock',
        'Product with small inventory',
        25.00,
        'USD',
        2000,
        'DECR-TEST-004',
        '8000000000004',
        10, -- small stock
        'low_stock',
        true,
        0,
        0,
        NOW(),
        NOW()
    ),
    -- 8004: Product with variants (parent)
    (
        8004,
        2000,
        'Test Product - With Variants',
        'Product that has size variants',
        100.00,
        'USD',
        2000,
        'DECR-TEST-005',
        '8000000000005',
        0, -- no stock at product level (has variants)
        'in_stock',
        true,
        0,
        0,
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
        'Test Product - Concurrent Stock',
        'Product for concurrent decrement testing',
        50.00,
        'USD',
        2000,
        'DECR-CONC-001',
        '8000000000006',
        100, -- enough for concurrent requests
        'in_stock',
        true,
        0,
        0,
        NOW(),
        NOW()
    ),
    -- 8006: Product for overselling prevention test
    (
        8006,
        2000,
        'Test Product - Overselling Test',
        'Product to test overselling prevention',
        30.00,
        'USD',
        2000,
        'DECR-OVER-001',
        '8000000000007',
        20, -- limited stock to trigger overselling scenario
        'low_stock',
        true,
        0,
        0,
        NOW(),
        NOW()
    ),
    -- 8007: Product for transaction isolation test
    (
        8007,
        2000,
        'Test Product - Transaction Isolation',
        'Product for testing transaction isolation',
        40.00,
        'USD',
        2000,
        'DECR-TX-001',
        '8000000000008',
        50,
        'in_stock',
        true,
        0,
        0,
        NOW(),
        NOW()
    ),
    -- 8008: Product for high-frequency test (large stock)
    (
        8008,
        2000,
        'Test Product - High Frequency',
        'Product with large stock for high-frequency testing',
        60.00,
        'USD',
        2000,
        'DECR-FREQ-001',
        '8000000000009',
        1000, -- large stock for many operations
        'in_stock',
        true,
        0,
        0,
        NOW(),
        NOW()
    )
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- Products for Large Batch Test (8010-8059)
-- ============================================================================

-- Insert 50 products for large batch testing
INSERT INTO b2c_products (
    id, storefront_id, name, description,
    price, currency, category_id,
    sku, barcode, stock_quantity, stock_status,
    is_active, view_count, sold_count,
    created_at, updated_at
)
SELECT
    8010 + seq AS id,
    2000 AS storefront_id,
    'Batch Test Product ' || seq AS name,
    'Product for batch testing' AS description,
    10.00 + seq AS price,
    'USD' AS currency,
    2001 AS category_id,
    'BATCH-' || LPAD(seq::text, 4, '0') AS sku,
    '8001' || LPAD(seq::text, 8, '0') AS barcode,
    10 AS stock_quantity, -- enough for batch test
    'in_stock' AS stock_status,
    true AS is_active,
    0 AS view_count,
    0 AS sold_count,
    NOW() AS created_at,
    NOW() AS updated_at
FROM generate_series(0, 49) AS seq
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- Product Variants for Variant-Level Tests
-- ============================================================================

-- Update product 8004 to have variants
UPDATE b2c_products
SET has_variants = true
WHERE id = 8004;

INSERT INTO b2c_product_variants (
    id, product_id, sku, barcode,
    price, stock_quantity, stock_status,
    variant_attributes, is_active, is_default,
    view_count, sold_count,
    created_at, updated_at
)
VALUES
    -- 9000: Size S variant for product 8004
    (
        9000,
        8004,
        'DECR-TEST-005-S',
        '8000000000005-S',
        100.00,
        50, -- sufficient stock
        'in_stock',
        '{"size": "S"}'::jsonb,
        true,
        true, -- default variant
        0,
        0,
        NOW(),
        NOW()
    ),
    -- 9001: Size M variant for product 8004
    (
        9001,
        8004,
        'DECR-TEST-005-M',
        '8000000000005-M',
        100.00,
        30,
        'in_stock',
        '{"size": "M"}'::jsonb,
        true,
        false,
        0,
        0,
        NOW(),
        NOW()
    ),
    -- 9002: Size L variant for product 8004
    (
        9003,
        8004,
        'DECR-TEST-005-L',
        '8000000000005-L',
        100.00,
        20,
        'low_stock',
        '{"size": "L"}'::jsonb,
        true,
        false,
        0,
        0,
        NOW(),
        NOW()
    )
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- Cleanup old test data (in case of re-runs)
-- ============================================================================

-- Clean any leftover inventory movements from previous test runs
DELETE FROM b2c_inventory_movements
WHERE storefront_product_id >= 8000 AND storefront_product_id < 9000;
