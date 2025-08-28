-- Исправляем конфликты is_required между таблицами category_attributes и category_attribute_mapping

-- 1. Condition не должен быть обязательным (это дополнительная информация)
UPDATE category_attribute_mapping
SET is_required = false
WHERE attribute_id = (SELECT id FROM category_attributes WHERE name = 'condition')
AND is_required = true;

-- 2. Brand не должен быть обязательным (может быть неизвестен или "другой")
UPDATE category_attribute_mapping  
SET is_required = false
WHERE attribute_id IN (
    SELECT id FROM category_attributes 
    WHERE name IN ('brand', 'car_make', 'motorcycle_make', 'truck_make', 'boat_make')
)
AND is_required = true;

-- 3. Fuel type и transmission должны быть обязательными для транспортных средств
UPDATE category_attribute_mapping
SET is_required = true
WHERE attribute_id IN (
    SELECT id FROM category_attributes 
    WHERE name IN ('fuel_type', 'transmission')
)
AND category_id IN (1302, 1303, 1304) -- Мотоциклы, Авто делови, Комерцијална возила
AND is_required = false;

-- 4. Year должен быть обязательным для сельхозтехники
UPDATE category_attribute_mapping
SET is_required = true
WHERE attribute_id = (SELECT id FROM category_attributes WHERE name = 'year')
AND category_id = 1601 -- Польопривредне машине
AND is_required = false;

-- Логирование изменений
DO $$
BEGIN
    RAISE NOTICE 'Fixed is_required conflicts between category_attributes and category_attribute_mapping tables';
END $$;