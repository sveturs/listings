-- Добавление новых вариативных атрибутов для различных категорий товаров

-- Сначала добавим уникальный индекс на name если его нет
CREATE UNIQUE INDEX IF NOT EXISTS product_variant_attributes_name_unique ON product_variant_attributes(name);

-- 1. Память (для электроники: телефоны, планшеты, компьютеры)
INSERT INTO product_variant_attributes (name, display_name, type, is_required, sort_order, affects_stock)
VALUES ('memory', 'Memory', 'memory', false, 3, true)
ON CONFLICT (name) DO NOTHING;

-- 2. Хранилище (для электроники)
INSERT INTO product_variant_attributes (name, display_name, type, is_required, sort_order, affects_stock)
VALUES ('storage', 'Storage', 'storage', false, 4, true)
ON CONFLICT (name) DO NOTHING;

-- 3. Материал (для одежды, обуви, аксессуаров, мебели)
INSERT INTO product_variant_attributes (name, display_name, type, is_required, sort_order, affects_stock)
VALUES ('material', 'Material', 'material', false, 5, false)
ON CONFLICT (name) DO NOTHING;

-- 4. Объем/Емкость (для кухонной утвари, бытовой техники)
INSERT INTO product_variant_attributes (name, display_name, type, is_required, sort_order, affects_stock)
VALUES ('capacity', 'Capacity', 'capacity', false, 6, true)
ON CONFLICT (name) DO NOTHING;

-- 5. Мощность (для бытовой техники)
INSERT INTO product_variant_attributes (name, display_name, type, is_required, sort_order, affects_stock)
VALUES ('power', 'Power', 'power', false, 7, false)
ON CONFLICT (name) DO NOTHING;

-- 6. Тип соединения (для электроники, аксессуаров)
INSERT INTO product_variant_attributes (name, display_name, type, is_required, sort_order, affects_stock)
VALUES ('connectivity', 'Connectivity', 'connectivity', false, 8, true)
ON CONFLICT (name) DO NOTHING;

-- 7. Стиль/Дизайн (для мебели, одежды)
INSERT INTO product_variant_attributes (name, display_name, type, is_required, sort_order, affects_stock)
VALUES ('style', 'Style', 'style', false, 9, false)
ON CONFLICT (name) DO NOTHING;

-- 8. Паттерн/Узор (для одежды, аксессуаров)
INSERT INTO product_variant_attributes (name, display_name, type, is_required, sort_order, affects_stock)
VALUES ('pattern', 'Pattern', 'pattern', false, 10, false)
ON CONFLICT (name) DO NOTHING;

-- 9. Вес (для спортивного оборудования, продуктов)
INSERT INTO product_variant_attributes (name, display_name, type, is_required, sort_order, affects_stock)
VALUES ('weight', 'Weight', 'weight', false, 11, true)
ON CONFLICT (name) DO NOTHING;

-- 10. Комплектация (для электроники, мебели)
INSERT INTO product_variant_attributes (name, display_name, type, is_required, sort_order, affects_stock)
VALUES ('bundle', 'Bundle', 'bundle', false, 12, true)
ON CONFLICT (name) DO NOTHING;

-- Добавление переводов для новых атрибутов
INSERT INTO translations (entity_type, entity_id, field_name, language, translated_text)
SELECT 
  'product_variant_attribute' as entity_type,
  pva.id as entity_id,
  'display_name' as field_name,
  lang.code as language,
  CASE 
    -- Memory
    WHEN pva.name = 'memory' AND lang.code = 'en' THEN 'Memory'
    WHEN pva.name = 'memory' AND lang.code = 'ru' THEN 'Память'
    WHEN pva.name = 'memory' AND lang.code = 'sr' THEN 'Memorija'
    -- Storage
    WHEN pva.name = 'storage' AND lang.code = 'en' THEN 'Storage'
    WHEN pva.name = 'storage' AND lang.code = 'ru' THEN 'Хранилище'
    WHEN pva.name = 'storage' AND lang.code = 'sr' THEN 'Skladište'
    -- Material
    WHEN pva.name = 'material' AND lang.code = 'en' THEN 'Material'
    WHEN pva.name = 'material' AND lang.code = 'ru' THEN 'Материал'
    WHEN pva.name = 'material' AND lang.code = 'sr' THEN 'Materijal'
    -- Capacity
    WHEN pva.name = 'capacity' AND lang.code = 'en' THEN 'Capacity'
    WHEN pva.name = 'capacity' AND lang.code = 'ru' THEN 'Емкость'
    WHEN pva.name = 'capacity' AND lang.code = 'sr' THEN 'Kapacitet'
    -- Power
    WHEN pva.name = 'power' AND lang.code = 'en' THEN 'Power'
    WHEN pva.name = 'power' AND lang.code = 'ru' THEN 'Мощность'
    WHEN pva.name = 'power' AND lang.code = 'sr' THEN 'Snaga'
    -- Connectivity
    WHEN pva.name = 'connectivity' AND lang.code = 'en' THEN 'Connectivity'
    WHEN pva.name = 'connectivity' AND lang.code = 'ru' THEN 'Подключение'
    WHEN pva.name = 'connectivity' AND lang.code = 'sr' THEN 'Povezivanje'
    -- Style
    WHEN pva.name = 'style' AND lang.code = 'en' THEN 'Style'
    WHEN pva.name = 'style' AND lang.code = 'ru' THEN 'Стиль'
    WHEN pva.name = 'style' AND lang.code = 'sr' THEN 'Stil'
    -- Pattern
    WHEN pva.name = 'pattern' AND lang.code = 'en' THEN 'Pattern'
    WHEN pva.name = 'pattern' AND lang.code = 'ru' THEN 'Узор'
    WHEN pva.name = 'pattern' AND lang.code = 'sr' THEN 'Uzorak'
    -- Weight
    WHEN pva.name = 'weight' AND lang.code = 'en' THEN 'Weight'
    WHEN pva.name = 'weight' AND lang.code = 'ru' THEN 'Вес'
    WHEN pva.name = 'weight' AND lang.code = 'sr' THEN 'Težina'
    -- Bundle
    WHEN pva.name = 'bundle' AND lang.code = 'en' THEN 'Bundle'
    WHEN pva.name = 'bundle' AND lang.code = 'ru' THEN 'Комплектация'
    WHEN pva.name = 'bundle' AND lang.code = 'sr' THEN 'Paket'
  END as translation
FROM product_variant_attributes pva
CROSS JOIN (VALUES ('en'), ('ru'), ('sr')) AS lang(code)
WHERE pva.name IN ('memory', 'storage', 'material', 'capacity', 'power', 'connectivity', 'style', 'pattern', 'weight', 'bundle')
ON CONFLICT (entity_type, entity_id, field_name, language) DO UPDATE
SET translated_text = EXCLUDED.translated_text,
    updated_at = NOW();

-- Связывание вариативных атрибутов с категориями через таблицу категорий-атрибутов
-- Эта функция создает вариативные атрибуты для категорий
CREATE OR REPLACE FUNCTION link_variant_attributes_to_categories()
RETURNS void AS $$
DECLARE
    v_color_id INTEGER;
    v_size_id INTEGER;
    v_memory_id INTEGER;
    v_storage_id INTEGER;
    v_material_id INTEGER;
    v_capacity_id INTEGER;
    v_power_id INTEGER;
    v_connectivity_id INTEGER;
    v_style_id INTEGER;
    v_pattern_id INTEGER;
    v_weight_id INTEGER;
    v_bundle_id INTEGER;
BEGIN
    -- Получаем ID атрибутов
    SELECT id INTO v_color_id FROM product_variant_attributes WHERE name = 'color';
    SELECT id INTO v_size_id FROM product_variant_attributes WHERE name = 'size';
    SELECT id INTO v_memory_id FROM product_variant_attributes WHERE name = 'memory';
    SELECT id INTO v_storage_id FROM product_variant_attributes WHERE name = 'storage';
    SELECT id INTO v_material_id FROM product_variant_attributes WHERE name = 'material';
    SELECT id INTO v_capacity_id FROM product_variant_attributes WHERE name = 'capacity';
    SELECT id INTO v_power_id FROM product_variant_attributes WHERE name = 'power';
    SELECT id INTO v_connectivity_id FROM product_variant_attributes WHERE name = 'connectivity';
    SELECT id INTO v_style_id FROM product_variant_attributes WHERE name = 'style';
    SELECT id INTO v_pattern_id FROM product_variant_attributes WHERE name = 'pattern';
    SELECT id INTO v_weight_id FROM product_variant_attributes WHERE name = 'weight';
    SELECT id INTO v_bundle_id FROM product_variant_attributes WHERE name = 'bundle';

    -- Примечание: вариативные атрибуты хранятся в отдельной системе и не связаны напрямую с category_attributes
    -- Они используются только для товаров витрин через поле variant_attributes в storefront_product_variants
END;
$$ LANGUAGE plpgsql;

-- Выполняем функцию
SELECT link_variant_attributes_to_categories();

-- Удаляем временную функцию
DROP FUNCTION link_variant_attributes_to_categories();

-- Добавляем комментарии к таблице для документирования системы
COMMENT ON TABLE product_variant_attributes IS 'Атрибуты для вариантов товаров витрин. Используется для создания комбинаций товаров (например, цвет + размер)';
COMMENT ON COLUMN product_variant_attributes.affects_stock IS 'Определяет, влияет ли этот атрибут на раздельный учет остатков. Если true, то каждая комбинация значений будет иметь свой остаток';

-- Примеры использования вариативных атрибутов для категорий:
-- Одежда (womens-clothing, mens-clothing, kids-clothing): color, size, material, pattern, style
-- Обувь (shoes): color, size, material, style
-- Электроника (smartphones, computers): color, memory, storage, connectivity
-- Бытовая техника (home-appliances): color, capacity, power
-- Мебель (furniture): color, material, style
-- Кухонная утварь (kitchenware): color, capacity, material