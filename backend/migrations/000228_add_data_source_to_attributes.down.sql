-- Откат изменений
ALTER TABLE category_attributes 
DROP COLUMN IF EXISTS data_source_config;

-- Восстановление оригинального CHECK constraint
ALTER TABLE category_attributes 
DROP CONSTRAINT IF EXISTS check_data_source_values;

ALTER TABLE category_attributes 
ADD CONSTRAINT check_data_source_values 
CHECK (data_source IN ('manual', 'api_external', 'ai_generated', 'imported', 'computed'));