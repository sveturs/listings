-- Настройка атрибутов которые могут использоваться как варианты товаров

-- Атрибуты которые могут использоваться как варианты и НЕ влияют на остатки
UPDATE category_attributes 
SET is_variant_compatible = true, affects_stock = false 
WHERE name IN (
    'color',           -- цвет
    'brand',           -- бренд/марка
    'material',        -- материал  
    'style',           -- стиль
    'pattern',         -- узор/рисунок
    'condition',       -- состояние (новое/б/у)
    'fuel_type',       -- тип топлива
    'transmission',    -- коробка передач
    'drivetrain',      -- привод
    'body_type',       -- тип кузова
    'doors',           -- количество дверей
    'seats'            -- количество мест
);

-- Атрибуты которые могут использоваться как варианты и ВЛИЯЮТ на остатки
UPDATE category_attributes 
SET is_variant_compatible = true, affects_stock = true 
WHERE name IN (
    'size',            -- размер (одежда, обувь)
    'storage',         -- объем памяти
    'ram',             -- оперативная память
    'display_size',    -- размер дисплея
    'engine_size',     -- объем двигателя
    'power_hp',        -- мощность
    'year',            -- год выпуска
    'mileage',         -- пробег
    'area',            -- площадь (недвижимость)
    'rooms',           -- количество комнат
    'floor'            -- этаж
);

-- Создаем индекс для быстрого поиска вариативных атрибутов
CREATE INDEX IF NOT EXISTS idx_category_attributes_variant_compatible 
ON category_attributes (is_variant_compatible) 
WHERE is_variant_compatible = true;