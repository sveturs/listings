-- Удаляем поле popularity_score
ALTER TABLE car_models
DROP COLUMN IF EXISTS popularity_score;