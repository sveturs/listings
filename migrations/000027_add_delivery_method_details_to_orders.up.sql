-- Migration: Add delivery_method_details column to orders table
-- Purpose: Store detailed delivery method information selected by customer at checkout
-- This enables WMS to know which delivery provider and method to use for shipment creation

-- Add JSONB column for delivery method details
ALTER TABLE orders
ADD COLUMN delivery_method_details JSONB;

-- Add comment explaining the structure
COMMENT ON COLUMN orders.delivery_method_details IS
  'Detailed delivery method information selected by customer at checkout. Contains: provider_id (string), method_type (string), display_name (string), price (number), currency (string), estimated_days_min (number), estimated_days_max (number), supports_cod (boolean), supports_tracking (boolean). This data is passed to WMS via order.confirmed event for shipment creation.';

-- Create index on provider_id for faster lookups
CREATE INDEX idx_orders_delivery_method_details_provider
ON orders ((delivery_method_details->>'provider_id'))
WHERE delivery_method_details IS NOT NULL;

-- Example data structure:
-- {
--   "provider_id": "post_express",
--   "method_type": "express",
--   "display_name": "Экспресс доставка",
--   "price": 450.0,
--   "currency": "RSD",
--   "estimated_days_min": 2,
--   "estimated_days_max": 3,
--   "supports_cod": true,
--   "supports_tracking": true
-- }
