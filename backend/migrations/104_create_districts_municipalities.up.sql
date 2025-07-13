-- Enable PostGIS if not already enabled
CREATE EXTENSION IF NOT EXISTS postgis;

-- Create table for districts (city districts like Belgrade districts)
CREATE TABLE IF NOT EXISTS districts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    city_id UUID, -- Can be NULL for top-level districts
    country_code VARCHAR(2) NOT NULL DEFAULT 'RS',
    -- PostGIS polygon for district boundaries
    boundary geometry(Polygon, 4326),
    -- Center point for display
    center_point geometry(Point, 4326),
    -- Additional metadata
    population INTEGER,
    area_km2 DECIMAL(10, 2),
    postal_codes TEXT[], -- Array of postal codes
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create table for municipalities  
CREATE TABLE IF NOT EXISTS municipalities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    district_id UUID REFERENCES districts(id) ON DELETE SET NULL,
    country_code VARCHAR(2) NOT NULL DEFAULT 'RS',
    -- PostGIS polygon for municipality boundaries
    boundary geometry(Polygon, 4326),
    -- Center point for display
    center_point geometry(Point, 4326),
    -- Additional metadata
    population INTEGER,
    area_km2 DECIMAL(10, 2),
    postal_codes TEXT[],
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Add indexes for spatial queries
CREATE INDEX idx_districts_boundary ON districts USING GIST (boundary);
CREATE INDEX idx_districts_center ON districts USING GIST (center_point);
CREATE INDEX idx_districts_name ON districts(name);
CREATE INDEX idx_districts_country ON districts(country_code);

CREATE INDEX idx_municipalities_boundary ON municipalities USING GIST (boundary);
CREATE INDEX idx_municipalities_center ON municipalities USING GIST (center_point);
CREATE INDEX idx_municipalities_name ON municipalities(name);
CREATE INDEX idx_municipalities_district ON municipalities(district_id);

-- Add foreign key to link listings to districts
ALTER TABLE listings_geo 
ADD COLUMN district_id UUID REFERENCES districts(id) ON DELETE SET NULL,
ADD COLUMN municipality_id UUID REFERENCES municipalities(id) ON DELETE SET NULL;

-- Create indexes for listing queries
CREATE INDEX idx_listings_geo_district ON listings_geo(district_id);
CREATE INDEX idx_listings_geo_municipality ON listings_geo(municipality_id);

-- Function to automatically assign district/municipality based on coordinates
CREATE OR REPLACE FUNCTION assign_district_municipality()
RETURNS TRIGGER AS $$
BEGIN
    -- Find district containing the point
    SELECT id INTO NEW.district_id
    FROM districts
    WHERE ST_Contains(boundary, NEW.location::geometry)
    LIMIT 1;
    
    -- Find municipality containing the point
    SELECT id INTO NEW.municipality_id
    FROM municipalities
    WHERE ST_Contains(boundary, NEW.location::geometry)
    LIMIT 1;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to auto-assign district/municipality
CREATE TRIGGER trigger_assign_district_municipality
BEFORE INSERT OR UPDATE OF location ON listings_geo
FOR EACH ROW
EXECUTE FUNCTION assign_district_municipality();

-- Update existing listings to assign districts/municipalities
UPDATE listings_geo lg
SET 
    district_id = d.id,
    municipality_id = m.id
FROM districts d, municipalities m
WHERE ST_Contains(d.boundary, lg.location::geometry)
   AND ST_Contains(m.boundary, lg.location::geometry);