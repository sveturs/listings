-- Название файла: backend/migrations/0016_fix_bed_availability.up.sql

-- Удаляем старую версию функции если она существует
DROP FUNCTION IF EXISTS update_available_beds CASCADE;

-- Создаем новую функцию
CREATE OR REPLACE FUNCTION update_available_beds()
RETURNS TRIGGER AS $$
BEGIN
    WITH bed_count AS (
        SELECT 
            b.room_id,
            COUNT(DISTINCT b.id) as total_beds,
            COUNT(DISTINCT CASE 
                WHEN NOT EXISTS (
                    SELECT 1 
                    FROM bed_bookings bb 
                    WHERE bb.bed_id = b.id 
                    AND bb.status = 'confirmed'
                    AND (
                        (bb.start_date <= CURRENT_DATE AND bb.end_date >= CURRENT_DATE) OR
                        (bb.start_date >= CURRENT_DATE AND bb.start_date <= CURRENT_DATE)
                    )
                ) THEN b.id 
            END) as available_beds
        FROM beds b
        GROUP BY b.room_id
    )
    UPDATE rooms r
    SET 
        total_beds = bc.total_beds,
        available_beds = bc.available_beds
    FROM bed_count bc
    WHERE r.id = bc.room_id
    AND r.accommodation_type = 'bed';

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Пересоздаем триггер
DROP TRIGGER IF EXISTS bed_bookings_update_trigger ON bed_bookings;
CREATE TRIGGER bed_bookings_update_trigger
    AFTER INSERT OR UPDATE OR DELETE ON bed_bookings
    FOR EACH ROW
    EXECUTE FUNCTION update_available_beds();

-- Создаем функцию для получения количества доступных мест на определенные даты
CREATE OR REPLACE FUNCTION get_beds_availability(
    p_room_id INTEGER,
    p_start_date DATE,
    p_end_date DATE
) RETURNS INTEGER AS $$
BEGIN
    RETURN (
        SELECT COUNT(DISTINCT b.id)
        FROM beds b
        WHERE b.room_id = p_room_id
        AND b.is_available = true
        AND NOT EXISTS (
            SELECT 1
            FROM bed_bookings bb
            WHERE bb.bed_id = b.id
            AND bb.status = 'confirmed'
            AND (
                (bb.start_date <= p_end_date AND bb.end_date >= p_start_date) OR
                (bb.start_date >= p_start_date AND bb.start_date <= p_end_date)
            )
        )
    );
END;
$$ LANGUAGE plpgsql;