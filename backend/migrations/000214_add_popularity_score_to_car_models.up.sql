-- Добавляем поле popularity_score для отслеживания популярности моделей в Сербии
ALTER TABLE car_models
ADD COLUMN IF NOT EXISTS popularity_score DECIMAL(5,2) DEFAULT 0;