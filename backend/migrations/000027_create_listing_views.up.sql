-- Create table for tracking listing views
CREATE TABLE IF NOT EXISTS listing_views (
    id SERIAL PRIMARY KEY,
    listing_id INTEGER NOT NULL REFERENCES marketplace_listings(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    ip_hash VARCHAR(255),
    view_time TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    CONSTRAINT listing_view_uniqueness UNIQUE (listing_id, user_id),
    CONSTRAINT at_least_one_identifier CHECK (user_id IS NOT NULL OR ip_hash IS NOT NULL)
);

-- Add index for faster lookup by listing and user
CREATE INDEX IF NOT EXISTS idx_listing_views_listing_user ON listing_views(listing_id, user_id);

-- Add index for faster lookup by listing and IP
CREATE INDEX IF NOT EXISTS idx_listing_views_listing_ip ON listing_views(listing_id, ip_hash);

-- Add index for time-based queries
CREATE INDEX IF NOT EXISTS idx_listing_views_time ON listing_views(view_time);