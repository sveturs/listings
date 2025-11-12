-- Phase 9.5.5: Add B2C Products and Inventory Tracking
-- Creates b2c_products table for storefront products
-- Creates b2c_product_variants table for product variations
-- Creates b2c_inventory_movements table for audit trail

-- =============================================================================
-- b2c_products (storefront products)
-- =============================================================================
CREATE TABLE IF NOT EXISTS b2c_products (
    id BIGSERIAL PRIMARY KEY,
    storefront_id INTEGER NOT NULL REFERENCES storefronts(id) ON DELETE CASCADE,

    -- Basic product information
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price NUMERIC(10, 2) NOT NULL CHECK (price >= 0),
    currency VARCHAR(3) DEFAULT 'USD' NOT NULL,
    category_id INTEGER,

    -- Inventory tracking
    sku VARCHAR(100) UNIQUE,
    barcode VARCHAR(100),
    stock_quantity INTEGER DEFAULT 0 NOT NULL CHECK (stock_quantity >= 0),
    stock_status VARCHAR(20) DEFAULT 'in_stock' NOT NULL CHECK (stock_status IN ('in_stock', 'out_of_stock', 'low_stock', 'preorder')),

    -- Product attributes & metadata
    attributes JSONB DEFAULT '{}'::jsonb,
    is_active BOOLEAN DEFAULT TRUE NOT NULL,

    -- Analytics
    view_count INTEGER DEFAULT 0 NOT NULL,
    sold_count INTEGER DEFAULT 0 NOT NULL,

    -- Individual product location (overrides storefront location)
    has_individual_location BOOLEAN DEFAULT FALSE,
    individual_address TEXT,
    individual_latitude NUMERIC(10,8),
    individual_longitude NUMERIC(11,8),
    location_privacy VARCHAR(20) CHECK (location_privacy IN ('exact', 'approximate', 'hidden')),
    show_on_map BOOLEAN DEFAULT TRUE,

    -- Variant management
    has_variants BOOLEAN DEFAULT FALSE,

    -- Timestamps
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW() NOT NULL,
    deleted_at TIMESTAMP
);

-- Indexes for b2c_products
CREATE INDEX IF NOT EXISTS idx_b2c_products_storefront_id ON b2c_products(storefront_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_b2c_products_category_id ON b2c_products(category_id) WHERE category_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_b2c_products_sku ON b2c_products(sku) WHERE sku IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_b2c_products_barcode ON b2c_products(barcode) WHERE barcode IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_b2c_products_is_active ON b2c_products(is_active, storefront_id) WHERE is_active = TRUE AND deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_b2c_products_stock_status ON b2c_products(stock_status, storefront_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_b2c_products_price ON b2c_products(price) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_b2c_products_created_at ON b2c_products(created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_b2c_products_location ON b2c_products USING GIST (
    ll_to_earth(individual_latitude::float, individual_longitude::float)
) WHERE individual_latitude IS NOT NULL AND individual_longitude IS NOT NULL AND has_individual_location = TRUE;

-- Update trigger for b2c_products
CREATE OR REPLACE FUNCTION update_b2c_products_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_b2c_products_updated_at
BEFORE UPDATE ON b2c_products
FOR EACH ROW
EXECUTE FUNCTION update_b2c_products_updated_at();

-- =============================================================================
-- b2c_product_variants (product variations like size, color)
-- =============================================================================
CREATE TABLE IF NOT EXISTS b2c_product_variants (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL REFERENCES b2c_products(id) ON DELETE CASCADE,

    -- Variant identifiers
    sku VARCHAR(100) UNIQUE,
    barcode VARCHAR(100),

    -- Pricing (optional override)
    price NUMERIC(10, 2) CHECK (price >= 0),
    compare_at_price NUMERIC(10, 2) CHECK (compare_at_price >= 0),
    cost_price NUMERIC(10, 2) CHECK (cost_price >= 0),

    -- Inventory
    stock_quantity INTEGER DEFAULT 0 NOT NULL CHECK (stock_quantity >= 0),
    stock_status VARCHAR(20) DEFAULT 'in_stock' NOT NULL CHECK (stock_status IN ('in_stock', 'out_of_stock', 'low_stock', 'preorder')),
    low_stock_threshold INTEGER DEFAULT 10,

    -- Variant attributes (e.g., {"size": "L", "color": "Red"})
    variant_attributes JSONB DEFAULT '{}'::jsonb,

    -- Physical properties
    weight NUMERIC(10, 3),
    dimensions JSONB, -- {"length": 10, "width": 5, "height": 3, "unit": "cm"}

    -- Status
    is_active BOOLEAN DEFAULT TRUE NOT NULL,
    is_default BOOLEAN DEFAULT FALSE, -- Only one variant can be default per product

    -- Analytics
    view_count INTEGER DEFAULT 0 NOT NULL,
    sold_count INTEGER DEFAULT 0 NOT NULL,

    -- Timestamps
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW() NOT NULL
);

-- Indexes for b2c_product_variants
CREATE INDEX IF NOT EXISTS idx_b2c_product_variants_product_id ON b2c_product_variants(product_id);
CREATE INDEX IF NOT EXISTS idx_b2c_product_variants_sku ON b2c_product_variants(sku) WHERE sku IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_b2c_product_variants_barcode ON b2c_product_variants(barcode) WHERE barcode IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_b2c_product_variants_is_default ON b2c_product_variants(product_id, is_default) WHERE is_default = TRUE;
CREATE INDEX IF NOT EXISTS idx_b2c_product_variants_is_active ON b2c_product_variants(product_id, is_active) WHERE is_active = TRUE;

-- Constraint: Only one default variant per product
CREATE UNIQUE INDEX idx_b2c_product_variants_unique_default
    ON b2c_product_variants(product_id)
    WHERE is_default = TRUE;

-- Update trigger for b2c_product_variants
CREATE OR REPLACE FUNCTION update_b2c_product_variants_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();

    -- Automatically update parent product's has_variants flag
    UPDATE b2c_products
    SET has_variants = EXISTS(
        SELECT 1 FROM b2c_product_variants
        WHERE product_id = NEW.product_id
    )
    WHERE id = NEW.product_id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_b2c_product_variants_updated_at
BEFORE UPDATE ON b2c_product_variants
FOR EACH ROW
EXECUTE FUNCTION update_b2c_product_variants_updated_at();

-- =============================================================================
-- b2c_inventory_movements (audit trail for stock changes)
-- =============================================================================
CREATE TABLE IF NOT EXISTS b2c_inventory_movements (
    id BIGSERIAL PRIMARY KEY,

    -- Product or variant reference
    storefront_product_id BIGINT NOT NULL REFERENCES b2c_products(id) ON DELETE CASCADE,
    variant_id BIGINT REFERENCES b2c_product_variants(id) ON DELETE CASCADE,

    -- Movement details
    type VARCHAR(20) NOT NULL CHECK (type IN ('in', 'out', 'adjustment')),
    quantity INTEGER NOT NULL,

    -- Context
    reason VARCHAR(100), -- e.g., 'restock', 'sale', 'inventory_count', 'damaged', 'returned'
    notes TEXT,

    -- User tracking
    user_id BIGINT NOT NULL,

    -- Timestamp
    created_at TIMESTAMP DEFAULT NOW() NOT NULL
);

-- Indexes for b2c_inventory_movements
CREATE INDEX IF NOT EXISTS idx_b2c_inventory_movements_product_id ON b2c_inventory_movements(storefront_product_id);
CREATE INDEX IF NOT EXISTS idx_b2c_inventory_movements_variant_id ON b2c_inventory_movements(variant_id) WHERE variant_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_b2c_inventory_movements_user_id ON b2c_inventory_movements(user_id);
CREATE INDEX IF NOT EXISTS idx_b2c_inventory_movements_created_at ON b2c_inventory_movements(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_b2c_inventory_movements_type ON b2c_inventory_movements(type);

-- Comments for documentation
COMMENT ON TABLE b2c_products IS 'B2C products sold by storefronts';
COMMENT ON COLUMN b2c_products.has_individual_location IS 'If true, product has its own location independent of storefront';
COMMENT ON COLUMN b2c_products.location_privacy IS 'Privacy level for product location: exact, approximate, or hidden';
COMMENT ON COLUMN b2c_products.stock_status IS 'Current stock status: in_stock, out_of_stock, low_stock, or preorder';

COMMENT ON TABLE b2c_product_variants IS 'Product variants (size, color, etc.) for b2c_products';
COMMENT ON COLUMN b2c_product_variants.is_default IS 'Only one variant can be marked as default per product';
COMMENT ON COLUMN b2c_product_variants.variant_attributes IS 'JSONB containing variant-specific attributes like size, color';

COMMENT ON TABLE b2c_inventory_movements IS 'Audit trail for all inventory stock changes';
COMMENT ON COLUMN b2c_inventory_movements.type IS 'Type of movement: in (add stock), out (remove stock), adjustment (set quantity)';
COMMENT ON COLUMN b2c_inventory_movements.reason IS 'Reason for movement: restock, sale, inventory_count, damaged, returned, etc.';
