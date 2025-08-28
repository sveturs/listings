-- Исправляем флаг is_required для марок автомобилей в таблице category_attributes
-- Марка не должна быть обязательной, так как может быть "Другая" или неизвестна

UPDATE category_attributes
SET is_required = false
WHERE name IN ('car_make', 'motorcycle_make', 'truck_make', 'boat_make')
AND is_required = true;

-- Логирование
DO $$
BEGIN
    RAISE NOTICE 'Fixed is_required flag for vehicle make attributes in category_attributes table';
END $$;