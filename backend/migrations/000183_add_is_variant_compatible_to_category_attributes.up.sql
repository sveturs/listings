-- Добавление поля is_variant_compatible в таблицу category_attributes
-- Это поле указывает, может ли атрибут категории использоваться как вариативный атрибут товара

ALTER TABLE category_attributes 
ADD COLUMN is_variant_compatible BOOLEAN DEFAULT FALSE;

-- Добавляем комментарий к полю
COMMENT ON COLUMN category_attributes.is_variant_compatible IS 'Указывает, может ли атрибут использоваться для вариантов товаров';

-- Автоматически отмечаем атрибуты, которые могут быть вариативными, на основе их названий
-- Это предварительная настройка, которую можно изменить через админку
UPDATE category_attributes 
SET is_variant_compatible = TRUE
WHERE LOWER(name) IN (
    'color', 'colour', 'цвет', 'boja',
    'size', 'размер', 'veličina', 
    'memory', 'память', 'memorija', 'ram',
    'storage', 'хранилище', 'складиште',
    'material', 'материал', 'materijal',
    'capacity', 'емкость', 'kapacitet', 'volume', 'объем',
    'power', 'мощность', 'snaga',
    'connectivity', 'подключение', 'povezivanje', 
    'style', 'стиль', 'stil',
    'pattern', 'узор', 'šablon',
    'weight', 'вес', 'težina',
    'bundle', 'комплект', 'paket'
);

-- Создаем индекс для быстрого поиска вариативных атрибутов
CREATE INDEX idx_category_attributes_is_variant_compatible 
ON category_attributes(is_variant_compatible) 
WHERE is_variant_compatible = TRUE;