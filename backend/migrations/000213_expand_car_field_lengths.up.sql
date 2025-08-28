-- Расширяем поля для более длинных значений
ALTER TABLE car_models 
ALTER COLUMN drive_type TYPE VARCHAR(50),
ALTER COLUMN segment TYPE VARCHAR(50),
ALTER COLUMN electric_range_standard TYPE VARCHAR(50);