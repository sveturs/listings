-- Rollback Migration: Drop Inventory Reservations Table
-- Phase: Phase 17 - Orders Migration (Day 3-4)
-- Date: 2025-11-14

-- Drop table
DROP TABLE IF EXISTS inventory_reservations CASCADE;

-- Drop trigger and function
DROP TRIGGER IF EXISTS trigger_inventory_reservations_updated_at ON inventory_reservations;
DROP FUNCTION IF EXISTS update_inventory_reservations_updated_at();
