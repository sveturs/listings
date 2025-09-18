-- Изменяем дефолтное значение поля is_active для витрин на true
-- Это позволит новым витринам быть активными сразу после создания

-- Изменяем дефолтное значение колонки is_active на true
ALTER TABLE storefronts
ALTER COLUMN is_active SET DEFAULT true;

-- Комментарий для документации
COMMENT ON COLUMN storefronts.is_active IS 'Активность витрины. По умолчанию true - витрины активны сразу после создания';