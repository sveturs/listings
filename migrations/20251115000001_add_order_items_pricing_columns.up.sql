-- Migration: Add subtotal, discount, and image_url columns to order_items
-- Purpose: Support full pricing breakdown and product snapshot
-- Author: System
-- Date: 2025-11-15

-- Add subtotal column (quantity * unit_price)
ALTER TABLE order_items
ADD COLUMN IF NOT EXISTS subtotal DECIMAL(10,2) NOT NULL DEFAULT 0;

-- Add discount column (item-level discount)
ALTER TABLE order_items
ADD COLUMN IF NOT EXISTS discount DECIMAL(10,2) NOT NULL DEFAULT 0;

-- Add image_url column (primary image snapshot at purchase time)
ALTER TABLE order_items
ADD COLUMN IF NOT EXISTS image_url VARCHAR(500);

-- Add check constraints for non-negative values
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'chk_order_items_subtotal_non_negative'
    ) THEN
        ALTER TABLE order_items ADD CONSTRAINT chk_order_items_subtotal_non_negative CHECK (subtotal >= 0);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'chk_order_items_discount_non_negative'
    ) THEN
        ALTER TABLE order_items ADD CONSTRAINT chk_order_items_discount_non_negative CHECK (discount >= 0);
    END IF;
END $$;

-- Add comment explaining the pricing logic
COMMENT ON COLUMN order_items.subtotal IS 'Subtotal before discount (quantity * unit_price)';
COMMENT ON COLUMN order_items.discount IS 'Item-level discount amount';
COMMENT ON COLUMN order_items.image_url IS 'Snapshot of primary image URL at purchase time';

-- Update existing rows to calculate subtotal from quantity * price
-- and ensure total = subtotal - discount (discount defaults to 0)
UPDATE order_items
SET subtotal = quantity * price,
    discount = 0
WHERE subtotal = 0;

-- Verify data integrity: total should equal subtotal - discount
DO $$
DECLARE
    invalid_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO invalid_count
    FROM order_items
    WHERE ABS(total - (subtotal - discount)) > 0.01; -- Allow for floating point rounding

    IF invalid_count > 0 THEN
        RAISE NOTICE 'Found % order_items with total != subtotal - discount', invalid_count;
    END IF;
END $$;
