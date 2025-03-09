-- backend/migrations/000009_add_auto_properties.up.sql
-- Создаем таблицу для хранения дополнительных свойств автомобилей
CREATE TABLE IF NOT EXISTS auto_properties (
    listing_id INT PRIMARY KEY REFERENCES marketplace_listings(id) ON DELETE CASCADE,
    brand VARCHAR(100),
    model VARCHAR(100),
    year INT CHECK (year BETWEEN 1900 AND 2100),
    mileage INT,
    fuel_type VARCHAR(50), -- бензин, дизель, электро, гибрид и т.д.
    transmission VARCHAR(50), -- механика, автомат, вариатор и т.д.
    engine_capacity DECIMAL(4,1), -- объем двигателя в литрах
    power INT, -- мощность в л.с.
    color VARCHAR(50),
    body_type VARCHAR(50), -- седан, хэтчбек, внедорожник и т.д.
    drive_type VARCHAR(50), -- передний, задний, полный
    number_of_doors INT,
    number_of_seats INT,
    additional_features TEXT, -- дополнительные особенности
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создаем индексы для оптимизации запросов
CREATE INDEX IF NOT EXISTS idx_auto_properties_brand ON auto_properties(brand);
CREATE INDEX IF NOT EXISTS idx_auto_properties_model ON auto_properties(model);
CREATE INDEX IF NOT EXISTS idx_auto_properties_year ON auto_properties(year);
CREATE INDEX IF NOT EXISTS idx_auto_properties_fuel_type ON auto_properties(fuel_type);
CREATE INDEX IF NOT EXISTS idx_auto_properties_transmission ON auto_properties(transmission);
CREATE INDEX IF NOT EXISTS idx_auto_properties_body_type ON auto_properties(body_type);

-- Функция для автоматического обновления даты изменения записи
CREATE OR REPLACE FUNCTION update_auto_properties_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Триггер для обновления даты изменения записи
CREATE TRIGGER update_auto_properties_timestamp
    BEFORE UPDATE ON auto_properties
    FOR EACH ROW
    EXECUTE FUNCTION update_auto_properties_updated_at();

-- Создаем представление для объединения основных данных объявлений с данными автомобилей
CREATE OR REPLACE VIEW auto_listings_view AS
SELECT 
    l.*,
    ap.brand,
    ap.model,
    ap.year,
    ap.mileage,
    ap.fuel_type,
    ap.transmission,
    ap.engine_capacity,
    ap.power,
    ap.color,
    ap.body_type,
    ap.drive_type,
    ap.number_of_doors,
    ap.number_of_seats,
    ap.additional_features
FROM 
    marketplace_listings l
LEFT JOIN 
    auto_properties ap ON l.id = ap.listing_id
WHERE 
    l.category_id IN (
        SELECT id FROM marketplace_categories 
        WHERE parent_id = 2000 
        OR id = 2000 
        OR parent_id IN (SELECT id FROM marketplace_categories WHERE parent_id = 2000)
    );