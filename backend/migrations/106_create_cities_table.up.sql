-- Create table for cities (major cities like Belgrade, Novi Sad, Niš, etc.)
CREATE TABLE IF NOT EXISTS cities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    -- Latinized name for URL-safe slugs
    slug VARCHAR(255) NOT NULL UNIQUE,
    country_code VARCHAR(2) NOT NULL DEFAULT 'RS',
    -- PostGIS polygon for city boundaries
    boundary geometry(Polygon, 4326),
    -- Center point for viewport calculations
    center_point geometry(Point, 4326) NOT NULL,
    -- Additional metadata
    population INTEGER,
    area_km2 DECIMAL(10, 2),
    postal_codes TEXT[], -- Array of postal codes
    -- Enable/disable district search for this city
    has_districts BOOLEAN NOT NULL DEFAULT FALSE,
    -- Priority for city selection (lower = higher priority)
    priority INTEGER NOT NULL DEFAULT 100,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Add spatial indexes for viewport queries
CREATE INDEX idx_cities_boundary ON cities USING GIST (boundary);
CREATE INDEX idx_cities_center ON cities USING GIST (center_point);
CREATE INDEX idx_cities_name ON cities(name);
CREATE INDEX idx_cities_slug ON cities(slug);
CREATE INDEX idx_cities_country ON cities(country_code);
CREATE INDEX idx_cities_has_districts ON cities(has_districts);
CREATE INDEX idx_cities_priority ON cities(priority);

-- Add foreign key constraint to districts table
ALTER TABLE districts 
ADD CONSTRAINT fk_districts_city_id 
FOREIGN KEY (city_id) REFERENCES cities(id) ON DELETE SET NULL;

-- Create index for district-city relationship
CREATE INDEX idx_districts_city_id ON districts(city_id);

-- Insert major Serbian cities with their coordinates
INSERT INTO cities (name, slug, center_point, has_districts, priority) VALUES
('Београд', 'beograd', ST_SetSRID(ST_MakePoint(20.4612, 44.8186), 4326), TRUE, 1),
('Нови Сад', 'novi-sad', ST_SetSRID(ST_MakePoint(19.8335, 45.2671), 4326), FALSE, 2),
('Ниш', 'nis', ST_SetSRID(ST_MakePoint(21.8958, 43.3209), 4326), FALSE, 3),
('Крагујевац', 'kragujevac', ST_SetSRID(ST_MakePoint(20.9133, 44.0165), 4326), FALSE, 4),
('Суботица', 'subotica', ST_SetSRID(ST_MakePoint(19.6636, 46.1000), 4326), FALSE, 5),
('Панчево', 'pancevo', ST_SetSRID(ST_MakePoint(20.6402, 44.8705), 4326), FALSE, 6),
('Чачак', 'cacak', ST_SetSRID(ST_MakePoint(20.3497, 43.8914), 4326), FALSE, 7),
('Нови Пазар', 'novi-pazar', ST_SetSRID(ST_MakePoint(20.5226, 43.1364), 4326), FALSE, 8),
('Кикинда', 'kikinda', ST_SetSRID(ST_MakePoint(20.4577, 45.8397), 4326), FALSE, 9),
('Лесковац', 'leskovac', ST_SetSRID(ST_MakePoint(21.9447, 42.9987), 4326), FALSE, 10),
('Смедерево', 'smederevo', ST_SetSRID(ST_MakePoint(20.9300, 44.6631), 4326), FALSE, 11);

-- Update existing Belgrade districts to link them to Belgrade city
UPDATE districts 
SET city_id = (SELECT id FROM cities WHERE slug = 'beograd')
WHERE city_id IS NULL;

-- Update Belgrade city to have proper boundary based on its districts
UPDATE cities 
SET boundary = (
    SELECT ST_ConvexHull(ST_Collect(boundary))
    FROM districts 
    WHERE city_id = cities.id
),
has_districts = TRUE
WHERE slug = 'beograd' AND EXISTS (
    SELECT 1 FROM districts WHERE city_id = cities.id
);