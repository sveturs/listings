-- Rollback: Revert entity_type values from c2c_/b2c_ back to marketplace_/storefronts_

-- Revert unified_geo table
UPDATE unified_geo
SET source_type = 'marketplace_listing'
WHERE source_type = 'c2c_listing';

UPDATE unified_geo
SET source_type = 'storefront'
WHERE source_type = 'b2c_store';

UPDATE unified_geo
SET source_type = 'storefront_product'
WHERE source_type = 'b2c_product';

-- Revert reviews table
UPDATE reviews
SET entity_type = 'marketplace_listing'
WHERE entity_type = 'c2c_listing';

UPDATE reviews
SET entity_type = 'storefront'
WHERE entity_type = 'b2c_store';

UPDATE reviews
SET entity_type = 'storefront_product'
WHERE entity_type = 'b2c_product';

-- Revert translations table
UPDATE translations
SET entity_type = 'marketplace_category'
WHERE entity_type = 'c2c_category';

UPDATE translations
SET entity_type = 'marketplace_listing'
WHERE entity_type = 'c2c_listing';

UPDATE translations
SET entity_type = 'storefront'
WHERE entity_type = 'b2c_store';

UPDATE translations
SET entity_type = 'storefront_product'
WHERE entity_type = 'b2c_product';

-- Revert notifications table
UPDATE notifications
SET entity_type = 'marketplace_listing'
WHERE entity_type = 'c2c_listing';

UPDATE notifications
SET entity_type = 'storefront'
WHERE entity_type = 'b2c_store';

UPDATE notifications
SET entity_type = 'storefront_product'
WHERE entity_type = 'b2c_product';

-- Revert activity_logs table
UPDATE activity_logs
SET entity_type = 'marketplace_listing'
WHERE entity_type = 'c2c_listing';

UPDATE activity_logs
SET entity_type = 'storefront'
WHERE entity_type = 'b2c_store';

UPDATE activity_logs
SET entity_type = 'storefront_product'
WHERE entity_type = 'b2c_product';
