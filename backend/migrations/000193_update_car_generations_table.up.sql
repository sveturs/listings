-- Обновляем существующую таблицу car_generations, созданную в миграции 184

-- Добавляем недостающие поля
ALTER TABLE car_generations 
ADD COLUMN IF NOT EXISTS body_types JSONB DEFAULT '[]',
ADD COLUMN IF NOT EXISTS engine_types JSONB DEFAULT '[]',
ADD COLUMN IF NOT EXISTS image_url VARCHAR(500),
ADD COLUMN IF NOT EXISTS is_active BOOLEAN DEFAULT true;

-- Создаем дополнительные индексы для производительности
CREATE INDEX IF NOT EXISTS idx_car_generations_model_id ON car_generations(model_id);
CREATE INDEX IF NOT EXISTS idx_car_generations_years ON car_generations(year_start, year_end);
CREATE INDEX IF NOT EXISTS idx_car_generations_active ON car_generations(is_active) WHERE is_active = true;

-- Создаем функцию для получения поколений с переводами
CREATE OR REPLACE FUNCTION get_car_generations_with_translations(
    p_model_id INTEGER,
    p_language VARCHAR DEFAULT 'ru'
) RETURNS TABLE (
    id INTEGER,
    name VARCHAR,
    year_from INTEGER,
    year_to INTEGER,
    facelift_year INTEGER,
    body_types JSONB,
    engine_types JSONB,
    specs JSONB,
    image_url VARCHAR
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        g.id,
        COALESCE(
            (SELECT translated_text FROM translations 
             WHERE entity_type = 'car_generation' 
             AND entity_id = g.id 
             AND field_name = 'name' 
             AND language = p_language),
            g.name
        ) as name,
        g.year_start as year_from,
        g.year_end as year_to,
        g.facelift_year,
        g.body_types,
        g.engine_types,
        g.specs,
        g.image_url
    FROM car_generations g
    WHERE g.model_id = p_model_id
    AND g.is_active = true
    ORDER BY g.year_start DESC;
END;
$$ LANGUAGE plpgsql;