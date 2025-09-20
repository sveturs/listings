-- ================================================================================
-- ДЕФОЛТНЫЕ АТРИБУТЫ ДОСТАВКИ ДЛЯ КАТЕГОРИЙ
-- Миграция 000019: Заполнение дефолтных значений для существующих категорий
-- ================================================================================

-- 1. ЗАПОЛНЕНИЕ ДЕФОЛТНЫХ АТРИБУТОВ ДЛЯ ОСНОВНЫХ КАТЕГОРИЙ
-- ================================================================================

-- Получаем ID категорий и заполняем дефолтные значения
-- Используем ON CONFLICT для безопасного обновления если запись уже существует

-- Электроника и техника
INSERT INTO delivery_category_defaults (
    category_id,
    default_weight_kg,
    default_length_cm,
    default_width_cm,
    default_height_cm,
    default_packaging_type,
    is_typically_fragile
)
SELECT
    id,
    0.5,  -- средний вес электроники
    25,   -- средняя длина
    20,   -- средняя ширина
    10,   -- средняя высота
    'box',
    true  -- обычно хрупкие
FROM marketplace_categories
WHERE slug IN ('electronics', 'elektronika', 'tehnika')
   OR name ILIKE '%electronic%'
   OR name ILIKE '%электрон%'
ON CONFLICT (category_id) DO UPDATE SET
    default_weight_kg = EXCLUDED.default_weight_kg,
    default_length_cm = EXCLUDED.default_length_cm,
    default_width_cm = EXCLUDED.default_width_cm,
    default_height_cm = EXCLUDED.default_height_cm,
    default_packaging_type = EXCLUDED.default_packaging_type,
    is_typically_fragile = EXCLUDED.is_typically_fragile,
    updated_at = NOW();

-- Одежда и обувь
INSERT INTO delivery_category_defaults (
    category_id,
    default_weight_kg,
    default_length_cm,
    default_width_cm,
    default_height_cm,
    default_packaging_type,
    is_typically_fragile
)
SELECT
    id,
    0.3,      -- средний вес одежды
    35,       -- средняя длина
    25,       -- средняя ширина
    5,        -- средняя высота (сложенная одежда)
    'envelope',
    false     -- не хрупкая
FROM marketplace_categories
WHERE slug IN ('clothing', 'odezhda', 'fashion', 'clothes')
   OR name ILIKE '%cloth%'
   OR name ILIKE '%fashion%'
   OR name ILIKE '%одежд%'
ON CONFLICT (category_id) DO UPDATE SET
    default_weight_kg = EXCLUDED.default_weight_kg,
    default_length_cm = EXCLUDED.default_length_cm,
    default_width_cm = EXCLUDED.default_width_cm,
    default_height_cm = EXCLUDED.default_height_cm,
    default_packaging_type = EXCLUDED.default_packaging_type,
    is_typically_fragile = EXCLUDED.is_typically_fragile,
    updated_at = NOW();

-- Мебель
INSERT INTO delivery_category_defaults (
    category_id,
    default_weight_kg,
    default_length_cm,
    default_width_cm,
    default_height_cm,
    default_packaging_type,
    is_typically_fragile
)
SELECT
    id,
    15.0,     -- средний вес мебели
    120,      -- средняя длина
    60,       -- средняя ширина
    80,       -- средняя высота
    'custom',
    false     -- обычно не хрупкая
FROM marketplace_categories
WHERE slug IN ('furniture', 'mebel', 'home', 'dom-i-sad')
   OR name ILIKE '%furniture%'
   OR name ILIKE '%мебел%'
ON CONFLICT (category_id) DO UPDATE SET
    default_weight_kg = EXCLUDED.default_weight_kg,
    default_length_cm = EXCLUDED.default_length_cm,
    default_width_cm = EXCLUDED.default_width_cm,
    default_height_cm = EXCLUDED.default_height_cm,
    default_packaging_type = EXCLUDED.default_packaging_type,
    is_typically_fragile = EXCLUDED.is_typically_fragile,
    updated_at = NOW();

-- Книги и медиа
INSERT INTO delivery_category_defaults (
    category_id,
    default_weight_kg,
    default_length_cm,
    default_width_cm,
    default_height_cm,
    default_packaging_type,
    is_typically_fragile
)
SELECT
    id,
    0.4,      -- средний вес книги
    22,       -- средняя длина
    15,       -- средняя ширина
    3,        -- средняя толщина
    'envelope',
    false     -- не хрупкие
FROM marketplace_categories
WHERE slug IN ('books', 'knigi', 'media', 'books-and-media')
   OR name ILIKE '%book%'
   OR name ILIKE '%книг%'
   OR name ILIKE '%media%'
ON CONFLICT (category_id) DO UPDATE SET
    default_weight_kg = EXCLUDED.default_weight_kg,
    default_length_cm = EXCLUDED.default_length_cm,
    default_width_cm = EXCLUDED.default_width_cm,
    default_height_cm = EXCLUDED.default_height_cm,
    default_packaging_type = EXCLUDED.default_packaging_type,
    is_typically_fragile = EXCLUDED.is_typically_fragile,
    updated_at = NOW();

-- Косметика и парфюмерия
INSERT INTO delivery_category_defaults (
    category_id,
    default_weight_kg,
    default_length_cm,
    default_width_cm,
    default_height_cm,
    default_packaging_type,
    is_typically_fragile
)
SELECT
    id,
    0.2,      -- средний вес косметики
    15,       -- средняя длина
    10,       -- средняя ширина
    8,        -- средняя высота
    'box',
    true      -- часто хрупкая (флаконы)
FROM marketplace_categories
WHERE slug IN ('beauty', 'krasota', 'cosmetics', 'kosmetika')
   OR name ILIKE '%beauty%'
   OR name ILIKE '%cosmetic%'
   OR name ILIKE '%косметик%'
   OR name ILIKE '%красот%'
ON CONFLICT (category_id) DO UPDATE SET
    default_weight_kg = EXCLUDED.default_weight_kg,
    default_length_cm = EXCLUDED.default_length_cm,
    default_width_cm = EXCLUDED.default_width_cm,
    default_height_cm = EXCLUDED.default_height_cm,
    default_packaging_type = EXCLUDED.default_packaging_type,
    is_typically_fragile = EXCLUDED.is_typically_fragile,
    updated_at = NOW();

-- Спорт и отдых
INSERT INTO delivery_category_defaults (
    category_id,
    default_weight_kg,
    default_length_cm,
    default_width_cm,
    default_height_cm,
    default_packaging_type,
    is_typically_fragile
)
SELECT
    id,
    2.0,      -- средний вес спорттоваров
    50,       -- средняя длина
    30,       -- средняя ширина
    20,       -- средняя высота
    'box',
    false     -- обычно не хрупкие
FROM marketplace_categories
WHERE slug IN ('sports', 'sport', 'fitness', 'sport-i-otdykh')
   OR name ILIKE '%sport%'
   OR name ILIKE '%спорт%'
   OR name ILIKE '%fitness%'
   OR name ILIKE '%фитнес%'
ON CONFLICT (category_id) DO UPDATE SET
    default_weight_kg = EXCLUDED.default_weight_kg,
    default_length_cm = EXCLUDED.default_length_cm,
    default_width_cm = EXCLUDED.default_width_cm,
    default_height_cm = EXCLUDED.default_height_cm,
    default_packaging_type = EXCLUDED.default_packaging_type,
    is_typically_fragile = EXCLUDED.is_typically_fragile,
    updated_at = NOW();

-- Автомобили и запчасти
INSERT INTO delivery_category_defaults (
    category_id,
    default_weight_kg,
    default_length_cm,
    default_width_cm,
    default_height_cm,
    default_packaging_type,
    is_typically_fragile
)
SELECT
    id,
    5.0,      -- средний вес запчастей
    40,       -- средняя длина
    30,       -- средняя ширина
    25,       -- средняя высота
    'box',
    false     -- обычно не хрупкие
FROM marketplace_categories
WHERE slug IN ('automotive', 'avtomobili', 'cars', 'auto')
   OR name ILIKE '%auto%'
   OR name ILIKE '%car%'
   OR name ILIKE '%авто%'
   OR name ILIKE '%машин%'
ON CONFLICT (category_id) DO UPDATE SET
    default_weight_kg = EXCLUDED.default_weight_kg,
    default_length_cm = EXCLUDED.default_length_cm,
    default_width_cm = EXCLUDED.default_width_cm,
    default_height_cm = EXCLUDED.default_height_cm,
    default_packaging_type = EXCLUDED.default_packaging_type,
    is_typically_fragile = EXCLUDED.is_typically_fragile,
    updated_at = NOW();

-- Детские товары
INSERT INTO delivery_category_defaults (
    category_id,
    default_weight_kg,
    default_length_cm,
    default_width_cm,
    default_height_cm,
    default_packaging_type,
    is_typically_fragile
)
SELECT
    id,
    1.0,      -- средний вес детских товаров
    30,       -- средняя длина
    25,       -- средняя ширина
    15,       -- средняя высота
    'box',
    false     -- обычно не хрупкие
FROM marketplace_categories
WHERE slug IN ('kids', 'deti', 'toys', 'igrushki', 'detskie-tovary')
   OR name ILIKE '%kid%'
   OR name ILIKE '%toy%'
   OR name ILIKE '%дет%'
   OR name ILIKE '%игрушк%'
ON CONFLICT (category_id) DO UPDATE SET
    default_weight_kg = EXCLUDED.default_weight_kg,
    default_length_cm = EXCLUDED.default_length_cm,
    default_width_cm = EXCLUDED.default_width_cm,
    default_height_cm = EXCLUDED.default_height_cm,
    default_packaging_type = EXCLUDED.default_packaging_type,
    is_typically_fragile = EXCLUDED.is_typically_fragile,
    updated_at = NOW();

-- Продукты питания
INSERT INTO delivery_category_defaults (
    category_id,
    default_weight_kg,
    default_length_cm,
    default_width_cm,
    default_height_cm,
    default_packaging_type,
    is_typically_fragile
)
SELECT
    id,
    1.0,      -- средний вес продуктов
    25,       -- средняя длина
    20,       -- средняя ширина
    15,       -- средняя высота
    'box',
    false     -- обычно не хрупкие (но может требоваться специальная обработка)
FROM marketplace_categories
WHERE slug IN ('food', 'produkty', 'groceries', 'eda')
   OR name ILIKE '%food%'
   OR name ILIKE '%grocer%'
   OR name ILIKE '%продукт%'
   OR name ILIKE '%еда%'
ON CONFLICT (category_id) DO UPDATE SET
    default_weight_kg = EXCLUDED.default_weight_kg,
    default_length_cm = EXCLUDED.default_length_cm,
    default_width_cm = EXCLUDED.default_width_cm,
    default_height_cm = EXCLUDED.default_height_cm,
    default_packaging_type = EXCLUDED.default_packaging_type,
    is_typically_fragile = EXCLUDED.is_typically_fragile,
    updated_at = NOW();

-- Недвижимость (для документов и ключей)
INSERT INTO delivery_category_defaults (
    category_id,
    default_weight_kg,
    default_length_cm,
    default_width_cm,
    default_height_cm,
    default_packaging_type,
    is_typically_fragile
)
SELECT
    id,
    0.1,      -- вес документов/ключей
    30,       -- A4 формат
    21,       -- A4 формат
    2,        -- толщина папки
    'envelope',
    false
FROM marketplace_categories
WHERE slug IN ('realestate', 'nedvizhimost', 'property', 'real-estate')
   OR name ILIKE '%realestate%'
   OR name ILIKE '%property%'
   OR name ILIKE '%недвижимост%'
ON CONFLICT (category_id) DO UPDATE SET
    default_weight_kg = EXCLUDED.default_weight_kg,
    default_length_cm = EXCLUDED.default_length_cm,
    default_width_cm = EXCLUDED.default_width_cm,
    default_height_cm = EXCLUDED.default_height_cm,
    default_packaging_type = EXCLUDED.default_packaging_type,
    is_typically_fragile = EXCLUDED.is_typically_fragile,
    updated_at = NOW();

-- 2. УНИВЕРСАЛЬНЫЕ ДЕФОЛТЫ ДЛЯ ВСЕХ ОСТАЛЬНЫХ КАТЕГОРИЙ
-- ================================================================================

-- Добавляем универсальные дефолты для категорий без специфических значений
INSERT INTO delivery_category_defaults (
    category_id,
    default_weight_kg,
    default_length_cm,
    default_width_cm,
    default_height_cm,
    default_packaging_type,
    is_typically_fragile
)
SELECT
    c.id,
    1.0,      -- универсальный средний вес
    30,       -- универсальная средняя длина
    20,       -- универсальная средняя ширина
    10,       -- универсальная средняя высота
    'box',
    false     -- по умолчанию не хрупкие
FROM marketplace_categories c
LEFT JOIN delivery_category_defaults dcd ON c.id = dcd.category_id
WHERE dcd.id IS NULL  -- только для категорий без дефолтов
ON CONFLICT (category_id) DO NOTHING;

-- 3. СТАТИСТИКА И ОТЧЕТ
-- ================================================================================

-- Выводим статистику по заполненным дефолтам
DO $$
DECLARE
    total_categories INTEGER;
    filled_categories INTEGER;
    coverage_percent NUMERIC;
BEGIN
    SELECT COUNT(*) INTO total_categories FROM marketplace_categories;
    SELECT COUNT(*) INTO filled_categories FROM delivery_category_defaults;

    IF total_categories > 0 THEN
        coverage_percent := (filled_categories::NUMERIC / total_categories::NUMERIC) * 100;
        RAISE NOTICE 'Дефолтные атрибуты доставки заполнены для % из % категорий (покрытие: %%%)',
            filled_categories, total_categories, ROUND(coverage_percent, 2);
    END IF;
END $$;

-- ================================================================================
-- КОНЕЦ МИГРАЦИИ
-- ================================================================================