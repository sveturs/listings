-- Add address_multilingual column to store localized addresses
ALTER TABLE marketplace_listings
ADD COLUMN IF NOT EXISTS address_multilingual JSONB;

-- Comment for documentation
COMMENT ON COLUMN marketplace_listings.address_multilingual IS 'Stores multilingual addresses in format {"en": "address", "ru": "адрес", "sr": "adresa"}';

-- Create index for better performance when querying specific languages
CREATE INDEX IF NOT EXISTS idx_marketplace_listings_address_multilingual
ON marketplace_listings USING gin (address_multilingual);

-- Migrate existing address translations from translations table
UPDATE marketplace_listings ml
SET address_multilingual = (
    SELECT jsonb_object_agg(t.language, t.translated_text)
    FROM translations t
    WHERE t.entity_type = 'listing'
    AND t.entity_id = ml.id
    AND t.field_name = 'address'
    AND t.translated_text IS NOT NULL
)
WHERE EXISTS (
    SELECT 1 FROM translations t2
    WHERE t2.entity_type = 'listing'
    AND t2.entity_id = ml.id
    AND t2.field_name = 'address'
);

-- Also set default address for listings that have location but no multilingual address
UPDATE marketplace_listings
SET address_multilingual = jsonb_build_object(
    'en', location,
    'ru', location,
    'sr', location
)
WHERE location IS NOT NULL
AND address_multilingual IS NULL;