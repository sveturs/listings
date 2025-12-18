-- Migration: 20251217030012_seed_global_attributes
-- Description: Seed global attributes (Phase 2, Task BE-2.4)
-- Author: Elite Full-Stack Architect (Claude AI)
-- Date: 2025-12-17
-- Note: ~15 global attributes that apply to all/most categories

-- ============================================================================
-- GLOBAL ATTRIBUTES
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

-- ===== BRAND =====
(
    'brand',
    '{"sr": "Brend", "en": "Brand", "ru": "Бренд"}'::jsonb,
    '{"sr": "Brend", "en": "Brand", "ru": "Бренд"}'::jsonb,
    'select',
    'regular',
    '{"allow_custom": true}'::jsonb,
    '{}'::jsonb,
    '{"display_as": "dropdown", "placeholder": {"sr": "Izaberite brend", "en": "Select brand", "ru": "Выберите бренд"}}'::jsonb,
    true,  -- is_searchable
    true,  -- is_filterable
    false, -- is_required
    false, -- is_variant_compatible
    false, -- affects_stock
    false, -- affects_price
    true,  -- show_in_card
    true,  -- is_active
    1      -- sort_order
),

-- ===== CONDITION =====
(
    'condition',
    '{"sr": "Stanje", "en": "Condition", "ru": "Состояние"}'::jsonb,
    '{"sr": "Stanje", "en": "Condition", "ru": "Состояние"}'::jsonb,
    'select',
    'regular',
    '[
        {"value": "new", "label": {"sr": "Novo", "en": "New", "ru": "Новое"}},
        {"value": "like_new", "label": {"sr": "Kao novo", "en": "Like New", "ru": "Как новое"}},
        {"value": "excellent", "label": {"sr": "Odlično", "en": "Excellent", "ru": "Отличное"}},
        {"value": "good", "label": {"sr": "Dobro", "en": "Good", "ru": "Хорошее"}},
        {"value": "fair", "label": {"sr": "Prihvatljivo", "en": "Fair", "ru": "Удовлетворительное"}},
        {"value": "poor", "label": {"sr": "Loše", "en": "Poor", "ru": "Плохое"}}
    ]'::jsonb,
    '{}'::jsonb,
    '{"display_as": "radio"}'::jsonb,
    false, -- is_searchable
    true,  -- is_filterable
    true,  -- is_required
    false, -- is_variant_compatible
    false, -- affects_stock
    false, -- affects_price
    true,  -- show_in_card
    true,  -- is_active
    2      -- sort_order
),

-- ===== COUNTRY OF ORIGIN =====
(
    'country_of_origin',
    '{"sr": "Zemlja porekla", "en": "Country of Origin", "ru": "Страна производства"}'::jsonb,
    '{"sr": "Zemlja porekla", "en": "Country of Origin", "ru": "Страна производства"}'::jsonb,
    'select',
    'regular',
    '[
        {"value": "RS", "label": {"sr": "Srbija", "en": "Serbia", "ru": "Сербия"}},
        {"value": "CN", "label": {"sr": "Kina", "en": "China", "ru": "Китай"}},
        {"value": "DE", "label": {"sr": "Nemačka", "en": "Germany", "ru": "Германия"}},
        {"value": "IT", "label": {"sr": "Italija", "en": "Italy", "ru": "Италия"}},
        {"value": "FR", "label": {"sr": "Francuska", "en": "France", "ru": "Франция"}},
        {"value": "US", "label": {"sr": "SAD", "en": "USA", "ru": "США"}},
        {"value": "GB", "label": {"sr": "Velika Britanija", "en": "UK", "ru": "Великобритания"}},
        {"value": "TR", "label": {"sr": "Turska", "en": "Turkey", "ru": "Турция"}},
        {"value": "JP", "label": {"sr": "Japan", "en": "Japan", "ru": "Япония"}},
        {"value": "KR", "label": {"sr": "Južna Koreja", "en": "South Korea", "ru": "Южная Корея"}}
    ]'::jsonb,
    '{}'::jsonb,
    '{"display_as": "dropdown"}'::jsonb,
    false, -- is_searchable
    true,  -- is_filterable
    false, -- is_required
    false, -- is_variant_compatible
    false, -- affects_stock
    false, -- affects_price
    false, -- show_in_card
    true,  -- is_active
    3      -- sort_order
),

-- ===== MATERIAL =====
(
    'material',
    '{"sr": "Materijal", "en": "Material", "ru": "Материал"}'::jsonb,
    '{"sr": "Materijal", "en": "Material", "ru": "Материал"}'::jsonb,
    'multiselect',
    'both',
    '[
        {"value": "cotton", "label": {"sr": "Pamuk", "en": "Cotton", "ru": "Хлопок"}},
        {"value": "polyester", "label": {"sr": "Poliester", "en": "Polyester", "ru": "Полиэстер"}},
        {"value": "leather", "label": {"sr": "Koža", "en": "Leather", "ru": "Кожа"}},
        {"value": "wool", "label": {"sr": "Vuna", "en": "Wool", "ru": "Шерсть"}},
        {"value": "silk", "label": {"sr": "Svila", "en": "Silk", "ru": "Шёлк"}},
        {"value": "metal", "label": {"sr": "Metal", "en": "Metal", "ru": "Металл"}},
        {"value": "plastic", "label": {"sr": "Plastika", "en": "Plastic", "ru": "Пластик"}},
        {"value": "wood", "label": {"sr": "Drvo", "en": "Wood", "ru": "Дерево"}},
        {"value": "glass", "label": {"sr": "Staklo", "en": "Glass", "ru": "Стекло"}},
        {"value": "rubber", "label": {"sr": "Guma", "en": "Rubber", "ru": "Резина"}}
    ]'::jsonb,
    '{}'::jsonb,
    '{"display_as": "checkbox"}'::jsonb,
    true,  -- is_searchable
    true,  -- is_filterable
    false, -- is_required
    false, -- is_variant_compatible
    false, -- affects_stock
    false, -- affects_price
    false, -- show_in_card
    true,  -- is_active
    4      -- sort_order
),

-- ===== WEIGHT =====
(
    'weight',
    '{"sr": "Težina", "en": "Weight", "ru": "Вес"}'::jsonb,
    '{"sr": "Težina (kg)", "en": "Weight (kg)", "ru": "Вес (кг)"}'::jsonb,
    'number',
    'regular',
    '{}'::jsonb,
    '{"min": 0.01, "max": 10000}'::jsonb,
    '{"unit": "kg", "step": 0.1, "placeholder": {"sr": "Unesite težinu", "en": "Enter weight", "ru": "Введите вес"}}'::jsonb,
    false, -- is_searchable
    false, -- is_filterable
    false, -- is_required
    false, -- is_variant_compatible
    false, -- affects_stock
    false, -- affects_price
    false, -- show_in_card
    true,  -- is_active
    5      -- sort_order
),

-- ===== DIMENSIONS =====
(
    'dimensions',
    '{"sr": "Dimenzije", "en": "Dimensions", "ru": "Размеры"}'::jsonb,
    '{"sr": "Dimenzije (DxŠxV)", "en": "Dimensions (LxWxH)", "ru": "Размеры (ДхШхВ)"}'::jsonb,
    'text',
    'regular',
    '{}'::jsonb,
    '{"pattern": "^\\\\d+(\\\\.\\\\d+)?\\\\s*x\\\\s*\\\\d+(\\\\.\\\\d+)?\\\\s*x\\\\s*\\\\d+(\\\\.\\\\d+)?\\\\s*(cm|mm|m)?$", "message": {"sr": "Format: 100x50x75 cm", "en": "Format: 100x50x75 cm", "ru": "Формат: 100x50x75 см"}}'::jsonb,
    '{"placeholder": {"sr": "npr. 100x50x75 cm", "en": "e.g. 100x50x75 cm", "ru": "напр. 100x50x75 см"}}'::jsonb,
    false, -- is_searchable
    false, -- is_filterable
    false, -- is_required
    false, -- is_variant_compatible
    false, -- affects_stock
    false, -- affects_price
    false, -- show_in_card
    true,  -- is_active
    6      -- sort_order
),

-- ===== WARRANTY =====
(
    'warranty_months',
    '{"sr": "Garancija", "en": "Warranty", "ru": "Гарантия"}'::jsonb,
    '{"sr": "Garancija (meseci)", "en": "Warranty (months)", "ru": "Гарантия (месяцев)"}'::jsonb,
    'number',
    'regular',
    '{}'::jsonb,
    '{"min": 0, "max": 120}'::jsonb,
    '{"unit": "months", "step": 1, "placeholder": {"sr": "Broj meseci garancije", "en": "Warranty months", "ru": "Месяцев гарантии"}}'::jsonb,
    false, -- is_searchable
    true,  -- is_filterable
    false, -- is_required
    false, -- is_variant_compatible
    false, -- affects_stock
    false, -- affects_price
    true,  -- show_in_card
    true,  -- is_active
    7      -- sort_order
),

-- ===== MODEL NUMBER =====
(
    'model_number',
    '{"sr": "Model", "en": "Model Number", "ru": "Номер модели"}'::jsonb,
    '{"sr": "Broj modela", "en": "Model Number", "ru": "Номер модели"}'::jsonb,
    'text',
    'regular',
    '{}'::jsonb,
    '{"pattern": "^[A-Z0-9-]+$", "message": {"sr": "Samo velika slova, brojevi i crtice", "en": "Only uppercase letters, numbers and dashes", "ru": "Только заглавные буквы, цифры и тире"}}'::jsonb,
    '{"placeholder": {"sr": "npr. ABC-123", "en": "e.g. ABC-123", "ru": "напр. ABC-123"}}'::jsonb,
    true,  -- is_searchable
    false, -- is_filterable
    false, -- is_required
    false, -- is_variant_compatible
    false, -- affects_stock
    false, -- affects_price
    false, -- show_in_card
    true,  -- is_active
    8      -- sort_order
),

-- ===== YEAR OF MANUFACTURE =====
(
    'year_of_manufacture',
    '{"sr": "Godina proizvodnje", "en": "Year of Manufacture", "ru": "Год выпуска"}'::jsonb,
    '{"sr": "Godina proizvodnje", "en": "Year of Manufacture", "ru": "Год выпуска"}'::jsonb,
    'number',
    'regular',
    '{}'::jsonb,
    '{"min": 1900, "max": 2026}'::jsonb,
    '{"step": 1, "placeholder": {"sr": "npr. 2023", "en": "e.g. 2023", "ru": "напр. 2023"}}'::jsonb,
    false, -- is_searchable
    true,  -- is_filterable
    false, -- is_required
    false, -- is_variant_compatible
    false, -- affects_stock
    false, -- affects_price
    false, -- show_in_card
    true,  -- is_active
    9      -- sort_order
),

-- ===== ENERGY CLASS =====
(
    'energy_class',
    '{"sr": "Energetska klasa", "en": "Energy Class", "ru": "Класс энергопотребления"}'::jsonb,
    '{"sr": "Energetska klasa", "en": "Energy Class", "ru": "Класс энергопотребления"}'::jsonb,
    'select',
    'regular',
    '[
        {"value": "A+++", "label": {"sr": "A+++", "en": "A+++", "ru": "A+++"}},
        {"value": "A++", "label": {"sr": "A++", "en": "A++", "ru": "A++"}},
        {"value": "A+", "label": {"sr": "A+", "en": "A+", "ru": "A+"}},
        {"value": "A", "label": {"sr": "A", "en": "A", "ru": "A"}},
        {"value": "B", "label": {"sr": "B", "en": "B", "ru": "B"}},
        {"value": "C", "label": {"sr": "C", "en": "C", "ru": "C"}},
        {"value": "D", "label": {"sr": "D", "en": "D", "ru": "D"}},
        {"value": "E", "label": {"sr": "E", "en": "E", "ru": "E"}},
        {"value": "F", "label": {"sr": "F", "en": "F", "ru": "F"}},
        {"value": "G", "label": {"sr": "G", "en": "G", "ru": "G"}}
    ]'::jsonb,
    '{}'::jsonb,
    '{"display_as": "dropdown"}'::jsonb,
    false, -- is_searchable
    true,  -- is_filterable
    false, -- is_required
    false, -- is_variant_compatible
    false, -- affects_stock
    false, -- affects_price
    true,  -- show_in_card
    true,  -- is_active
    10     -- sort_order
);

-- ============================================================================
-- PROGRESS NOTIFICATION
-- ============================================================================

DO $$
DECLARE
    inserted_count INT;
BEGIN
    SELECT COUNT(*) INTO inserted_count FROM attributes WHERE code IN (
        'brand', 'condition', 'country_of_origin', 'material', 'weight',
        'dimensions', 'warranty_months', 'model_number', 'year_of_manufacture', 'energy_class'
    );

    RAISE NOTICE '';
    RAISE NOTICE '✅ Global attributes seed complete!';
    RAISE NOTICE '   Total attributes inserted: %', inserted_count;
    RAISE NOTICE '';
END $$;

-- ============================================================================
-- END OF MIGRATION
-- ============================================================================
