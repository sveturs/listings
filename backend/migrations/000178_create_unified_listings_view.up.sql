-- Migration: 000178_create_unified_listings_view
-- Description: Создает VIEW для unified списка C2C + B2C объявлений БЕЗ дублирования данных
-- New approach: UNION query вместо костыля с дублированием в БД
-- Date: 2025-10-11

BEGIN;

-- Создать VIEW для unified listings
CREATE OR REPLACE VIEW unified_listings AS
-- C2C listings
SELECT
    l.id,
    'c2c' AS source_type,
    l.user_id,
    l.category_id,
    l.title,
    l.description,
    l.price,
    l.condition,
    l.status,
    l.location,
    l.latitude,
    l.longitude,
    l.address_city,
    l.address_country,
    l.views_count,
    l.show_on_map,
    l.original_language,
    l.created_at,
    l.updated_at,
    NULL::integer AS storefront_id,
    l.external_id,
    l.metadata,
    l.needs_reindex,
    l.address_multilingual,
    -- Агрегация изображений из c2c_images
    COALESCE(
        (
            SELECT jsonb_agg(
                jsonb_build_object(
                    'id', img.id,
                    'file_path', img.file_path,
                    'file_name', img.file_name,
                    'public_url', img.public_url,
                    'is_main', img.is_main,
                    'storage_type', img.storage_type
                )
                ORDER BY img.is_main DESC, img.id
            )
            FROM c2c_images img
            WHERE img.listing_id = l.id
        ),
        '[]'::jsonb
    ) AS images
FROM c2c_listings l

UNION ALL

-- B2C products
SELECT
    p.id,
    'b2c' AS source_type,
    s.user_id,
    p.category_id,
    p.name AS title,
    p.description,
    p.price,
    'new' AS condition,
    CASE
        WHEN p.is_active THEN 'active'::character varying
        ELSE 'inactive'::character varying
    END AS status,
    COALESCE(p.individual_address, s.address) AS location,
    COALESCE(p.individual_latitude, s.latitude) AS latitude,
    COALESCE(p.individual_longitude, s.longitude) AS longitude,
    s.city AS address_city,
    s.country AS address_country,
    p.view_count AS views_count,
    COALESCE(p.show_on_map, true) AS show_on_map,
    'sr' AS original_language,
    p.created_at::timestamp without time zone AS created_at,
    p.updated_at::timestamp without time zone AS updated_at,
    p.storefront_id,
    p.sku AS external_id,
    jsonb_build_object(
        'source', 'storefront',
        'storefront_id', p.storefront_id,
        'stock_quantity', p.stock_quantity,
        'stock_status', p.stock_status,
        'currency', p.currency,
        'barcode', p.barcode,
        'attributes', p.attributes,
        'has_variants', p.has_variants
    ) AS metadata,
    false AS needs_reindex,
    NULL::jsonb AS address_multilingual,
    -- Агрегация изображений из b2c_product_images
    COALESCE(
        (
            SELECT jsonb_agg(
                jsonb_build_object(
                    'id', img.id,
                    'image_url', img.image_url,
                    'thumbnail_url', img.thumbnail_url,
                    'is_default', img.is_default,
                    'display_order', img.display_order
                )
                ORDER BY img.is_default DESC, img.display_order, img.id
            )
            FROM b2c_product_images img
            WHERE img.storefront_product_id = p.id
        ),
        '[]'::jsonb
    ) AS images
FROM b2c_products p
JOIN b2c_stores s ON s.id = p.storefront_id;

-- Комментарий к VIEW
COMMENT ON VIEW unified_listings IS
'Unified VIEW для отображения C2C и B2C объявлений вместе БЕЗ дублирования данных.
- source_type: "c2c" или "b2c"
- images: JSONB массив изображений из соответствующей таблицы
- Каждый товар находится ТОЛЬКО в одной таблице (источник истины)
- Используйте WHERE source_type = ''c2c'' для фильтрации только C2C
- Используйте WHERE source_type = ''b2c'' для фильтрации только B2C';

COMMIT;
