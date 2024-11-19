-- Удаляем старый триггер и функцию
DROP TRIGGER IF EXISTS bed_bookings_update_trigger ON bed_bookings;
DROP FUNCTION IF EXISTS update_available_beds();

-- Создаем улучшенную функцию обновления доступных мест
CREATE OR REPLACE FUNCTION update_available_beds()
RETURNS TRIGGER AS $$
BEGIN
    -- Обновляем количество доступных мест для комнаты, связанной с бронированием
    WITH bed_count AS (
        SELECT b.room_id,
               COUNT(DISTINCT bb.bed_id) as booked_beds
        FROM beds b
        LEFT JOIN bed_bookings bb ON b.id = bb.bed_id
        WHERE bb.status = 'confirmed'
        AND bb.start_date <= CURRENT_DATE 
        AND bb.end_date >= CURRENT_DATE
        GROUP BY b.room_id
    )
    UPDATE rooms r
    SET available_beds = r.total_beds - COALESCE(bc.booked_beds, 0)
    FROM bed_count bc
    WHERE r.id = bc.room_id
    AND r.accommodation_type = 'bed';

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Создаем новый триггер
CREATE TRIGGER bed_bookings_update_trigger
AFTER INSERT OR UPDATE OR DELETE ON bed_bookings
FOR EACH ROW
EXECUTE FUNCTION update_available_beds();

-- Создаем функцию для получения количества доступных мест на конкретные даты
CREATE OR REPLACE FUNCTION get_available_beds(
    room_id_param INTEGER,
    start_date_param DATE,
    end_date_param DATE
)
RETURNS INTEGER AS $$
DECLARE
    total_beds INTEGER;
    booked_beds INTEGER;
BEGIN
    -- Получаем общее количество кроватей
    SELECT total_beds INTO total_beds
    FROM rooms
    WHERE id = room_id_param;

    -- Считаем количество забронированных кроватей на указанные даты
    SELECT COUNT(DISTINCT bb.bed_id) INTO booked_beds
    FROM bed_bookings bb
    JOIN beds b ON bb.bed_id = b.id
    WHERE b.room_id = room_id_param
    AND bb.status = 'confirmed'
    AND bb.start_date <= end_date_param
    AND bb.end_date >= start_date_param;

    -- Возвращаем разницу
    RETURN total_beds - COALESCE(booked_beds, 0);
END;
$$ LANGUAGE plpgsql;

-- Функция инициализации доступных мест
CREATE OR REPLACE FUNCTION initialize_available_beds()
RETURNS void AS $$
BEGIN
    WITH bed_count AS (
        SELECT b.room_id,
               COUNT(DISTINCT bb.bed_id) as booked_beds
        FROM beds b
        LEFT JOIN bed_bookings bb ON b.id = bb.bed_id
        WHERE bb.status = 'confirmed'
        AND bb.start_date <= CURRENT_DATE 
        AND bb.end_date >= CURRENT_DATE
        GROUP BY b.room_id
    )
    UPDATE rooms r
    SET available_beds = r.total_beds - COALESCE(bc.booked_beds, 0)
    FROM bed_count bc
    WHERE r.id = bc.room_id
    AND r.accommodation_type = 'bed';
END;
$$ LANGUAGE plpgsql;

-- Выполняем начальную инициализацию
SELECT initialize_available_beds();