-- Связываем атрибут мощности с категориями автомобилей
DO $$
DECLARE
    car_category_id INT;
    power_attr_id INT;
BEGIN
    -- Получаем ID категории автомобилей
    SELECT id INTO car_category_id FROM marketplace_categories WHERE name = 'Автомобили' OR name = 'Cars' LIMIT 1;
    
    -- Получаем ID атрибута мощности
    SELECT id INTO power_attr_id FROM category_attributes WHERE name = 'power';
    
    -- Добавляем связь
    IF car_category_id IS NOT NULL AND power_attr_id IS NOT NULL THEN
        -- Добавляем для основной категории
        INSERT INTO category_attribute_mapping (category_id, attribute_id, is_enabled, is_required)
        VALUES (car_category_id, power_attr_id, true, false)
        ON CONFLICT DO NOTHING;
        
        -- Добавляем для подкатегорий
        INSERT INTO category_attribute_mapping (category_id, attribute_id, is_enabled, is_required)
        SELECT c.id, power_attr_id, true, false
        FROM marketplace_categories c 
        WHERE c.parent_id = car_category_id OR EXISTS (
            SELECT 1 FROM marketplace_categories c2 
            WHERE c2.parent_id = car_category_id AND c.parent_id = c2.id
        )
        ON CONFLICT DO NOTHING;
        
        RAISE NOTICE 'Linked power attribute to car categories: category_id=%, attribute_id=%', car_category_id, power_attr_id;
    ELSE
        RAISE NOTICE 'Could not find car category or power attribute: car_category_id=%, power_attr_id=%', car_category_id, power_attr_id;
    END IF;
END $$;


