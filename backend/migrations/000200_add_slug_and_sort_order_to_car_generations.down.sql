-- Удаляем добавленные колонки
ALTER TABLE car_generations 
DROP COLUMN IF EXISTS slug,
DROP COLUMN IF EXISTS sort_order;