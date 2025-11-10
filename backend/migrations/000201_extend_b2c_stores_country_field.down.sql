-- Revert country field back to ISO codes only (VARCHAR(2))
-- WARNING: This will truncate data if any full country names are stored
ALTER TABLE b2c_stores
    ALTER COLUMN country TYPE VARCHAR(2);

-- Restore original default
ALTER TABLE b2c_stores
    ALTER COLUMN country SET DEFAULT 'RS';
