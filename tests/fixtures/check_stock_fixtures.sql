-- CheckStockAvailability Integration Test Fixtures
-- Extended test data for comprehensive stock availability testing
-- Complements b2c_inventory_fixtures.sql with additional edge cases

-- Additional test storefronts for edge cases
INSERT INTO storefronts (id, user_id, name, slug, description, is_active, created_at, updated_at)
VALUES
    (1002, 1002, 'Edge Case Storefront', 'edge-case-store', 'Storefront for edge case testing', true, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Products for specific CheckStockAvailability scenarios
INSERT INTO b2c_products (
    id, storefront_id, name, description,
    price, currency, category_id,
    sku, barcode, stock_quantity, stock_status,
    is_active, view_count, sold_count,
    created_at, updated_at
)
VALUES
    -- Product for exact match test (quantity = 10)
    (
        5100,
        1000,
        'Exact Match Product',
        'Product with exactly 10 units for exact match testing',
        100.00,
        'USD',
        2000,
        'TEST-EXACT-001',
        '1234567890100',
        10, -- exact quantity for testing
        'in_stock',
        true,
        0,
        0,
        NOW(),
        NOW()
    ),
    -- Product for boundary test (quantity = 1)
    (
        5101,
        1000,
        'Single Unit Product',
        'Product with only 1 unit',
        50.00,
        'USD',
        2000,
        'TEST-SINGLE-001',
        '1234567890101',
        1, -- minimum stock
        'low_stock',
        true,
        0,
        0,
        NOW(),
        NOW()
    ),
    -- Product for large quantity test (quantity = 10000)
    (
        5102,
        1000,
        'Large Stock Product',
        'Product with large stock quantity',
        25.00,
        'USD',
        2001,
        'TEST-LARGE-001',
        '1234567890102',
        10000, -- large stock
        'in_stock',
        true,
        0,
        0,
        NOW(),
        NOW()
    ),
    -- Product for concurrent testing (quantity = 1000)
    (
        5103,
        1000,
        'Concurrent Test Product',
        'Product for testing concurrent availability checks',
        75.00,
        'USD',
        2000,
        'TEST-CONCURRENT-001',
        '1234567890103',
        1000,
        'in_stock',
        true,
        0,
        0,
        NOW(),
        NOW()
    ),
    -- Additional zero stock product
    (
        5104,
        1000,
        'Zero Stock Product 2',
        'Another product with zero stock',
        60.00,
        'USD',
        2000,
        'TEST-ZERO-002',
        '1234567890104',
        0,
        'out_of_stock',
        true,
        0,
        0,
        NOW(),
        NOW()
    ),
    -- Product for variant testing (has multiple variants)
    (
        5105,
        1000,
        'Multi-Variant Product',
        'Product with multiple variants for testing',
        120.00,
        'USD',
        2001,
        'TEST-VARIANT-001',
        '1234567890105',
        200, -- base product stock
        'in_stock',
        true,
        0,
        0,
        NOW(),
        NOW()
    ),
    -- Inactive product with stock (should still check availability)
    (
        5106,
        1000,
        'Inactive Stock Product',
        'Inactive product but has stock',
        40.00,
        'USD',
        2000,
        'TEST-INACTIVE-STOCK-001',
        '1234567890106',
        50,
        'in_stock',
        false, -- inactive
        0,
        0,
        NOW(),
        NOW()
    )
ON CONFLICT (id) DO NOTHING;

-- Additional variants for variant-level testing
INSERT INTO b2c_product_variants (
    id, product_id, sku, barcode,
    price, stock_quantity, stock_status,
    variant_attributes, is_active, is_default,
    view_count, sold_count,
    created_at, updated_at
)
VALUES
    -- Variant with zero stock
    (
        6100,
        5105,
        'TEST-VARIANT-001-XS',
        '1234567890105-XS',
        120.00,
        0, -- out of stock
        'out_of_stock',
        '{"size": "XS"}'::jsonb,
        true,
        false,
        0,
        0,
        NOW(),
        NOW()
    ),
    -- Variant with low stock
    (
        6101,
        5105,
        'TEST-VARIANT-001-L',
        '1234567890105-L',
        120.00,
        3, -- low stock
        'low_stock',
        '{"size": "L"}'::jsonb,
        true,
        false,
        0,
        0,
        NOW(),
        NOW()
    ),
    -- Variant with exact match quantity
    (
        6102,
        5105,
        'TEST-VARIANT-001-XL',
        '1234567890105-XL',
        120.00,
        15, -- for exact match tests
        'in_stock',
        '{"size": "XL"}'::jsonb,
        true,
        false,
        0,
        0,
        NOW(),
        NOW()
    ),
    -- Variant with large stock
    (
        6103,
        5105,
        'TEST-VARIANT-001-XXL',
        '1234567890105-XXL',
        120.00,
        500, -- large stock
        'in_stock',
        '{"size": "XXL"}'::jsonb,
        true,
        false,
        0,
        0,
        NOW(),
        NOW()
    ),
    -- Inactive variant with stock
    (
        6104,
        5105,
        'TEST-VARIANT-001-INACTIVE',
        '1234567890105-INACTIVE',
        120.00,
        100,
        'in_stock',
        '{"size": "Custom", "status": "discontinued"}'::jsonb,
        false, -- inactive
        false,
        0,
        0,
        NOW(),
        NOW()
    ),
    -- Variant for product 5000 (Size L) - additional to existing S and M
    (
        6105,
        5000,
        'TEST-STOCK-001-L',
        '1234567890001-L',
        100.00,
        20,
        'in_stock',
        '{"size": "L"}'::jsonb,
        true,
        false,
        0,
        0,
        NOW(),
        NOW()
    )
ON CONFLICT (id) DO NOTHING;

-- Update has_variants flag for product 5105
UPDATE b2c_products
SET has_variants = true
WHERE id = 5105;

-- Performance test data: Multiple products for batch testing
INSERT INTO b2c_products (
    id, storefront_id, name, description,
    price, currency, category_id,
    sku, barcode, stock_quantity, stock_status,
    is_active, view_count, sold_count,
    created_at, updated_at
)
SELECT
    5200 + i,
    1000,
    'Batch Test Product ' || i,
    'Product for batch performance testing',
    50.00 + (i * 5.00),
    'USD',
    2000 + (i % 3),
    'TEST-BATCH-' || LPAD(i::text, 3, '0'),
    '1234567890200' || LPAD(i::text, 3, '0'),
    100 + (i * 10), -- varying stock quantities
    CASE
        WHEN i % 10 = 0 THEN 'out_of_stock'
        WHEN i % 5 = 0 THEN 'low_stock'
        ELSE 'in_stock'
    END,
    true,
    0,
    0,
    NOW(),
    NOW()
FROM generate_series(1, 50) AS i
ON CONFLICT (id) DO NOTHING;

-- Comment explaining test data structure
COMMENT ON TABLE b2c_products IS 'B2C products table with test fixtures for integration testing';

-- Verification queries (commented out, for reference)
-- SELECT COUNT(*) FROM b2c_products WHERE id >= 5100 AND id < 5200; -- Should return ~7
-- SELECT COUNT(*) FROM b2c_products WHERE id >= 5200 AND id < 5300; -- Should return 50
-- SELECT COUNT(*) FROM b2c_product_variants WHERE id >= 6100; -- Should return ~6
