-- Возвращаем старые размеры полей
ALTER TABLE car_models 
ALTER COLUMN drive_type TYPE VARCHAR(20),
ALTER COLUMN segment TYPE VARCHAR(20),
ALTER COLUMN electric_range_standard TYPE VARCHAR(20);