-- Rollback migration: Revert payment_transaction_id type change
-- WARNING: This rollback may fail if there are values that don't fit in integer type

-- 1. Drop the foreign key constraint
ALTER TABLE marketplace_orders 
DROP CONSTRAINT IF EXISTS marketplace_orders_payment_transaction_id_fkey;

-- 2. Change the column type back to integer
-- Note: This will fail if there are bigint values that don't fit in integer
ALTER TABLE marketplace_orders 
ALTER COLUMN payment_transaction_id TYPE integer;

-- 3. Re-add the foreign key constraint (will fail if types don't match)
ALTER TABLE marketplace_orders 
ADD CONSTRAINT marketplace_orders_payment_transaction_id_fkey 
FOREIGN KEY (payment_transaction_id) 
REFERENCES payment_transactions(id) 
ON DELETE SET NULL;

-- 4. Re-create the duplicate indexes (not recommended, but for rollback completeness)
-- These were removed as duplicates, so we don't recreate them

-- Remove comment
COMMENT ON COLUMN marketplace_orders.payment_transaction_id IS NULL;