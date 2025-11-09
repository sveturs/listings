-- Migration: 000015_create_product_variants
-- Description: Create b2c_product_variants table for listings microservice
-- Date: 2025-11-09
-- Author: Phase 13.1.10
--
-- This migration creates the product variants table to support variant management
-- in the listings microservice. It mirrors the structure from the monolith.

-- Create sequence for product variants
CREATE SEQUENCE IF NOT EXISTS b2c_product_variants_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

-- Create b2c_product_variants table
CREATE TABLE IF NOT EXISTS b2c_product_variants (
    id INTEGER NOT NULL DEFAULT nextval('b2c_product_variants_id_seq'::regclass),
    product_id INTEGER NOT NULL,
    sku VARCHAR(100),
    barcode VARCHAR(100),
    price NUMERIC(15,2),
    compare_at_price NUMERIC(15,2),
    cost_price NUMERIC(15,2),
    stock_quantity INTEGER DEFAULT 0 NOT NULL,
    stock_status VARCHAR(20) DEFAULT 'in_stock' NOT NULL,
    low_stock_threshold INTEGER DEFAULT 5,
    variant_attributes JSONB DEFAULT '{}'::jsonb NOT NULL,
    weight NUMERIC(10,3),
    dimensions JSONB,
    is_active BOOLEAN DEFAULT true NOT NULL,
    is_default BOOLEAN DEFAULT false NOT NULL,
    view_count INTEGER DEFAULT 0 NOT NULL,
    sold_count INTEGER DEFAULT 0 NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,

    -- Constraints
    CONSTRAINT b2c_product_variants_pkey PRIMARY KEY (id),
    CONSTRAINT b2c_product_variants_stock_quantity_check CHECK (stock_quantity >= 0),
    CONSTRAINT b2c_product_variants_stock_status_check CHECK (stock_status IN ('in_stock', 'low_stock', 'out_of_stock'))
);

-- Add foreign key constraint
-- Note: We reference listings table (unified table in microservice)
ALTER TABLE b2c_product_variants
ADD CONSTRAINT fk_b2c_product_variants_product_id
FOREIGN KEY (product_id)
REFERENCES listings(id)
ON DELETE CASCADE
ON UPDATE CASCADE;

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_b2c_product_variants_product_id
    ON b2c_product_variants(product_id);

CREATE INDEX IF NOT EXISTS idx_b2c_product_variants_is_active
    ON b2c_product_variants(is_active);

CREATE INDEX IF NOT EXISTS idx_b2c_product_variants_stock_status
    ON b2c_product_variants(stock_status);

-- Add comment
COMMENT ON TABLE b2c_product_variants IS 'Product variants for B2C products - supports multiple SKUs, prices, and stock levels per product';
