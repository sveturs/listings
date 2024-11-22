CREATE OR REPLACE FUNCTION update_available_beds()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE rooms r
    SET available_beds = (
        SELECT COUNT(*)
        FROM beds b
        WHERE b.room_id = NEW.room_id
        AND b.is_available = true
        AND b.id NOT IN (
            SELECT bed_id
            FROM bed_bookings
            WHERE status = 'confirmed'
        )
    )
    WHERE id = NEW.room_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;