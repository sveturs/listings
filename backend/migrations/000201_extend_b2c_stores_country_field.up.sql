-- Extend country field to support full country names (not just ISO codes)
ALTER TABLE b2c_stores
    ALTER COLUMN country TYPE VARCHAR(100);

-- Update default value
ALTER TABLE b2c_stores
    ALTER COLUMN country SET DEFAULT 'RS';
