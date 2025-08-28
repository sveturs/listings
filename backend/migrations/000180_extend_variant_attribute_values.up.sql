-- Расширить таблицу product_variant_attribute_values для лучшей поддержки различных типов атрибутов
ALTER TABLE product_variant_attribute_values
ADD COLUMN IF NOT EXISTS metadata JSONB DEFAULT '{}',
ADD COLUMN IF NOT EXISTS is_popular BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN IF NOT EXISTS usage_count INTEGER NOT NULL DEFAULT 0;

-- Добавить индекс для популярных значений
CREATE INDEX idx_product_variant_attribute_values_popular
ON product_variant_attribute_values(attribute_id, is_popular, sort_order)
WHERE is_active = true;

-- Добавить индекс для поиска по metadata
CREATE INDEX idx_product_variant_attribute_values_metadata_gin
ON product_variant_attribute_values USING GIN(metadata);

-- Добавить несколько стандартных размеров одежды
INSERT INTO product_variant_attribute_values (attribute_id, value, display_name, sort_order, is_active, is_popular)
SELECT 
    pva.id,
    size.value,
    size.display,
    size.order_num,
    true,
    size.popular
FROM product_variant_attributes pva
CROSS JOIN (VALUES 
    ('xs', 'XS', 1, true),
    ('s', 'S', 2, true),
    ('m', 'M', 3, true),
    ('l', 'L', 4, true),
    ('xl', 'XL', 5, true),
    ('xxl', 'XXL', 6, true),
    ('xxxl', 'XXXL', 7, false),
    ('36', '36', 10, false),
    ('38', '38', 11, false),
    ('40', '40', 12, false),
    ('42', '42', 13, false),
    ('44', '44', 14, false),
    ('46', '46', 15, false),
    ('48', '48', 16, false),
    ('50', '50', 17, false),
    ('52', '52', 18, false)
) AS size(value, display, order_num, popular)
WHERE pva.name = 'size'
AND NOT EXISTS (
    SELECT 1 FROM product_variant_attribute_values 
    WHERE attribute_id = pva.id AND value = size.value
);

-- Добавить стандартные цвета
INSERT INTO product_variant_attribute_values (attribute_id, value, display_name, color_hex, sort_order, is_active, is_popular, metadata)
SELECT 
    pva.id,
    color.value,
    color.display,
    color.hex,
    color.order_num,
    true,
    color.popular,
    jsonb_build_object('group', color.color_group)
FROM product_variant_attributes pva
CROSS JOIN (VALUES 
    ('black', 'Black', '#000000', 1, true, 'dark'),
    ('white', 'White', '#FFFFFF', 2, true, 'light'),
    ('gray', 'Gray', '#808080', 3, true, 'neutral'),
    ('red', 'Red', '#FF0000', 4, true, 'warm'),
    ('blue', 'Blue', '#0000FF', 5, true, 'cool'),
    ('green', 'Green', '#00FF00', 6, true, 'cool'),
    ('yellow', 'Yellow', '#FFFF00', 7, false, 'warm'),
    ('orange', 'Orange', '#FFA500', 8, false, 'warm'),
    ('purple', 'Purple', '#800080', 9, false, 'cool'),
    ('pink', 'Pink', '#FFC0CB', 10, false, 'warm'),
    ('brown', 'Brown', '#964B00', 11, true, 'warm'),
    ('navy', 'Navy', '#000080', 12, true, 'cool'),
    ('beige', 'Beige', '#F5F5DC', 13, true, 'neutral'),
    ('gold', 'Gold', '#FFD700', 14, false, 'warm'),
    ('silver', 'Silver', '#C0C0C0', 15, false, 'neutral')
) AS color(value, display, hex, order_num, popular, color_group)
WHERE pva.name = 'color'
AND NOT EXISTS (
    SELECT 1 FROM product_variant_attribute_values 
    WHERE attribute_id = pva.id AND value = color.value
);

-- Добавить переводы для значений атрибутов
INSERT INTO translations (entity_type, entity_id, field_name, language, translated_text)
SELECT 
    'product_variant_attribute_value' as entity_type,
    pvav.id as entity_id,
    'display_name' as field_name,
    lang.code as language,
    CASE 
        -- Размеры остаются без перевода (международные)
        WHEN pva.name = 'size' THEN pvav.display_name
        -- Цвета
        WHEN pva.name = 'color' AND pvav.value = 'black' AND lang.code = 'ru' THEN 'Черный'
        WHEN pva.name = 'color' AND pvav.value = 'black' AND lang.code = 'sr' THEN 'Crna'
        WHEN pva.name = 'color' AND pvav.value = 'white' AND lang.code = 'ru' THEN 'Белый'
        WHEN pva.name = 'color' AND pvav.value = 'white' AND lang.code = 'sr' THEN 'Bela'
        WHEN pva.name = 'color' AND pvav.value = 'gray' AND lang.code = 'ru' THEN 'Серый'
        WHEN pva.name = 'color' AND pvav.value = 'gray' AND lang.code = 'sr' THEN 'Siva'
        WHEN pva.name = 'color' AND pvav.value = 'red' AND lang.code = 'ru' THEN 'Красный'
        WHEN pva.name = 'color' AND pvav.value = 'red' AND lang.code = 'sr' THEN 'Crvena'
        WHEN pva.name = 'color' AND pvav.value = 'blue' AND lang.code = 'ru' THEN 'Синий'
        WHEN pva.name = 'color' AND pvav.value = 'blue' AND lang.code = 'sr' THEN 'Plava'
        WHEN pva.name = 'color' AND pvav.value = 'green' AND lang.code = 'ru' THEN 'Зеленый'
        WHEN pva.name = 'color' AND pvav.value = 'green' AND lang.code = 'sr' THEN 'Zelena'
        WHEN pva.name = 'color' AND pvav.value = 'yellow' AND lang.code = 'ru' THEN 'Желтый'
        WHEN pva.name = 'color' AND pvav.value = 'yellow' AND lang.code = 'sr' THEN 'Žuta'
        WHEN pva.name = 'color' AND pvav.value = 'brown' AND lang.code = 'ru' THEN 'Коричневый'
        WHEN pva.name = 'color' AND pvav.value = 'brown' AND lang.code = 'sr' THEN 'Braon'
        ELSE pvav.display_name
    END as translation
FROM product_variant_attribute_values pvav
JOIN product_variant_attributes pva ON pva.id = pvav.attribute_id
CROSS JOIN (VALUES ('en'), ('ru'), ('sr')) AS lang(code)
WHERE pva.name IN ('size', 'color')
ON CONFLICT (entity_type, entity_id, field_name, language) DO UPDATE
SET translated_text = EXCLUDED.translated_text,
    updated_at = NOW();