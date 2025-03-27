-- Объединённая миграция атрибутов категорий

-- Таблица для определения атрибутов категорий
CREATE TABLE IF NOT EXISTS category_attributes (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    attribute_type VARCHAR(50) NOT NULL, -- text, number, select, boolean, etc.
    options JSONB, -- Для хранения вариантов выбора, диапазонов и т.д.
    validation_rules JSONB, -- Правила валидации
    is_searchable BOOLEAN DEFAULT true,
    is_filterable BOOLEAN DEFAULT true,
    is_required BOOLEAN DEFAULT false,
    sort_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица для связи атрибутов с категориями
CREATE TABLE IF NOT EXISTS category_attribute_mapping (
    category_id INT NOT NULL REFERENCES marketplace_categories(id) ON DELETE CASCADE,
    attribute_id INT NOT NULL REFERENCES category_attributes(id) ON DELETE CASCADE,
    is_enabled BOOLEAN DEFAULT true,
    is_required BOOLEAN DEFAULT false,
    PRIMARY KEY (category_id, attribute_id)
);

-- Таблица для хранения значений атрибутов объявлений
CREATE TABLE IF NOT EXISTS listing_attribute_values (
    id SERIAL PRIMARY KEY,
    listing_id INT NOT NULL REFERENCES marketplace_listings(id) ON DELETE CASCADE,
    attribute_id INT NOT NULL REFERENCES category_attributes(id) ON DELETE CASCADE,
    text_value TEXT,
    numeric_value NUMERIC(20,5),
    boolean_value BOOLEAN,
    json_value JSONB
);

-- Создаем индекс для быстрого поиска атрибутов конкретного объявления
CREATE INDEX IF NOT EXISTS idx_listing_attr_listing_id ON listing_attribute_values(listing_id);

-- Индексы для ускорения поиска по атрибутам
CREATE INDEX IF NOT EXISTS idx_listing_attr_text ON listing_attribute_values(attribute_id, text_value) WHERE text_value IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_listing_attr_numeric ON listing_attribute_values(attribute_id, numeric_value) WHERE numeric_value IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_listing_attr_boolean ON listing_attribute_values(attribute_id, boolean_value) WHERE boolean_value IS NOT NULL;

-- Вставка начальных атрибутов для различных категорий товаров

-- Автомобили
INSERT INTO category_attributes (name, display_name, attribute_type, options, is_searchable, is_filterable, is_required) VALUES
('make', 'Make', 'select', '{"values": ["Audi", "BMW", "Mercedes", "Toyota", "Honda", "Ford", "Chevrolet", "Volkswagen", "Nissan", "Other"]}', true, true, true),
('model', 'Model', 'text', NULL, true, true, true),
('year', 'Year', 'number', '{"min": 1990, "max": 2025}', true, true, true),
('mileage', 'Mileage', 'number', '{"min": 0}', true, true, false),
('engine_capacity', 'Engine capacity', 'number', '{"min": 0.1, "max": 10, "step": 0.1}', true, true, false),
('fuel_type', 'Fuel type', 'select', '{"values": ["Бензин", "Дизель", "Гибрид", "Электро", "Газ"]}', true, true, false),
('transmission', 'Transmission', 'select', '{"values": ["Механика", "Автомат", "Робот", "Вариатор"]}', true, true, false),
('body_type', 'Body type', 'select', '{"values": ["Седан", "Хэтчбек", "Универсал", "Внедорожник", "Купе", "Кабриолет", "Минивэн", "Пикап"]}', true, true, false),
('color', 'Color', 'select', '{"values": ["Белый", "Черный", "Серый", "Серебристый", "Красный", "Синий", "Зеленый", "Желтый", "Коричневый", "Другой"]}', true, true, false),
('power', 'Power', 'number', '{"min": 1, "max": 2000}', true, true, false),
('drive_type', 'Drive type', 'select', '{"values": ["Передний", "Задний", "Полный", "Другой"]}', true, true, false),
('number_of_doors', 'Number of doors', 'select', '{"values": ["2", "3", "4", "5", "6+"]}', true, true, false),
('number_of_seats', 'Number of seats', 'select', '{"values": ["1", "2", "3", "4", "5", "6", "7", "8+"]}', true, true, false)
ON CONFLICT (id) DO NOTHING;

-- Недвижимость
INSERT INTO category_attributes (name, display_name, attribute_type, options, is_searchable, is_filterable, is_required) VALUES
('property_type', 'Тип недвижимости', 'select', '{"values": ["Квартира", "Дом", "Комната", "Земельный участок", "Коммерческая недвижимость", "Гараж"]}', true, true, true),
('rooms', 'Количество комнат', 'select', '{"values": ["Студия", "1", "2", "3", "4", "5+"]}', true, true, false),
('floor', 'Этаж', 'number', '{"min": 0, "max": 100}', true, true, false),
('total_floors', 'Этажей в доме', 'number', '{"min": 1, "max": 100}', true, true, false),
('area', 'Площадь (м²)', 'number', '{"min": 1}', true, true, true),
('land_area', 'Площадь участка (сот.)', 'number', '{"min": 0}', true, true, false),
('building_type', 'Тип дома', 'select', '{"values": ["Панельный", "Кирпичный", "Монолитный", "Деревянный", "Блочный", "Другой"]}', true, true, false),
('has_balcony', 'Балкон', 'boolean', NULL, true, true, false),
('has_elevator', 'Лифт', 'boolean', NULL, true, true, false),
('has_parking', 'Парковка', 'boolean', NULL, true, true, false)
ON CONFLICT (id) DO NOTHING;

-- Электроника - Телефоны
INSERT INTO category_attributes (name, display_name, attribute_type, options, is_searchable, is_filterable, is_required) VALUES
('brand', 'Бренд', 'select', '{"values": ["Apple", "Samsung", "Xiaomi", "Huawei", "Google", "OnePlus", "Sony", "Nokia", "LG", "Other"]}', true, true, true),
('model_phone', 'Модель', 'text', NULL, true, true, true),
('memory', 'Память (ГБ)', 'select', '{"values": ["8", "16", "32", "64", "128", "256", "512", "1024"]}', true, true, false),
('ram', 'ОЗУ (ГБ)', 'select', '{"values": ["1", "2", "3", "4", "6", "8", "12", "16"]}', true, true, false),
('os', 'Операционная система', 'select', '{"values": ["iOS", "Android", "Windows", "Другая"]}', true, true, false),
('screen_size', 'Размер экрана (дюймы)', 'number', '{"min": 1, "max": 15, "step": 0.1}', true, true, false),
('camera', 'Камера (МП)', 'number', '{"min": 1}', true, true, false),
('has_5g', '5G', 'boolean', NULL, true, true, false)
ON CONFLICT (id) DO NOTHING;

-- Компьютеры
INSERT INTO category_attributes (name, display_name, attribute_type, options, is_searchable, is_filterable, is_required) VALUES
('pc_brand', 'Бренд', 'select', '{"values": ["Apple", "Dell", "HP", "Lenovo", "Asus", "Acer", "MSI", "Gigabyte", "Сборка"]}', true, true, true),
('pc_type', 'Тип', 'select', '{"values": ["Ноутбук", "Настольный ПК", "Моноблок", "Сервер", "Другое"]}', true, true, true),
('cpu', 'Процессор', 'text', NULL, true, true, false),
('gpu', 'Видеокарта', 'text', NULL, true, true, false),
('ram_pc', 'ОЗУ (ГБ)', 'select', '{"values": ["2", "4", "8", "16", "32", "64", "128"]}', true, true, false),
('storage_type', 'Тип накопителя', 'select', '{"values": ["HDD", "SSD", "HDD+SSD"]}', true, true, false),
('storage_capacity', 'Объем накопителя (ГБ)', 'number', '{"min": 1}', true, true, false),
('os_pc', 'Операционная система', 'select', '{"values": ["Windows", "macOS", "Linux", "Без ОС"]}', true, true, false)
ON CONFLICT (id) DO NOTHING;

-- Добавляем переводы для атрибутов на английский
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified, created_at, updated_at) 
SELECT 'attribute', id, 'en', 'display_name', 
    CASE 
        -- Автомобили
        WHEN name = 'make' THEN 'Make'
        WHEN name = 'model' THEN 'Model'
        WHEN name = 'year' THEN 'Year'
        WHEN name = 'mileage' THEN 'Mileage (km)'
        WHEN name = 'engine_capacity' THEN 'Engine capacity (L)'
        WHEN name = 'fuel_type' THEN 'Fuel type'
        WHEN name = 'transmission' THEN 'Transmission'
        WHEN name = 'body_type' THEN 'Body type'
        WHEN name = 'color' THEN 'Color'
        WHEN name = 'power' THEN 'Power (hp)'
        WHEN name = 'drive_type' THEN 'Drive Type'
        WHEN name = 'number_of_doors' THEN 'Number of Doors'
        WHEN name = 'number_of_seats' THEN 'Number of Seats'
        -- Недвижимость
        WHEN name = 'property_type' THEN 'Property type'
        WHEN name = 'rooms' THEN 'Rooms'
        WHEN name = 'floor' THEN 'Floor'
        WHEN name = 'total_floors' THEN 'Total floors'
        WHEN name = 'area' THEN 'Area (m²)'
        WHEN name = 'land_area' THEN 'Land area'
        WHEN name = 'building_type' THEN 'Building type'
        WHEN name = 'has_balcony' THEN 'Balcony'
        WHEN name = 'has_elevator' THEN 'Elevator'
        WHEN name = 'has_parking' THEN 'Parking'
        -- Телефоны
        WHEN name = 'brand' THEN 'Brand'
        WHEN name = 'model_phone' THEN 'Model'
        WHEN name = 'memory' THEN 'Memory (GB)'
        WHEN name = 'ram' THEN 'RAM (GB)'
        WHEN name = 'os' THEN 'Operating system'
        WHEN name = 'screen_size' THEN 'Screen size (inches)'
        WHEN name = 'camera' THEN 'Camera (MP)'
        WHEN name = 'has_5g' THEN '5G'
        -- Компьютеры
        WHEN name = 'pc_brand' THEN 'Brand'
        WHEN name = 'pc_type' THEN 'Type'
        WHEN name = 'cpu' THEN 'Processor'
        WHEN name = 'gpu' THEN 'Graphics card'
        WHEN name = 'ram_pc' THEN 'RAM (GB)'
        WHEN name = 'storage_type' THEN 'Storage type'
        WHEN name = 'storage_capacity' THEN 'Storage capacity (GB)'
        WHEN name = 'os_pc' THEN 'Operating system'
        ELSE display_name
    END, 
    false, true, NOW(), NOW()
FROM category_attributes
ON CONFLICT DO NOTHING;

-- Добавляем переводы для атрибутов на сербский
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified, created_at, updated_at) 
SELECT 'attribute', id, 'sr', 'display_name', 
    CASE 
        -- Автомобили
        WHEN name = 'make' THEN 'Marka'
        WHEN name = 'model' THEN 'Model'
        WHEN name = 'year' THEN 'Godina proizvodnje'
        WHEN name = 'mileage' THEN 'Kilometraža'
        WHEN name = 'engine_capacity' THEN 'Zapremina motora (L)'
        WHEN name = 'fuel_type' THEN 'Vrsta goriva'
        WHEN name = 'transmission' THEN 'Menjač'
        WHEN name = 'body_type' THEN 'Tip karoserije'
        WHEN name = 'color' THEN 'Boja'
        WHEN name = 'power' THEN 'Snaga (ks)'
        WHEN name = 'drive_type' THEN 'Pogon'
        WHEN name = 'number_of_doors' THEN 'Broj vrata'
        WHEN name = 'number_of_seats' THEN 'Broj sedišta'
        -- Недвижимость
        WHEN name = 'property_type' THEN 'Tip nekretnine'
        WHEN name = 'rooms' THEN 'Broj soba'
        WHEN name = 'floor' THEN 'Sprat'
        WHEN name = 'total_floors' THEN 'Ukupno spratova'
        WHEN name = 'area' THEN 'Površina (m²)'
        WHEN name = 'land_area' THEN 'Površina zemljišta'
        WHEN name = 'building_type' THEN 'Tip zgrade'
        WHEN name = 'has_balcony' THEN 'Balkon'
        WHEN name = 'has_elevator' THEN 'Lift'
        WHEN name = 'has_parking' THEN 'Parking'
        -- Телефоны
        WHEN name = 'brand' THEN 'Brend'
        WHEN name = 'model_phone' THEN 'Model'
        WHEN name = 'memory' THEN 'Memorija (GB)'
        WHEN name = 'ram' THEN 'RAM (GB)'
        WHEN name = 'os' THEN 'Operativni sistem'
        WHEN name = 'screen_size' THEN 'Veličina ekrana (inči)'
        WHEN name = 'camera' THEN 'Kamera (MP)'
        WHEN name = 'has_5g' THEN '5G'
        -- Компьютеры
        WHEN name = 'pc_brand' THEN 'Brend'
        WHEN name = 'pc_type' THEN 'Tip'
        WHEN name = 'cpu' THEN 'Procesor'
        WHEN name = 'gpu' THEN 'Grafička kartica'
        WHEN name = 'ram_pc' THEN 'RAM (GB)'
        WHEN name = 'storage_type' THEN 'Tip skladišta'
        WHEN name = 'storage_capacity' THEN 'Kapacitet skladišta (GB)'
        WHEN name = 'os_pc' THEN 'Operativni sistem'
        ELSE display_name
    END, 
    false, true, NOW(), NOW()
FROM category_attributes
ON CONFLICT DO NOTHING;

-- Очистка существующих связей (чтобы избежать дублирования)
DELETE FROM category_attribute_mapping;

-- Создание связей для категории Автомобили (2000) и всех её подкатегорий
INSERT INTO category_attribute_mapping (category_id, attribute_id, is_enabled, is_required)
SELECT c.id, a.id, true, false
FROM marketplace_categories c
CROSS JOIN category_attributes a
WHERE (c.id = 2000 OR c.parent_id = 2000 OR EXISTS (
    SELECT 1 FROM marketplace_categories c2 
    WHERE c2.parent_id = 2000 AND c.parent_id = c2.id
))
AND a.name IN ('make', 'model', 'year', 'mileage', 'engine_capacity', 'fuel_type', 
               'transmission', 'body_type', 'color', 'power', 'drive_type', 
               'number_of_doors', 'number_of_seats');

-- Создание связей для категории Недвижимость (1000) и подкатегорий
INSERT INTO category_attribute_mapping (category_id, attribute_id, is_enabled, is_required)
SELECT c.id, a.id, true, false
FROM marketplace_categories c
CROSS JOIN category_attributes a
WHERE (c.id = 1000 OR c.parent_id = 1000 OR EXISTS (
    SELECT 1 FROM marketplace_categories c2 
    WHERE c2.parent_id = 1000 AND c.parent_id = c2.id
))
AND a.name IN ('property_type', 'rooms', 'floor', 'total_floors', 'area', 'land_area', 
               'building_type', 'has_balcony', 'has_elevator', 'has_parking');

-- Создание связей для категории Телефоны (3100) и подкатегорий
INSERT INTO category_attribute_mapping (category_id, attribute_id, is_enabled, is_required)
SELECT c.id, a.id, true, false
FROM marketplace_categories c
CROSS JOIN category_attributes a
WHERE (c.id = 3100 OR c.parent_id = 3100 OR EXISTS (
    SELECT 1 FROM marketplace_categories c2 
    WHERE c2.parent_id = 3100 AND c.parent_id = c2.id
))
AND a.name IN ('brand', 'model_phone', 'memory', 'ram', 'os', 'screen_size', 'camera', 'has_5g');

-- Создание связей для категории Планшеты (3810)
INSERT INTO category_attribute_mapping (category_id, attribute_id, is_enabled, is_required)
SELECT c.id, a.id, true, false
FROM marketplace_categories c
CROSS JOIN category_attributes a
WHERE (c.id = 3810 OR c.parent_id = 3810)
AND a.name IN ('brand', 'model_phone', 'memory', 'ram', 'os', 'screen_size', 'camera', 'has_5g');

-- Создание связей для категории Системные блоки (3310), Моноблоки (3320) и их подкатегорий
INSERT INTO category_attribute_mapping (category_id, attribute_id, is_enabled, is_required)
SELECT c.id, a.id, true, false
FROM marketplace_categories c
CROSS JOIN category_attributes a
WHERE (c.id IN (3310, 3320) OR c.parent_id IN (3310, 3320))
AND a.name IN ('pc_brand', 'pc_type', 'cpu', 'gpu', 'ram_pc', 'storage_type', 
               'storage_capacity', 'os_pc');

-- Создание связей для категории Ноутбуки (3600) и подкатегорий
INSERT INTO category_attribute_mapping (category_id, attribute_id, is_enabled, is_required)
SELECT c.id, a.id, true, false
FROM marketplace_categories c
CROSS JOIN category_attributes a
WHERE (c.id = 3600 OR c.parent_id = 3600)
AND a.name IN ('pc_brand', 'pc_type', 'cpu', 'gpu', 'ram_pc', 'storage_type', 
               'storage_capacity', 'os_pc');

-- Метка обязательных атрибутов для автомобилей
UPDATE category_attribute_mapping 
SET is_required = true
WHERE attribute_id IN (
    SELECT id FROM category_attributes 
    WHERE name IN ('make', 'model', 'year')
)
AND category_id IN (
    SELECT id FROM marketplace_categories 
    WHERE id = 2000 OR parent_id = 2000 OR EXISTS (
        SELECT 1 FROM marketplace_categories c2 
        WHERE c2.parent_id = 2000 AND parent_id = c2.id
    )
);

-- Метка обязательных атрибутов для недвижимости
UPDATE category_attribute_mapping 
SET is_required = true
WHERE attribute_id IN (
    SELECT id FROM category_attributes 
    WHERE name IN ('property_type', 'area')
)
AND category_id IN (
    SELECT id FROM marketplace_categories 
    WHERE id = 1000 OR parent_id = 1000 OR EXISTS (
        SELECT 1 FROM marketplace_categories c2 
        WHERE c2.parent_id = 1000 AND parent_id = c2.id
    )
);

-- Метка обязательных атрибутов для телефонов
UPDATE category_attribute_mapping 
SET is_required = true
WHERE attribute_id IN (
    SELECT id FROM category_attributes 
    WHERE name IN ('brand', 'model_phone')
)
AND category_id IN (
    SELECT id FROM marketplace_categories 
    WHERE id = 3100 OR parent_id = 3100 OR EXISTS (
        SELECT 1 FROM marketplace_categories c2 
        WHERE c2.parent_id = 3100 AND parent_id = c2.id
    )
);

-- Метка обязательных атрибутов для компьютеров
UPDATE category_attribute_mapping 
SET is_required = true
WHERE attribute_id IN (
    SELECT id FROM category_attributes 
    WHERE name IN ('pc_brand', 'pc_type')
)
AND category_id IN (
    SELECT id FROM marketplace_categories 
    WHERE id IN (3310, 3320, 3600) OR parent_id IN (3310, 3320, 3600)
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_listing_attr_unique ON listing_attribute_values (listing_id, attribute_id);
-- Обновляем материализованное представление категорий, если оно существует
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_class WHERE relname = 'category_listing_counts') THEN
        REFRESH MATERIALIZED VIEW category_listing_counts;
    END IF;
END
$$;