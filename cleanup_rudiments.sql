-- ============================================================================
-- Listings Service Database Cleanup Script
-- ============================================================================
-- Удаление рудиментов после аудита базы данных 2025-12-16
--
-- ВАЖНО: Выполнять ТОЛЬКО после backup БД!
--
-- Экономия: ~1.5-2 MB дискового пространства + 60-80 индексов
-- ============================================================================

BEGIN;

-- ============================================================================
-- ФАЗА 1: УДАЛЕНИЕ ДУБЛИРУЮЩИХСЯ ИНДЕКСОВ
-- ============================================================================
-- Экономия: ~400-500 KB + ускорение INSERT/UPDATE

-- listings (2 индекса)
DROP INDEX IF EXISTS idx_listings_uuid;              -- оставляем listings_uuid_key (UNIQUE)
DROP INDEX IF EXISTS idx_listings_slug_all;          -- оставляем idx_listings_slug (UNIQUE WHERE is_deleted = false)

-- listing_locations (1 индекс)
DROP INDEX IF EXISTS idx_listing_locations_listing_id; -- оставляем listing_locations_listing_id_key (UNIQUE)

-- listing_favorites (2 индекса)
DROP INDEX IF EXISTS idx_listing_favorites_unique;    -- оставляем listing_favorites_pkey
DROP INDEX IF EXISTS listing_favorites_listing_id_idx; -- оставляем idx_listing_favorites_listing_id

-- storefronts (1 индекс)
DROP INDEX IF EXISTS idx_storefronts_slug;           -- оставляем storefronts_slug_key (UNIQUE)

-- categories (1 индекс)
DROP INDEX IF EXISTS idx_categories_slug;            -- оставляем categories_slug_key (UNIQUE)

-- attributes (1 индекс)
DROP INDEX IF EXISTS idx_attributes_code;            -- оставляем attributes_code_key (UNIQUE)

-- attribute_search_cache (1 индекс)
DROP INDEX IF EXISTS idx_attr_search_cache_listing;  -- оставляем attribute_search_cache_listing_id_key (UNIQUE)

-- shopping_carts (2 индекса)
DROP INDEX IF EXISTS idx_shopping_carts_user_storefront;    -- оставляем idx_shopping_carts_unique_user_per_storefront
DROP INDEX IF EXISTS idx_shopping_carts_session_storefront; -- оставляем idx_shopping_carts_unique_session_per_storefront

-- orders (1 индекс)
DROP INDEX IF EXISTS idx_orders_order_number;        -- оставляем orders_order_number_key (UNIQUE)

-- storefront_invitations (1 индекс)
DROP INDEX IF EXISTS idx_storefront_invitations_code; -- оставляем storefront_invitations_invite_code_key (UNIQUE)

-- c2c_chats (2 индекса)
DROP INDEX IF EXISTS c2c_chats_least_greatest_idx;   -- оставляем c2c_chats_least_greatest_idx1 (UNIQUE WHERE)
DROP INDEX IF EXISTS c2c_chats_listing_id_buyer_id_seller_id_idx; -- оставляем c2c_chats_listing_id_buyer_id_seller_id_key (UNIQUE)

-- indexing_queue (1 индекс)
DROP INDEX IF EXISTS idx_indexing_queue_listing_id;  -- оставляем idx_indexing_queue_listing_id_pending (UNIQUE WHERE status='pending')

-- ============================================================================
-- ФАЗА 2: УДАЛЕНИЕ КОЛОНОК-РУДИМЕНТОВ
-- ============================================================================
-- Экономия: ~50-100 KB + упрощение схемы

-- attributes: legacy колонка из старой архитектуры
ALTER TABLE attributes DROP COLUMN IF EXISTS legacy_product_variant_attribute_id;

-- categories: неиспользуемый external_id
ALTER TABLE categories DROP COLUMN IF EXISTS external_id;

-- category_attributes: 5 неиспользуемых колонок
ALTER TABLE category_attributes
    DROP COLUMN IF EXISTS category_specific_options,
    DROP COLUMN IF EXISTS custom_ui_settings,
    DROP COLUMN IF EXISTS custom_validation_rules,
    DROP COLUMN IF EXISTS is_filterable,
    DROP COLUMN IF EXISTS is_searchable;

-- orders: дублирующиеся колонки
ALTER TABLE orders
    DROP COLUMN IF EXISTS notes,              -- дублирует customer_notes
    DROP COLUMN IF EXISTS shipping_method;    -- заменено на shipping_method_id

-- ============================================================================
-- ФАЗА 3: УДАЛЕНИЕ ПУСТЫХ ТАБЛИЦ-РУДИМЕНТОВ
-- ============================================================================
-- ВАЖНО: Эти таблицы полностью пустые и не используются в коде!
-- Экономия: ~700 KB + 60+ индексов

-- Старая аналитика (заменена на listing_stats)
DROP TABLE IF EXISTS analytics_events CASCADE;

-- Устаревшая система атрибутов
DROP TABLE IF EXISTS attribute_options CASCADE;
DROP TABLE IF EXISTS attribute_search_cache CASCADE;

-- Неиспользуемая B2C функциональность
DROP TABLE IF EXISTS b2c_product_variants CASCADE;

-- Дублирующие таблицы
DROP TABLE IF EXISTS category_variant_attributes CASCADE;
DROP TABLE IF EXISTS listing_attribute_values CASCADE;
DROP TABLE IF EXISTS variant_attribute_values CASCADE;

-- ============================================================================
-- ФАЗА 4: ОПТИМИЗАЦИЯ ПОСЛЕ ОЧИСТКИ
-- ============================================================================

-- Обновить статистику планировщика
ANALYZE;

-- Освободить место (опционально, может быть долгим на больших БД)
-- VACUUM FULL ANALYZE;

COMMIT;

-- ============================================================================
-- ПРОВЕРКА РЕЗУЛЬТАТОВ
-- ============================================================================

-- Размер БД до и после
SELECT pg_size_pretty(pg_database_size('listings_dev_db')) AS db_size;

-- Количество индексов
SELECT count(*) AS total_indexes FROM pg_indexes WHERE schemaname = 'public';

-- Количество таблиц
SELECT count(*) AS total_tables FROM pg_tables WHERE schemaname = 'public';

-- ============================================================================
-- ROLLBACK (если что-то пошло не так)
-- ============================================================================
-- ROLLBACK;
