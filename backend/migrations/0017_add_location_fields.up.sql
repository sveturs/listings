-- Добавляем поля для координат и форматированного адреса
ALTER TABLE rooms
    ADD COLUMN latitude DECIMAL(10, 8),
    ADD COLUMN longitude DECIMAL(11, 8),
    ADD COLUMN formatted_address TEXT;

-- Индекс для географического поиска
CREATE INDEX idx_rooms_location ON rooms(latitude, longitude);

-- Функция для валидации координат
CREATE OR REPLACE FUNCTION validate_coordinates()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.latitude IS NOT NULL AND (NEW.latitude < -90 OR NEW.latitude > 90) THEN
        RAISE EXCEPTION 'Широта должна быть между -90 и 90';
    END IF;
    IF NEW.longitude IS NOT NULL AND (NEW.longitude < -180 OR NEW.longitude > 180) THEN
        RAISE EXCEPTION 'Долгота должна быть между -180 и 180';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Триггер для валидации координат
CREATE TRIGGER validate_coordinates_trigger
    BEFORE INSERT OR UPDATE ON rooms
    FOR EACH ROW
    EXECUTE FUNCTION validate_coordinates();