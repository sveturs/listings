-- Добавляем настройки отображения для атрибутов категорий
ALTER TABLE category_attributes 
ADD COLUMN IF NOT EXISTS show_in_card BOOLEAN DEFAULT true,
ADD COLUMN IF NOT EXISTS show_in_list BOOLEAN DEFAULT false;

-- Добавляем те же настройки в таблицу маппинга для переопределения на уровне категории
ALTER TABLE category_attribute_mapping 
ADD COLUMN IF NOT EXISTS show_in_card BOOLEAN DEFAULT NULL,
ADD COLUMN IF NOT EXISTS show_in_list BOOLEAN DEFAULT NULL;

-- Комментарии для понимания логики
COMMENT ON COLUMN category_attributes.show_in_card IS 'Показывать атрибут на странице детального просмотра товара';
COMMENT ON COLUMN category_attributes.show_in_list IS 'Показывать атрибут в списке товаров на странице поиска';
COMMENT ON COLUMN category_attribute_mapping.show_in_card IS 'Переопределение настройки show_in_card для конкретной категории (NULL = использовать значение из category_attributes)';
COMMENT ON COLUMN category_attribute_mapping.show_in_list IS 'Переопределение настройки show_in_list для конкретной категории (NULL = использовать значение из category_attributes)';

-- Обновляем существующие атрибуты с разумными значениями по умолчанию
-- Обязательные атрибуты всегда показываем в карточке
UPDATE category_attributes 
SET show_in_card = true 
WHERE is_required = true;

-- Основные атрибуты показываем в списке для удобства пользователей
UPDATE category_attributes 
SET show_in_list = true 
WHERE name IN (
    -- Автомобили
    'make', 'model', 'year', 'mileage', 'engine_capacity',
    -- Недвижимость  
    'property_type', 'rooms', 'area', 'floor',
    -- Телефоны
    'brand', 'model_phone', 'memory', 'ram',
    -- Компьютеры
    'pc_brand', 'pc_type', 'ram_pc', 'storage_capacity'
);

-- Все searchable и filterable атрибуты по умолчанию показываем в карточке
UPDATE category_attributes 
SET show_in_card = true 
WHERE (is_searchable = true OR is_filterable = true) 
AND show_in_card IS NOT true;