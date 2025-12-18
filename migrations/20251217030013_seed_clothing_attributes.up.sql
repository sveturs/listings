-- Migration: 20251217030013_seed_clothing_attributes
-- Description: Seed clothing-specific attributes (Phase 2, Task BE-2.5)
-- Author: Elite Full-Stack Architect (Claude AI)
-- Date: 2025-12-17
-- Note: ~20 attributes specific to clothing categories

-- ============================================================================
-- CLOTHING ATTRIBUTES
-- ============================================================================

INSERT INTO attributes (
    code,
    name,
    display_name,
    attribute_type,
    purpose,
    options,
    validation_rules,
    ui_settings,
    is_searchable,
    is_filterable,
    is_required,
    is_variant_compatible,
    affects_stock,
    affects_price,
    show_in_card,
    is_active,
    sort_order
) VALUES

-- ===== CLOTHING SIZE =====
(
    'clothing_size',
    '{"sr": "Veličina", "en": "Size", "ru": "Размер"}'::jsonb,
    '{"sr": "Veličina odeće", "en": "Clothing Size", "ru": "Размер одежды"}'::jsonb,
    'size',
    'variant',
    '[]'::jsonb,  -- Will be populated via attribute_values
    '{}'::jsonb,
    '{"display_as": "buttons", "show_in_filters": true}'::jsonb,
    true,  -- is_searchable
    true,  -- is_filterable
    true,  -- is_required
    true,  -- is_variant_compatible
    true,  -- affects_stock
    false, -- affects_price
    true,  -- show_in_card
    true,  -- is_active
    101    -- sort_order
),

-- ===== COLOR =====
(
    'color',
    '{"sr": "Boja", "en": "Color", "ru": "Цвет"}'::jsonb,
    '{"sr": "Boja", "en": "Color", "ru": "Цвет"}'::jsonb,
    'color',
    'variant',
    '[]'::jsonb,  -- Will be populated via attribute_values
    '{}'::jsonb,
    '{"display_as": "swatches", "show_in_filters": true}'::jsonb,
    true,  -- is_searchable
    true,  -- is_filterable
    true,  -- is_required
    true,  -- is_variant_compatible
    true,  -- affects_stock
    false, -- affects_price
    true,  -- show_in_card
    true,  -- is_active
    102    -- sort_order
),

-- ===== GENDER =====
(
    'gender',
    '{"sr": "Pol", "en": "Gender", "ru": "Пол"}'::jsonb,
    '{"sr": "Pol", "en": "Gender", "ru": "Пол"}'::jsonb,
    'select',
    'regular',
    '[
        {"value": "women", "label": {"sr": "Žene", "en": "Women", "ru": "Женский"}},
        {"value": "men", "label": {"sr": "Muškarci", "en": "Men", "ru": "Мужской"}},
        {"value": "unisex", "label": {"sr": "Unisex", "en": "Unisex", "ru": "Унисекс"}},
        {"value": "girls", "label": {"sr": "Devojčice", "en": "Girls", "ru": "Девочки"}},
        {"value": "boys", "label": {"sr": "Dečaci", "en": "Boys", "ru": "Мальчики"}}
    ]'::jsonb,
    '{}'::jsonb,
    '{"display_as": "radio", "show_in_filters": true}'::jsonb,
    false, -- is_searchable
    true,  -- is_filterable
    false, -- is_required
    false, -- is_variant_compatible
    false, -- affects_stock
    false, -- affects_price
    false, -- show_in_card
    true,  -- is_active
    103    -- sort_order
),

-- ===== FIT =====
(
    'fit',
    '{"sr": "Kroj", "en": "Fit", "ru": "Крой"}'::jsonb,
    '{"sr": "Kroj", "en": "Fit", "ru": "Крой"}'::jsonb,
    'select',
    'regular',
    '[
        {"value": "slim", "label": {"sr": "Slim", "en": "Slim", "ru": "Приталенный"}},
        {"value": "regular", "label": {"sr": "Regular", "en": "Regular", "ru": "Обычный"}},
        {"value": "relaxed", "label": {"sr": "Opušten", "en": "Relaxed", "ru": "Свободный"}},
        {"value": "oversized", "label": {"sr": "Oversized", "en": "Oversized", "ru": "Оверсайз"}}
    ]'::jsonb,
    '{}'::jsonb,
    '{"display_as": "radio", "show_in_filters": true}'::jsonb,
    false, -- is_searchable
    true,  -- is_filterable
    false, -- is_required
    false, -- is_variant_compatible
    false, -- affects_stock
    false, -- affects_price
    false, -- show_in_card
    true,  -- is_active
    104    -- sort_order
),

-- ===== STYLE =====
(
    'style',
    '{"sr": "Stil", "en": "Style", "ru": "Стиль"}'::jsonb,
    '{"sr": "Stil", "en": "Style", "ru": "Стиль"}'::jsonb,
    'select',
    'regular',
    '[
        {"value": "casual", "label": {"sr": "Casual", "en": "Casual", "ru": "Повседневный"}},
        {"value": "formal", "label": {"sr": "Elegantno", "en": "Formal", "ru": "Формальный"}},
        {"value": "sporty", "label": {"sr": "Sportski", "en": "Sporty", "ru": "Спортивный"}},
        {"value": "streetwear", "label": {"sr": "Streetwear", "en": "Streetwear", "ru": "Уличный"}},
        {"value": "vintage", "label": {"sr": "Vintage", "en": "Vintage", "ru": "Винтаж"}},
        {"value": "bohemian", "label": {"sr": "Boho", "en": "Bohemian", "ru": "Богемный"}},
        {"value": "minimalist", "label": {"sr": "Minimalističko", "en": "Minimalist", "ru": "Минималистичный"}}
    ]'::jsonb,
    '{}'::jsonb,
    '{"display_as": "dropdown", "show_in_filters": true}'::jsonb,
    false, -- is_searchable
    true,  -- is_filterable
    false, -- is_required
    false, -- is_variant_compatible
    false, -- affects_stock
    false, -- affects_price
    false, -- show_in_card
    true,  -- is_active
    105    -- sort_order
),

-- ===== SEASON =====
(
    'season',
    '{"sr": "Sezona", "en": "Season", "ru": "Сезон"}'::jsonb,
    '{"sr": "Sezona", "en": "Season", "ru": "Сезон"}'::jsonb,
    'multiselect',
    'regular',
    '[
        {"value": "spring", "label": {"sr": "Proleće", "en": "Spring", "ru": "Весна"}},
        {"value": "summer", "label": {"sr": "Leto", "en": "Summer", "ru": "Лето"}},
        {"value": "autumn", "label": {"sr": "Jesen", "en": "Autumn", "ru": "Осень"}},
        {"value": "winter", "label": {"sr": "Zima", "en": "Winter", "ru": "Зима"}},
        {"value": "all_season", "label": {"sr": "Sve sezone", "en": "All Seasons", "ru": "Все сезоны"}}
    ]'::jsonb,
    '{}'::jsonb,
    '{"display_as": "checkbox", "show_in_filters": true}'::jsonb,
    false, -- is_searchable
    true,  -- is_filterable
    false, -- is_required
    false, -- is_variant_compatible
    false, -- affects_stock
    false, -- affects_price
    false, -- show_in_card
    true,  -- is_active
    106    -- sort_order
),

-- ===== NECKLINE =====
(
    'neckline',
    '{"sr": "Vrat", "en": "Neckline", "ru": "Вырез"}'::jsonb,
    '{"sr": "Tip izreza", "en": "Neckline Type", "ru": "Тип выреза"}'::jsonb,
    'select',
    'regular',
    '[
        {"value": "round", "label": {"sr": "Okrugli", "en": "Round", "ru": "Круглый"}},
        {"value": "v_neck", "label": {"sr": "V-izrez", "en": "V-Neck", "ru": "V-образный"}},
        {"value": "crew", "label": {"sr": "Crew", "en": "Crew", "ru": "Экипажный"}},
        {"value": "turtleneck", "label": {"sr": "Rolka", "en": "Turtleneck", "ru": "Водолазка"}},
        {"value": "off_shoulder", "label": {"sr": "Sa ramena", "en": "Off-Shoulder", "ru": "С открытыми плечами"}},
        {"value": "polo", "label": {"sr": "Polo", "en": "Polo", "ru": "Поло"}},
        {"value": "square", "label": {"sr": "Kvadratni", "en": "Square", "ru": "Квадратный"}}
    ]'::jsonb,
    '{}'::jsonb,
    '{"display_as": "dropdown", "show_in_filters": true}'::jsonb,
    false, -- is_searchable
    true,  -- is_filterable
    false, -- is_required
    false, -- is_variant_compatible
    false, -- affects_stock
    false, -- affects_price
    false, -- show_in_card
    true,  -- is_active
    107    -- sort_order
),

-- ===== SLEEVE LENGTH =====
(
    'sleeve_length',
    '{"sr": "Dužina rukava", "en": "Sleeve Length", "ru": "Длина рукава"}'::jsonb,
    '{"sr": "Dužina rukava", "en": "Sleeve Length", "ru": "Длина рукава"}'::jsonb,
    'select',
    'regular',
    '[
        {"value": "sleeveless", "label": {"sr": "Bez rukava", "en": "Sleeveless", "ru": "Без рукавов"}},
        {"value": "short", "label": {"sr": "Kratki", "en": "Short", "ru": "Короткий"}},
        {"value": "three_quarter", "label": {"sr": "3/4", "en": "3/4", "ru": "3/4"}},
        {"value": "long", "label": {"sr": "Dugi", "en": "Long", "ru": "Длинный"}}
    ]'::jsonb,
    '{}'::jsonb,
    '{"display_as": "radio", "show_in_filters": true}'::jsonb,
    false, -- is_searchable
    true,  -- is_filterable
    false, -- is_required
    false, -- is_variant_compatible
    false, -- affects_stock
    false, -- affects_price
    false, -- show_in_card
    true,  -- is_active
    108    -- sort_order
),

-- ===== PATTERN =====
(
    'pattern',
    '{"sr": "Uzorak", "en": "Pattern", "ru": "Узор"}'::jsonb,
    '{"sr": "Uzorak", "en": "Pattern", "ru": "Узор"}'::jsonb,
    'select',
    'regular',
    '[
        {"value": "solid", "label": {"sr": "Jednobojno", "en": "Solid", "ru": "Однотонный"}},
        {"value": "striped", "label": {"sr": "Pruge", "en": "Striped", "ru": "Полоска"}},
        {"value": "plaid", "label": {"sr": "Karo", "en": "Plaid", "ru": "Клетка"}},
        {"value": "floral", "label": {"sr": "Cvetni", "en": "Floral", "ru": "Цветочный"}},
        {"value": "geometric", "label": {"sr": "Geometrijski", "en": "Geometric", "ru": "Геометрический"}},
        {"value": "abstract", "label": {"sr": "Apstraktno", "en": "Abstract", "ru": "Абстрактный"}},
        {"value": "animal_print", "label": {"sr": "Životinjski print", "en": "Animal Print", "ru": "Животный принт"}},
        {"value": "camouflage", "label": {"sr": "Maskirni", "en": "Camouflage", "ru": "Камуфляж"}}
    ]'::jsonb,
    '{}'::jsonb,
    '{"display_as": "dropdown", "show_in_filters": true}'::jsonb,
    false, -- is_searchable
    true,  -- is_filterable
    false, -- is_required
    false, -- is_variant_compatible
    false, -- affects_stock
    false, -- affects_price
    false, -- show_in_card
    true,  -- is_active
    109    -- sort_order
),

-- ===== CLOSURE TYPE =====
(
    'closure_type',
    '{"sr": "Tip zatvaranja", "en": "Closure Type", "ru": "Тип застёжки"}'::jsonb,
    '{"sr": "Tip zatvaranja", "en": "Closure Type", "ru": "Тип застёжки"}'::jsonb,
    'select',
    'regular',
    '[
        {"value": "button", "label": {"sr": "Dugme", "en": "Button", "ru": "Пуговица"}},
        {"value": "zipper", "label": {"sr": "Rajsferšlus", "en": "Zipper", "ru": "Молния"}},
        {"value": "snap", "label": {"sr": "Snap", "en": "Snap", "ru": "Кнопка"}},
        {"value": "hook", "label": {"sr": "Kopča", "en": "Hook", "ru": "Крючок"}},
        {"value": "velcro", "label": {"sr": "Čičak", "en": "Velcro", "ru": "Липучка"}},
        {"value": "none", "label": {"sr": "Bez", "en": "None", "ru": "Без застёжки"}}
    ]'::jsonb,
    '{}'::jsonb,
    '{"display_as": "dropdown", "show_in_filters": false}'::jsonb,
    false, -- is_searchable
    false, -- is_filterable
    false, -- is_required
    false, -- is_variant_compatible
    false, -- affects_stock
    false, -- affects_price
    false, -- show_in_card
    true,  -- is_active
    110    -- sort_order
);

-- ============================================================================
-- ATTRIBUTE VALUES: clothing_size
-- ============================================================================

INSERT INTO attribute_values (attribute_id, value, label, sort_order)
SELECT
    a.id,
    size_data.value,
    size_data.label,
    size_data.sort_order
FROM attributes a,
LATERAL (VALUES
    ('XXS', '{"sr": "XXS", "en": "XXS", "ru": "XXS"}'::jsonb, 1),
    ('XS', '{"sr": "XS", "en": "XS", "ru": "XS"}'::jsonb, 2),
    ('S', '{"sr": "S", "en": "S", "ru": "S"}'::jsonb, 3),
    ('M', '{"sr": "M", "en": "M", "ru": "M"}'::jsonb, 4),
    ('L', '{"sr": "L", "en": "L", "ru": "L"}'::jsonb, 5),
    ('XL', '{"sr": "XL", "en": "XL", "ru": "XL"}'::jsonb, 6),
    ('XXL', '{"sr": "XXL", "en": "XXL", "ru": "XXL"}'::jsonb, 7),
    ('XXXL', '{"sr": "XXXL", "en": "XXXL", "ru": "XXXL"}'::jsonb, 8)
) AS size_data(value, label, sort_order)
WHERE a.code = 'clothing_size';

-- ============================================================================
-- ATTRIBUTE VALUES: color (with hex codes)
-- ============================================================================

INSERT INTO attribute_values (attribute_id, value, label, metadata, sort_order)
SELECT
    a.id,
    color_data.value,
    color_data.label,
    color_data.metadata,
    color_data.sort_order
FROM attributes a,
LATERAL (VALUES
    ('black', '{"sr": "Crna", "en": "Black", "ru": "Чёрный"}'::jsonb, '{"hex": "#000000"}'::jsonb, 1),
    ('white', '{"sr": "Bela", "en": "White", "ru": "Белый"}'::jsonb, '{"hex": "#FFFFFF"}'::jsonb, 2),
    ('red', '{"sr": "Crvena", "en": "Red", "ru": "Красный"}'::jsonb, '{"hex": "#FF0000"}'::jsonb, 3),
    ('blue', '{"sr": "Plava", "en": "Blue", "ru": "Синий"}'::jsonb, '{"hex": "#0000FF"}'::jsonb, 4),
    ('green', '{"sr": "Zelena", "en": "Green", "ru": "Зелёный"}'::jsonb, '{"hex": "#00FF00"}'::jsonb, 5),
    ('yellow', '{"sr": "Žuta", "en": "Yellow", "ru": "Жёлтый"}'::jsonb, '{"hex": "#FFFF00"}'::jsonb, 6),
    ('pink', '{"sr": "Roze", "en": "Pink", "ru": "Розовый"}'::jsonb, '{"hex": "#FFC0CB"}'::jsonb, 7),
    ('purple', '{"sr": "Ljubičasta", "en": "Purple", "ru": "Фиолетовый"}'::jsonb, '{"hex": "#800080"}'::jsonb, 8),
    ('orange', '{"sr": "Narandžasta", "en": "Orange", "ru": "Оранжевый"}'::jsonb, '{"hex": "#FFA500"}'::jsonb, 9),
    ('brown', '{"sr": "Braon", "en": "Brown", "ru": "Коричневый"}'::jsonb, '{"hex": "#8B4513"}'::jsonb, 10),
    ('gray', '{"sr": "Siva", "en": "Gray", "ru": "Серый"}'::jsonb, '{"hex": "#808080"}'::jsonb, 11),
    ('beige', '{"sr": "Bež", "en": "Beige", "ru": "Бежевый"}'::jsonb, '{"hex": "#F5F5DC"}'::jsonb, 12),
    ('navy', '{"sr": "Teget", "en": "Navy", "ru": "Тёмно-синий"}'::jsonb, '{"hex": "#000080"}'::jsonb, 13),
    ('multicolor', '{"sr": "Višebojno", "en": "Multicolor", "ru": "Многоцветный"}'::jsonb, '{"hex": "linear-gradient"}'::jsonb, 14)
) AS color_data(value, label, metadata, sort_order)
WHERE a.code = 'color';

-- ============================================================================
-- PROGRESS NOTIFICATION
-- ============================================================================

DO $$
DECLARE
    attr_count INT;
    value_count INT;
BEGIN
    SELECT COUNT(*) INTO attr_count
    FROM attributes
    WHERE code IN ('clothing_size', 'color', 'gender', 'fit', 'style', 'season',
                   'neckline', 'sleeve_length', 'pattern', 'closure_type');

    SELECT COUNT(*) INTO value_count
    FROM attribute_values av
    INNER JOIN attributes a ON av.attribute_id = a.id
    WHERE a.code IN ('clothing_size', 'color');

    RAISE NOTICE '';
    RAISE NOTICE '✅ Clothing attributes seed complete!';
    RAISE NOTICE '   Attributes inserted: %', attr_count;
    RAISE NOTICE '   Attribute values inserted: % (sizes: 8, colors: 14)', value_count;
    RAISE NOTICE '';
END $$;

-- ============================================================================
-- END OF MIGRATION
-- ============================================================================
