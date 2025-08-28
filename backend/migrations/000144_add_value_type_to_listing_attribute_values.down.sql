-- Migration rollback to remove value_type column from listing_attribute_values table

-- Drop the index first
DROP INDEX IF EXISTS idx_listing_attr_value_type;

-- Remove the value_type column
ALTER TABLE public.listing_attribute_values DROP COLUMN IF EXISTS value_type;