-- Add method_type column to distinguish between delivery methods of same provider
-- (e.g., post_express:standard vs post_express:express)

ALTER TABLE storefront_delivery_options
ADD COLUMN method_type VARCHAR(50);

-- Add index for faster lookups by provider+method_type combination
CREATE INDEX idx_storefront_delivery_options_provider_method
ON storefront_delivery_options (storefront_id, provider, method_type)
WHERE provider IS NOT NULL AND method_type IS NOT NULL;

-- Update existing records with default method_type based on name
UPDATE storefront_delivery_options
SET method_type = CASE
    WHEN name ILIKE '%экспресс%' OR name ILIKE '%express%' THEN 'express'
    WHEN name ILIKE '%самовывоз%' OR name ILIKE '%pickup%' THEN 'pickup'
    WHEN name ILIKE '%стандарт%' OR name ILIKE '%standard%' THEN 'standard'
    ELSE 'standard' -- default
END
WHERE method_type IS NULL;

-- Add comment
COMMENT ON COLUMN storefront_delivery_options.method_type IS
  'Delivery method type: standard, express, overnight, economy, pickup, pickup_point. Distinguishes methods within same provider.';
