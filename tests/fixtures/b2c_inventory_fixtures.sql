-- B2C Inventory Integration Test Fixtures
-- This file provides test data for B2C inventory integration tests
-- Uses unified listings table with source_type='b2c'

-- Test storefronts (from migration 000003)
INSERT INTO storefronts (id, user_id, name, slug, description, is_active, created_at, updated_at)
VALUES
    (1000, 1000, 'Test Inventory Storefront 1', 'test-inventory-store-1', 'Test storefront for inventory tests', true, NOW(), NOW()),
    (1001, 1001, 'Test Inventory Storefront 2', 'test-inventory-store-2', 'Another test storefront', true, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Test listings for inventory operations
INSERT INTO listings (
    id, user_id, storefront_id, title, description,
    price, currency, category_id,
    sku, quantity, status,
    source_type,
    created_at, updated_at
)
VALUES
    -- Product with sufficient stock
    (
        5000, -- id
        1000, -- user_id
        1000, -- storefront_id
        'Test Product - Sufficient Stock',
        'Product with enough stock for testing',
        100.00, -- price
        'USD',  -- currency
        1301,   -- category_id (Laptops)
        'TEST-STOCK-001', -- sku
        100, -- quantity
        'active', -- status
        'b2c', -- source_type
        NOW(),
        NOW()
    ),
    -- Product with low stock
    (
        5001,
        1000,
        1000,
        'Test Product - Low Stock',
        'Product with low stock',
        50.00,
        'USD',
        1301,
        'TEST-STOCK-002',
        5, -- low stock
        'active',
        'b2c',
        NOW(),
        NOW()
    ),
    -- Product out of stock
    (
        5002,
        1000,
        1000,
        'Test Product - Out of Stock',
        'Product with zero stock',
        75.00,
        'USD',
        1301,
        'TEST-STOCK-003',
        0, -- out of stock
        'inactive',
        'b2c',
        NOW(),
        NOW()
    ),
    -- Product for batch update test
    (
        5003,
        1000,
        1000,
        'Test Product - Batch Update 1',
        'First product for batch testing',
        25.00,
        'USD',
        1302,
        'TEST-BATCH-001',
        50,
        'active',
        'b2c',
        NOW(),
        NOW()
    ),
    (
        5004,
        1000,
        1000,
        'Test Product - Batch Update 2',
        'Second product for batch testing',
        30.00,
        'USD',
        1302,
        'TEST-BATCH-002',
        75,
        'active',
        'b2c',
        NOW(),
        NOW()
    ),
    -- Product for stats testing (inactive)
    (
        5005,
        1000,
        1000,
        'Test Product - Inactive',
        'Inactive product for stats',
        40.00,
        'USD',
        1301,
        'TEST-INACTIVE-001',
        20,
        'inactive',
        'b2c',
        NOW(),
        NOW()
    ),
    -- Product for view count testing
    (
        5006,
        1000,
        1000,
        'Test Product - View Counter',
        'Product for testing view increments',
        60.00,
        'USD',
        1301,
        'TEST-VIEW-001',
        30,
        'active',
        'b2c',
        NOW(),
        NOW()
    ),
    -- Product in second storefront
    (
        5007,
        1001,
        1001,
        'Test Product - Storefront 2',
        'Product in second storefront',
        80.00,
        'EUR',
        1301,
        'TEST-STORE2-001',
        40,
        'active',
        'b2c',
        NOW(),
        NOW()
    )
ON CONFLICT (id) DO NOTHING;

-- Insert initial inventory movements (history)
INSERT INTO inventory_movements (
    id, listing_id, variant_id, movement_type, quantity, reason, notes, user_id, created_at
)
VALUES
    -- Initial stock for product 5000
    (
        7000,
        5000, -- listing_id (было storefront_product_id)
        NULL, -- variant_id
        'in', -- movement_type (было type)
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
