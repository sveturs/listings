-- Add icon field to category_attributes table

ALTER TABLE category_attributes ADD COLUMN icon VARCHAR(10) DEFAULT '';

-- Add comment
COMMENT ON COLUMN category_attributes.icon IS 'Unicode emoji icon for attribute display';