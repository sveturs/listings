-- Откат миграции 103: Удаление таблицы кэша геокодирования

-- Удаляем функции
DROP FUNCTION IF EXISTS get_geocoding_cache_stats();
DROP FUNCTION IF EXISTS cleanup_expired_geocoding_cache();
DROP FUNCTION IF EXISTS trigger_cleanup_geocoding_cache();
DROP FUNCTION IF EXISTS update_geocoding_cache_updated_at();

-- Удаляем триггеры
DROP TRIGGER IF EXISTS trigger_cleanup_geocoding_cache ON geocoding_cache;
DROP TRIGGER IF EXISTS trigger_geocoding_cache_updated_at ON geocoding_cache;

-- Удаляем индексы
DROP INDEX IF EXISTS idx_geocoding_cache_input_address;
DROP INDEX IF EXISTS idx_geocoding_cache_normalized;
DROP INDEX IF EXISTS idx_geocoding_cache_location;
DROP INDEX IF EXISTS idx_geocoding_cache_expires_at;
DROP INDEX IF EXISTS idx_geocoding_cache_provider;
DROP INDEX IF EXISTS idx_geocoding_cache_confidence;
DROP INDEX IF EXISTS idx_geocoding_cache_cache_hits;
DROP INDEX IF EXISTS idx_geocoding_cache_country_lang;
DROP INDEX IF EXISTS idx_geocoding_cache_address_components;

-- Удаляем таблицу
DROP TABLE IF EXISTS geocoding_cache;