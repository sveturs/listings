-- ============================================================================
-- SEED: Test products with variants (Nike Air Max 90)
-- ============================================================================

-- Тестовый товар: Nike Air Max 90
-- Note: listings.id is bigint, product_variants use uuid
WITH new_product AS (
    INSERT INTO listings (
        user_id, title, description, category_id, price, currency,
        status, created_at, updated_at
    ) VALUES (
        1,  -- system user
        'Nike Air Max 90',
        'Classic sneakers with Air Max cushioning',
        (SELECT id FROM categories WHERE slug = 'muska-obuca' LIMIT 1),
        12990,
        'RSD',
        'active',
        NOW(),
        NOW()
    )
    RETURNING id, uuid
)
INSERT INTO product_variants (id, product_id, sku, price, stock_quantity, reserved_quantity, is_default, status)
SELECT
    variant_id::uuid,
    new_product.uuid,
    sku,
    price,
    stock,
    0,
    is_default,
    'active'
FROM new_product
CROSS JOIN (VALUES
    ('11111111-1111-1111-1111-000000000042', 'NAM90-42-BLK', 12990, 5, true),
    ('22222222-2222-2222-2222-000000000042', 'NAM90-42-WHT', 12990, 3, false),
    ('33333333-3333-3333-3333-000000000043', 'NAM90-43-BLK', 12990, 8, false),
    ('44444444-4444-4444-4444-000000000043', 'NAM90-43-WHT', 12990, 0, false),  -- out of stock
    ('55555555-5555-5555-5555-000000000044', 'NAM90-44-BLK', 13490, 2, false)   -- разная цена
) AS v(variant_id, sku, price, stock, is_default);

-- Атрибуты вариантов (Size + Color)
-- attribute_id = 11 (clothing_size), 12 (color)
-- Size 42 (Black/White)
INSERT INTO variant_attribute_values (variant_id, attribute_id, value_text)
VALUES
('11111111-1111-1111-1111-000000000042'::uuid, 11, '42'),
('22222222-2222-2222-2222-000000000042'::uuid, 11, '42')
ON CONFLICT (variant_id, attribute_id) DO NOTHING;

-- Size 43 (Black/White)
INSERT INTO variant_attribute_values (variant_id, attribute_id, value_text)
VALUES
('33333333-3333-3333-3333-000000000043'::uuid, 11, '43'),
('44444444-4444-4444-4444-000000000043'::uuid, 11, '43')
ON CONFLICT (variant_id, attribute_id) DO NOTHING;

-- Size 44 (Black)
INSERT INTO variant_attribute_values (variant_id, attribute_id, value_text)
VALUES
('55555555-5555-5555-5555-000000000044'::uuid, 11, '44')
ON CONFLICT (variant_id, attribute_id) DO NOTHING;

-- Colors
INSERT INTO variant_attribute_values (variant_id, attribute_id, value_text)
VALUES
('11111111-1111-1111-1111-000000000042'::uuid, 12, 'black'),
('22222222-2222-2222-2222-000000000042'::uuid, 12, 'white'),
('33333333-3333-3333-3333-000000000043'::uuid, 12, 'black'),
('44444444-4444-4444-4444-000000000043'::uuid, 12, 'white'),
('55555555-5555-5555-5555-000000000044'::uuid, 12, 'black')
ON CONFLICT (variant_id, attribute_id) DO NOTHING;
