-- Drop triggers
DROP TRIGGER IF EXISTS update_car_generations_updated_at ON car_generations;
DROP TRIGGER IF EXISTS update_car_models_updated_at ON car_models;
DROP TRIGGER IF EXISTS update_car_makes_updated_at ON car_makes;

-- Drop tables in reverse order due to foreign key constraints
DROP TABLE IF EXISTS car_generations;
DROP TABLE IF EXISTS car_models;
DROP TABLE IF EXISTS car_makes;