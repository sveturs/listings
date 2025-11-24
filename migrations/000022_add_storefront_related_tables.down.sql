-- Rollback Migration: Drop Storefront Related Tables
-- Phase: Storefront Merge - Delivery to Listings
-- Date: 2025-11-12

-- Drop tables in reverse order (respect FK constraints)
DROP TABLE IF EXISTS storefront_delivery_options CASCADE;
DROP TABLE IF EXISTS storefront_payment_methods CASCADE;
DROP TABLE IF EXISTS storefront_hours CASCADE;
DROP TABLE IF EXISTS storefront_staff CASCADE;

-- Drop triggers and functions
DROP TRIGGER IF EXISTS trigger_storefront_delivery_options_updated_at ON storefront_delivery_options;
DROP FUNCTION IF EXISTS update_storefront_delivery_options_updated_at();

DROP TRIGGER IF EXISTS trigger_storefront_staff_updated_at ON storefront_staff;
DROP FUNCTION IF EXISTS update_storefront_staff_updated_at();
