-- Inventory Integration Test Fixtures
-- This file provides test data for inventory integration tests
-- Note: This is a listings microservice, so we only use tables that exist in migrations

-- Test storefronts (from migration 000003)
INSERT INTO storefronts (id, user_id, name, slug, description, is_active, created_at, updated_at)
VALUES
    (1000, 1000, 'Test Inventory Storefront 1', 'test-inventory-store-1', 'Test storefront for inventory tests', true, NOW(), NOW()),
    (1001, 1001, 'Test Inventory Storefront 2', 'test-inventory-store-2', 'Another test storefront', true, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Test products (b2b listings) for inventory operations
INSERT INTO listings (
    id, uuid, user_id, storefront_id, title, description,
    price, currency, category_id, status,
    quantity, sku, is_b2c, created_at, updated_at
)
VALUES
    -- Product with sufficient stock
    (
        5000,
        gen_random_uuid(),
        1000,
        1000,
        'Test Product - Sufficient Stock',
        'Product with enough stock for testing',
        100.00,
        'USD',
        2000,
        'active',
        100,
        'TEST-STOCK-001',
        false, -- B2B listing
        NOW(),
        NOW()
    ),
    -- Product with low stock
    (
        5001,
        gen_random_uuid(),
        1000,
        1000,
        'Test Product - Low Stock',
        'Product with low stock',
        50.00,
        'USD',
        2000,
        'active',
        5,
        'TEST-STOCK-002',
        false,
        NOW(),
        NOW()
    ),
    -- Product out of stock
    (
        5002,
        gen_random_uuid(),
        1000,
        1000,
        'Test Product - Out of Stock',
        'Product with zero stock',
        75.00,
        'USD',
        2000,
        'active',
        0,
        'TEST-STOCK-003',
        false,
        NOW(),
        NOW()
    ),
    -- Product for batch update test
    (
        5003,
        gen_random_uuid(),
        1000,
        1000,
        'Test Product - Batch Update 1',
        'First product for batch testing',
        25.00,
        'USD',
        2001,
        'active',
        50,
        'TEST-BATCH-001',
        false,
        NOW(),
        NOW()
    ),
    (
        5004,
        gen_random_uuid(),
        1000,
        1000,
        'Test Product - Batch Update 2',
        'Second product for batch testing',
        30.00,
        'USD',
        2001,
        'active',
        75,
        'TEST-BATCH-002',
        false,
        NOW(),
        NOW()
    ),
    -- Product for stats testing (inactive)
    (
        5005,
        gen_random_uuid(),
        1000,
        1000,
        'Test Product - Inactive',
        'Inactive product for stats',
        40.00,
        'USD',
        2000,
        'inactive',
        20,
        'TEST-INACTIVE-001',
        false,
        NOW(),
        NOW()
    ),
    -- Product for view count testing
    (
        5006,
        gen_random_uuid(),
        1000,
        1000,
        'Test Product - View Counter',
        'Product for testing view increments',
        60.00,
        'USD',
        2000,
        'active',
        30,
        'TEST-VIEW-001',
        false,
        NOW(),
        NOW()
    ),
    -- Product in second storefront
    (
        5007,
        gen_random_uuid(),
        1001,
        1001,
        'Test Product - Storefront 2',
        'Product in second storefront',
        80.00,
        'EUR',
        2000,
        'active',
        40,
        'TEST-STORE2-001',
        false,
        NOW(),
        NOW()
    )
ON CONFLICT (id) DO NOTHING;

-- Initialize listing_stats for view count testing
INSERT INTO listing_stats (listing_id, views, favorites, shares, created_at, updated_at)
VALUES
    (5000, 0, 0, 0, NOW(), NOW()),
    (5001, 0, 0, 0, NOW(), NOW()),
    (5002, 0, 0, 0, NOW(), NOW()),
    (5003, 0, 0, 0, NOW(), NOW()),
    (5004, 0, 0, 0, NOW(), NOW()),
    (5005, 0, 0, 0, NOW(), NOW()),
    (5006, 10, 0, 0, NOW(), NOW()), -- Already has 10 views
    (5007, 0, 0, 0, NOW(), NOW())
ON CONFLICT (listing_id) DO NOTHING;

-- Note: inventory_movements, product_variants tables don't exist in listings microservice
-- They will need to be added via migrations if needed for inventory tracking
-- For now, tests will focus on listing quantity management
