-- Откат миграции - возвращаем дефолтное значение is_active обратно на false

-- Возвращаем дефолтное значение колонки is_active на false
ALTER TABLE storefronts
ALTER COLUMN is_active SET DEFAULT false;

-- Убираем комментарий
COMMENT ON COLUMN storefronts.is_active IS NULL;