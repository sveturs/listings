-- Добавление поля is_variant_compatible в таблицу unified_attributes
-- Это поле указывает, может ли атрибут использоваться как вариант товара

ALTER TABLE unified_attributes 
ADD COLUMN IF NOT EXISTS is_variant_compatible BOOLEAN DEFAULT FALSE;

-- Комментарии для полей
COMMENT ON COLUMN unified_attributes.is_variant_compatible IS 'Может ли атрибут использоваться для создания вариантов товаров';
COMMENT ON COLUMN unified_attributes.affects_stock IS 'Влияет ли вариант на учет остатков товара';

-- Установка флагов для существующих атрибутов, которые часто используются как варианты
UPDATE unified_attributes 
SET is_variant_compatible = TRUE
WHERE code IN (
    'color',
    'size',
    'material',
    'memory',
    'storage',
    'ram_size',
    'storage_size',
    'screen_size',
    'processor_type',
    'book_language',
    'book_binding_type',
    'pet_age',
    'pet_gender'
);

-- Обновляем атрибуты, которые влияют на остатки
UPDATE unified_attributes 
SET affects_stock = TRUE
WHERE code IN (
    'size',
    'memory',
    'storage',
    'ram_size', 
    'storage_size'
);

-- Создание таблицы для связи вариативных атрибутов с категориями
CREATE TABLE IF NOT EXISTS variant_attribute_mappings (
    id SERIAL PRIMARY KEY,
    variant_attribute_id INTEGER NOT NULL REFERENCES unified_attributes(id) ON DELETE CASCADE,
    category_id INTEGER NOT NULL REFERENCES marketplace_categories(id) ON DELETE CASCADE,
    sort_order INTEGER DEFAULT 0,
    is_required BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(variant_attribute_id, category_id)
);

-- Индекс для быстрого поиска вариативных атрибутов по категории
CREATE INDEX IF NOT EXISTS idx_variant_mappings_category 
ON variant_attribute_mappings(category_id);

-- Индекс для быстрого поиска категорий по атрибуту
CREATE INDEX IF NOT EXISTS idx_variant_mappings_attribute
ON variant_attribute_mappings(variant_attribute_id);

-- Комментарий к таблице
COMMENT ON TABLE variant_attribute_mappings IS 'Связь между вариативными атрибутами и категориями, где они могут использоваться';