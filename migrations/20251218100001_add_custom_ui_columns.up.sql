-- Add custom UI columns to categories table
-- Required by categories_repository.go

ALTER TABLE categories
ADD COLUMN IF NOT EXISTS has_custom_ui BOOLEAN DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS custom_ui_component VARCHAR(100);

-- Add comment
COMMENT ON COLUMN categories.has_custom_ui IS 'Whether category has custom UI component';
COMMENT ON COLUMN categories.custom_ui_component IS 'Name of custom UI component for this category';
