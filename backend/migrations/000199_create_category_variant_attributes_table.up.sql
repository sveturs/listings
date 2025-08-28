-- Создание таблицы связей между категориями и доступными вариативными атрибутами
-- Эта таблица определяет, какие вариативные атрибуты доступны для каждой категории

CREATE TABLE IF NOT EXISTS category_variant_attributes (
    id SERIAL PRIMARY KEY,
    category_id INTEGER NOT NULL REFERENCES marketplace_categories(id) ON DELETE CASCADE ON UPDATE CASCADE,
    variant_attribute_name VARCHAR(100) NOT NULL,
    sort_order INTEGER DEFAULT 0,
    is_required BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Индекс для быстрого поиска по категории
CREATE INDEX idx_category_variant_attributes_category 
ON category_variant_attributes(category_id);

-- Индекс для быстрого поиска по вариативному атрибуту
CREATE INDEX idx_category_variant_attributes_variant_name 
ON category_variant_attributes(variant_attribute_name);

-- Индекс для сортировки
CREATE INDEX idx_category_variant_attributes_sort_order 
ON category_variant_attributes(category_id, sort_order);

-- Уникальный индекс для предотвращения дублирования
CREATE UNIQUE INDEX idx_category_variant_attributes_unique 
ON category_variant_attributes(category_id, variant_attribute_name);

-- Комментарии для документации
COMMENT ON TABLE category_variant_attributes IS 'Таблица связей между категориями и доступными им вариативными атрибутами';
COMMENT ON COLUMN category_variant_attributes.category_id IS 'ID категории';
COMMENT ON COLUMN category_variant_attributes.variant_attribute_name IS 'Системное имя вариативного атрибута (например: color, size, memory)';
COMMENT ON COLUMN category_variant_attributes.sort_order IS 'Порядок отображения вариативного атрибута в категории';
COMMENT ON COLUMN category_variant_attributes.is_required IS 'Обязателен ли данный вариативный атрибут для товаров этой категории';