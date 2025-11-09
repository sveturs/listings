-- Add stock_status column to listings table
-- This field tracks inventory availability status for products

-- Add stock_status column with default 'in_stock'
ALTER TABLE listings
ADD COLUMN IF NOT EXISTS stock_status VARCHAR(50) DEFAULT 'in_stock'
CHECK (stock_status IN ('in_stock', 'out_of_stock', 'low_stock', 'discontinued'));

-- Update existing records based on quantity
UPDATE listings
SET stock_status = CASE
    WHEN quantity > 0 THEN 'in_stock'
    WHEN quantity = 0 THEN 'out_of_stock'
    ELSE 'in_stock'
END
WHERE stock_status IS NULL OR stock_status = 'in_stock';

-- Create index for stock_status queries
CREATE INDEX IF NOT EXISTS idx_listings_stock_status ON listings(stock_status) WHERE is_deleted = false;

-- Add comment
COMMENT ON COLUMN listings.stock_status IS 'Inventory status: in_stock (available), out_of_stock (no stock), low_stock (below threshold), discontinued (no longer available)';
