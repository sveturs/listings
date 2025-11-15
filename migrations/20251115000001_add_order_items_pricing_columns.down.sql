-- Migration Rollback: Remove subtotal, discount, and image_url columns from order_items
-- Purpose: Revert to original schema
-- Author: System
-- Date: 2025-11-15

-- Drop check constraints first
ALTER TABLE order_items
DROP CONSTRAINT IF EXISTS chk_order_items_subtotal_non_negative;

ALTER TABLE order_items
DROP CONSTRAINT IF EXISTS chk_order_items_discount_non_negative;

-- Drop columns
ALTER TABLE order_items
DROP COLUMN IF EXISTS subtotal;

ALTER TABLE order_items
DROP COLUMN IF EXISTS discount;

ALTER TABLE order_items
DROP COLUMN IF EXISTS image_url;
