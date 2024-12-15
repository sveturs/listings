-- backend/migrations/0031_add_car_rating.up.sql

ALTER TABLE cars
ADD COLUMN rating DECIMAL(3,2) DEFAULT 0,
ADD COLUMN reviews_count INT DEFAULT 0;

-- Создаем триггер для обновления рейтинга
CREATE OR REPLACE FUNCTION update_car_rating()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.entity_type = 'car' THEN
        UPDATE cars
        SET 
            rating = (
                SELECT COALESCE(AVG(rating)::DECIMAL(3,2), 0)
                FROM reviews 
                WHERE entity_type = 'car' 
                AND entity_id = NEW.entity_id
                AND status = 'published'
            ),
            reviews_count = (
                SELECT COUNT(*)
                FROM reviews 
                WHERE entity_type = 'car' 
                AND entity_id = NEW.entity_id
                AND status = 'published'
            )
        WHERE id = NEW.entity_id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER reviews_car_rating_update
AFTER INSERT OR UPDATE OR DELETE ON reviews
FOR EACH ROW
EXECUTE FUNCTION update_car_rating();