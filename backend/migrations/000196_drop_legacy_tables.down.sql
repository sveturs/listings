-- Rollback migration: Restore legacy tables
-- WARNING: This rollback recreates table structure but WITHOUT data
-- Data is lost after UP migration. Use backup to restore data if needed.

-- Restore listing_attribute_values
CREATE TABLE IF NOT EXISTS listing_attribute_values (
    id SERIAL PRIMARY KEY,
    listing_id INTEGER,
    attribute_id INTEGER REFERENCES unified_attributes(id) ON DELETE CASCADE,
    value_text TEXT,
    value_int INTEGER,
    value_float FLOAT,
    value_bool BOOLEAN,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Restore districts backup tables (structure only, data lost)
CREATE TABLE IF NOT EXISTS districts_leskovac_backup (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    slug VARCHAR(255),
    city VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS districts_novi_sad_backup_20250715 (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    slug VARCHAR(255),
    city VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Note: To fully restore data, use database backup:
-- PGPASSWORD=mX3g1XGhMRUZEX3l pg_restore -h localhost -p 5433 -U postgres -d svetubd /tmp/backup_before_sprint_2.1_*.sql
