-- Create storefront_products table
CREATE TABLE IF NOT EXISTS storefront_products (
    id SERIAL PRIMARY KEY,
    storefront_id INTEGER NOT NULL REFERENCES storefronts(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    price DECIMAL(10, 2) NOT NULL CHECK (price >= 0),
    currency CHAR(3) NOT NULL DEFAULT 'USD',
    category_id INTEGER NOT NULL REFERENCES marketplace_categories(id),
    sku VARCHAR(100),
    barcode VARCHAR(100),
    stock_quantity INTEGER NOT NULL DEFAULT 0 CHECK (stock_quantity >= 0),
    stock_status VARCHAR(20) NOT NULL DEFAULT 'in_stock' CHECK (stock_status IN ('in_stock', 'low_stock', 'out_of_stock')),
    is_active BOOLEAN NOT NULL DEFAULT true,
    attributes JSONB DEFAULT '{}',
    view_count INTEGER NOT NULL DEFAULT 0,
    sold_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for storefront_products
CREATE INDEX idx_storefront_products_storefront_id ON storefront_products(storefront_id);
CREATE INDEX idx_storefront_products_category_id ON storefront_products(category_id);
CREATE INDEX idx_storefront_products_stock_status ON storefront_products(stock_status);
CREATE INDEX idx_storefront_products_is_active ON storefront_products(is_active);
CREATE INDEX idx_storefront_products_sku ON storefront_products(sku) WHERE sku IS NOT NULL;
CREATE INDEX idx_storefront_products_barcode ON storefront_products(barcode) WHERE barcode IS NOT NULL;
CREATE INDEX idx_storefront_products_name_gin ON storefront_products USING gin(to_tsvector('simple', name));

-- Create storefront_product_images table
CREATE TABLE IF NOT EXISTS storefront_product_images (
    id SERIAL PRIMARY KEY,
    storefront_product_id INTEGER NOT NULL REFERENCES storefront_products(id) ON DELETE CASCADE,
    image_url TEXT NOT NULL,
    thumbnail_url TEXT NOT NULL,
    display_order INTEGER NOT NULL DEFAULT 0,
    is_default BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for storefront_product_images
CREATE INDEX idx_storefront_product_images_product_id ON storefront_product_images(storefront_product_id);
CREATE INDEX idx_storefront_product_images_display_order ON storefront_product_images(display_order);

-- Create storefront_product_variants table
CREATE TABLE IF NOT EXISTS storefront_product_variants (
    id SERIAL PRIMARY KEY,
    storefront_product_id INTEGER NOT NULL REFERENCES storefront_products(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    sku VARCHAR(100),
    price DECIMAL(10, 2) NOT NULL CHECK (price >= 0),
    stock_quantity INTEGER NOT NULL DEFAULT 0 CHECK (stock_quantity >= 0),
    attributes JSONB DEFAULT '{}',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for storefront_product_variants
CREATE INDEX idx_storefront_product_variants_product_id ON storefront_product_variants(storefront_product_id);
CREATE INDEX idx_storefront_product_variants_sku ON storefront_product_variants(sku) WHERE sku IS NOT NULL;
CREATE INDEX idx_storefront_product_variants_is_active ON storefront_product_variants(is_active);

-- Create storefront_inventory_movements table
CREATE TABLE IF NOT EXISTS storefront_inventory_movements (
    id SERIAL PRIMARY KEY,
    storefront_product_id INTEGER NOT NULL REFERENCES storefront_products(id) ON DELETE CASCADE,
    variant_id INTEGER REFERENCES storefront_product_variants(id) ON DELETE CASCADE,
    type VARCHAR(20) NOT NULL CHECK (type IN ('in', 'out', 'adjustment')),
    quantity INTEGER NOT NULL,
    reason VARCHAR(50) NOT NULL,
    order_id INTEGER,
    notes TEXT,
    user_id INTEGER NOT NULL REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for storefront_inventory_movements
CREATE INDEX idx_storefront_inventory_movements_product_id ON storefront_inventory_movements(storefront_product_id);
CREATE INDEX idx_storefront_inventory_movements_variant_id ON storefront_inventory_movements(variant_id) WHERE variant_id IS NOT NULL;
CREATE INDEX idx_storefront_inventory_movements_type ON storefront_inventory_movements(type);
CREATE INDEX idx_storefront_inventory_movements_created_at ON storefront_inventory_movements(created_at);

-- Add unique constraints using partial indexes
CREATE UNIQUE INDEX unique_storefront_product_sku ON storefront_products (storefront_id, sku) WHERE sku IS NOT NULL;
CREATE UNIQUE INDEX unique_storefront_product_barcode ON storefront_products (storefront_id, barcode) WHERE barcode IS NOT NULL;
CREATE UNIQUE INDEX unique_variant_sku ON storefront_product_variants (storefront_product_id, sku) WHERE sku IS NOT NULL;

-- Add trigger to update stock_status automatically
CREATE OR REPLACE FUNCTION update_product_stock_status()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.stock_quantity = 0 THEN
        NEW.stock_status = 'out_of_stock';
    ELSIF NEW.stock_quantity <= 5 THEN
        NEW.stock_status = 'low_stock';
    ELSE
        NEW.stock_status = 'in_stock';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_stock_status_trigger
BEFORE INSERT OR UPDATE OF stock_quantity ON storefront_products
FOR EACH ROW
EXECUTE FUNCTION update_product_stock_status();

-- Add trigger to update updated_at timestamp
CREATE TRIGGER update_storefront_products_updated_at
BEFORE UPDATE ON storefront_products
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_storefront_product_variants_updated_at
BEFORE UPDATE ON storefront_product_variants
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Add comments
COMMENT ON TABLE storefront_products IS 'Products sold in storefronts';
COMMENT ON TABLE storefront_product_images IS 'Images for storefront products';
COMMENT ON TABLE storefront_product_variants IS 'Product variants (e.g., sizes, colors)';
COMMENT ON TABLE storefront_inventory_movements IS 'Track inventory changes for audit';

COMMENT ON COLUMN storefront_products.stock_status IS 'Calculated based on stock_quantity: out_of_stock (0), low_stock (1-5), in_stock (>5)';
COMMENT ON COLUMN storefront_inventory_movements.type IS 'Type of movement: in (increase), out (decrease), adjustment';
COMMENT ON COLUMN storefront_inventory_movements.reason IS 'Reason for movement: sale, return, damage, restock, adjustment';