-- Create car_makes table for storing car manufacturers
CREATE TABLE car_makes (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    logo_url VARCHAR(500),
    country VARCHAR(50),
    is_active BOOLEAN DEFAULT true,
    sort_order INTEGER DEFAULT 0,
    is_domestic BOOLEAN DEFAULT false, -- for Serbian brands
    popularity_rs INTEGER DEFAULT 0, -- popularity in Serbia
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create car_models table for storing car models
CREATE TABLE car_models (
    id SERIAL PRIMARY KEY,
    make_id INTEGER REFERENCES car_makes(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL,
    generation VARCHAR(50),
    production_start INTEGER,
    production_end INTEGER,
    body_types JSONB, -- ["sedan", "hatchback", "wagon"]
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(make_id, slug)
);

-- Create car_generations table for storing model generations/facelifts
CREATE TABLE car_generations (
    id SERIAL PRIMARY KEY,
    model_id INTEGER REFERENCES car_models(id) ON DELETE CASCADE,
    name VARCHAR(100),
    year_start INTEGER,
    year_end INTEGER,
    facelift_year INTEGER,
    specs JSONB, -- technical specifications
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for fast search
CREATE INDEX idx_car_makes_slug ON car_makes(slug);
CREATE INDEX idx_car_makes_popularity ON car_makes(popularity_rs DESC);
CREATE INDEX idx_car_makes_is_domestic ON car_makes(is_domestic);
CREATE INDEX idx_car_models_make_slug ON car_models(make_id, slug);
CREATE INDEX idx_car_models_search ON car_models USING gin(to_tsvector('simple', name));
CREATE INDEX idx_car_generations_model_id ON car_generations(model_id);

-- Insert popular car makes in Serbia
INSERT INTO car_makes (name, slug, country, is_domestic, popularity_rs, sort_order) VALUES
-- Serbian domestic brands
('Zastava', 'zastava', 'Serbia', true, 100, 1),
('Yugo', 'yugo', 'Serbia', true, 95, 2),
('FAP', 'fap', 'Serbia', true, 80, 3),
('IMT', 'imt', 'Serbia', true, 75, 4),

-- Most popular foreign brands in Serbia
('Volkswagen', 'volkswagen', 'Germany', false, 90, 5),
('Fiat', 'fiat', 'Italy', false, 88, 6),
('Opel', 'opel', 'Germany', false, 85, 7),
('Peugeot', 'peugeot', 'France', false, 83, 8),
('Renault', 'renault', 'France', false, 82, 9),
('Škoda', 'skoda', 'Czech Republic', false, 81, 10),
('Ford', 'ford', 'USA', false, 80, 11),
('Citroën', 'citroen', 'France', false, 78, 12),
('Mercedes-Benz', 'mercedes-benz', 'Germany', false, 77, 13),
('BMW', 'bmw', 'Germany', false, 76, 14),
('Audi', 'audi', 'Germany', false, 75, 15),
('Toyota', 'toyota', 'Japan', false, 70, 16),
('Hyundai', 'hyundai', 'South Korea', false, 68, 17),
('Mazda', 'mazda', 'Japan', false, 65, 18),
('Nissan', 'nissan', 'Japan', false, 63, 19),
('Honda', 'honda', 'Japan', false, 60, 20),
('Kia', 'kia', 'South Korea', false, 58, 21),
('Seat', 'seat', 'Spain', false, 55, 22),
('Suzuki', 'suzuki', 'Japan', false, 53, 23),
('Mitsubishi', 'mitsubishi', 'Japan', false, 50, 24),
('Chevrolet', 'chevrolet', 'USA', false, 48, 25),
('Dacia', 'dacia', 'Romania', false, 45, 26),
('Volvo', 'volvo', 'Sweden', false, 43, 27),
('Alfa Romeo', 'alfa-romeo', 'Italy', false, 40, 28),
('Lada', 'lada', 'Russia', false, 38, 29),
('Lancia', 'lancia', 'Italy', false, 35, 30);

-- Insert some popular models for Serbian brands
INSERT INTO car_models (make_id, name, slug, production_start, production_end, body_types) VALUES
-- Zastava models
((SELECT id FROM car_makes WHERE slug = 'zastava'), '750 (Fića)', '750', 1955, 1985, '["sedan"]'),
((SELECT id FROM car_makes WHERE slug = 'zastava'), '101', '101', 1971, 2008, '["sedan", "hatchback"]'),
((SELECT id FROM car_makes WHERE slug = 'zastava'), '128', '128', 1988, 2008, '["sedan", "hatchback"]'),
((SELECT id FROM car_makes WHERE slug = 'zastava'), 'Koral (Yugo)', 'koral', 1980, 2008, '["hatchback"]'),
((SELECT id FROM car_makes WHERE slug = 'zastava'), 'Florida', 'florida', 1988, 2008, '["sedan", "hatchback"]'),
((SELECT id FROM car_makes WHERE slug = 'zastava'), '10', '10', 2008, 2008, '["hatchback"]'),

-- Yugo models (separate brand entry for classics)
((SELECT id FROM car_makes WHERE slug = 'yugo'), '45', '45', 1980, 1992, '["hatchback"]'),
((SELECT id FROM car_makes WHERE slug = 'yugo'), '55', '55', 1980, 1992, '["hatchback"]'),
((SELECT id FROM car_makes WHERE slug = 'yugo'), '65', '65', 1985, 1992, '["hatchback"]'),
((SELECT id FROM car_makes WHERE slug = 'yugo'), 'Koral', 'koral', 1988, 2008, '["hatchback"]'),
((SELECT id FROM car_makes WHERE slug = 'yugo'), 'Tempo', 'tempo', 1990, 1998, '["hatchback"]'),
((SELECT id FROM car_makes WHERE slug = 'yugo'), 'Sana', 'sana', 1989, 1995, '["hatchback"]'),

-- FAP trucks
((SELECT id FROM car_makes WHERE slug = 'fap'), '1314', '1314', 1975, NULL, '["truck"]'),
((SELECT id FROM car_makes WHERE slug = 'fap'), '1620', '1620', 1980, NULL, '["truck"]'),
((SELECT id FROM car_makes WHERE slug = 'fap'), '2023', '2023', 1985, NULL, '["truck"]'),

-- IMT tractors
((SELECT id FROM car_makes WHERE slug = 'imt'), '533', '533', 1960, NULL, '["tractor"]'),
((SELECT id FROM car_makes WHERE slug = 'imt'), '539', '539', 1975, NULL, '["tractor"]'),
((SELECT id FROM car_makes WHERE slug = 'imt'), '577', '577', 1985, NULL, '["tractor"]');

-- Add some popular foreign models (most common in Serbia)
INSERT INTO car_models (make_id, name, slug, production_start, body_types) VALUES
-- Volkswagen
((SELECT id FROM car_makes WHERE slug = 'volkswagen'), 'Golf', 'golf', 1974, '["hatchback", "wagon"]'),
((SELECT id FROM car_makes WHERE slug = 'volkswagen'), 'Passat', 'passat', 1973, '["sedan", "wagon"]'),
((SELECT id FROM car_makes WHERE slug = 'volkswagen'), 'Polo', 'polo', 1975, '["hatchback", "sedan"]'),
((SELECT id FROM car_makes WHERE slug = 'volkswagen'), 'Jetta', 'jetta', 1979, '["sedan"]'),

-- Fiat
((SELECT id FROM car_makes WHERE slug = 'fiat'), 'Punto', 'punto', 1993, '["hatchback"]'),
((SELECT id FROM car_makes WHERE slug = 'fiat'), 'Stilo', 'stilo', 2001, '["hatchback", "wagon"]'),
((SELECT id FROM car_makes WHERE slug = 'fiat'), 'Bravo', 'bravo', 1995, '["hatchback"]'),
((SELECT id FROM car_makes WHERE slug = 'fiat'), '500', '500', 2007, '["hatchback"]'),

-- Opel
((SELECT id FROM car_makes WHERE slug = 'opel'), 'Astra', 'astra', 1991, '["hatchback", "sedan", "wagon"]'),
((SELECT id FROM car_makes WHERE slug = 'opel'), 'Corsa', 'corsa', 1982, '["hatchback"]'),
((SELECT id FROM car_makes WHERE slug = 'opel'), 'Vectra', 'vectra', 1988, '["sedan", "hatchback", "wagon"]'),
((SELECT id FROM car_makes WHERE slug = 'opel'), 'Insignia', 'insignia', 2008, '["sedan", "hatchback", "wagon"]'),

-- Peugeot
((SELECT id FROM car_makes WHERE slug = 'peugeot'), '206', '206', 1998, '["hatchback", "sedan", "wagon"]'),
((SELECT id FROM car_makes WHERE slug = 'peugeot'), '307', '307', 2001, '["hatchback", "wagon"]'),
((SELECT id FROM car_makes WHERE slug = 'peugeot'), '308', '308', 2007, '["hatchback", "wagon"]'),
((SELECT id FROM car_makes WHERE slug = 'peugeot'), '407', '407', 2004, '["sedan", "wagon"]'),

-- Renault
((SELECT id FROM car_makes WHERE slug = 'renault'), 'Clio', 'clio', 1990, '["hatchback"]'),
((SELECT id FROM car_makes WHERE slug = 'renault'), 'Megane', 'megane', 1995, '["hatchback", "sedan", "wagon"]'),
((SELECT id FROM car_makes WHERE slug = 'renault'), 'Laguna', 'laguna', 1993, '["hatchback", "wagon"]'),
((SELECT id FROM car_makes WHERE slug = 'renault'), 'Scenic', 'scenic', 1996, '["minivan"]'),

-- Škoda
((SELECT id FROM car_makes WHERE slug = 'skoda'), 'Octavia', 'octavia', 1996, '["sedan", "wagon"]'),
((SELECT id FROM car_makes WHERE slug = 'skoda'), 'Fabia', 'fabia', 1999, '["hatchback", "wagon"]'),
((SELECT id FROM car_makes WHERE slug = 'skoda'), 'Superb', 'superb', 2001, '["sedan", "wagon"]');

-- Create triggers for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_car_makes_updated_at BEFORE UPDATE ON car_makes
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_car_models_updated_at BEFORE UPDATE ON car_models
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_car_generations_updated_at BEFORE UPDATE ON car_generations
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();