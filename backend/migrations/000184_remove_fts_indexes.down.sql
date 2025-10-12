-- Восстановление Full-Text Search индексов
-- (на случай если потребуется вернуть FTS поиск в PostgreSQL)

-- ВНИМАНИЕ: Эти индексы используют расширения pg_trgm и требуют значительного времени
-- для построения на больших таблицах

-- 1. Car Models - GIN индекс для поиска
CREATE INDEX IF NOT EXISTS idx_car_models_search ON car_models 
USING gin(to_tsvector('english', coalesce(name, '') || ' ' || coalesce(slug, '')));

-- 2. C2C Listings - GIN индекс на description
CREATE INDEX IF NOT EXISTS c2c_listings_description_idx ON c2c_listings 
USING gin(to_tsvector('english', coalesce(description, '')));

-- 3. Unified Attributes - Trigram индекс на name
CREATE INDEX IF NOT EXISTS idx_unified_attributes_name_trgm ON unified_attributes 
USING gin(name gin_trgm_ops);

-- 4. C2C Listings - GIN индекс на title
CREATE INDEX IF NOT EXISTS c2c_listings_title_idx ON c2c_listings 
USING gin(to_tsvector('english', coalesce(title, '')));

-- 5. B2C Products - tsvector индекс
CREATE INDEX IF NOT EXISTS idx_b2c_products_text_search ON b2c_products 
USING gin(to_tsvector('english', coalesce(name, '') || ' ' || coalesce(description, '')));

-- 6. C2C Listings - общий tsvector индекс
CREATE INDEX IF NOT EXISTS idx_c2c_listings_text_search ON c2c_listings 
USING gin(to_tsvector('english', coalesce(title, '') || ' ' || coalesce(description, '')));

-- 7. B2C Products - старый tsvector индекс
CREATE INDEX IF NOT EXISTS b2c_products_to_tsvector_idx ON b2c_products 
USING gin(to_tsvector('english', coalesce(name, '')));

-- 8. Unified Attribute Values - Trigram индекс
CREATE INDEX IF NOT EXISTS idx_unified_attr_values_text_trgm ON unified_attribute_values 
USING gin(value gin_trgm_ops);

-- 9. C2C Listings - старый tsvector индекс
CREATE INDEX IF NOT EXISTS c2c_listings_to_tsvector_idx ON c2c_listings 
USING gin(to_tsvector('english', coalesce(title, '')));
