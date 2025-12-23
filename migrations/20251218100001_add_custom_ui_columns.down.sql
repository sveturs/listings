-- Remove custom UI columns from categories table

ALTER TABLE categories
DROP COLUMN IF EXISTS custom_ui_component,
DROP COLUMN IF EXISTS has_custom_ui;
