-- Create categories VIEW on c2c_categories
-- This allows the microservice to use 'categories' table name
-- while keeping compatibility with monolith's c2c_categories

BEGIN;

-- Drop existing view if any
DROP VIEW IF EXISTS categories CASCADE;

-- Create view that maps c2c_categories to categories
CREATE OR REPLACE VIEW categories AS
SELECT
    id,
    parent_id,
    name,
    slug,
    level,
    sort_order,
    is_active,
    icon,
    description,
    has_custom_ui,
    custom_ui_component,
    created_at,
    created_at as updated_at
FROM c2c_categories;

COMMENT ON VIEW categories IS 'View that maps c2c_categories to categories for microservice compatibility';

COMMIT;
