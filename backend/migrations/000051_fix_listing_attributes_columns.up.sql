-- Добавляем недостающие колонки для таблицы listing_attribute_values
ALTER TABLE listing_attribute_values
ADD COLUMN IF NOT EXISTS value_type VARCHAR(50);

-- Заполняем value_type на основе существующих данных
UPDATE listing_attribute_values
SET value_type =
    CASE
        WHEN text_value IS NOT NULL THEN 'text'
        WHEN numeric_value IS NOT NULL THEN 'number'
        WHEN boolean_value IS NOT NULL THEN 'boolean'
        WHEN date_value IS NOT NULL THEN 'date'
        WHEN json_value IS NOT NULL THEN 'json'
        ELSE 'text'
    END
WHERE value_type IS NULL;

-- Добавляем NOT NULL constraint после заполнения
ALTER TABLE listing_attribute_values
ALTER COLUMN value_type SET NOT NULL;

-- Добавляем колонку show_in_card в таблицу unified_attributes если её нет
ALTER TABLE unified_attributes
ADD COLUMN IF NOT EXISTS show_in_card BOOLEAN DEFAULT false;

-- Пересоздаем view category_attributes с новой колонкой
DROP VIEW IF EXISTS category_attributes;
CREATE VIEW category_attributes AS
SELECT unified_attributes.id,
    unified_attributes.code AS name,
    unified_attributes.display_name,
    unified_attributes.attribute_type,
    COALESCE(unified_attributes.icon, ''::character varying) AS icon,
    unified_attributes.options,
    unified_attributes.validation_rules,
    unified_attributes.is_searchable,
    unified_attributes.is_filterable,
    unified_attributes.is_required,
    unified_attributes.is_required AS is_mandatory,
    unified_attributes.is_active,
    unified_attributes.sort_order,
    unified_attributes.created_at,
    unified_attributes.updated_at,
    ''::text AS custom_component,
    unified_attributes.is_variant_compatible,
    unified_attributes.affects_stock,
    unified_attributes.show_in_card
FROM unified_attributes
WHERE unified_attributes.is_active = true;

-- Устанавливаем show_in_card для важных атрибутов автомобилей
UPDATE unified_attributes
SET show_in_card = true
WHERE code IN ('year', 'mileage', 'fuel_type', 'transmission', 'engine_size', 'body_type');

-- Создаем индекс для ускорения запросов
CREATE INDEX IF NOT EXISTS idx_listing_attribute_values_value_type
ON listing_attribute_values(value_type);

CREATE INDEX IF NOT EXISTS idx_unified_attributes_show_in_card
ON unified_attributes(show_in_card)
WHERE show_in_card = true;