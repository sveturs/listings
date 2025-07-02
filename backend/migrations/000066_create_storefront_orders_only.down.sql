-- Drop all triggers first
DROP TRIGGER IF EXISTS calculate_escrow_release_date_trigger ON storefront_orders;
DROP TRIGGER IF EXISTS set_order_number_trigger ON storefront_orders;
DROP TRIGGER IF EXISTS update_storefront_orders_updated_at ON storefront_orders;

-- Drop functions
DROP FUNCTION IF EXISTS calculate_escrow_release_date();
DROP FUNCTION IF EXISTS set_order_number();
DROP FUNCTION IF EXISTS generate_order_number();

-- Drop indexes
DROP INDEX IF EXISTS idx_storefront_orders_escrow_date;
DROP INDEX IF EXISTS idx_storefront_orders_status;
DROP INDEX IF EXISTS idx_storefront_orders_storefront;
DROP INDEX IF EXISTS idx_storefront_orders_customer;
DROP INDEX IF EXISTS idx_payment_transactions_source;
DROP INDEX IF EXISTS idx_inventory_reservations_expires;
DROP INDEX IF EXISTS idx_inventory_reservations_product;

-- Remove columns from payment_transactions
ALTER TABLE payment_transactions 
DROP COLUMN IF EXISTS storefront_id,
DROP COLUMN IF EXISTS source_id,
DROP COLUMN IF EXISTS source_type;

-- Drop tables
DROP TABLE IF EXISTS inventory_reservations;
DROP TABLE IF EXISTS storefront_order_items;
DROP TABLE IF EXISTS storefront_orders;