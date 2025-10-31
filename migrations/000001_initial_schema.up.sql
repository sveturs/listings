-- Initial schema for listings microservice
-- This migration creates the core tables needed for managing marketplace listings

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Listings table (main entity)
CREATE TABLE IF NOT EXISTS listings (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID NOT NULL DEFAULT uuid_generate_v4() UNIQUE,

    -- User & ownership
    user_id BIGINT NOT NULL,
    storefront_id BIGINT,

    -- Core fields
    title VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(15,2) NOT NULL CHECK (price >= 0),
    currency VARCHAR(3) NOT NULL DEFAULT 'RSD',

    -- Categorization
    category_id BIGINT NOT NULL,

    -- Status & visibility
    status VARCHAR(50) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'active', 'inactive', 'sold', 'archived')),
    visibility VARCHAR(50) NOT NULL DEFAULT 'public' CHECK (visibility IN ('public', 'private', 'unlisted')),

    -- Inventory
    quantity INTEGER NOT NULL DEFAULT 1 CHECK (quantity >= 0),
    sku VARCHAR(100),

    -- Metadata
    views_count INTEGER NOT NULL DEFAULT 0,
    favorites_count INTEGER NOT NULL DEFAULT 0,

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    published_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,

    -- Soft delete support
    is_deleted BOOLEAN NOT NULL DEFAULT false
);

-- Indexes for listings
CREATE INDEX idx_listings_user_id ON listings(user_id) WHERE is_deleted = false;
CREATE INDEX idx_listings_storefront_id ON listings(storefront_id) WHERE is_deleted = false;
CREATE INDEX idx_listings_category_id ON listings(category_id) WHERE is_deleted = false;
CREATE INDEX idx_listings_status ON listings(status) WHERE is_deleted = false;
CREATE INDEX idx_listings_created_at ON listings(created_at DESC) WHERE is_deleted = false;
CREATE INDEX idx_listings_price ON listings(price) WHERE is_deleted = false;
CREATE INDEX idx_listings_uuid ON listings(uuid) WHERE is_deleted = false;

-- Full-text search index
CREATE INDEX idx_listings_search ON listings USING gin(to_tsvector('english', title || ' ' || COALESCE(description, '')));

-- Listing attributes (flexible key-value storage)
CREATE TABLE IF NOT EXISTS listing_attributes (
    id BIGSERIAL PRIMARY KEY,
    listing_id BIGINT NOT NULL REFERENCES listings(id) ON DELETE CASCADE,

    attribute_key VARCHAR(100) NOT NULL,
    attribute_value TEXT,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(listing_id, attribute_key)
);

CREATE INDEX idx_listing_attributes_listing_id ON listing_attributes(listing_id);
CREATE INDEX idx_listing_attributes_key ON listing_attributes(attribute_key);

-- Listing images
CREATE TABLE IF NOT EXISTS listing_images (
    id BIGSERIAL PRIMARY KEY,
    listing_id BIGINT NOT NULL REFERENCES listings(id) ON DELETE CASCADE,

    url TEXT NOT NULL,
    storage_path TEXT,
    thumbnail_url TEXT,

    display_order INTEGER NOT NULL DEFAULT 0,
    is_primary BOOLEAN NOT NULL DEFAULT false,

    width INTEGER,
    height INTEGER,
    file_size BIGINT,
    mime_type VARCHAR(100),

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_listing_images_listing_id ON listing_images(listing_id);
CREATE INDEX idx_listing_images_display_order ON listing_images(listing_id, display_order);

-- Listing tags
CREATE TABLE IF NOT EXISTS listing_tags (
    id BIGSERIAL PRIMARY KEY,
    listing_id BIGINT NOT NULL REFERENCES listings(id) ON DELETE CASCADE,
    tag VARCHAR(100) NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(listing_id, tag)
);

CREATE INDEX idx_listing_tags_listing_id ON listing_tags(listing_id);
CREATE INDEX idx_listing_tags_tag ON listing_tags(tag);

-- Listing location (optional)
CREATE TABLE IF NOT EXISTS listing_locations (
    id BIGSERIAL PRIMARY KEY,
    listing_id BIGINT NOT NULL UNIQUE REFERENCES listings(id) ON DELETE CASCADE,

    -- Address fields
    country VARCHAR(100),
    city VARCHAR(100),
    postal_code VARCHAR(20),
    address_line1 TEXT,
    address_line2 TEXT,

    -- Geo coordinates
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_listing_locations_listing_id ON listing_locations(listing_id);
CREATE INDEX idx_listing_locations_geo ON listing_locations(latitude, longitude) WHERE latitude IS NOT NULL AND longitude IS NOT NULL;

-- Listing statistics (cached aggregations)
CREATE TABLE IF NOT EXISTS listing_stats (
    listing_id BIGINT PRIMARY KEY REFERENCES listings(id) ON DELETE CASCADE,

    views_count INTEGER NOT NULL DEFAULT 0,
    favorites_count INTEGER NOT NULL DEFAULT 0,
    inquiries_count INTEGER NOT NULL DEFAULT 0,

    last_viewed_at TIMESTAMP WITH TIME ZONE,

    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexing queue (for async OpenSearch indexing)
CREATE TABLE IF NOT EXISTS indexing_queue (
    id BIGSERIAL PRIMARY KEY,
    listing_id BIGINT NOT NULL REFERENCES listings(id) ON DELETE CASCADE,

    operation VARCHAR(20) NOT NULL CHECK (operation IN ('index', 'update', 'delete')),
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'completed', 'failed')),

    retry_count INTEGER NOT NULL DEFAULT 0,
    max_retries INTEGER NOT NULL DEFAULT 3,

    error_message TEXT,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    processed_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_indexing_queue_status ON indexing_queue(status) WHERE status != 'completed';
CREATE INDEX idx_indexing_queue_listing_id ON indexing_queue(listing_id);
CREATE INDEX idx_indexing_queue_created_at ON indexing_queue(created_at);

-- Trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply trigger to tables
CREATE TRIGGER update_listings_updated_at BEFORE UPDATE ON listings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_listing_images_updated_at BEFORE UPDATE ON listing_images
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_listing_locations_updated_at BEFORE UPDATE ON listing_locations
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_listing_stats_updated_at BEFORE UPDATE ON listing_stats
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_indexing_queue_updated_at BEFORE UPDATE ON indexing_queue
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
