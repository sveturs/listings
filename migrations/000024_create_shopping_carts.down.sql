-- Rollback Migration: Drop Shopping Carts Table
-- Phase: Phase 17 - Orders Migration (Day 3-4)
-- Date: 2025-11-14

-- Drop table (CASCADE will handle dependent cart_items)
DROP TABLE IF EXISTS shopping_carts CASCADE;

-- Drop trigger and function
DROP TRIGGER IF EXISTS trigger_shopping_carts_updated_at ON shopping_carts;
DROP FUNCTION IF EXISTS update_shopping_carts_updated_at();
