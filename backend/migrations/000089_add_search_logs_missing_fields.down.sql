-- Откат изменений в таблице search_logs

-- Удаляем добавленные индексы
DROP INDEX IF EXISTS idx_search_logs_composite;
DROP INDEX IF EXISTS idx_search_logs_price_range;
DROP INDEX IF EXISTS idx_search_logs_search_type;
DROP INDEX IF EXISTS idx_search_logs_language;
DROP INDEX IF EXISTS idx_search_logs_device_type;

-- Имя колонки не меняли, поэтому ничего не делаем

-- Возвращаем тип ip_address обратно на INET
ALTER TABLE search_logs ALTER COLUMN ip_address TYPE INET USING ip_address::INET;

-- Удаляем добавленные колонки
ALTER TABLE search_logs 
DROP COLUMN IF EXISTS device_type,
DROP COLUMN IF EXISTS language,
DROP COLUMN IF EXISTS search_type,
DROP COLUMN IF EXISTS has_spell_correct,
DROP COLUMN IF EXISTS clicked_items,
DROP COLUMN IF EXISTS price_min,
DROP COLUMN IF EXISTS price_max;