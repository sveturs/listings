-- Добавить поле affects_stock к атрибутам вариантов
ALTER TABLE product_variant_attributes 
ADD COLUMN affects_stock BOOLEAN NOT NULL DEFAULT false;

-- Добавить комментарий к колонке
COMMENT ON COLUMN product_variant_attributes.affects_stock IS 'Определяет, влияет ли этот атрибут на раздельный учет остатков';

-- Обновить существующие атрибуты, которые влияют на остатки
UPDATE product_variant_attributes 
SET affects_stock = true 
WHERE name IN ('size', 'color', 'memory', 'storage', 'capacity', 'volume', 'weight');

-- Добавить переводы для атрибутов вариантов
INSERT INTO translations (entity_type, entity_id, field_name, language, translated_text)
SELECT 
  'product_variant_attribute' as entity_type,
  id as entity_id,
  'display_name' as field_name,
  lang.code as language,
  CASE 
    -- Size
    WHEN name = 'size' AND lang.code = 'en' THEN 'Size'
    WHEN name = 'size' AND lang.code = 'ru' THEN 'Размер'
    WHEN name = 'size' AND lang.code = 'sr' THEN 'Veličina'
    -- Color
    WHEN name = 'color' AND lang.code = 'en' THEN 'Color'
    WHEN name = 'color' AND lang.code = 'ru' THEN 'Цвет'
    WHEN name = 'color' AND lang.code = 'sr' THEN 'Boja'
    -- Memory
    WHEN name = 'memory' AND lang.code = 'en' THEN 'Memory'
    WHEN name = 'memory' AND lang.code = 'ru' THEN 'Память'
    WHEN name = 'memory' AND lang.code = 'sr' THEN 'Memorija'
    -- Storage
    WHEN name = 'storage' AND lang.code = 'en' THEN 'Storage'
    WHEN name = 'storage' AND lang.code = 'ru' THEN 'Хранилище'
    WHEN name = 'storage' AND lang.code = 'sr' THEN 'Skladište'
    -- Capacity
    WHEN name = 'capacity' AND lang.code = 'en' THEN 'Capacity'
    WHEN name = 'capacity' AND lang.code = 'ru' THEN 'Емкость'
    WHEN name = 'capacity' AND lang.code = 'sr' THEN 'Kapacitet'
    -- Volume
    WHEN name = 'volume' AND lang.code = 'en' THEN 'Volume'
    WHEN name = 'volume' AND lang.code = 'ru' THEN 'Объем'
    WHEN name = 'volume' AND lang.code = 'sr' THEN 'Volumen'
    -- Weight
    WHEN name = 'weight' AND lang.code = 'en' THEN 'Weight'
    WHEN name = 'weight' AND lang.code = 'ru' THEN 'Вес'
    WHEN name = 'weight' AND lang.code = 'sr' THEN 'Težina'
    -- Material
    WHEN name = 'material' AND lang.code = 'en' THEN 'Material'
    WHEN name = 'material' AND lang.code = 'ru' THEN 'Материал'
    WHEN name = 'material' AND lang.code = 'sr' THEN 'Materijal'
    -- Brand
    WHEN name = 'brand' AND lang.code = 'en' THEN 'Brand'
    WHEN name = 'brand' AND lang.code = 'ru' THEN 'Бренд'
    WHEN name = 'brand' AND lang.code = 'sr' THEN 'Brend'
    -- Model
    WHEN name = 'model' AND lang.code = 'en' THEN 'Model'
    WHEN name = 'model' AND lang.code = 'ru' THEN 'Модель'
    WHEN name = 'model' AND lang.code = 'sr' THEN 'Model'
  END as translation
FROM product_variant_attributes
CROSS JOIN (VALUES ('en'), ('ru'), ('sr')) AS lang(code)
WHERE name IN ('size', 'color', 'memory', 'storage', 'capacity', 'volume', 'weight', 'material', 'brand', 'model')
AND CASE 
    WHEN name = 'size' AND lang.code = 'en' THEN 'Size'
    WHEN name = 'size' AND lang.code = 'ru' THEN 'Размер'
    WHEN name = 'size' AND lang.code = 'sr' THEN 'Veličina'
    WHEN name = 'color' AND lang.code = 'en' THEN 'Color'
    WHEN name = 'color' AND lang.code = 'ru' THEN 'Цвет'
    WHEN name = 'color' AND lang.code = 'sr' THEN 'Boja'
    WHEN name = 'memory' AND lang.code = 'en' THEN 'Memory'
    WHEN name = 'memory' AND lang.code = 'ru' THEN 'Память'
    WHEN name = 'memory' AND lang.code = 'sr' THEN 'Memorija'
    WHEN name = 'storage' AND lang.code = 'en' THEN 'Storage'
    WHEN name = 'storage' AND lang.code = 'ru' THEN 'Хранилище'
    WHEN name = 'storage' AND lang.code = 'sr' THEN 'Skladište'
    WHEN name = 'capacity' AND lang.code = 'en' THEN 'Capacity'
    WHEN name = 'capacity' AND lang.code = 'ru' THEN 'Емкость'
    WHEN name = 'capacity' AND lang.code = 'sr' THEN 'Kapacitet'
    WHEN name = 'volume' AND lang.code = 'en' THEN 'Volume'
    WHEN name = 'volume' AND lang.code = 'ru' THEN 'Объем'
    WHEN name = 'volume' AND lang.code = 'sr' THEN 'Volumen'
    WHEN name = 'weight' AND lang.code = 'en' THEN 'Weight'
    WHEN name = 'weight' AND lang.code = 'ru' THEN 'Вес'
    WHEN name = 'weight' AND lang.code = 'sr' THEN 'Težina'
    WHEN name = 'material' AND lang.code = 'en' THEN 'Material'
    WHEN name = 'material' AND lang.code = 'ru' THEN 'Материал'
    WHEN name = 'material' AND lang.code = 'sr' THEN 'Materijal'
    WHEN name = 'brand' AND lang.code = 'en' THEN 'Brand'
    WHEN name = 'brand' AND lang.code = 'ru' THEN 'Бренд'
    WHEN name = 'brand' AND lang.code = 'sr' THEN 'Brend'
    WHEN name = 'model' AND lang.code = 'en' THEN 'Model'
    WHEN name = 'model' AND lang.code = 'ru' THEN 'Модель'
    WHEN name = 'model' AND lang.code = 'sr' THEN 'Model'
END IS NOT NULL
ON CONFLICT (entity_type, entity_id, field_name, language) DO UPDATE
SET translated_text = EXCLUDED.translated_text,
    updated_at = NOW();

-- Добавить индекс для быстрого поиска атрибутов, влияющих на остатки
CREATE INDEX idx_product_variant_attributes_affects_stock 
ON product_variant_attributes(affects_stock) 
WHERE affects_stock = true;