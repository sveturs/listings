-- Migration: Create product_variants table (Phase 3: Variants System)
-- Description: Product variants with SKU, pricing, inventory management
-- Author: Claude (Elite Architect)
-- Date: 2025-12-17
-- Phase: 3 (Variants)
-- Task: BE-3.1

BEGIN;

-- ============================================================================
-- TABLE: product_variants - Product Variants
-- ============================================================================
CREATE TABLE IF NOT EXISTS product_variants (
    -- Primary key
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Reference to parent product
    product_id UUID NOT NULL,

    -- Unique identifier
    sku VARCHAR(100) NOT NULL,

    -- Pricing (NULL = use product's base_price)
    price DECIMAL(12,2),
    compare_at_price DECIMAL(12,2),

    -- Inventory management
    stock_quantity INTEGER NOT NULL DEFAULT 0,
    reserved_quantity INTEGER NOT NULL DEFAULT 0,
    low_stock_alert INTEGER DEFAULT 5,

    -- Physical characteristics
    weight_grams DECIMAL(8,3),
    barcode VARCHAR(50),

    -- Display settings
    is_default BOOLEAN NOT NULL DEFAULT false,
    position INTEGER NOT NULL DEFAULT 0,

    -- Status
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (
        status IN ('active', 'out_of_stock', 'discontinued')
    ),

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Constraints
    CONSTRAINT unique_variant_sku UNIQUE (sku),
    CONSTRAINT positive_stock CHECK (stock_quantity >= 0),
    CONSTRAINT positive_reserved CHECK (reserved_quantity >= 0),
    CONSTRAINT reserved_not_exceed_stock CHECK (reserved_quantity <= stock_quantity),
    CONSTRAINT valid_price CHECK (price IS NULL OR price >= 0),
    CONSTRAINT valid_compare_price CHECK (compare_at_price IS NULL OR compare_at_price >= 0)
);

-- ============================================================================
-- INDEXES for product_variants
-- ============================================================================

-- Primary lookups
CREATE INDEX idx_variants_product ON product_variants(product_id, position);
CREATE INDEX idx_variants_sku ON product_variants(sku);

-- Stock queries
CREATE INDEX idx_variants_stock_status ON product_variants(product_id, status, stock_quantity)
WHERE status = 'active';

-- Available quantity (calculated: stock - reserved)
CREATE INDEX idx_variants_available ON product_variants(product_id, (stock_quantity - reserved_quantity))
WHERE status = 'active' AND (stock_quantity - reserved_quantity) > 0;

-- Default variant lookup
CREATE INDEX idx_variants_default ON product_variants(product_id)
WHERE is_default = true;

-- Barcode lookup
CREATE INDEX idx_variants_barcode ON product_variants(barcode)
WHERE barcode IS NOT NULL;

-- Low stock alerts
CREATE INDEX idx_variants_low_stock ON product_variants(product_id, stock_quantity, low_stock_alert)
WHERE status = 'active' AND stock_quantity <= low_stock_alert;

-- ============================================================================
-- TRIGGERS
-- ============================================================================

-- Auto-update updated_at timestamp
CREATE TRIGGER trigger_variants_updated_at
    BEFORE UPDATE ON product_variants
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Auto-update status based on stock
CREATE OR REPLACE FUNCTION auto_update_variant_status()
RETURNS TRIGGER AS $$
BEGIN
    -- Set status to out_of_stock when stock reaches zero
    IF NEW.stock_quantity = 0 AND OLD.stock_quantity > 0 THEN
        NEW.status = 'out_of_stock';
    END IF;

    -- Set status to active when stock becomes available again
    IF NEW.stock_quantity > 0 AND OLD.status = 'out_of_stock' THEN
        NEW.status = 'active';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_variants_auto_status
    BEFORE UPDATE OF stock_quantity ON product_variants
    FOR EACH ROW
    EXECUTE FUNCTION auto_update_variant_status();

-- Ensure only one default variant per product
CREATE OR REPLACE FUNCTION enforce_single_default_variant()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.is_default = true THEN
        -- Unset is_default for all other variants of the same product
        UPDATE product_variants
        SET is_default = false
        WHERE product_id = NEW.product_id AND id != NEW.id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_variants_single_default
    BEFORE INSERT OR UPDATE OF is_default ON product_variants
    FOR EACH ROW
    WHEN (NEW.is_default = true)
    EXECUTE FUNCTION enforce_single_default_variant();

-- ============================================================================
-- COMMENTS
-- ============================================================================

COMMENT ON TABLE product_variants IS 'Product variants with unique SKU, pricing, and inventory';
COMMENT ON COLUMN product_variants.product_id IS 'Reference to parent product (foreign key added later)';
COMMENT ON COLUMN product_variants.sku IS 'Stock Keeping Unit - unique identifier for this variant';
COMMENT ON COLUMN product_variants.price IS 'Variant-specific price (NULL = use product base_price)';
COMMENT ON COLUMN product_variants.compare_at_price IS 'Original price for discount display';
COMMENT ON COLUMN product_variants.stock_quantity IS 'Total available stock';
COMMENT ON COLUMN product_variants.reserved_quantity IS 'Stock reserved in active orders';
COMMENT ON COLUMN product_variants.low_stock_alert IS 'Threshold for low stock notifications';
COMMENT ON COLUMN product_variants.weight_grams IS 'Variant weight in grams';
COMMENT ON COLUMN product_variants.barcode IS 'EAN/UPC barcode';
COMMENT ON COLUMN product_variants.is_default IS 'Default variant for product (only one allowed)';
COMMENT ON COLUMN product_variants.position IS 'Display order for variant selector';
COMMENT ON COLUMN product_variants.status IS 'Variant status: active, out_of_stock, discontinued';

COMMIT;
