-- Migration: Create product variants system
-- This migration creates tables for product variants with different attributes, prices, and stock

-- Create product_variant_attributes table for defining variant attributes (color, size, etc.)
CREATE TABLE IF NOT EXISTS product_variant_attributes (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL, -- e.g., "color", "size", "model"
    display_name VARCHAR(255) NOT NULL, -- e.g., "Цвет", "Размер", "Модель"
    type VARCHAR(50) NOT NULL DEFAULT 'text', -- text, color, image, number
    is_required BOOLEAN NOT NULL DEFAULT false,
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create product_variant_attribute_values table for possible values of attributes
CREATE TABLE IF NOT EXISTS product_variant_attribute_values (
    id SERIAL PRIMARY KEY,
    attribute_id INTEGER NOT NULL REFERENCES product_variant_attributes(id) ON DELETE CASCADE,
    value VARCHAR(255) NOT NULL, -- e.g., "red", "XL", "iPhone 15 Pro"
    display_name VARCHAR(255) NOT NULL, -- e.g., "Красный", "XL", "iPhone 15 Pro"
    color_hex VARCHAR(7), -- for color attributes: #FF0000
    image_url TEXT, -- for image-based attributes
    sort_order INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create storefront_product_variants table for actual product variants
CREATE TABLE IF NOT EXISTS storefront_product_variants (
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL REFERENCES storefront_products(id) ON DELETE CASCADE,
    sku VARCHAR(100) UNIQUE, -- unique SKU for this variant
    barcode VARCHAR(100),

    -- Pricing (can override parent product price)
    price NUMERIC(15, 2), -- if NULL, uses parent product price
    compare_at_price NUMERIC(15, 2), -- original price for discounts
    cost_price NUMERIC(15, 2), -- cost for profit calculations

    -- Stock management
    stock_quantity INTEGER NOT NULL DEFAULT 0 CHECK (stock_quantity >= 0),
    stock_status VARCHAR(20) NOT NULL DEFAULT 'in_stock' CHECK (stock_status IN ('in_stock', 'low_stock', 'out_of_stock')),
    low_stock_threshold INTEGER DEFAULT 5,

    -- Variant attributes (JSON for flexibility)
    variant_attributes JSONB NOT NULL DEFAULT '{}', -- {"color": "red", "size": "XL"}

    -- Physical properties
    weight NUMERIC(10, 3), -- in kg
    dimensions JSONB, -- {"length": 10, "width": 5, "height": 2} in cm

    -- Status
    is_active BOOLEAN NOT NULL DEFAULT true,
    is_default BOOLEAN NOT NULL DEFAULT false, -- one variant should be default

    -- Tracking
    view_count INTEGER NOT NULL DEFAULT 0,
    sold_count INTEGER NOT NULL DEFAULT 0,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create storefront_product_variant_images table for variant-specific images
CREATE TABLE IF NOT EXISTS storefront_product_variant_images (
    id SERIAL PRIMARY KEY,
    variant_id INTEGER NOT NULL REFERENCES storefront_product_variants(id) ON DELETE CASCADE,
    image_url TEXT NOT NULL,
    thumbnail_url TEXT,
    alt_text VARCHAR(255),
    display_order INTEGER NOT NULL DEFAULT 0,
    is_main BOOLEAN NOT NULL DEFAULT false, -- main image for this variant
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_product_variant_attributes_name ON product_variant_attributes(name);
CREATE INDEX IF NOT EXISTS idx_product_variant_attribute_values_attribute_id ON product_variant_attribute_values(attribute_id);
CREATE INDEX IF NOT EXISTS idx_product_variant_attribute_values_value ON product_variant_attribute_values(value);

CREATE INDEX IF NOT EXISTS idx_storefront_product_variants_product_id ON storefront_product_variants(product_id);
CREATE INDEX IF NOT EXISTS idx_storefront_product_variants_sku ON storefront_product_variants(sku) WHERE sku IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_storefront_product_variants_barcode ON storefront_product_variants(barcode) WHERE barcode IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_storefront_product_variants_stock_status ON storefront_product_variants(stock_status);
CREATE INDEX IF NOT EXISTS idx_storefront_product_variants_is_active ON storefront_product_variants(is_active);
CREATE INDEX IF NOT EXISTS idx_storefront_product_variants_is_default ON storefront_product_variants(is_default);
CREATE INDEX IF NOT EXISTS idx_storefront_product_variants_attributes_gin ON storefront_product_variants USING gin(variant_attributes);

CREATE INDEX IF NOT EXISTS idx_storefront_product_variant_images_variant_id ON storefront_product_variant_images(variant_id);
CREATE INDEX IF NOT EXISTS idx_storefront_product_variant_images_is_main ON storefront_product_variant_images(is_main);

-- Add unique constraint to ensure only one default variant per product
CREATE UNIQUE INDEX IF NOT EXISTS idx_storefront_product_variants_default_unique
ON storefront_product_variants(product_id)
WHERE is_default = true;

-- Add trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_product_variants_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_product_variant_attributes_updated_at
    BEFORE UPDATE ON product_variant_attributes
    FOR EACH ROW
    EXECUTE FUNCTION update_product_variants_updated_at();

CREATE TRIGGER trigger_update_product_variant_attribute_values_updated_at
    BEFORE UPDATE ON product_variant_attribute_values
    FOR EACH ROW
    EXECUTE FUNCTION update_product_variants_updated_at();

CREATE TRIGGER trigger_update_storefront_product_variants_updated_at
    BEFORE UPDATE ON storefront_product_variants
    FOR EACH ROW
    EXECUTE FUNCTION update_product_variants_updated_at();
