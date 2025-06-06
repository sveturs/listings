-- Remove icon field from category_attributes table

ALTER TABLE category_attributes DROP COLUMN IF EXISTS icon;