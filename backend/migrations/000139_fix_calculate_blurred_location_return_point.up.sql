-- Исправление функции calculate_blurred_location для возврата Point вместо Polygon

CREATE OR REPLACE FUNCTION calculate_blurred_location(
    exact_location geography,
    privacy_level location_privacy_level
) RETURNS geography AS $$
BEGIN
    CASE privacy_level
        WHEN 'exact' THEN
            RETURN exact_location;
        WHEN 'street' THEN
            -- Смещаем точку на случайное расстояние в пределах 200м
            RETURN ST_Project(
                exact_location,
                200 + (random() * 100 - 50), -- 200м ± 50м
                radians(random() * 360)       -- случайное направление
            )::geography;
        WHEN 'district' THEN
            -- Смещаем точку на случайное расстояние в пределах 1км
            RETURN ST_Project(
                exact_location,
                1000 + (random() * 200 - 100), -- 1км ± 100м
                radians(random() * 360)        -- случайное направление
            )::geography;
        WHEN 'city' THEN
            -- Смещаем точку на случайное расстояние в пределах 5км
            RETURN ST_Project(
                exact_location,
                5000 + (random() * 1000 - 500), -- 5км ± 500м
                radians(random() * 360)         -- случайное направление
            )::geography;
        ELSE
            RETURN NULL;
    END CASE;
END;
$$ LANGUAGE plpgsql IMMUTABLE;