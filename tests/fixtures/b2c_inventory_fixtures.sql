-- B2C Inventory Integration Test Fixtures
-- This file provides test data for B2C inventory integration tests
-- Uses b2c_products tables created in migration 000004

-- Test storefronts (from migration 000003)
INSERT INTO storefronts (id, user_id, name, slug, description, is_active, created_at, updated_at)
VALUES
    (1000, 1000, 'Test Inventory Storefront 1', 'test-inventory-store-1', 'Test storefront for inventory tests', true, NOW(), NOW()),
    (1001, 1001, 'Test Inventory Storefront 2', 'test-inventory-store-2', 'Another test storefront', true, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Test b2c_products for inventory operations
INSERT INTO b2c_products (
    id, storefront_id, name, description,
    price, currency, category_id,
    sku, barcode, stock_quantity, stock_status,
    is_active, view_count, sold_count,
    created_at, updated_at
)
VALUES
    -- Product with sufficient stock
    (
        5000, -- id
        1000, -- storefront_id
        'Test Product - Sufficient Stock',
        'Product with enough stock for testing',
        100.00, -- price
        'USD',  -- currency
        2000,   -- category_id
        'TEST-STOCK-001', -- sku
        '1234567890001',  -- barcode
        100, -- stock_quantity
        'in_stock', -- stock_status
        true, -- is_active
        0,    -- view_count
        0,    -- sold_count
        NOW(),
        NOW()
    ),
    -- Product with low stock
    (
        5001,
        1000,
        'Test Product - Low Stock',
        'Product with low stock',
        50.00,
        'USD',
        2000,
        'TEST-STOCK-002',
        '1234567890002',
        5, -- low stock
        'low_stock',
        true,
        0,
        0,
        NOW(),
        NOW()
    ),
    -- Product out of stock
    (
        5002,
        1000,
        'Test Product - Out of Stock',
        'Product with zero stock',
        75.00,
        'USD',
        2000,
        'TEST-STOCK-003',
        '1234567890003',
        0, -- out of stock
        'out_of_stock',
        true,
        0,
        0,
        NOW(),
        NOW()
    ),
    -- Product for batch update test
    (
        5003,
        1000,
        'Test Product - Batch Update 1',
        'First product for batch testing',
        25.00,
        'USD',
        2001,
        'TEST-BATCH-001',
        '1234567890004',
        50,
        'in_stock',
        true,
        0,
        0,
        NOW(),
        NOW()
    ),
    (
        5004,
        1000,
        'Test Product - Batch Update 2',
        'Second product for batch testing',
        30.00,
        'USD',
        2001,
        'TEST-BATCH-002',
        '1234567890005',
        75,
        'in_stock',
        true,
        0,
        0,
        NOW(),
        NOW()
    ),
    -- Product for stats testing (inactive)
    (
        5005,
        1000,
        'Test Product - Inactive',
        'Inactive product for stats',
        40.00,
        'USD',
        2000,
        'TEST-INACTIVE-001',
        '1234567890006',
        20,
        'in_stock',
        false, -- inactive
        0,
        0,
        NOW(),
        NOW()
    ),
    -- Product for view count testing
    (
        5006,
        1000,
        'Test Product - View Counter',
        'Product for testing view increments',
        60.00,
        'USD',
        2000,
        'TEST-VIEW-001',
        '1234567890007',
        30,
        'in_stock',
        true,
        10, -- already has 10 views
        0,
        NOW(),
        NOW()
    ),
    -- Product in second storefront
    (
        5007,
        1001,
        'Test Product - Storefront 2',
        'Product in second storefront',
        80.00,
        'EUR',
        2000,
        'TEST-STORE2-001',
        '1234567890008',
        40,
        'in_stock',
        true,
        0,
        0,
        NOW(),
        NOW()
    )
ON CONFLICT (id) DO NOTHING;

-- Test b2c_product_variants
INSERT INTO b2c_product_variants (
    id, product_id, sku, barcode,
    price, stock_quantity, stock_status,
    variant_attributes, is_active, is_default,
    view_count, sold_count,
    created_at, updated_at
)
VALUES
    -- Variant for product 5000 (Size S)
    (
        6000, -- id
        5000, -- product_id
        'TEST-STOCK-001-S', -- sku
        '1234567890001-S',  -- barcode
        100.00, -- price
        50,     -- stock_quantity
        'in_stock',
        '{"size": "S"}'::jsonb, -- variant_attributes
        true,   -- is_active
        true,   -- is_default
        0,      -- view_count
        0,      -- sold_count
        NOW(),
        NOW()
    ),
    -- Variant for product 5000 (Size M)
    (
        6001,
        5000,
        'TEST-STOCK-001-M',
        '1234567890001-M',
        100.00,
        30,
        'in_stock',
        '{"size": "M"}'::jsonb,
        true,
        false, -- not default
        0,
        0,
        NOW(),
        NOW()
    )
ON CONFLICT (id) DO NOTHING;

-- Update has_variants flag for product 5000
UPDATE b2c_products
SET has_variants = true
WHERE id = 5000;

-- Insert initial inventory movements (history)
INSERT INTO b2c_inventory_movements (
    id, storefront_product_id, variant_id, type, quantity, reason, notes, user_id, created_at
)
VALUES
    -- Initial stock for product 5000
    (
        7000,
        5000, -- storefront_product_id
        NULL, -- variant_id
        'in',
        100,
        'initial_stock',
        'Initial product creation',
        1000, -- user_id
        NOW() - INTERVAL '7 days'
    ),
    -- Restock for product 5001
    (
        7001,
        5001,
        NULL,
        'in',
        5,
        'restock',
        'Supplier delivery',
        1000,
        NOW() - INTERVAL '5 days'
    ),
    -- Sale from product 5006
    (
        7002,
        5006,
        NULL,
        'out',
        5,
        'sale',
        'Customer purchase',
        1000,
        NOW() - INTERVAL '2 days'
    )
ON CONFLICT (id) DO NOTHING;
