-- Индексы для поиска по атрибутам
CREATE INDEX IF NOT EXISTS idx_attr_name_text_val ON listing_attribute_values (attribute_id, text_value) 
WHERE text_value IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_attr_name_num_val ON listing_attribute_values (attribute_id, numeric_value) 
WHERE numeric_value IS NOT NULL;

-- Индекс для поиска по единицам измерения
CREATE INDEX IF NOT EXISTS idx_attr_unit ON listing_attribute_values (unit)
WHERE unit IS NOT NULL;

-- Составной индекс для категорий атрибутов
CREATE INDEX IF NOT EXISTS idx_category_attr_mapping ON category_attribute_mapping (category_id, is_enabled)
WHERE is_enabled = true;