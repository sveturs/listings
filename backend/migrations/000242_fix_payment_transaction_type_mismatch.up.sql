-- Fix type mismatch between marketplace_orders.payment_transaction_id (integer) and payment_transactions.id (bigint)
-- This migration corrects the data type to ensure referential integrity

-- 1. Drop the existing foreign key constraint if it exists
ALTER TABLE marketplace_orders 
DROP CONSTRAINT IF EXISTS marketplace_orders_payment_transaction_id_fkey;

-- 2. Change the column type from integer to bigint
ALTER TABLE marketplace_orders 
ALTER COLUMN payment_transaction_id TYPE bigint;

-- 3. Re-add the foreign key constraint with correct type
ALTER TABLE marketplace_orders 
ADD CONSTRAINT marketplace_orders_payment_transaction_id_fkey 
FOREIGN KEY (payment_transaction_id) 
REFERENCES payment_transactions(id) 
ON DELETE SET NULL;

-- 4. Remove duplicate indexes in marketplace_orders table
DROP INDEX IF EXISTS idx_marketplace_orders_buyer_id_duplicate;
DROP INDEX IF EXISTS idx_marketplace_orders_seller_id_duplicate; 
DROP INDEX IF EXISTS idx_marketplace_orders_listing_id_duplicate;
DROP INDEX IF EXISTS idx_marketplace_orders_created_at_duplicate;
DROP INDEX IF EXISTS idx_marketplace_orders_status_duplicate;
DROP INDEX IF EXISTS idx_marketplace_orders_payment_status_duplicate;

-- Add comment to document the fix
COMMENT ON COLUMN marketplace_orders.payment_transaction_id IS 'Reference to payment_transactions table (bigint type to match payment_transactions.id)';