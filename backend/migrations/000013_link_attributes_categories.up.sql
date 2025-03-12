-- backend/migrations/000013_link_attributes_categories.up.sql
-- Находим ID категории "Автомобили"
DO $$
DECLARE
    car_category_id INT;
    realty_category_id INT;
    phones_category_id INT;
    computers_category_id INT;
BEGIN
    -- Получаем ID для основных категорий
    SELECT id INTO car_category_id FROM marketplace_categories WHERE name = 'Автомобили' OR name = 'Cars' LIMIT 1;
    SELECT id INTO realty_category_id FROM marketplace_categories WHERE name = 'Недвижимость' OR name = 'Real Estate' LIMIT 1;
    SELECT id INTO phones_category_id FROM marketplace_categories WHERE name = 'Телефоны' OR name = 'Phones' LIMIT 1;
    SELECT id INTO computers_category_id FROM marketplace_categories WHERE name = 'Компьютеры' OR name = 'Computers' LIMIT 1;

    -- Связываем атрибуты с категориями автомобилей
    IF car_category_id IS NOT NULL THEN
        INSERT INTO category_attribute_mapping (category_id, attribute_id)
        SELECT c.id, a.id 
        FROM marketplace_categories c 
        CROSS JOIN category_attributes a
        WHERE (c.id = car_category_id OR c.parent_id = car_category_id OR EXISTS (
            SELECT 1 FROM marketplace_categories c2 
            WHERE c2.parent_id = car_category_id AND c.parent_id = c2.id
        ))
        AND a.name IN ('make', 'model', 'year', 'mileage', 'engine_capacity', 'fuel_type', 'transmission', 'body_type', 'color')
        ON CONFLICT DO NOTHING;
    END IF;

    -- Связываем атрибуты с категориями недвижимости
    IF realty_category_id IS NOT NULL THEN
        INSERT INTO category_attribute_mapping (category_id, attribute_id)
        SELECT c.id, a.id 
        FROM marketplace_categories c 
        CROSS JOIN category_attributes a
        WHERE (c.id = realty_category_id OR c.parent_id = realty_category_id OR EXISTS (
            SELECT 1 FROM marketplace_categories c2 
            WHERE c2.parent_id = realty_category_id AND c.parent_id = c2.id
        ))
        AND a.name IN ('property_type', 'rooms', 'floor', 'total_floors', 'area', 'land_area', 'building_type', 'has_balcony', 'has_elevator', 'has_parking')
        ON CONFLICT DO NOTHING;
    END IF;

    -- Связываем атрибуты с категориями телефонов
    IF phones_category_id IS NOT NULL THEN
        INSERT INTO category_attribute_mapping (category_id, attribute_id)
        SELECT c.id, a.id 
        FROM marketplace_categories c 
        CROSS JOIN category_attributes a
        WHERE (c.id = phones_category_id OR c.parent_id = phones_category_id)
        AND a.name IN ('brand', 'model_phone', 'memory', 'ram', 'os', 'screen_size', 'camera', 'has_5g')
        ON CONFLICT DO NOTHING;
    END IF;

    -- Связываем атрибуты с категориями компьютеров
    IF computers_category_id IS NOT NULL THEN
        INSERT INTO category_attribute_mapping (category_id, attribute_id)
        SELECT c.id, a.id 
        FROM marketplace_categories c 
        CROSS JOIN category_attributes a
        WHERE (c.id = computers_category_id OR c.parent_id = computers_category_id)
        AND a.name IN ('pc_brand', 'pc_type', 'cpu', 'gpu', 'ram_pc', 'storage_type', 'storage_capacity', 'os_pc')
        ON CONFLICT DO NOTHING;
    END IF;
END $$;

-- Связываем атрибуты компьютеров с нужными категориями
INSERT INTO category_attribute_mapping (category_id, attribute_id)
SELECT 
  c.id, 
  a.id
FROM 
  marketplace_categories c
CROSS JOIN 
  category_attributes a
WHERE 
  c.id IN (3310, 3320,
 3600)
  AND a.name IN ('pc_brand', 'pc_type', 'cpu', 'gpu', 'ram_pc', 'storage_type', 'storage_capacity', 'os_pc')
ON CONFLICT DO NOTHING;

-- Связываем атрибуты телефонов с категорией планшетов (3810)
INSERT INTO category_attribute_mapping (category_id, attribute_id)
SELECT 
  3810, 
  a.id
FROM 
  category_attributes a
WHERE 
  a.name IN ('brand', 'model_phone', 'memory', 'ram', 'os', 'screen_size', 'camera', 'has_5g');