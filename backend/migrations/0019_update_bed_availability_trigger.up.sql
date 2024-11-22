-- Обновляем функцию подсчета доступных мест
CREATE OR REPLACE FUNCTION update_available_beds()
RETURNS TRIGGER AS $$
BEGIN
    -- Обновляем количество доступных мест для комнаты с учетом текущих броней
    WITH current_bookings AS (
        SELECT b.room_id,
               COUNT(DISTINCT bb.bed_id) as booked_beds
        FROM beds b
        JOIN bed_bookings bb ON b.id = bb.bed_id
        WHERE bb.status = 'confirmed'
        AND bb.start_date <= CURRENT_DATE
        AND bb.end_date >= CURRENT_DATE
        GROUP BY b.room_id
    )
    UPDATE rooms r
    SET available_beds = COALESCE(r.total_beds, 0) - COALESCE(cb.booked_beds, 0)
    FROM current_bookings cb
    WHERE r.id = cb.room_id
    AND r.accommodation_type = 'bed';

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;