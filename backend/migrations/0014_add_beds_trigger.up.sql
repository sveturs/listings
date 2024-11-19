-- Создаем функцию для подсчета доступных кроватей
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

-- Создаем триггер для таблицы bed_bookings
CREATE TRIGGER bed_bookings_update_trigger
AFTER INSERT OR UPDATE OR DELETE ON bed_bookings
FOR EACH ROW
EXECUTE FUNCTION update_available_beds();