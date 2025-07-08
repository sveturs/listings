-- Добавление недостающих полей в таблицу search_logs

-- Добавляем поля device_type, language, search_type, has_spell_correct, clicked_items, timestamp
ALTER TABLE search_logs 
ADD COLUMN IF NOT EXISTS device_type VARCHAR(50) DEFAULT 'desktop',
ADD COLUMN IF NOT EXISTS language VARCHAR(10) DEFAULT 'ru',
ADD COLUMN IF NOT EXISTS search_type VARCHAR(50) DEFAULT 'listings',
ADD COLUMN IF NOT EXISTS has_spell_correct BOOLEAN DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS clicked_items JSONB DEFAULT '[]'::jsonb,
ADD COLUMN IF NOT EXISTS price_min NUMERIC(10,2),
ADD COLUMN IF NOT EXISTS price_max NUMERIC(10,2);

-- Не переименовываем created_at, оставляем для совместимости с существующим кодом

-- Изменяем тип ip_address на TEXT для совместимости
ALTER TABLE search_logs ALTER COLUMN ip_address TYPE TEXT USING ip_address::TEXT;

-- Добавляем псевдонимы для совместимости с кодом
-- (некоторые поля в коде называются иначе)
-- query_text -> query (в запросах используем правильное имя)
-- results_count -> result_count (в запросах используем правильное имя)

-- Добавляем индексы для новых полей
CREATE INDEX IF NOT EXISTS idx_search_logs_device_type ON search_logs(device_type);
CREATE INDEX IF NOT EXISTS idx_search_logs_language ON search_logs(language);
CREATE INDEX IF NOT EXISTS idx_search_logs_search_type ON search_logs(search_type);
CREATE INDEX IF NOT EXISTS idx_search_logs_price_range ON search_logs(price_min, price_max) WHERE price_min IS NOT NULL OR price_max IS NOT NULL;

-- Добавляем составной индекс для частых запросов
CREATE INDEX IF NOT EXISTS idx_search_logs_composite ON search_logs(created_at DESC, search_type, language);

-- Комментарии для документации
COMMENT ON COLUMN search_logs.device_type IS 'Тип устройства: desktop, mobile, tablet, bot';
COMMENT ON COLUMN search_logs.language IS 'Язык интерфейса пользователя';
COMMENT ON COLUMN search_logs.search_type IS 'Тип поиска: listings, advanced, suggestions, fuzzy';
COMMENT ON COLUMN search_logs.has_spell_correct IS 'Было ли применено исправление орфографии';
COMMENT ON COLUMN search_logs.clicked_items IS 'Массив ID объявлений, на которые кликнул пользователь';
COMMENT ON COLUMN search_logs.price_min IS 'Минимальная цена в фильтре';
COMMENT ON COLUMN search_logs.price_max IS 'Максимальная цена в фильтре';