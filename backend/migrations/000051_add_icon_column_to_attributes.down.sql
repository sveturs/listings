-- Пересоздаем view category_attributes без колонки icon
DROP VIEW IF EXISTS category_attributes;

CREATE VIEW category_attributes AS
SELECT 
    id,
    code AS name,
    display_name,
    attribute_type,
    options,
    validation_rules,
    is_searchable,
    is_filterable,
    is_required AS is_mandatory,
    is_active,
    sort_order,
    created_at,
    updated_at
FROM unified_attributes
WHERE is_active = true;

-- Удаление колонки icon из таблицы unified_attributes
ALTER TABLE unified_attributes 
DROP COLUMN IF EXISTS icon;