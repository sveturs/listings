-- Добавляем колонку unit в таблицу listing_attribute_values
ALTER TABLE listing_attribute_values
ADD COLUMN IF NOT EXISTS unit VARCHAR(50);

-- Добавляем индекс для unit
CREATE INDEX IF NOT EXISTS idx_listing_attribute_values_unit
ON listing_attribute_values(unit)
WHERE unit IS NOT NULL;

-- Комментарий к колонке
COMMENT ON COLUMN listing_attribute_values.unit IS 'Единица измерения для числовых атрибутов (kg, m, l, etc)';