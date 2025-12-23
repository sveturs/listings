-- Rollback migration: Revert category_id from UUID back to bigint
-- WARNING: This will lose all category_id values!

BEGIN;

-- Step 1: Drop FK constraint
ALTER TABLE listings DROP CONSTRAINT IF EXISTS listings_category_id_fkey;

-- Step 2: Drop index
DROP INDEX IF EXISTS idx_listings_category_id;

-- Step 3: Add temporary bigint column
ALTER TABLE listings ADD COLUMN category_id_bigint BIGINT;

-- Step 4: Set all values to NULL (cannot convert UUID to bigint)
UPDATE listings SET category_id_bigint = NULL;

-- Step 5: Drop UUID column
ALTER TABLE listings DROP COLUMN category_id;

-- Step 6: Rename bigint column back
ALTER TABLE listings RENAME COLUMN category_id_bigint TO category_id;

-- Step 7: Recreate old check constraint
ALTER TABLE listings ADD CONSTRAINT listings_category_id_check CHECK (category_id > 0);

-- Step 8: Recreate old index
CREATE INDEX idx_listings_category_id ON listings(category_id) WHERE is_deleted = false;

-- Step 9: Remove comment
COMMENT ON COLUMN listings.category_id IS NULL;

COMMIT;
