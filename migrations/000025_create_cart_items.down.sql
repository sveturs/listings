-- Rollback Migration: Drop Cart Items Table
-- Phase: Phase 17 - Orders Migration (Day 3-4)
-- Date: 2025-11-14

-- Drop table
DROP TABLE IF EXISTS cart_items CASCADE;

-- Drop trigger and function
DROP TRIGGER IF EXISTS trigger_cart_items_updated_at ON cart_items;
DROP FUNCTION IF EXISTS update_cart_items_updated_at();
