-- Добавляем дополнительные технические атрибуты для автомобилей

-- 1. Привод (если не существует)
INSERT INTO category_attributes (
    name,
    display_name,
    attribute_type,
    options,
    is_searchable,
    is_filterable,
    is_required,
    sort_order,
    show_in_card,
    show_in_list
) VALUES (
    'drivetrain',
    'Pogon',
    'select',
    '{"values": ["fwd", "rwd", "awd", "4wd"]}',
    true,
    true,
    false,
    31,
    true,
    false
)
ON CONFLICT (name) DO NOTHING;

-- Добавляем все технические атрибуты к категории Автомобили (ID: 1003)
INSERT INTO category_attribute_mapping (category_id, attribute_id)
SELECT 1003, id FROM category_attributes WHERE name IN ('body_type', 'drivetrain', 'doors', 'seats', 'power_hp')
ON CONFLICT (category_id, attribute_id) DO NOTHING;