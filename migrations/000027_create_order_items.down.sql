-- Rollback Migration: Drop Order Items Table
-- Phase: Phase 17 - Orders Migration (Day 3-4)
-- Date: 2025-11-14

-- Drop table
DROP TABLE IF EXISTS order_items CASCADE;
