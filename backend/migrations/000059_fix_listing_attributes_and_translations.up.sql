-- Создание таблицы listing_attribute_values для хранения атрибутов объявлений
CREATE TABLE IF NOT EXISTS listing_attribute_values (
    id SERIAL PRIMARY KEY,
    listing_id INTEGER NOT NULL REFERENCES marketplace_listings(id) ON DELETE CASCADE,
    attribute_id INTEGER NOT NULL REFERENCES unified_attributes(id) ON DELETE CASCADE,
    text_value TEXT,
    numeric_value NUMERIC(15, 2),
    boolean_value BOOLEAN,
    date_value DATE,
    json_value JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(listing_id, attribute_id)
);

-- Индексы для производительности
CREATE INDEX IF NOT EXISTS idx_listing_attribute_values_listing ON listing_attribute_values(listing_id);
CREATE INDEX IF NOT EXISTS idx_listing_attribute_values_attribute ON listing_attribute_values(attribute_id);
CREATE INDEX IF NOT EXISTS idx_listing_attribute_values_numeric ON listing_attribute_values(numeric_value) WHERE numeric_value IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_listing_attribute_values_text ON listing_attribute_values(text_value) WHERE text_value IS NOT NULL;

-- Добавление недостающей колонки в таблицу translations
ALTER TABLE translations 
ADD COLUMN IF NOT EXISTS last_modified_by INTEGER REFERENCES users(id);

-- Таблица map_items_cache уже существует, пропускаем создание view
-- Это таблица используется для кеширования данных карты

-- Создаем функцию для обновления кеша если её нет
CREATE OR REPLACE FUNCTION refresh_map_items_cache() RETURNS void AS $$
BEGIN
    -- Функция-заглушка для совместимости
    -- В будущем здесь можно добавить логику обновления кеша
    RETURN;
END;
$$ LANGUAGE plpgsql;

-- Альтернативный VIEW для получения данных карты
CREATE OR REPLACE VIEW map_items_view AS
SELECT 
    ml.id,
    ml.title,
    ml.description,
    ml.price,
    ml.condition,
    ml.location,
    ml.latitude,
    ml.longitude,
    ml.status,
    ml.created_at,
    ml.updated_at,
    ml.user_id,
    ml.category_id,
    mc.name as category_name,
    mc.slug as category_slug,
    (
        SELECT mi.public_url 
        FROM marketplace_images mi 
        WHERE mi.listing_id = ml.id 
          AND mi.is_main = true 
        LIMIT 1
    ) as main_image_url,
    u.name as user_name,
    ml.show_on_map
FROM marketplace_listings ml
LEFT JOIN marketplace_categories mc ON ml.category_id = mc.id
LEFT JOIN users u ON ml.user_id = u.id
WHERE ml.status = 'active' 
  AND ml.show_on_map = true
  AND ml.latitude IS NOT NULL 
  AND ml.longitude IS NOT NULL;