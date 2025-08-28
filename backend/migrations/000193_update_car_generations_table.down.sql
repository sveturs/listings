-- Удаляем функцию
DROP FUNCTION IF EXISTS get_car_generations_with_translations(INTEGER, VARCHAR);

-- Удаляем индексы
DROP INDEX IF EXISTS idx_car_generations_model_id;
DROP INDEX IF EXISTS idx_car_generations_years;
DROP INDEX IF EXISTS idx_car_generations_active;

-- Удаляем добавленные колонки
ALTER TABLE car_generations 
DROP COLUMN IF EXISTS body_types,
DROP COLUMN IF EXISTS engine_types,
DROP COLUMN IF EXISTS image_url,
DROP COLUMN IF EXISTS is_active;