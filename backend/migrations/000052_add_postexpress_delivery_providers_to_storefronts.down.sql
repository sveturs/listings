-- Remove delivery_providers from storefronts settings
UPDATE storefronts
SET settings = settings - 'delivery_providers'
WHERE settings ? 'delivery_providers';