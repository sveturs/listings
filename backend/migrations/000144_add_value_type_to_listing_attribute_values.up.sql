-- Migration to add value_type column to listing_attribute_values table

ALTER TABLE public.listing_attribute_values 
ADD COLUMN value_type VARCHAR(20) DEFAULT 'text' NOT NULL;

-- Add index on value_type for better query performance
CREATE INDEX idx_listing_attr_value_type ON public.listing_attribute_values USING btree (value_type);

-- Update existing records based on which value fields are populated
UPDATE public.listing_attribute_values 
SET value_type = CASE 
    WHEN numeric_value IS NOT NULL THEN 'numeric'
    WHEN boolean_value IS NOT NULL THEN 'boolean'
    WHEN json_value IS NOT NULL THEN 'json'
    WHEN text_value IS NOT NULL THEN 'text'
    ELSE 'text'
END;