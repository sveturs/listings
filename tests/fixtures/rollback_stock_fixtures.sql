-- RollbackStock Integration Test Fixtures
-- This file provides test data specifically for RollbackStock testing
-- Focuses on compensating transaction scenarios and idempotency

-- Test storefronts (reuse storefront 1000 from b2c_inventory_fixtures if exists)
INSERT INTO storefronts (id, user_id, name, slug, description, is_active, created_at, updated_at)
VALUES
    (1000, 1000, 'Test Inventory Storefront 1', 'test-inventory-store-1', 'Test storefront for inventory tests', true, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Test products for rollback scenarios
INSERT INTO b2c_products (
    id, storefront_id, name, description,
    price, currency, category_id,
    sku, barcode, stock_quantity, stock_status,
    is_active, view_count, sold_count,
    created_at, updated_at
)
VALUES
    -- Product for single rollback (initially decremented)
    (
        8000, -- id
        1000, -- storefront_id
        'Test Product - Single Rollback',
        'Product for testing single item rollback',
        50.00, -- price
        'USD',  -- currency
        2000,   -- category_id
        'TEST-ROLLBACK-001', -- sku
        '8000000000001',  -- barcode
        90, -- stock_quantity (was 100, decremented by 10)
        'in_stock', -- stock_status
        true, -- is_active
        0,    -- view_count
        10,   -- sold_count (10 items "sold")
        NOW(),
        NOW()
    ),
    -- Product for batch rollback test
    (
        8001,
        1000,
        'Test Product - Batch Rollback 1',
        'First product for batch rollback',
        30.00,
        'USD',
        2000,
        'TEST-ROLLBACK-002',
        '8000000000002',
        80, -- was 100, decremented by 20
        'in_stock',
        true,
        0,
        20,
        NOW(),
        NOW()
    ),
    (
        8002,
        1000,
        'Test Product - Batch Rollback 2',
        'Second product for batch rollback',
        40.00,
        'USD',
        2000,
        'TEST-ROLLBACK-003',
        '8000000000003',
        85, -- was 100, decremented by 15
        'in_stock',
        true,
        0,
        15,
        NOW(),
        NOW()
    ),
    (
        8003,
        1000,
        'Test Product - Batch Rollback 3',
        'Third product for batch rollback',
        35.00,
        'USD',
        2000,
        'TEST-ROLLBACK-004',
        '8000000000004',
        95, -- was 100, decremented by 5
        'in_stock',
        true,
        0,
        5,
        NOW(),
        NOW()
    ),
    -- Product for idempotency test (double rollback protection)
    (
        8004,
        1000,
        'Test Product - Idempotency',
        'Product for testing rollback idempotency',
        60.00,
        'USD',
        2000,
        'TEST-ROLLBACK-005',
        '8000000000005',
        70, -- was 100, decremented by 30
        'in_stock',
        true,
        0,
        30,
        NOW(),
        NOW()
    ),
    -- Product for partial rollback (rollback less than decremented)
    (
        8005,
        1000,
        'Test Product - Partial Rollback',
        'Product for partial quantity rollback',
        45.00,
        'USD',
        2000,
        'TEST-ROLLBACK-006',
        '8000000000006',
        50, -- was 100, decremented by 50
        'in_stock',
        true,
        0,
        50,
        NOW(),
        NOW()
    ),
    -- Product for concurrent rollback test
    (
        8006,
        1000,
        'Test Product - Concurrent Rollback',
        'Product for testing concurrent rollback operations',
        55.00,
        'USD',
        2000,
        'TEST-ROLLBACK-007',
        '8000000000007',
        60, -- was 100, decremented by 40
        'in_stock',
        true,
        0,
        40,
        NOW(),
        NOW()
    ),
    -- Product that doesn't exist in decrement history (error case)
    (
        8007,
        1000,
        'Test Product - No Decrement History',
        'Product without prior decrement (should fail rollback)',
        70.00,
        'USD',
        2000,
        'TEST-ROLLBACK-008',
        '8000000000008',
        100, -- full stock, never decremented
        'in_stock',
        true,
        0,
        0,
        NOW(),
        NOW()
    ),
    -- Product for E2E workflow test
    (
        8008,
        1000,
        'Test Product - E2E Workflow',
        'Product for end-to-end stock workflow',
        80.00,
        'USD',
        2000,
        'TEST-ROLLBACK-009',
        '8000000000009',
        100, -- fresh stock for complete workflow
        'in_stock',
        true,
        0,
        0,
        NOW(),
        NOW()
    ),
    -- Product for variant rollback test
    (
        8009,
        1000,
        'Test Product - Variant Rollback',
        'Product with variants for rollback testing',
        90.00,
        'USD',
        2000,
        'TEST-ROLLBACK-010',
        '8000000000010',
        100, -- product-level stock
        'in_stock',
        true,
        0,
        0,
        NOW(),
        NOW()
    )
ON CONFLICT (id) DO NOTHING;

-- Test variants for rollback testing
INSERT INTO b2c_product_variants (
    id, product_id, sku, barcode,
    price, stock_quantity, stock_status,
    variant_attributes, is_active, is_default,
    view_count, sold_count,
    created_at, updated_at
)
VALUES
    -- Variant for product 8009 (Size S - decremented)
    (
        9000, -- id
        8009, -- product_id
        'TEST-ROLLBACK-010-S', -- sku
        '8000000000010-S',  -- barcode
        90.00, -- price
        40,    -- stock_quantity (was 50, decremented by 10)
        'in_stock',
        '{"size": "S"}'::jsonb, -- variant_attributes
        true,   -- is_active
        true,   -- is_default
        0,      -- view_count
        10,     -- sold_count
        NOW(),
        NOW()
    ),
    -- Variant for product 8009 (Size M - decremented)
    (
        9001,
        8009,
        'TEST-ROLLBACK-010-M',
        '8000000000010-M',
        90.00,
        25, -- was 50, decremented by 25
        'low_stock',
        '{"size": "M"}'::jsonb,
        true,
        false, -- not default
        0,
        25,
        NOW(),
        NOW()
    ),
    -- Variant for product 8009 (Size L - never decremented)
    (
        9002,
        8009,
        'TEST-ROLLBACK-010-L',
        '8000000000010-L',
        90.00,
        50, -- full stock
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

-- Update has_variants flag
UPDATE b2c_products
SET has_variants = true
WHERE id = 8009;

-- Simulate inventory movements (decrement history)
-- These represent the original decrements that we will rollback
INSERT INTO b2c_inventory_movements (
    id, storefront_product_id, variant_id, type, quantity, reason, notes, user_id, created_at
)
VALUES
    -- Decrement for product 8000 (order ORDER-001)
    (
        10000,
        8000, -- storefront_product_id
        NULL, -- variant_id
        'out', -- type (decrement)
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
