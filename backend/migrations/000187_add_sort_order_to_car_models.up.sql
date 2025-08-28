-- Add missing sort_order column to car_models table
ALTER TABLE car_models ADD COLUMN IF NOT EXISTS sort_order INTEGER DEFAULT 0;

-- Update sort_order for existing models
UPDATE car_models SET sort_order = id;