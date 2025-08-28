-- Удаление пустых и неиспользуемых таблиц
DROP TABLE IF EXISTS car_trims CASCADE;
DROP TABLE IF EXISTS car_sync_log CASCADE;
DROP TABLE IF EXISTS carapi_usage CASCADE;
DROP TABLE IF EXISTS vin_decode_cache CASCADE;

-- Удаление неиспользуемой системы атрибутов
DROP TABLE IF EXISTS category_variant_attributes CASCADE;

-- Очистка variant_attribute_mappings (только 2 тестовые записи)
TRUNCATE TABLE variant_attribute_mappings;

-- Добавление недостающих external_id для синхронизации
UPDATE car_generations SET external_id = 'local_' || id::text WHERE external_id IS NULL;
UPDATE car_models SET external_id = 'local_' || id::text WHERE external_id IS NULL OR external_id = '';
UPDATE car_makes SET external_id = 'local_' || id::text WHERE external_id IS NULL OR external_id = '';

-- Комментарии для документации
COMMENT ON TABLE car_makes IS 'Марки автомобилей';
COMMENT ON TABLE car_models IS 'Модели автомобилей';
COMMENT ON TABLE car_generations IS 'Поколения моделей автомобилей';

-- Удаление неиспользуемых индексов и оптимизация существующих
-- Анализ таблиц для обновления статистики
ANALYZE car_makes;
ANALYZE car_models;
ANALYZE car_generations;