-- Create enum for reservation status if it doesn't exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'reservation_status') THEN
        CREATE TYPE reservation_status AS ENUM ('active', 'committed', 'released', 'expired');
    END IF;
END$$;

-- Create inventory_reservations table
CREATE TABLE IF NOT EXISTS inventory_reservations (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL,
    variant_id BIGINT NULL,
    order_id BIGINT NOT NULL,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    status reservation_status NOT NULL DEFAULT 'active',
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Add foreign key constraints
ALTER TABLE inventory_reservations 
    ADD CONSTRAINT fk_inventory_reservations_product_id 
    FOREIGN KEY (product_id) REFERENCES storefront_products(id) ON DELETE CASCADE;

ALTER TABLE inventory_reservations 
    ADD CONSTRAINT fk_inventory_reservations_variant_id 
    FOREIGN KEY (variant_id) REFERENCES storefront_product_variants(id) ON DELETE CASCADE;

ALTER TABLE inventory_reservations 
    ADD CONSTRAINT fk_inventory_reservations_order_id 
    FOREIGN KEY (order_id) REFERENCES storefront_orders(id) ON DELETE CASCADE;

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_inventory_reservations_product_id ON inventory_reservations(product_id);
CREATE INDEX IF NOT EXISTS idx_inventory_reservations_variant_id ON inventory_reservations(variant_id);
CREATE INDEX IF NOT EXISTS idx_inventory_reservations_order_id ON inventory_reservations(order_id);
CREATE INDEX IF NOT EXISTS idx_inventory_reservations_status ON inventory_reservations(status);
CREATE INDEX IF NOT EXISTS idx_inventory_reservations_expires_at ON inventory_reservations(expires_at);

-- Create update trigger for updated_at
CREATE OR REPLACE FUNCTION update_inventory_reservations_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER trigger_update_inventory_reservations_updated_at
    BEFORE UPDATE ON inventory_reservations
    FOR EACH ROW
    EXECUTE FUNCTION update_inventory_reservations_updated_at();