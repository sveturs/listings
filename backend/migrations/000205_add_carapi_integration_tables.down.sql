-- Откат миграции для интеграции CarAPI

-- Удаляем триггеры
DROP TRIGGER IF EXISTS update_car_trims_updated_at ON car_trims;
DROP FUNCTION IF EXISTS update_car_trims_updated_at();

-- Удаляем новые таблицы
DROP TABLE IF EXISTS carapi_usage;
DROP TABLE IF EXISTS car_sync_log;
DROP TABLE IF EXISTS vin_decode_cache;
DROP TABLE IF EXISTS car_trims;

-- Удаляем добавленные колонки из car_generations
ALTER TABLE car_generations
DROP COLUMN IF EXISTS external_id,
DROP COLUMN IF EXISTS platform,
DROP COLUMN IF EXISTS production_country,
DROP COLUMN IF EXISTS metadata,
DROP COLUMN IF EXISTS last_sync_at;

-- Удаляем добавленные колонки из car_models
ALTER TABLE car_models
DROP COLUMN IF EXISTS external_id,
DROP COLUMN IF EXISTS body_type,
DROP COLUMN IF EXISTS segment,
DROP COLUMN IF EXISTS years_range,
DROP COLUMN IF EXISTS metadata,
DROP COLUMN IF EXISTS last_sync_at;

-- Удаляем добавленные колонки из car_makes
ALTER TABLE car_makes
DROP COLUMN IF EXISTS external_id,
DROP COLUMN IF EXISTS manufacturer_id,
DROP COLUMN IF EXISTS last_sync_at,
DROP COLUMN IF EXISTS metadata;

-- Удаляем индексы
DROP INDEX IF EXISTS idx_car_makes_external_id;
DROP INDEX IF EXISTS idx_car_models_external_id;