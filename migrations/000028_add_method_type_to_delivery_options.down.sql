-- Rollback: remove method_type column

DROP INDEX IF EXISTS idx_storefront_delivery_options_provider_method;

ALTER TABLE storefront_delivery_options
DROP COLUMN IF EXISTS method_type;
