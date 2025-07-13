-- Drop trigger first
DROP TRIGGER IF EXISTS trigger_assign_district_municipality ON listings_geo;
DROP FUNCTION IF EXISTS assign_district_municipality();

-- Remove columns from listings_geo
ALTER TABLE listings_geo 
DROP COLUMN IF EXISTS district_id,
DROP COLUMN IF EXISTS municipality_id;

-- Drop indexes
DROP INDEX IF EXISTS idx_municipalities_boundary;
DROP INDEX IF EXISTS idx_municipalities_center;
DROP INDEX IF EXISTS idx_municipalities_name;
DROP INDEX IF EXISTS idx_municipalities_district;

DROP INDEX IF EXISTS idx_districts_boundary;
DROP INDEX IF EXISTS idx_districts_center;
DROP INDEX IF EXISTS idx_districts_name;
DROP INDEX IF EXISTS idx_districts_country;

-- Drop tables
DROP TABLE IF EXISTS municipalities;
DROP TABLE IF EXISTS districts;