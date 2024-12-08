ALTER TABLE cars 
ADD COLUMN transmission VARCHAR(20),
ADD COLUMN fuel_type VARCHAR(20),
ADD COLUMN seats INT DEFAULT 4,
ADD COLUMN features TEXT[];

-- Добавляем ограничения
ALTER TABLE cars 
ADD CONSTRAINT check_transmission CHECK (transmission IN ('automatic', 'manual')),
ADD CONSTRAINT check_fuel_type CHECK (fuel_type IN ('petrol', 'diesel', 'electric', 'hybrid'));