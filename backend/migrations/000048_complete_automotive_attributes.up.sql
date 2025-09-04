-- Добавление базовых атрибутов для всех автомобильных категорий без атрибутов
-- Migration 000048: Complete automotive category attributes

-- Создание универсальных автомобильных атрибутов
INSERT INTO unified_attributes (code, name, display_name, attribute_type, purpose, options, validation_rules, ui_settings, is_searchable, is_filterable, is_required, is_variant_compatible, affects_stock, affects_price, sort_order, is_active) VALUES
    ('auto_part_brand', 'Part Brand', 'Бренд запчасти', 'text', 'regular', NULL, '{"minLength": 2, "maxLength": 50}', '{"showInCard": true}', true, true, false, false, false, false, 10, true),
    ('auto_part_oem', 'OEM Number', 'OEM номер', 'text', 'regular', NULL, '{"maxLength": 50}', '{"showInCard": true}', true, true, false, false, false, false, 20, true),
    ('auto_part_condition', 'Part Condition', 'Состояние', 'select', 'regular', '["New", "Used - Like New", "Used - Good", "Used - Fair", "Refurbished", "For Parts"]', '{"required": true}', '{"showInCard": true}', true, true, true, false, false, false, 30, true),
    ('auto_compatibility', 'Compatibility', 'Совместимость', 'text', 'regular', NULL, '{"maxLength": 200}', '{"showInCard": false}', true, true, false, false, false, false, 40, true),
    ('auto_warranty', 'Warranty', 'Гарантия', 'select', 'regular', '["No Warranty", "1 Month", "3 Months", "6 Months", "1 Year", "2 Years", "Lifetime"]', NULL, '{"showInCard": true}', false, true, false, false, false, false, 50, true),
    ('auto_year_from', 'Year From', 'Год от', 'number', 'regular', NULL, '{"min": 1900, "max": 2025}', '{"showInCard": false}', true, true, false, false, false, false, 60, true),
    ('auto_year_to', 'Year To', 'Год до', 'number', 'regular', NULL, '{"min": 1900, "max": 2025}', '{"showInCard": false}', true, true, false, false, false, false, 70, true),
    ('auto_installation', 'Installation Available', 'Установка', 'boolean', 'regular', NULL, NULL, '{"showInCard": false}', false, true, false, false, false, false, 80, true)
ON CONFLICT (code) DO NOTHING;

-- Создание специфичных атрибутов для разных типов автозапчастей
INSERT INTO unified_attributes (code, name, display_name, attribute_type, purpose, options, validation_rules, ui_settings, is_searchable, is_filterable, is_required, is_variant_compatible, affects_stock, affects_price, sort_order, is_active) VALUES
    -- Для шин
    ('tire_width', 'Tire Width', 'Ширина шины', 'number', 'regular', NULL, '{"min": 100, "max": 400}', '{"showInCard": true}', true, true, false, true, true, true, 90, true),
    ('tire_profile', 'Tire Profile', 'Профиль шины', 'number', 'regular', NULL, '{"min": 20, "max": 100}', '{"showInCard": true}', true, true, false, true, true, true, 100, true),
    ('tire_diameter', 'Tire Diameter', 'Диаметр', 'number', 'regular', NULL, '{"min": 10, "max": 30}', '{"showInCard": true}', true, true, false, true, true, true, 110, true),
    ('tire_season', 'Season', 'Сезон', 'select', 'regular', '["Summer", "Winter", "All-Season"]', '{"required": true}', '{"showInCard": true}', true, true, true, false, false, false, 120, true),
    ('tire_speed_index', 'Speed Index', 'Индекс скорости', 'text', 'regular', NULL, '{"maxLength": 5}', '{"showInCard": false}', true, true, false, false, false, false, 130, true),
    ('tire_load_index', 'Load Index', 'Индекс нагрузки', 'number', 'regular', NULL, '{"min": 50, "max": 150}', '{"showInCard": false}', true, true, false, false, false, false, 140, true),
    
    -- Для дисков
    ('rim_diameter', 'Rim Diameter', 'Диаметр диска', 'number', 'regular', NULL, '{"min": 10, "max": 30}', '{"showInCard": true}', true, true, false, true, true, true, 150, true),
    ('rim_width', 'Rim Width', 'Ширина диска', 'number', 'regular', NULL, '{"min": 4, "max": 15}', '{"showInCard": true}', true, true, false, true, true, true, 160, true),
    ('rim_bolt_pattern', 'Bolt Pattern', 'Разболтовка', 'text', 'regular', NULL, '{"maxLength": 20}', '{"showInCard": true}', true, true, false, false, false, false, 170, true),
    ('rim_offset', 'Offset (ET)', 'Вылет (ET)', 'number', 'regular', NULL, '{"min": -50, "max": 100}', '{"showInCard": false}', true, true, false, false, false, false, 180, true),
    ('rim_center_bore', 'Center Bore', 'Центральное отверстие', 'number', 'regular', NULL, '{"min": 50, "max": 150}', '{"showInCard": false}', true, true, false, false, false, false, 190, true),
    
    -- Для двигателя
    ('engine_volume', 'Engine Volume', 'Объем двигателя', 'number', 'regular', NULL, '{"min": 0.5, "max": 10}', '{"showInCard": true}', true, true, false, false, false, false, 200, true),
    ('engine_power', 'Engine Power', 'Мощность', 'number', 'regular', NULL, '{"min": 50, "max": 1000}', '{"showInCard": true}', true, true, false, false, false, false, 210, true),
    ('engine_type', 'Engine Type', 'Тип двигателя', 'select', 'regular', '["Petrol", "Diesel", "Hybrid", "Electric", "Gas", "Other"]', NULL, '{"showInCard": true}', true, true, false, false, false, false, 220, true),
    
    -- Для коммерческих и специальных транспортных средств
    ('vehicle_capacity', 'Capacity', 'Грузоподъемность', 'number', 'regular', NULL, '{"min": 100, "max": 50000}', '{"showInCard": true}', true, true, false, false, false, false, 230, true),
    ('vehicle_seats', 'Number of Seats', 'Количество мест', 'number', 'regular', NULL, '{"min": 1, "max": 100}', '{"showInCard": true}', true, true, false, false, false, false, 240, true),
    ('vehicle_axles', 'Number of Axles', 'Количество осей', 'number', 'regular', NULL, '{"min": 1, "max": 10}', '{"showInCard": false}', true, true, false, false, false, false, 250, true)
ON CONFLICT (code) DO NOTHING;

-- Функция для массового добавления базовых атрибутов к автомобильным категориям
DO $$
DECLARE
    cat_id INTEGER;
    base_attrs INTEGER[];
BEGIN
    -- Базовые атрибуты для всех автомобильных категорий
    SELECT array_agg(id) INTO base_attrs 
    FROM unified_attributes 
    WHERE code IN ('auto_part_brand', 'auto_part_oem', 'auto_part_condition', 'auto_compatibility', 'auto_warranty');
    
    -- Добавляем базовые атрибуты ко всем автомобильным категориям без атрибутов
    FOR cat_id IN 
        SELECT c.id 
        FROM marketplace_categories c
        LEFT JOIN unified_category_attributes uca ON c.id = uca.category_id
        WHERE c.slug LIKE '%auto%' 
           OR c.slug LIKE '%car%'
           OR c.slug LIKE '%tire%'
           OR c.slug LIKE '%rim%'
           OR c.slug LIKE '%engine%'
           OR c.slug LIKE '%brake%'
           OR c.slug LIKE '%suspension%'
           OR c.slug IN ('electrical-parts', 'body-parts', 'cooling-system', 'transmission-parts', 'interior-parts')
        GROUP BY c.id
        HAVING COUNT(uca.id) = 0
    LOOP
        -- Добавляем базовые атрибуты
        INSERT INTO unified_category_attributes (category_id, attribute_id, is_enabled, is_required, sort_order)
        SELECT cat_id, unnest(base_attrs), true, false, generate_series(10, 50, 10)
        ON CONFLICT (category_id, attribute_id) DO NOTHING;
    END LOOP;
END $$;

-- Добавление специфичных атрибутов для категорий шин
INSERT INTO unified_category_attributes (category_id, attribute_id, is_enabled, is_required, sort_order)
SELECT c.id, ua.id, true, (ua.code = 'tire_season'), ua.sort_order
FROM marketplace_categories c
CROSS JOIN unified_attributes ua
WHERE c.slug IN ('summer-tires', 'winter-tires', 'all-season-tires', 
                  'passenger-summer-tires', 'passenger-winter-tires',
                  'truck-summer-tires', 'truck-winter-tires', 
                  'suv-summer-tires', 'suv-winter-tires')
  AND ua.code IN ('tire_width', 'tire_profile', 'tire_diameter', 'tire_season', 'tire_speed_index', 'tire_load_index')
ON CONFLICT (category_id, attribute_id) DO NOTHING;

-- Добавление специфичных атрибутов для категорий дисков
INSERT INTO unified_category_attributes (category_id, attribute_id, is_enabled, is_required, sort_order)
SELECT c.id, ua.id, true, false, ua.sort_order
FROM marketplace_categories c
CROSS JOIN unified_attributes ua
WHERE c.slug IN ('rims', 'aluminum-rims', 'steel-rims', 'sport-rims', 'complete-wheels')
  AND ua.code IN ('rim_diameter', 'rim_width', 'rim_bolt_pattern', 'rim_offset', 'rim_center_bore')
ON CONFLICT (category_id, attribute_id) DO NOTHING;

-- Добавление специфичных атрибутов для категорий двигателя
INSERT INTO unified_category_attributes (category_id, attribute_id, is_enabled, is_required, sort_order)
SELECT c.id, ua.id, true, false, ua.sort_order
FROM marketplace_categories c
CROSS JOIN unified_attributes ua
WHERE c.slug IN ('engine-and-parts', 'filters', 'belts-and-chains', 'spark-plugs', 'oils-and-fluids')
  AND ua.code IN ('engine_volume', 'engine_power', 'engine_type')
ON CONFLICT (category_id, attribute_id) DO NOTHING;

-- Добавление атрибутов для коммерческого транспорта
INSERT INTO unified_category_attributes (category_id, attribute_id, is_enabled, is_required, sort_order)
SELECT c.id, ua.id, true, false, ua.sort_order
FROM marketplace_categories c
CROSS JOIN unified_attributes ua
WHERE c.slug IN ('kamioni', 'autobusi', 'furgoni', 'prikolice', 'specijalna-vozila', 'komercijalna-vozila')
  AND ua.code IN ('vehicle_capacity', 'vehicle_seats', 'vehicle_axles', 'auto_year_from', 'auto_year_to')
ON CONFLICT (category_id, attribute_id) DO NOTHING;

-- Добавление атрибутов для водного транспорта
INSERT INTO unified_category_attributes (category_id, attribute_id, is_enabled, is_required, sort_order)
SELECT c.id, ua.id, true, false, ua.sort_order
FROM marketplace_categories c
CROSS JOIN unified_attributes ua
WHERE c.slug IN ('vodni-transport', 'camci', 'jahte', 'jet-ski', 'motori-za-camce', 'prikolice-za-camce')
  AND ua.code IN ('auto_part_brand', 'auto_part_condition', 'auto_year_from', 'auto_year_to', 'engine_power')
ON CONFLICT (category_id, attribute_id) DO NOTHING;

-- Добавление атрибутов для альтернативного транспорта
INSERT INTO unified_category_attributes (category_id, attribute_id, is_enabled, is_required, sort_order)
SELECT c.id, ua.id, true, false, ua.sort_order
FROM marketplace_categories c
CROSS JOIN unified_attributes ua
WHERE c.slug IN ('alternativni-transport', 'elektricni-bicikli', 'elektricni-skuteri', 'kvadovi', 'golf-vozila', 'motorne-sanke')
  AND ua.code IN ('auto_part_brand', 'auto_part_condition', 'engine_power', 'auto_warranty')
ON CONFLICT (category_id, attribute_id) DO NOTHING;

-- Добавление атрибутов для сельхозтехники
INSERT INTO unified_category_attributes (category_id, attribute_id, is_enabled, is_required, sort_order)
SELECT c.id, ua.id, true, false, ua.sort_order
FROM marketplace_categories c
CROSS JOIN unified_attributes ua
WHERE c.slug IN ('poljoprivredna-tehnika', 'traktori', 'kombajni', 'prikljucne-masine', 'ostala-poljoprivredna-tehnika')
  AND ua.code IN ('auto_part_brand', 'auto_part_condition', 'engine_power', 'vehicle_capacity', 'auto_year_from')
ON CONFLICT (category_id, attribute_id) DO NOTHING;

-- Добавление атрибутов для классических автомобилей
INSERT INTO unified_category_attributes (category_id, attribute_id, is_enabled, is_required, sort_order)
SELECT c.id, ua.id, true, false, ua.sort_order
FROM marketplace_categories c
CROSS JOIN unified_attributes ua
WHERE c.slug IN ('klasicna-vozila', 'oldtajmeri', 'youngtajmeri', 'kolekcijski-automobili', 'restaurirani-automobili', 'delovi-za-oldtajmere')
  AND ua.code IN ('auto_part_brand', 'auto_part_condition', 'auto_year_from', 'auto_year_to', 'auto_compatibility')
ON CONFLICT (category_id, attribute_id) DO NOTHING;

-- Обновление статистики
UPDATE marketplace_categories SET updated_at = NOW() 
WHERE slug LIKE '%auto%' OR slug LIKE '%car%' OR slug LIKE '%tire%' OR slug LIKE '%rim%' 
   OR slug LIKE '%engine%' OR slug LIKE '%brake%' OR slug LIKE '%suspension%'
   OR slug IN ('electrical-parts', 'body-parts', 'cooling-system', 'transmission-parts', 'interior-parts',
               'vodni-transport', 'camci', 'jahte', 'jet-ski', 'alternativni-transport', 
               'poljoprivredna-tehnika', 'klasicna-vozila');