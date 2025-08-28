-- Возвращаем флаг is_required для марок автомобилей

UPDATE category_attributes
SET is_required = true
WHERE name IN ('car_make', 'motorcycle_make', 'truck_make', 'boat_make');