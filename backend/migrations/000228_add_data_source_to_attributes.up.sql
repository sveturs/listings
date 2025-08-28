-- Добавление поля data_source_config для конфигурации источников данных
ALTER TABLE category_attributes 
ADD COLUMN IF NOT EXISTS data_source_config JSONB;

-- Добавление нового типа источника 'database' в существующий CHECK constraint
ALTER TABLE category_attributes 
DROP CONSTRAINT IF EXISTS check_data_source_values;

ALTER TABLE category_attributes 
ADD CONSTRAINT check_data_source_values 
CHECK (data_source IN ('manual', 'api_external', 'ai_generated', 'imported', 'computed', 'database', 'internal'));

-- Комментарии для документации
COMMENT ON COLUMN category_attributes.data_source IS 'Источник данных для атрибута: internal (встроенный список), database (из таблицы БД), api (внешний API)';
COMMENT ON COLUMN category_attributes.data_source_config IS 'Конфигурация источника данных в JSON формате (table, valueField, labelField, apiEndpoint и т.д.)';

-- Настройка связи car_make_id с таблицей car_makes
UPDATE category_attributes 
SET data_source = 'database',
    data_source_config = jsonb_build_object(
        'table', 'car_makes',
        'valueField', 'id',
        'labelField', 'name',
        'sortField', 'popularity_rs',
        'cacheTime', 3600
    )
WHERE name = 'car_make_id';

-- Настройка связи car_model_id с таблицей car_models
UPDATE category_attributes 
SET data_source = 'database',
    data_source_config = jsonb_build_object(
        'table', 'car_models',
        'valueField', 'id',
        'labelField', 'name',
        'filterField', 'car_make_id',
        'filterBy', 'car_make_id',
        'cacheTime', 3600
    )
WHERE name = 'car_model_id';

-- Другие атрибуты остаются с internal источником по умолчанию