-- Migration Down: 000023_create_attributes_schema
-- Description: Rollback unified attributes system
-- WARNING: This will drop all attribute data!

-- Drop triggers first
DROP TRIGGER IF EXISTS trigger_attribute_options_updated_at ON attribute_options;
DROP TRIGGER IF EXISTS trigger_variant_attr_values_updated_at ON variant_attribute_values;
DROP TRIGGER IF EXISTS trigger_category_variant_attrs_updated_at ON category_variant_attributes;
DROP TRIGGER IF EXISTS trigger_listing_attr_values_updated_at ON listing_attribute_values;
DROP TRIGGER IF EXISTS trigger_category_attributes_updated_at ON category_attributes;
DROP TRIGGER IF EXISTS trigger_attributes_search_vector ON attributes;
DROP TRIGGER IF EXISTS trigger_attributes_updated_at ON attributes;

-- Drop functions
DROP FUNCTION IF EXISTS update_attributes_search_vector();
DROP FUNCTION IF EXISTS update_attributes_timestamp();

-- Drop tables in reverse order (respecting foreign keys)
DROP TABLE IF EXISTS attribute_search_cache CASCADE;
DROP TABLE IF EXISTS attribute_options CASCADE;
DROP TABLE IF EXISTS variant_attribute_values CASCADE;
DROP TABLE IF EXISTS category_variant_attributes CASCADE;
DROP TABLE IF EXISTS listing_attribute_values CASCADE;
DROP TABLE IF EXISTS category_attributes CASCADE;
DROP TABLE IF EXISTS attributes CASCADE;

-- Note: uuid-ossp extension is left intact as it may be used by other tables
