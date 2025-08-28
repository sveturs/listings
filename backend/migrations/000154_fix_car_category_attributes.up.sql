-- Создаем новый атрибут car_make для автомобильных марок
INSERT INTO category_attributes (id, name, display_name, attribute_type, options, is_searchable, is_filterable, is_required, sort_order)
VALUES (
    2210, 
    'car_make', 
    'Marka',
    'select',
    '{"values": ["Audi", "BMW", "Mercedes-Benz", "Volkswagen", "Toyota", "Honda", "Nissan", "Mazda", "Ford", "Chevrolet", "Opel", "Renault", "Peugeot", "Citroen", "Fiat", "Alfa Romeo", "Volvo", "Škoda", "Seat", "Hyundai", "Kia", "Mitsubishi", "Suzuki", "Subaru", "Lexus", "Porsche", "Ferrari", "Lamborghini", "Tesla", "Other"]}'::jsonb,
    true,
    true,
    true,
    3
) ON CONFLICT (id) DO NOTHING;

-- Обновляем связь категории автомобилей - удаляем общий brand и добавляем car_make
DELETE FROM category_attribute_mapping 
WHERE category_id = 1301 AND attribute_id = 2003;

INSERT INTO category_attribute_mapping (category_id, attribute_id, is_enabled, is_required, sort_order)
VALUES (1301, 2210, true, true, 3)
ON CONFLICT (category_id, attribute_id) DO UPDATE
SET is_enabled = true, is_required = true, sort_order = 3;

-- Сдвигаем порядок следующих атрибутов
UPDATE category_attribute_mapping 
SET sort_order = sort_order + 1
WHERE category_id = 1301 AND sort_order >= 3 AND attribute_id != 2210;

-- Также давайте создадим специфичные атрибуты для других категорий, которые используют общий brand

-- Для категории мотоциклов (1302)
INSERT INTO category_attributes (id, name, display_name, attribute_type, options, is_searchable, is_filterable, is_required, sort_order)
VALUES (
    2211, 
    'motorcycle_make', 
    'Marka',
    'select',
    '{"values": ["Yamaha", "Honda", "Kawasaki", "Suzuki", "BMW", "Harley-Davidson", "Ducati", "KTM", "Aprilia", "Triumph", "Indian", "Vespa", "Piaggio", "Other"]}'::jsonb,
    true,
    true,
    true,
    3
) ON CONFLICT (id) DO NOTHING;

DELETE FROM category_attribute_mapping 
WHERE category_id = 1302 AND attribute_id = 2003;

INSERT INTO category_attribute_mapping (category_id, attribute_id, is_enabled, is_required, sort_order)
VALUES (1302, 2211, true, true, 3)
ON CONFLICT (category_id, attribute_id) DO UPDATE
SET is_enabled = true, is_required = true, sort_order = 3;

UPDATE category_attribute_mapping 
SET sort_order = sort_order + 1
WHERE category_id = 1302 AND sort_order >= 3 AND attribute_id != 2211;

-- Для грузовиков и спецтехники (1303)
INSERT INTO category_attributes (id, name, display_name, attribute_type, options, is_searchable, is_filterable, is_required, sort_order)
VALUES (
    2212, 
    'truck_make', 
    'Marka',
    'select',
    '{"values": ["MAN", "Mercedes-Benz", "Volvo", "Scania", "DAF", "Iveco", "Renault", "Kamaz", "MAZ", "Tatra", "Caterpillar", "JCB", "Komatsu", "Liebherr", "Other"]}'::jsonb,
    true,
    true,
    true,
    3
) ON CONFLICT (id) DO NOTHING;

DELETE FROM category_attribute_mapping 
WHERE category_id = 1303 AND attribute_id = 2003;

INSERT INTO category_attribute_mapping (category_id, attribute_id, is_enabled, is_required, sort_order)
VALUES (1303, 2212, true, true, 3)
ON CONFLICT (category_id, attribute_id) DO UPDATE
SET is_enabled = true, is_required = true, sort_order = 3;

UPDATE category_attribute_mapping 
SET sort_order = sort_order + 1
WHERE category_id = 1303 AND sort_order >= 3 AND attribute_id != 2212;

-- Для водного транспорта (1304)
INSERT INTO category_attributes (id, name, display_name, attribute_type, options, is_searchable, is_filterable, is_required, sort_order)
VALUES (
    2213, 
    'boat_make', 
    'Marka',
    'select',
    '{"values": ["Yamaha", "Mercury", "Suzuki", "Honda", "Sea-Doo", "Bayliner", "Sea Ray", "Azimut", "Princess", "Sunseeker", "Other"]}'::jsonb,
    true,
    true,
    true,
    3
) ON CONFLICT (id) DO NOTHING;

DELETE FROM category_attribute_mapping 
WHERE category_id = 1304 AND attribute_id = 2003;

INSERT INTO category_attribute_mapping (category_id, attribute_id, is_enabled, is_required, sort_order)
VALUES (1304, 2213, true, true, 3)
ON CONFLICT (category_id, attribute_id) DO UPDATE
SET is_enabled = true, is_required = true, sort_order = 3;

UPDATE category_attribute_mapping 
SET sort_order = sort_order + 1
WHERE category_id = 1304 AND sort_order >= 3 AND attribute_id != 2213;