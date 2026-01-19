-- Add social_links JSONB column to storefronts table
-- Purpose: Store social media links (Instagram, Facebook, Telegram, WhatsApp)
-- Expected format: {"instagram": "url", "facebook": "url", "telegram": "url", "whatsapp": "number"}

ALTER TABLE storefronts
ADD COLUMN social_links JSONB DEFAULT '{}'::jsonb;

-- Add index for faster queries when filtering by social links
CREATE INDEX idx_storefronts_social_links_gin ON storefronts USING gin(social_links);

-- Comment
COMMENT ON COLUMN storefronts.social_links IS 'Social media links (Instagram, Facebook, Telegram, WhatsApp)';
