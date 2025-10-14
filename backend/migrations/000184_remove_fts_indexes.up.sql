-- Удаление Full-Text Search индексов
-- Поиск в проекте осуществляется через OpenSearch, не PostgreSQL FTS
-- Эти индексы замедляют INSERT/UPDATE операции без реальной пользы
-- Освобождает ~852 KB места

-- 1. Car Models - GIN индекс для поиска (312 KB)
DROP INDEX IF EXISTS idx_car_models_search;

-- 2. C2C Listings - GIN индекс на description (224 KB)
DROP INDEX IF EXISTS c2c_listings_description_idx;

-- 3. Unified Attributes - Trigram индекс на name (96 KB)
DROP INDEX IF EXISTS idx_unified_attributes_name_trgm;

-- 4. C2C Listings - GIN индекс на title (88 KB)
DROP INDEX IF EXISTS c2c_listings_title_idx;

-- 5. B2C Products - tsvector индекс (32 KB)
DROP INDEX IF EXISTS idx_b2c_products_text_search;

-- 6. C2C Listings - общий tsvector индекс (24 KB)
DROP INDEX IF EXISTS idx_c2c_listings_text_search;

-- 7. B2C Products - старый tsvector индекс (24 KB)
DROP INDEX IF EXISTS b2c_products_to_tsvector_idx;

-- 8. Unified Attribute Values - Trigram индекс (24 KB)
DROP INDEX IF EXISTS idx_unified_attr_values_text_trgm;

-- 9. C2C Listings - старый tsvector индекс (24 KB)
DROP INDEX IF EXISTS c2c_listings_to_tsvector_idx;

-- Итого освобождено: ~852 KB
