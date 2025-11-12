-- RollbackStock Integration Test Fixtures
-- This file provides test data specifically for RollbackStock testing
-- Focuses on compensating transaction scenarios and idempotency

-- Test storefronts (reuse storefront 1000 from b2c_inventory_fixtures if exists)
INSERT INTO storefronts (id, user_id, name, slug, description, is_active, created_at, updated_at)
VALUES
    (1000, 1000, 'Test Inventory Storefront 1', 'test-inventory-store-1', 'Test storefront for inventory tests', true, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Test products for rollback scenarios
INSERT INTO listings (
    id, user_id, storefront_id, title, description,
    price, currency, category_id,
    sku, quantity, status,
    source_type,
    created_at, updated_at
)
VALUES
    -- Product for single rollback (initially decremented)
    (
        8000, -- id
        1000, -- user_id
        1000, -- storefront_id
        'Test Product - Single Rollback',
        'Product for testing single item rollback',
        50.00, -- price
        'USD',  -- currency
        2000,   -- category_id
        'TEST-ROLLBACK-001', -- sku
        90, -- quantity (was 100, decremented by 10)
        'active', -- status
        'b2c', -- source_type
        NOW(),
        NOW()
    ),
    -- Product for batch rollback test
    (
        8001,
        1000,
        1000,
        'Test Product - Batch Rollback 1',
        'First product for batch rollback',
        30.00,
        'USD',
        2000,
        'TEST-ROLLBACK-002',
        80, -- was 100, decremented by 20
        'active',
        'b2c',
        NOW(),
        NOW()
    ),
    (
        8002,
        1000,
        1000,
        'Test Product - Batch Rollback 2',
        'Second product for batch rollback',
        40.00,
        'USD',
        2000,
        'TEST-ROLLBACK-003',
        85, -- was 100, decremented by 15
        'active',
        'b2c',
        NOW(),
        NOW()
    ),
    (
        8003,
        1000,
        1000,
        'Test Product - Batch Rollback 3',
        'Third product for batch rollback',
        35.00,
        'USD',
        2000,
        'TEST-ROLLBACK-004',
        95, -- was 100, decremented by 5
        'active',
        'b2c',
        NOW(),
        NOW()
    ),
    -- Product for idempotency test (double rollback protection)
    (
        8004,
        1000,
        1000,
        'Test Product - Idempotency',
        'Product for testing rollback idempotency',
        60.00,
        'USD',
        2000,
        'TEST-ROLLBACK-005',
        70, -- was 100, decremented by 30
        'active',
        'b2c',
        NOW(),
        NOW()
    ),
    -- Product for partial rollback (rollback less than decremented)
    (
        8005,
        1000,
        1000,
        'Test Product - Partial Rollback',
        'Product for partial quantity rollback',
        45.00,
        'USD',
        2000,
        'TEST-ROLLBACK-006',
        50, -- was 100, decremented by 50
        'active',
        'b2c',
        NOW(),
        NOW()
    ),
    -- Product for concurrent rollback test
    (
        8006,
        1000,
        1000,
        'Test Product - Concurrent Rollback',
        'Product for testing concurrent rollback operations',
        55.00,
        'USD',
        2000,
        'TEST-ROLLBACK-007',
        60, -- was 100, decremented by 40
        'active',
        'b2c',
        NOW(),
        NOW()
    ),
    -- Product that doesn't exist in decrement history (error case)
    (
        8007,
        1000,
        1000,
        'Test Product - No Decrement History',
        'Product without prior decrement (should fail rollback)',
        70.00,
        'USD',
        2000,
        'TEST-ROLLBACK-008',
        100, -- full stock, never decremented
        'active',
        'b2c',
        NOW(),
        NOW()
    ),
    -- Product for E2E workflow test
    (
        8008,
        1000,
        1000,
        'Test Product - E2E Workflow',
        'Product for end-to-end stock workflow',
        80.00,
        'USD',
        2000,
        'TEST-ROLLBACK-009',
        100, -- fresh stock for complete workflow
        'active',
        'b2c',
        NOW(),
        NOW()
    ),
    -- Product for variant rollback test
    (
        8009,
        1000,
        1000,
        'Test Product - Variant Rollback',
        'Product with variants for rollback testing',
        90.00,
        'USD',
        2000,
        'TEST-ROLLBACK-010',
        100, -- product-level stock
        'active',
        'b2c',
        NOW(),
        NOW()
    )
ON CONFLICT (id) DO NOTHING;

-- Simulate inventory movements (decrement history)
-- These represent the original decrements that we will rollback
INSERT INTO inventory_movements (
    id, listing_id, variant_id, movement_type, quantity, reason, notes, user_id, created_at
)
VALUES
    -- Decrement for product 8000 (order ORDER-001)
    (
        10000,
        8000, -- listing_id (было storefront_product_id)
        NULL, -- variant_id
        'out', -- movement_type (было type)
        10,    -- quantity
        'order_created',
        'Order ORDER-001 created',
        1000, -- user_id
        NOW() - INTERVAL '1 hour'
    ),
    -- Decrement for product 8001 (order ORDER-002)
    (
        10001,
        8001,
        NULL,
        'out',
        20,
        'order_created',
        'Order ORDER-002 created',
        1000,
        NOW() - INTERVAL '2 hours'
    ),
    -- Decrement for product 8002 (order ORDER-002)
    (
        10002,
        8002,
        NULL,
        'out',
        15,
        'order_created',
        'Order ORDER-002 created',
        1000,
        NOW() - INTERVAL '2 hours'
    ),
    -- Decrement for product 8003 (order ORDER-002)
    (
        10003,
        8003,
        NULL,
        'out',
        5,
        'order_created',
        'Order ORDER-002 created',
        1000,
        NOW() - INTERVAL '2 hours'
    ),
    -- Decrement for product 8004 (order ORDER-003)
    (
        10004,
        8004,
        NULL,
        'out',
        30,
        'order_created',
        'Order ORDER-003 created',
        1000,
        NOW() - INTERVAL '3 hours'
    ),
    -- Decrement for product 8005 (order ORDER-004)
    (
        10005,
        8005,
        NULL,
        'out',
        50,
        'order_created',
        'Order ORDER-004 created',
        1000,
        NOW() - INTERVAL '4 hours'
    ),
    -- Decrement for product 8006 (order ORDER-005)
    (
        10006,
        8006,
        NULL,
        'out',
        40,
        'order_created',
        'Order ORDER-005 created',
        1000,
        NOW() - INTERVAL '5 hours'
    ),
    -- Decrement for variant 9000 (Size S, order ORDER-006)
    (
        10007,
        8009,
        9000, -- variant_id
        'out',
        10,
        'order_created',
        'Order ORDER-006 created (Size S)',
        1000,
        NOW() - INTERVAL '6 hours'
    ),
    -- Decrement for variant 9001 (Size M, order ORDER-007)
    (
        10008,
        8009,
        9001, -- variant_id
        'out',
        25,
        'order_created',
        'Order ORDER-007 created (Size M)',
        1000,
        NOW() - INTERVAL '7 hours'
    )
ON CONFLICT (id) DO NOTHING;
