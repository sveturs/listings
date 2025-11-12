-- CheckStockAvailability Integration Test Fixtures
-- Extended test data for comprehensive stock availability testing
-- Complements b2c_inventory_fixtures.sql with additional edge cases

-- Additional test storefronts for edge cases
INSERT INTO storefronts (id, user_id, name, slug, description, is_active, created_at, updated_at)
VALUES
    (1002, 1002, 'Edge Case Storefront', 'edge-case-store', 'Storefront for edge case testing', true, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Products for specific CheckStockAvailability scenarios
INSERT INTO listings (
    id, user_id, storefront_id, title, description,
    price, currency, category_id,
    sku, quantity, status,
    source_type,
    created_at, updated_at
)
VALUES
    -- Product for exact match test (quantity = 10)
    (
        5100,
        1000,
        1000,
        'Exact Match Product',
        'Product with exactly 10 units for exact match testing',
        100.00,
        'USD',
        2000,
        'TEST-EXACT-001',
        10,
        'active',
        'b2c',
        NOW(),
        NOW()
    ),
    -- Product for boundary test (quantity = 1)
    (
        5101,
        1000,
        1000,
        'Single Unit Product',
        'Product with only 1 unit',
        50.00,
        'USD',
        2000,
        'TEST-SINGLE-001',
        1,
        'active',
        'b2c',
        NOW(),
        NOW()
    ),
    -- Product for large quantity test (quantity = 10000)
    (
        5102,
        1000,
        1000,
        'Large Stock Product',
        'Product with large stock quantity',
        25.00,
        'USD',
        2001,
        'TEST-LARGE-001',
        10000,
        'active',
        'b2c',
        NOW(),
        NOW()
    ),
    -- Product for concurrent testing (quantity = 1000)
    (
        5103,
        1000,
        1000,
        'Concurrent Test Product',
        'Product for testing concurrent availability checks',
        75.00,
        'USD',
        2000,
        'TEST-CONCURRENT-001',
        1000,
        'active',
        'b2c',
        NOW(),
        NOW()
    ),
    -- Additional zero stock product
    (
        5104,
        1000,
        1000,
        'Zero Stock Product 2',
        'Another product with zero stock',
        60.00,
        'USD',
        2000,
        'TEST-ZERO-002',
        0,
        'inactive',
        'b2c',
        NOW(),
        NOW()
    ),
    -- Product for variant testing
    (
        5105,
        1000,
        1000,
        'Multi-Variant Product',
        'Product with multiple variants for testing',
        120.00,
        'USD',
        2001,
        'TEST-VARIANT-001',
        200,
        'active',
        'b2c',
        NOW(),
        NOW()
    ),
    -- Inactive product with stock
    (
        5106,
        1000,
        1000,
        'Inactive Stock Product',
        'Inactive product but has stock',
        40.00,
        'USD',
        2000,
        'TEST-INACTIVE-STOCK-001',
        50,
        'inactive',
        'b2c',
        NOW(),
        NOW()
    )
ON CONFLICT (id) DO NOTHING;

-- Performance test data: Multiple products for batch testing
INSERT INTO listings (
    id, user_id, storefront_id, title, description,
    price, currency, category_id,
    sku, quantity, status,
    source_type,
    created_at, updated_at
)
SELECT
    5200 + i,
    1000,
    1000,
    'Batch Test Product ' || i,
    'Product for batch performance testing',
    50.00 + (i * 5.00),
    'USD',
    2000 + (i % 3),
    'TEST-BATCH-' || LPAD(i::text, 3, '0'),
    100 + (i * 10),
    CASE
        WHEN i % 10 = 0 THEN 'inactive'
        ELSE 'active'
    END,
    'b2c',
    NOW(),
    NOW()
FROM generate_series(1, 50) AS i
ON CONFLICT (id) DO NOTHING;

-- Verification queries (commented out, for reference)
-- SELECT COUNT(*) FROM listings WHERE id >= 5100 AND id < 5200; -- Should return ~7
-- SELECT COUNT(*) FROM listings WHERE id >= 5200 AND id < 5300; -- Should return 50
