-- Rollback migration: Restore rudiment tables
-- NOTE: This is a data-loss migration. Rollback recreates empty tables only.

-- Restore address change log
CREATE TABLE IF NOT EXISTS address_change_log (
    id SERIAL PRIMARY KEY,
    entity_type VARCHAR(50),
    entity_id INTEGER,
    old_address TEXT,
    new_address TEXT,
    changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Restore GIS cache tables
CREATE TABLE IF NOT EXISTS gis_poi_cache (
    id SERIAL PRIMARY KEY,
    poi_type VARCHAR(100),
    lat DOUBLE PRECISION,
    lon DOUBLE PRECISION,
    data JSONB,
    cached_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS gis_filter_analytics (
    id SERIAL PRIMARY KEY,
    filter_type VARCHAR(100),
    usage_count INTEGER DEFAULT 0,
    last_used TIMESTAMP
);

CREATE TABLE IF NOT EXISTS gis_isochrone_cache (
    id SERIAL PRIMARY KEY,
    center_lat DOUBLE PRECISION,
    center_lon DOUBLE PRECISION,
    time_minutes INTEGER,
    polygon JSONB,
    cached_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Restore listings_geo (old table)
CREATE TABLE IF NOT EXISTS listings_geo (
    listing_id INTEGER PRIMARY KEY,
    geom GEOMETRY(Point, 4326),
    city VARCHAR(255),
    district VARCHAR(255)
);

-- Restore import sources
CREATE TABLE IF NOT EXISTS import_sources (
    id SERIAL PRIMARY KEY,
    source_name VARCHAR(255),
    source_type VARCHAR(50),
    config JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Restore custom UI tables
CREATE TABLE IF NOT EXISTS custom_ui_components (
    id SERIAL PRIMARY KEY,
    component_name VARCHAR(255),
    component_type VARCHAR(100),
    config JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS custom_ui_templates (
    id SERIAL PRIMARY KEY,
    template_name VARCHAR(255),
    template_html TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS custom_ui_component_usage (
    id SERIAL PRIMARY KEY,
    component_id INTEGER REFERENCES custom_ui_components(id),
    entity_type VARCHAR(50),
    entity_id INTEGER,
    used_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
