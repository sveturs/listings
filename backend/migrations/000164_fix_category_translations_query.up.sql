-- This migration doesn't change the database schema
-- It's a placeholder to remind that the code needs to be fixed
-- The issue is in the GetAllCategories function where it uses 'marketplace_category' instead of 'category'

-- Add a comment to remind about this
COMMENT ON TABLE translations IS 'Entity types: category (for marketplace categories), attribute, etc. NOT marketplace_category!';