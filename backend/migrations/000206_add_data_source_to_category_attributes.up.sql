-- Добавление поля data_source в таблицу category_attributes
-- Это поле указывает источник данных для атрибута (manual, api_external, ai_generated и т.д.)

-- Добавляем поле data_source
ALTER TABLE category_attributes 
ADD COLUMN IF NOT EXISTS data_source VARCHAR(50) DEFAULT 'manual';

-- Добавляем комментарий к полю
COMMENT ON COLUMN category_attributes.data_source IS 'Источник данных для атрибута: manual (ручной ввод), api_external (внешний API), ai_generated (сгенерировано AI), imported (импортировано)';

-- Создаем индекс для быстрого поиска по источнику данных
CREATE INDEX IF NOT EXISTS idx_category_attributes_data_source 
ON category_attributes(data_source);

-- Обновляем существующие атрибуты автомобилей с предполагаемым источником
UPDATE category_attributes 
SET data_source = 'api_external'
WHERE name IN ('car_make_id', 'car_model_id', 'car_generation_id')
AND EXISTS (
    SELECT 1 FROM category_attribute_mapping cam
    JOIN marketplace_categories mc ON mc.id = cam.category_id
    WHERE cam.attribute_id = category_attributes.id
    AND (mc.id >= 1301 AND mc.id <= 1302 OR mc.id >= 10100 AND mc.id < 10200)
);

-- Обновляем атрибуты, которые могут быть из AI
UPDATE category_attributes 
SET data_source = 'ai_generated'
WHERE name IN ('suggested_price', 'condition_assessment', 'market_demand')
AND data_source = 'manual';

-- Добавляем проверочное ограничение для допустимых значений
ALTER TABLE category_attributes 
ADD CONSTRAINT check_data_source_values 
CHECK (data_source IN ('manual', 'api_external', 'ai_generated', 'imported', 'computed'));