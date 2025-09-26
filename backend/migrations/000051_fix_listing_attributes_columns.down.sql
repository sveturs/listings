-- Удаляем индексы
DROP INDEX IF EXISTS idx_listing_attribute_values_value_type;
DROP INDEX IF EXISTS idx_unified_attributes_show_in_card;

-- Пересоздаем view без новой колонки
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
    unified_attributes.affects_stock
FROM unified_attributes
WHERE unified_attributes.is_active = true;

-- Удаляем колонки
ALTER TABLE listing_attribute_values
DROP COLUMN IF EXISTS value_type;

ALTER TABLE unified_attributes
DROP COLUMN IF EXISTS show_in_card;