-- Возврат к предыдущей версии функции calculate_blurred_location с Buffer

CREATE OR REPLACE FUNCTION calculate_blurred_location(
    exact_location geography,
    privacy_level location_privacy_level
) RETURNS geography AS $$
BEGIN
    CASE privacy_level
        WHEN 'exact' THEN
            RETURN exact_location;
        WHEN 'street' THEN
            -- Размываем до ~200 метров
            RETURN ST_Buffer(
                ST_Centroid(exact_location::geometry)::geography,
                200 + (random() * 100 - 50) -- 200м ± 50м
            );
        WHEN 'district' THEN
            -- Размываем до ~1000 метров
            RETURN ST_Buffer(
                ST_Centroid(exact_location::geometry)::geography,
                1000 + (random() * 200 - 100) -- 1км ± 100м
            );
        WHEN 'city' THEN
            -- Размываем до ~5000 метров
            RETURN ST_Buffer(
                ST_Centroid(exact_location::geometry)::geography,
                5000 + (random() * 1000 - 500) -- 5км ± 500м
            );
        ELSE
            RETURN NULL;
    END CASE;
END;
$$ LANGUAGE plpgsql IMMUTABLE;