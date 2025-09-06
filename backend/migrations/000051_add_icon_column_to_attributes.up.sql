-- Добавление колонки icon в таблицу unified_attributes
ALTER TABLE unified_attributes 
ADD COLUMN IF NOT EXISTS icon VARCHAR(255) DEFAULT '';

-- Комментарий к новой колонке
COMMENT ON COLUMN unified_attributes.icon IS 'Иконка для отображения атрибута в UI';

-- Пересоздаем view category_attributes чтобы включить новую колонку icon
DROP VIEW IF EXISTS category_attributes;

CREATE VIEW category_attributes AS
SELECT 
    id,
    code AS name,
    display_name,
    attribute_type,
    COALESCE(icon, '') AS icon,  -- Добавляем колонку icon
    options,
    validation_rules AS validation_rules,
    is_searchable,
    is_filterable,
    is_required,
    is_required AS is_mandatory,
    is_active,
    sort_order,
    created_at,
    updated_at,
    '' AS custom_component,  -- Добавляем для совместимости
    is_variant_compatible,
    affects_stock
FROM unified_attributes
WHERE is_active = true;