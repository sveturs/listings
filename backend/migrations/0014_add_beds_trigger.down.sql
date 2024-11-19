-- Удаляем триггер
DROP TRIGGER IF EXISTS bed_bookings_update_trigger ON bed_bookings;

-- Удаляем функцию
DROP FUNCTION IF EXISTS update_available_beds();