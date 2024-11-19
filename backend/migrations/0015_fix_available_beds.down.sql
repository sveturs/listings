-- Удаляем созданные функции и триггеры
DROP TRIGGER IF EXISTS bed_bookings_update_trigger ON bed_bookings;
DROP FUNCTION IF EXISTS update_available_beds();
DROP FUNCTION IF EXISTS get_available_beds(INTEGER, DATE, DATE);
DROP FUNCTION IF EXISTS initialize_available_beds();

-- Возвращаем простой триггер для обновления available_beds
CREATE OR REPLACE FUNCTION update_available_beds()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE rooms
    SET available_beds = (
        SELECT COUNT(*)
        FROM beds b
        WHERE b.room_id = NEW.room_id
        AND b.is_available = true
        AND b.id NOT IN (
            SELECT bed_id
            FROM bed_bookings
            WHERE status = 'confirmed'
            AND (
                (start_date <= CURRENT_DATE AND end_date >= CURRENT_DATE)
            )
        )
    )
    WHERE id = NEW.room_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER bed_bookings_update_trigger
AFTER INSERT OR UPDATE OR DELETE ON bed_bookings
FOR EACH ROW
EXECUTE FUNCTION update_available_beds();