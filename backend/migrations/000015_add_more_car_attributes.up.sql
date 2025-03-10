-- Добавляем дополнительные атрибуты для автомобилей
INSERT INTO category_attributes (name, display_name, attribute_type, options, is_searchable, is_filterable, is_required) VALUES
('drive_type', 'Привод', 'select', '{"values": ["Передний", "Задний", "Полный", "Другой"]}', true, true, false),
('number_of_doors', 'Количество дверей', 'select', '{"values": ["2", "3", "4", "5", "6+"]}', true, true, false),
('number_of_seats', 'Количество мест', 'select', '{"values": ["1", "2", "3", "4", "5", "6", "7", "8+"]}', true, true, false);

-- Связываем атрибуты с категориями автомобилей
DO $$
DECLARE
    car_category_id INT;
BEGIN
    -- Получаем ID категории автомобилей
    SELECT id INTO car_category_id FROM marketplace_categories WHERE name = 'Автомобили' OR name = 'Cars' LIMIT 1;

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
        AND a.name IN ('drive_type', 'number_of_doors', 'number_of_seats')
        ON CONFLICT DO NOTHING;
    END IF;
END $$;

-- Добавляем переводы для новых атрибутов
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified, created_at, updated_at)
SELECT 
    'attribute', 
    id, 
    'en', 
    'display_name', 
    CASE 
        WHEN name = 'drive_type' THEN 'Drive Type'
        WHEN name = 'number_of_doors' THEN 'Number of Doors'
        WHEN name = 'number_of_seats' THEN 'Number of Seats'
    END, 
    false, 
    true, 
    NOW(), 
    NOW()
FROM category_attributes
WHERE name IN ('drive_type', 'number_of_doors', 'number_of_seats')
ON CONFLICT DO NOTHING;

-- Добавляем переводы на сербский
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified, created_at, updated_at)
SELECT 
    'attribute', 
    id, 
    'sr', 
    'display_name', 
    CASE 
        WHEN name = 'drive_type' THEN 'Pogon'
        WHEN name = 'number_of_doors' THEN 'Broj vrata'
        WHEN name = 'number_of_seats' THEN 'Broj sedišta'
    END, 
    false, 
    true, 
    NOW(), 
    NOW()
FROM category_attributes
WHERE name IN ('drive_type', 'number_of_doors', 'number_of_seats')
ON CONFLICT DO NOTHING;