-- Добавляем недостающие колонки slug и sort_order в таблицу car_generations
ALTER TABLE car_generations 
ADD COLUMN slug VARCHAR(100),
ADD COLUMN sort_order INTEGER DEFAULT 0;

-- Создаем индекс для slug
CREATE INDEX idx_car_generations_slug ON car_generations(slug);

-- Обновляем существующие записи, генерируя slug из name
UPDATE car_generations 
SET slug = LOWER(REPLACE(REPLACE(REPLACE(name, ' ', '-'), '(', ''), ')', '')),
    sort_order = CASE 
        WHEN year_end IS NULL THEN 1  -- Текущие поколения первые
        ELSE 1000 - year_start  -- Остальные сортируются по году начала (новые первые)
    END
WHERE slug IS NULL;

-- Делаем slug обязательным после заполнения
ALTER TABLE car_generations 
ALTER COLUMN slug SET NOT NULL;