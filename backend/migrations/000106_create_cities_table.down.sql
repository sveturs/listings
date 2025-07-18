-- Remove foreign key constraint
ALTER TABLE districts DROP CONSTRAINT IF EXISTS fk_districts_city_id;

-- Drop indexes
DROP INDEX IF EXISTS idx_cities_boundary;
DROP INDEX IF EXISTS idx_cities_center;
DROP INDEX IF EXISTS idx_cities_name;
DROP INDEX IF EXISTS idx_cities_slug;
DROP INDEX IF EXISTS idx_cities_country;
DROP INDEX IF EXISTS idx_cities_has_districts;
DROP INDEX IF EXISTS idx_cities_priority;
DROP INDEX IF EXISTS idx_districts_city_id;

-- Drop cities table
DROP TABLE IF EXISTS cities;