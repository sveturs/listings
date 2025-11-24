-- Rollback Migration: Drop Orders Table
-- Phase: Phase 17 - Orders Migration (Day 3-4)
-- Date: 2025-11-14

-- Drop table (CASCADE will handle dependent order_items)
DROP TABLE IF EXISTS orders CASCADE;

-- Drop trigger and function
DROP TRIGGER IF EXISTS trigger_orders_updated_at ON orders;
DROP FUNCTION IF EXISTS update_orders_updated_at();
