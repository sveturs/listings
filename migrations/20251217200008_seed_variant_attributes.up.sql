-- ============================================================================
-- VARIANT ATTRIBUTES SEED (Purpose = 'variant')
-- ============================================================================

-- Clothing Size (для одежды)
INSERT INTO attributes (id, code, attribute_type, purpose, name, display_name, is_required, is_filterable, is_searchable, is_variant_compatible) VALUES
(
    3001,
    'clothing_size',
    'select',
    'variant',
    '{"en": "Size", "ru": "Размер", "sr": "Veličina"}',
    '{"en": "Size", "ru": "Размер", "sr": "Veličina"}',
    true,
    true,
    false,
    true
) ON CONFLICT (id) DO NOTHING;

-- Значения размеров
INSERT INTO attribute_values (id, attribute_id, value, label, sort_order) VALUES
(30011, 3001, 'xs', '{"en": "XS", "ru": "XS", "sr": "XS"}', 1),
(30012, 3001, 's', '{"en": "S", "ru": "S", "sr": "S"}', 2),
(30013, 3001, 'm', '{"en": "M", "ru": "M", "sr": "M"}', 3),
(30014, 3001, 'l', '{"en": "L", "ru": "L", "sr": "L"}', 4),
(30015, 3001, 'xl', '{"en": "XL", "ru": "XL", "sr": "XL"}', 5),
(30042, 3001, '42', '{"en": "42", "ru": "42", "sr": "42"}', 42),
(30043, 3001, '43', '{"en": "43", "ru": "43", "sr": "43"}', 43),
(30044, 3001, '44', '{"en": "44", "ru": "44", "sr": "44"}', 44)
ON CONFLICT (id) DO NOTHING;

-- Color (универсальный)
INSERT INTO attributes (id, code, attribute_type, purpose, name, display_name, is_required, is_filterable, is_searchable, is_variant_compatible, ui_settings) VALUES
(
    3003,
    'color',
    'color',
    'variant',
    '{"en": "Color", "ru": "Цвет", "sr": "Boja"}',
    '{"en": "Color", "ru": "Цвет", "sr": "Boja"}',
    true,
    true,
    false,
    true,
    '{"display_type": "swatch"}'
) ON CONFLICT (id) DO NOTHING;

-- Значения цветов с HEX кодами
INSERT INTO attribute_values (id, attribute_id, value, label, metadata, sort_order) VALUES
(30031, 3003, 'black', '{"en": "Black", "ru": "Черный", "sr": "Crna"}', '{"hex": "#000000"}', 1),
(30032, 3003, 'white', '{"en": "White", "ru": "Белый", "sr": "Bela"}', '{"hex": "#FFFFFF"}', 2),
(30033, 3003, 'red', '{"en": "Red", "ru": "Красный", "sr": "Crvena"}', '{"hex": "#FF0000"}', 3),
(30034, 3003, 'blue', '{"en": "Blue", "ru": "Синий", "sr": "Plava"}', '{"hex": "#0000FF"}', 4)
ON CONFLICT (id) DO NOTHING;

-- Storage Capacity (для электроники)
INSERT INTO attributes (id, code, attribute_type, purpose, name, display_name, is_required, is_filterable, is_variant_compatible) VALUES
(
    3004,
    'storage_capacity',
    'select',
    'variant',
    '{"en": "Storage", "ru": "Память", "sr": "Memorija"}',
    '{"en": "Storage", "ru": "Память", "sr": "Memorija"}',
    true,
    true,
    true
) ON CONFLICT (id) DO NOTHING;

INSERT INTO attribute_values (id, attribute_id, value, label, sort_order) VALUES
(30041, 3004, '64gb', '{"en": "64 GB", "ru": "64 ГБ", "sr": "64 GB"}', 64),
(30042, 3004, '128gb', '{"en": "128 GB", "ru": "128 ГБ", "sr": "128 GB"}', 128),
(30043, 3004, '256gb', '{"en": "256 GB", "ru": "256 ГБ", "sr": "256 GB"}', 256)
ON CONFLICT (id) DO NOTHING;
