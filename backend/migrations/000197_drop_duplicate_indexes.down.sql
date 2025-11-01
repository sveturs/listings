-- Rollback migration: Restore dropped indexes
-- Note: Recreating indexes can take significant time on large tables

-- 1. Restore duplicate UNIQUE indexes
CREATE UNIQUE INDEX IF NOT EXISTS ai_category_decisions_unique_hash
    ON ai_category_decisions (title_hash, entity_type);

CREATE INDEX IF NOT EXISTS idx_ai_decisions_title_hash
    ON ai_category_decisions (title_hash);

-- 2. Restore duplicate b2c_products indexes
CREATE INDEX IF NOT EXISTS b2c_products_storefront_id_idx
    ON b2c_products (storefront_id);

-- 3. Restore duplicate b2c_product_variants indexes
CREATE INDEX IF NOT EXISTS b2c_product_variants_sku_idx
    ON b2c_product_variants (sku) WHERE (sku IS NOT NULL);

-- 4. Restore car models unused indexes
CREATE INDEX IF NOT EXISTS idx_car_models_drive_type
    ON car_models (drive_type);

CREATE INDEX IF NOT EXISTS idx_car_models_transmission_type
    ON car_models (transmission_type);

CREATE INDEX IF NOT EXISTS idx_car_models_engine_type
    ON car_models (engine_type);

CREATE INDEX IF NOT EXISTS idx_car_models_serbia_popularity
    ON car_models (serbia_popularity);

-- 5. Restore category keywords unused index
CREATE INDEX IF NOT EXISTS idx_category_weight
    ON category_keywords (category_id, weight DESC);

-- Note: translations_pkey cannot be manually restored as it's auto-managed by PostgreSQL
-- It will be recreated automatically if needed
