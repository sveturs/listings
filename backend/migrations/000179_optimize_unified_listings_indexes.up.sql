-- Migration: 000179_optimize_unified_listings_indexes
-- Description: Создание оптимизированных индексов для unified listings VIEW
-- Created: 2025-10-11

-- ============================================
-- C2C LISTINGS INDEXES
-- ============================================

-- Индекс для фильтрации активных объявлений с сортировкой по дате
CREATE INDEX IF NOT EXISTS idx_c2c_listings_active_created
ON c2c_listings(status, created_at DESC)
WHERE status = 'active';

-- Индекс для фильтрации по категории
CREATE INDEX IF NOT EXISTS idx_c2c_listings_category_active
ON c2c_listings(category_id, status)
WHERE status = 'active';

-- Индекс для фильтрации по цене
CREATE INDEX IF NOT EXISTS idx_c2c_listings_price
ON c2c_listings(price)
WHERE status = 'active' AND price IS NOT NULL;

-- Индекс для геолокации
CREATE INDEX IF NOT EXISTS idx_c2c_listings_location
ON c2c_listings(latitude, longitude)
WHERE status = 'active' AND latitude IS NOT NULL AND longitude IS NOT NULL;

-- Индекс для полнотекстового поиска по title и description
CREATE INDEX IF NOT EXISTS idx_c2c_listings_text_search
ON c2c_listings USING gin(to_tsvector('russian', COALESCE(title, '') || ' ' || COALESCE(description, '')))
WHERE status = 'active';

-- ============================================
-- B2C PRODUCTS INDEXES
-- ============================================

-- Индекс для фильтрации активных товаров с сортировкой по дате
CREATE INDEX IF NOT EXISTS idx_b2c_products_active_created
ON b2c_products(is_active, created_at DESC)
WHERE is_active = true;

-- Индекс для фильтрации по категории
CREATE INDEX IF NOT EXISTS idx_b2c_products_category_active
ON b2c_products(category_id, is_active)
WHERE is_active = true;

-- Индекс для фильтрации по цене
CREATE INDEX IF NOT EXISTS idx_b2c_products_price
ON b2c_products(price)
WHERE is_active = true AND price IS NOT NULL;

-- Индекс для фильтрации по витрине
CREATE INDEX IF NOT EXISTS idx_b2c_products_storefront
ON b2c_products(storefront_id, is_active)
WHERE is_active = true;

-- Индекс для полнотекстового поиска по name и description
CREATE INDEX IF NOT EXISTS idx_b2c_products_text_search
ON b2c_products USING gin(to_tsvector('russian', COALESCE(name, '') || ' ' || COALESCE(description, '')))
WHERE is_active = true;

-- ============================================
-- IMAGES INDEXES
-- ============================================

-- Индекс для быстрого получения изображений C2C listing
CREATE INDEX IF NOT EXISTS idx_c2c_images_listing_main
ON c2c_images(listing_id, is_main DESC);

-- Индекс для быстрого получения изображений B2C product
CREATE INDEX IF NOT EXISTS idx_b2c_images_product_main
ON b2c_product_images(storefront_product_id, is_default DESC, display_order ASC);

-- ============================================
-- STATISTICS UPDATE
-- ============================================

-- Обновить статистику для всех затронутых таблиц
ANALYZE c2c_listings;
ANALYZE c2c_images;
ANALYZE b2c_products;
ANALYZE b2c_product_images;
ANALYZE b2c_stores;

COMMENT ON INDEX idx_c2c_listings_active_created IS
'Оптимизированный индекс для выборки активных C2C listings с сортировкой по дате';

COMMENT ON INDEX idx_b2c_products_active_created IS
'Оптимизированный индекс для выборки активных B2C products с сортировкой по дате';

COMMENT ON INDEX idx_c2c_images_listing_main IS
'Индекс для быстрого получения изображений C2C listing с приоритетом главного изображения';

COMMENT ON INDEX idx_b2c_images_product_main IS
'Индекс для быстрого получения изображений B2C product с приоритетом главного изображения';
