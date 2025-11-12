-- Migration: Add Storefront Related Tables
-- Phase: Storefront Merge - Delivery to Listings
-- Date: 2025-11-12
--
-- Purpose: Add staff, hours, payment methods, and delivery options tables
--          to support full storefront functionality from delivery microservice
--
-- Related tables:
-- 1. storefront_staff - Staff members (owners, managers, cashiers)
-- 2. storefront_hours - Working hours (regular and special dates)
-- 3. storefront_payment_methods - Payment methods supported by storefront
-- 4. storefront_delivery_options - Delivery options configured by storefront

-- =====================================================
-- 1. STOREFRONT_STAFF (Staff Members)
-- =====================================================

CREATE TABLE IF NOT EXISTS storefront_staff (
    id SERIAL PRIMARY KEY,
    storefront_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    role VARCHAR(50) DEFAULT 'staff' NOT NULL,
    permissions JSONB DEFAULT '{}'::JSONB,
    last_active_at TIMESTAMP,
    actions_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    -- Foreign Keys
    CONSTRAINT fk_storefront_staff_storefront FOREIGN KEY (storefront_id)
        REFERENCES storefronts(id) ON DELETE CASCADE
);

-- Indexes for storefront_staff
CREATE INDEX IF NOT EXISTS idx_storefront_staff_storefront_id ON storefront_staff(storefront_id);
CREATE INDEX IF NOT EXISTS idx_storefront_staff_user_id ON storefront_staff(user_id);
CREATE INDEX IF NOT EXISTS idx_storefront_staff_role ON storefront_staff(role);
CREATE UNIQUE INDEX IF NOT EXISTS idx_storefront_staff_unique_user_per_store
    ON storefront_staff(storefront_id, user_id);

-- Update trigger for storefront_staff
CREATE OR REPLACE FUNCTION update_storefront_staff_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_storefront_staff_updated_at
BEFORE UPDATE ON storefront_staff
FOR EACH ROW
EXECUTE FUNCTION update_storefront_staff_updated_at();

-- Comments
COMMENT ON TABLE storefront_staff IS 'Storefront staff members (owners, managers, cashiers)';
COMMENT ON COLUMN storefront_staff.role IS 'Values: owner, manager, cashier, support, moderator';
COMMENT ON COLUMN storefront_staff.permissions IS 'JSON with detailed permissions';

-- =====================================================
-- 2. STOREFRONT_HOURS (Working Hours)
-- =====================================================

CREATE TABLE IF NOT EXISTS storefront_hours (
    id SERIAL PRIMARY KEY,
    storefront_id INTEGER NOT NULL,
    day_of_week INTEGER NOT NULL,  -- 0=Sunday, 6=Saturday
    open_time TIME,
    close_time TIME,
    is_closed BOOLEAN DEFAULT FALSE,
    special_date DATE,
    special_note VARCHAR(255),

    -- Foreign Keys
    CONSTRAINT fk_storefront_hours_storefront FOREIGN KEY (storefront_id)
        REFERENCES storefronts(id) ON DELETE CASCADE,

    -- Constraints
    CONSTRAINT storefront_hours_day_of_week_check
        CHECK (day_of_week >= 0 AND day_of_week <= 6)
);

-- Indexes for storefront_hours
CREATE INDEX IF NOT EXISTS idx_storefront_hours_storefront_id ON storefront_hours(storefront_id);
CREATE INDEX IF NOT EXISTS idx_storefront_hours_day_of_week ON storefront_hours(day_of_week);
CREATE INDEX IF NOT EXISTS idx_storefront_hours_special_date ON storefront_hours(special_date)
    WHERE special_date IS NOT NULL;

-- Comments
COMMENT ON TABLE storefront_hours IS 'Storefront working hours (regular and special dates)';
COMMENT ON COLUMN storefront_hours.day_of_week IS '0=Sunday, 1=Monday, ..., 6=Saturday';
COMMENT ON COLUMN storefront_hours.special_date IS 'For special hours on specific dates (holidays, etc)';
COMMENT ON COLUMN storefront_hours.is_closed IS 'If true, store is closed on this day/date';

-- =====================================================
-- 3. STOREFRONT_PAYMENT_METHODS (Payment Methods)
-- =====================================================

CREATE TABLE IF NOT EXISTS storefront_payment_methods (
    id SERIAL PRIMARY KEY,
    storefront_id INTEGER NOT NULL,
    method_type VARCHAR(50) NOT NULL,
    is_enabled BOOLEAN DEFAULT TRUE,
    provider VARCHAR(50),
    settings JSONB DEFAULT '{}'::JSONB,
    transaction_fee NUMERIC(5,2) DEFAULT 0.00,
    min_amount NUMERIC(10,2),
    max_amount NUMERIC(10,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    -- Foreign Keys
    CONSTRAINT fk_storefront_payment_methods_storefront FOREIGN KEY (storefront_id)
        REFERENCES storefronts(id) ON DELETE CASCADE
);

-- Indexes for storefront_payment_methods
CREATE INDEX IF NOT EXISTS idx_storefront_payment_methods_storefront_id
    ON storefront_payment_methods(storefront_id);
CREATE INDEX IF NOT EXISTS idx_storefront_payment_methods_method_type
    ON storefront_payment_methods(method_type);
CREATE INDEX IF NOT EXISTS idx_storefront_payment_methods_is_enabled
    ON storefront_payment_methods(is_enabled) WHERE is_enabled = TRUE;

-- Comments
COMMENT ON TABLE storefront_payment_methods IS 'Payment methods supported by storefront';
COMMENT ON COLUMN storefront_payment_methods.method_type IS 'Values: cash, cod, card, bank_transfer, paypal, crypto, postanska, keks_pay, ips';
COMMENT ON COLUMN storefront_payment_methods.transaction_fee IS 'Transaction fee in percent';
COMMENT ON COLUMN storefront_payment_methods.provider IS 'Payment provider name (e.g., stripe, paypal)';

-- =====================================================
-- 4. STOREFRONT_DELIVERY_OPTIONS (Delivery Options)
-- =====================================================

CREATE TABLE IF NOT EXISTS storefront_delivery_options (
    id SERIAL PRIMARY KEY,
    storefront_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    base_price NUMERIC(10,2) DEFAULT 0.00,
    price_per_km NUMERIC(10,2) DEFAULT 0.00,
    price_per_kg NUMERIC(10,2) DEFAULT 0.00,
    free_above_amount NUMERIC(10,2),
    min_order_amount NUMERIC(10,2),
    max_weight_kg NUMERIC(10,2),
    max_distance_km NUMERIC(10,2),
    estimated_days_min INTEGER DEFAULT 1,
    estimated_days_max INTEGER DEFAULT 3,
    zones JSONB DEFAULT '[]'::JSONB,
    available_days JSONB DEFAULT '[1, 2, 3, 4, 5]'::JSONB,
    cutoff_time TIME,
    provider VARCHAR(50),
    provider_config JSONB DEFAULT '{}'::JSONB,
    is_active BOOLEAN DEFAULT TRUE,
    display_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    -- Foreign Keys
    CONSTRAINT fk_storefront_delivery_options_storefront FOREIGN KEY (storefront_id)
        REFERENCES storefronts(id) ON DELETE CASCADE
);

-- Indexes for storefront_delivery_options
CREATE INDEX IF NOT EXISTS idx_storefront_delivery_options_storefront_id
    ON storefront_delivery_options(storefront_id);
CREATE INDEX IF NOT EXISTS idx_storefront_delivery_options_provider
    ON storefront_delivery_options(provider) WHERE provider IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_storefront_delivery_options_is_active
    ON storefront_delivery_options(is_active) WHERE is_active = TRUE;
CREATE INDEX IF NOT EXISTS idx_storefront_delivery_options_display_order
    ON storefront_delivery_options(storefront_id, display_order);

-- Update trigger for storefront_delivery_options
CREATE OR REPLACE FUNCTION update_storefront_delivery_options_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_storefront_delivery_options_updated_at
BEFORE UPDATE ON storefront_delivery_options
FOR EACH ROW
EXECUTE FUNCTION update_storefront_delivery_options_updated_at();

-- Comments
COMMENT ON TABLE storefront_delivery_options IS 'Delivery options configured by storefront';
COMMENT ON COLUMN storefront_delivery_options.provider IS 'Values: posta_srbije, aks, bex, d_express, city_express, self_pickup, own_delivery';
COMMENT ON COLUMN storefront_delivery_options.zones IS 'JSON array of delivery zones with pricing';
COMMENT ON COLUMN storefront_delivery_options.available_days IS 'JSON array of weekdays (1=Monday, 5=Friday)';
COMMENT ON COLUMN storefront_delivery_options.cutoff_time IS 'Order cutoff time for same-day delivery';
COMMENT ON COLUMN storefront_delivery_options.display_order IS 'Display order in frontend (lower = higher priority)';
