-- Удаляем триггер и функцию
DROP TRIGGER IF EXISTS validate_coordinates_trigger ON rooms;
DROP FUNCTION IF EXISTS validate_coordinates;

-- Удаляем индекс
DROP INDEX IF EXISTS idx_rooms_location;

-- Удаляем колонки
ALTER TABLE rooms
    DROP COLUMN IF EXISTS latitude,
    DROP COLUMN IF EXISTS longitude,
    DROP COLUMN IF EXISTS formatted_address;