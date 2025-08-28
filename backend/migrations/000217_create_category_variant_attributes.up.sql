-- Создаем таблицу для связи категорий с вариативными атрибутами
-- Эта таблица определяет, какие вариативные атрибуты доступны для товаров в определенной категории
CREATE TABLE IF NOT EXISTS category_variant_attributes (
    id SERIAL PRIMARY KEY,
    category_id INTEGER NOT NULL REFERENCES marketplace_categories(id) ON DELETE CASCADE,
    variant_attribute_name VARCHAR(100) NOT NULL,
    sort_order INTEGER DEFAULT 0,
    is_required BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Уникальный индекс для предотвращения дублирования
    CONSTRAINT unique_category_variant_attribute UNIQUE (category_id, variant_attribute_name)
);

-- Добавляем индексы для производительности
CREATE INDEX IF NOT EXISTS idx_category_variant_attributes_category_id 
ON category_variant_attributes(category_id);

CREATE INDEX IF NOT EXISTS idx_category_variant_attributes_variant_name 
ON category_variant_attributes(variant_attribute_name);

-- Добавляем комментарии к таблице и полям
COMMENT ON TABLE category_variant_attributes IS 
'Связь между категориями и доступными для них вариативными атрибутами';

COMMENT ON COLUMN category_variant_attributes.category_id IS 
'ID категории из marketplace_categories';

COMMENT ON COLUMN category_variant_attributes.variant_attribute_name IS 
'Имя вариативного атрибута из product_variant_attributes.name';

COMMENT ON COLUMN category_variant_attributes.sort_order IS 
'Порядок отображения вариативных атрибутов для данной категории';

COMMENT ON COLUMN category_variant_attributes.is_required IS 
'Обязателен ли этот вариант для товаров в данной категории';

-- Создаем триггер для автоматического обновления updated_at
CREATE TRIGGER update_category_variant_attributes_updated_at 
BEFORE UPDATE ON category_variant_attributes 
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();